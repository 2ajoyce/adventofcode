package simulation

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

/////////////////////////////////////////////////////////////////////////////////////
// ERRORS
/////////////////////////////////////////////////////////////////////////////////////

type CellEmptyError struct{}

func (e CellEmptyError) Error() string {
	return "cell is empty"
}

/////////////////////////////////////////////////////////////////////////////////////
// SPATIAL MAP CELLS
/////////////////////////////////////////////////////////////////////////////////////

type SpatialMapCell interface {
	GetEntityIds() ([]uuid.UUID)
	AddEntityId(entityId uuid.UUID) (bool, error)
	RemoveEntityId(entityId uuid.UUID) (bool, error)
}

type spatialMapCell struct {
	entityIds []uuid.UUID
	mu        sync.RWMutex // To handle concurrent access
}

func NewSpatialMapCell() SpatialMapCell {
	c := new(spatialMapCell)
	c.entityIds = []uuid.UUID{}
	return c
}

func (c *spatialMapCell) GetEntityIds() ([]uuid.UUID) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to prevent external modification
	idsCopy := make([]uuid.UUID, len(c.entityIds))
	copy(idsCopy, c.entityIds)
	return idsCopy
}

func (c *spatialMapCell) AddEntityId(entityId uuid.UUID) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Check if the entity already exists in the cell
	for _, id := range c.entityIds {
		if id == entityId {
			return false, fmt.Errorf("entity %s already exists in the cell", entityId)
		}
	}
	c.entityIds = append(c.entityIds, entityId)
	return true, nil
}

func (c *spatialMapCell) RemoveEntityId(entityId uuid.UUID) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, id := range c.entityIds {
		if id == entityId {
			// Remove the entityId from the slice
			c.entityIds = append(c.entityIds[:i], c.entityIds[i+1:]...)
			return true, nil
		}
	}
	return false, fmt.Errorf("entity %s not found in the cell", entityId)
}

/////////////////////////////////////////////////////////////////////////////////////
// SPATIAL MAP
/////////////////////////////////////////////////////////////////////////////////////

type SpatialMap interface {
	GetCell(x, y int) (SpatialMapCell, error)
	GetHeight() int
	GetWidth() int
	GetIndex(x, y int) int
	removeEntity(x, y int, entityId uuid.UUID) (bool, error) // Mutation functions are private
	addEntity(x, y int, entityId uuid.UUID) (bool, error)    // Renamed from setEntity
	ValidateCoord(x, y int) bool
}

type spatialMap struct {
	width  int
	height int
	cells  []SpatialMapCell
}

func NewSpatialMap(width, height int) SpatialMap {
	m := new(spatialMap)
	m.width = width
	m.height = height
	m.cells = make([]SpatialMapCell, width*height)
	for i := range m.cells {
		m.cells[i] = NewSpatialMapCell()
	}
	return m
}

func (m *spatialMap) GetWidth() int {
	return m.width
}

func (m *spatialMap) GetHeight() int {
	return m.height
}

func (m *spatialMap) GetIndex(x, y int) int {
	return y*m.width + x
}

func (m *spatialMap) GetCell(x, y int) (SpatialMapCell, error) {
	if valid := m.ValidateCoord(x, y); !valid {
		return nil, fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}
	index := m.GetIndex(x, y)
	return m.cells[index], nil
}

func (m *spatialMap) ValidateCoord(x, y int) bool {
	// Check if the position is within bounds
	if x < 0 || x >= m.GetWidth() || y < 0 || y >= m.GetHeight() {
		return false
	}
	return true
}

func (m *spatialMap) addEntity(x, y int, entityId uuid.UUID) (bool, error) {
	cell, err := m.GetCell(x, y)
	if err != nil {
		return false, err
	}
	return cell.AddEntityId(entityId)
}

func (m *spatialMap) removeEntity(x, y int, entityId uuid.UUID) (bool, error) {
	cell, err := m.GetCell(x, y)
	if err != nil {
		return false, err
	}

	return cell.RemoveEntityId(entityId)
}

/////////////////////////////////////////////////////////////////////////////////////
// ENTITY
/////////////////////////////////////////////////////////////////////////////////////

type Entity interface {
	GetId() uuid.UUID
	GetPosition() (x int, y int)
	GetVelocity() (xv float64, yv float64)
	setPosition(x, y int)       // Mutation functions are private
	setVelocity(xv, yv float64) // Mutation functions are private
}

type entity struct {
	id uuid.UUID
	x  int
	y  int
	xv float64
	yv float64
}

func NewEntity() (Entity, error) {
	e := new(entity)
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	e.id = id
	return e, nil
}

func (e *entity) GetId() uuid.UUID {
	return e.id
}

func (e *entity) GetPosition() (x int, y int) {
	x = e.x
	y = e.y
	return
}

func (e *entity) setPosition(x, y int) {
	e.x = x
	e.y = y
	return
}

func (e *entity) GetVelocity() (xv float64, yv float64) {
	xv = e.xv
	yv = e.yv
	return
}

func (e *entity) setVelocity(xv, yv float64) {
	e.xv = xv
	e.yv = yv
}

/////////////////////////////////////////////////////////////////////////////////////
// SIMULATION
/////////////////////////////////////////////////////////////////////////////////////

type Simulation interface {
	AddEntity(e Entity, x, y int, xv, yv float64) (Entity, error)
	GetEntity(entityId uuid.UUID) (Entity, error)
	GetEntities() []Entity
	MoveEntity(entityId uuid.UUID, newX, newY int) (bool, error)
	RemoveEntity(entityId uuid.UUID) (bool, error)
	GetMap() SpatialMap
}

type simulation struct {
	updateMutex sync.Mutex
	spatialMap  SpatialMap
	entities    []Entity
	entityMap   map[uuid.UUID]int // EntityID -> index in Entities slice
}

func NewSimulation(width, height int) Simulation {
	s := new(simulation)
	s.spatialMap = NewSpatialMap(width, height)
	s.entities = make([]Entity, 0)
	s.entityMap = make(map[uuid.UUID]int)
	return s
}

func (s *simulation) GetMap() SpatialMap {
	return s.spatialMap
}

func (s *simulation) GetEntity(entityId uuid.UUID) (Entity, error) {
	// Get the entities index from the entityMap
	index, ok := s.entityMap[entityId]
	if !ok {
		return nil, fmt.Errorf("entity with ID %s not found", entityId)
	}
	return s.entities[index], nil
}

// GetEntities returns a copy of the current list of entities.
func (s *simulation) GetEntities() []Entity {
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()
	entitiesCopy := make([]Entity, len(s.entities))
	copy(entitiesCopy, s.entities)
	return entitiesCopy
}

func (s *simulation) AddEntity(e Entity, x, y int, xv, yv float64) (Entity, error) {
	// Lock the mutex to ensure thread safety when adding entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Validate the coordinates before adding the entity
	if valid := s.spatialMap.ValidateCoord(x, y); !valid {
		return nil, fmt.Errorf("invalid coordinates for new entity")
	}

	// Add the entity to the spatial map at the specified coordinates
	success, err := s.spatialMap.addEntity(x, y, e.GetId())
	if err != nil || !success {
		return nil, err
	}

	// Update the entity's position
	e.setPosition(x, y)

	// Update the entity's velocity
	e.setVelocity(xv, yv)

	// Add the entity to the slice and map
	s.entities = append(s.entities, e)
	s.entityMap[e.GetId()] = len(s.entities) - 1

	return e, nil
}

func (s *simulation) RemoveEntity(entityId uuid.UUID) (bool, error) {
	// Lock the mutex to ensure thread safety when removing entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Find the index of the entity in the slice
	index, exists := s.entityMap[entityId]
	if !exists {
		return false, fmt.Errorf("entity with ID %v not found", entityId)
	}

	// Get the entity from the slice using the index
	entityToRemove := s.entities[index]
	x, y := entityToRemove.GetPosition()

	// Remove the entity from the spatial map
	success, err := s.spatialMap.removeEntity(x, y, entityId)
	if err != nil || !success {
		return false, err
	}

	// Remove the entity from the slice and update the map
	lastIndex := len(s.entities) - 1
	if index != lastIndex {
		s.entities[index] = s.entities[lastIndex]
		s.entityMap[s.entities[lastIndex].GetId()] = index
	}
	delete(s.entityMap, entityId)
	s.entities = s.entities[:lastIndex]

	return true, nil
}

func (s *simulation) MoveEntity(entityId uuid.UUID, newX, newY int) (bool, error) {
	// Lock the mutex to ensure thread safety when moving entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Validate the new coordinates
	if success := s.spatialMap.ValidateCoord(newX, newY); !success {
		return false, fmt.Errorf("invalid coordinates for moving entity")
	}

	// Find the index of the entity
	index, ok := s.entityMap[entityId]
	if !ok {
		return false, fmt.Errorf("entity %s not found", entityId)
	}

	// Get the current position of the entity
	currentX, currentY := s.entities[index].GetPosition()

	// Remove the entity from its current cell
	success, err := s.spatialMap.removeEntity(currentX, currentY, entityId)
	if err != nil || !success {
		return false, err
	}

	// Add the entity to the new cell
	success, err = s.spatialMap.addEntity(newX, newY, entityId)
	if err != nil || !success {
		// Attempt to re-add the entity to its original cell in case of failure
		success, rollbackErr := s.spatialMap.addEntity(currentX, currentY, entityId)
		if !success || rollbackErr != nil {
			return false, fmt.Errorf("failed to move entity and failed to rollback: %v", rollbackErr)
		}
		return false, err
	}

	// Update the entity's position
	s.entities[index].setPosition(newX, newY)

	return true, nil
}
