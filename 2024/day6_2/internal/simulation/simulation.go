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
	GetEntityIds() []uuid.UUID
	IsEmpty() bool
	AddEntityId(entityId uuid.UUID) error
	RemoveEntityId(entityId uuid.UUID) error
	Clone() SpatialMapCell
}

type spatialMapCell struct {
	entityIds []uuid.UUID
	mu        sync.RWMutex // To handle concurrent access
}

func NewSpatialMapCell() SpatialMapCell {
	return &spatialMapCell{
		entityIds: make([]uuid.UUID, 0),
	}
}

func (c *spatialMapCell) IsEmpty() bool {
	return len(c.entityIds) == 0
}

func (c *spatialMapCell) GetEntityIds() []uuid.UUID {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to prevent external modification
	idsCopy := make([]uuid.UUID, 0)
	idsCopy = append(idsCopy, c.entityIds...)
	return idsCopy
}

func (c *spatialMapCell) AddEntityId(entityId uuid.UUID) error {
	if entityId == (uuid.UUID{}) {
		return fmt.Errorf("cannot add zero-value UUID to cell")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, id := range c.entityIds {
		if id == entityId {
			return fmt.Errorf("entity %s already exists in the cell", entityId)
		}
	}

	c.entityIds = append(c.entityIds, entityId)
	return nil
}

func (c *spatialMapCell) RemoveEntityId(entityId uuid.UUID) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, id := range c.entityIds {
		if id == entityId {
			// Remove the entityId from the slice
			c.entityIds = append(c.entityIds[:i], c.entityIds[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("entity %s not found in the cell", entityId)
}

func (c *spatialMapCell) clone() SpatialMapCell {
	c.mu.RLock()
	defer c.mu.RUnlock()

	clonedEntityIds := append([]uuid.UUID(nil), c.entityIds...)

	return &spatialMapCell{
		entityIds: clonedEntityIds,
		mu:        sync.RWMutex{},
	}
}

func (c *spatialMapCell) Clone() SpatialMapCell {
	return c.clone()
}

/////////////////////////////////////////////////////////////////////////////////////
// SPATIAL MAP
/////////////////////////////////////////////////////////////////////////////////////

type SpatialMap interface {
	GetCell(x, y int) (SpatialMapCell, error)
	GetHeight() int
	GetWidth() int
	GetIndex(x, y int) int
	removeEntity(x, y int, entityId uuid.UUID) error // Mutation functions are private
	addEntity(x, y int, entityId uuid.UUID) error    // Renamed from setEntity
	ValidateCoord(x, y int) bool
	Clone() SpatialMap
}

type spatialMap struct {
	width  int
	height int
	cells  []SpatialMapCell
	mu     sync.RWMutex
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

func (m *spatialMap) addEntity(x, y int, entityId uuid.UUID) error {
	cell, err := m.GetCell(x, y)
	if err != nil {
		return fmt.Errorf("error accessing cell at coordinates (%d, %d): %w", x, y, err)
	}
	return cell.AddEntityId(entityId)
}

func (m *spatialMap) removeEntity(x, y int, entityId uuid.UUID) error {
	cell, err := m.GetCell(x, y)
	if err != nil {
		return fmt.Errorf("error accessing cell at coordinates (%d, %d): %w", x, y, err)
	}

	return cell.RemoveEntityId(entityId)
}

func (m *spatialMap) clone() *spatialMap {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mClone := &spatialMap{
		width:  m.width,
		height: m.height,
		cells:  make([]SpatialMapCell, len(m.cells)),
	}

	for i, cell := range m.cells {
		if cell != nil {
			clonedCell := cell.Clone() // Clone the individual cell
			if clonedCell == nil {
				panic(fmt.Sprintf("Cell %d clone failed; nil returned.\n", i))
			}
			mClone.cells[i] = clonedCell
		}
	}

	return mClone
}

// Implementing Clone() for spatialMap to satisfy SpatialMap interface
func (m *spatialMap) Clone() SpatialMap {
	return m.clone()
}

/////////////////////////////////////////////////////////////////////////////////////
// ENTITY
/////////////////////////////////////////////////////////////////////////////////////

type Entity interface {
	GetId() uuid.UUID
	GetEntityType() string
	GetPosition() (x int, y int)
	GetVector() (xv int, yv int)
	setPosition(x, y int) // Mutation functions are private
	setVector(xv, yv int) // Mutation functions are private
	Clone() Entity        // Added Clone() method
}

type entity struct {
	id         uuid.UUID
	entityType string
	x          int
	y          int
	xv         int
	yv         int
}

func NewEntity(entityType string) (Entity, error) {
	e := new(entity)
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	e.id = id
	e.entityType = entityType
	e.setPosition(0, 0) // Default position is the origin
	e.setVector(0, 0)   // Default vector is zero
	return e, nil
}

func (e *entity) GetId() uuid.UUID {
	return e.id
}

func (e *entity) GetPosition() (x, y int) {
	x = e.x
	y = e.y
	return
}

func (e *entity) setPosition(x, y int) {
	e.x = x
	e.y = y
}

type Vector struct {
	vx, vy int
}

func (e *entity) GetVector() (xv, yv int) {
	xv = e.xv
	yv = e.yv
	return
}

func (e *entity) setVector(xv, yv int) {
	e.xv = xv
	e.yv = yv
}

func (e *entity) GetEntityType() string {
	return e.entityType
}

// Clone method for entity
func (e *entity) clone() Entity {
	clone := new(entity)
	clone.id = e.id
	clone.entityType = e.entityType
	clone.setPosition(e.x, e.y)
	clone.setVector(e.xv, e.yv)
	return clone
}

// Implementing Clone() for entity to satisfy Entity interface
func (e *entity) Clone() Entity {
	return e.clone()
}

/////////////////////////////////////////////////////////////////////////////////////
// SIMULATION
/////////////////////////////////////////////////////////////////////////////////////

type Simulation interface {
	AddEntity(e Entity, x, y int, xv, yv int) (Entity, error)
	GetEntity(entityId uuid.UUID) (Entity, error)
	GetEntities() []Entity
	MoveEntity(entityId uuid.UUID, newX, newY int, wrapping bool) error
	SetEntityVector(entityId uuid.UUID, newXv, newVy int) error
	RemoveEntity(entityId uuid.UUID) error
	GetMap() SpatialMap
	Clone() Simulation
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
	for i, e := range s.entities {
		entitiesCopy[i] = e.Clone()
	}
	return entitiesCopy
}

func (s *simulation) AddEntity(e Entity, x, y int, xv, yv int) (Entity, error) {
	// Lock the mutex to ensure thread safety when adding entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Validate the coordinates before adding the entity
	if valid := s.spatialMap.ValidateCoord(x, y); !valid {
		return nil, fmt.Errorf("invalid coordinates for new entity")
	}

	// Add the entity to the spatial map at the specified coordinates
	err := s.spatialMap.addEntity(x, y, e.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to add entity at coordinates (%d, %d)", x, y)
	}

	// Update the entity's position
	e.setPosition(x, y)

	// Update the entity's vector
	e.setVector(xv, yv)

	// Add the entity to the slice and map
	s.entities = append(s.entities, e)
	s.entityMap[e.GetId()] = len(s.entities) - 1

	return e, nil
}

func (s *simulation) RemoveEntity(entityId uuid.UUID) error {
	// Lock the mutex to ensure thread safety when removing entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Find the index of the entity in the slice
	index, exists := s.entityMap[entityId]
	if !exists {
		return fmt.Errorf("entity with ID %v not found", entityId)
	}

	// Get the entity from the slice using the index
	entityToRemove := s.entities[index]
	x, y := entityToRemove.GetPosition()

	// Remove the entity from the spatial map
	err := s.spatialMap.removeEntity(x, y, entityId)
	if err != nil {
		return fmt.Errorf("error removing entity from spatial map: %v", err)
	}

	// Remove the entity from the slice and update the map
	lastIndex := len(s.entities) - 1
	if index != lastIndex {
		s.entities[index] = s.entities[lastIndex]
		s.entityMap[s.entities[lastIndex].GetId()] = index
	}
	delete(s.entityMap, entityId)
	s.entities = s.entities[:lastIndex]

	return nil
}

func (s *simulation) MoveEntity(entityId uuid.UUID, newX, newY int, wrapping bool) error {
	// Lock the mutex to ensure thread safety when moving entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Get the map's dimensions
	width := s.spatialMap.GetWidth()
	height := s.spatialMap.GetHeight()

	if wrapping {
		// Wrap the new coordinates so that when entities leave the map they re-enter on the other side.
		newX = ((newX % width) + width) % width
		newY = ((newY % height) + height) % height
	}

	// Validate the new coordinates
	if success := s.spatialMap.ValidateCoord(newX, newY); !success {
		return fmt.Errorf("invalid coordinates (%d, %d) for moving entity", newX, newY)
	}

	// Find the index of the entity
	index, ok := s.entityMap[entityId]
	if !ok {
		fmt.Println(s.entityMap)
		return fmt.Errorf("entity %s not found", entityId)
	}

	// Get the current position of the entity
	currentX, currentY := s.entities[index].GetPosition()

	// Remove the entity from its current cell
	err := s.spatialMap.removeEntity(currentX, currentY, entityId)
	if err != nil {
		fmt.Printf("Looking for: %s\n", entityId.String())
		fmt.Printf("EntityMap: %v\n", s.entityMap)
		fmt.Printf("Entities: [")
		for i := 0; i < len(s.entities); i++ {
			fmt.Printf("%v, ", s.entities[i].GetId().String())
		}
		fmt.Printf("]\n")
		cell, _ := s.spatialMap.GetCell(currentX, currentY)
		lastCell, _ := s.spatialMap.GetCell(7, 9)
		fmt.Printf("Entities in cell: %v\n", cell.GetEntityIds())
		fmt.Printf("Entities in cell at (7,9): %v\n", lastCell.GetEntityIds())
		return fmt.Errorf("failed to remove entity from current cell at (%d, %d): %v", currentX, currentY, err)
	}

	// Add the entity to the new cell
	err = s.spatialMap.addEntity(newX, newY, entityId)
	if err != nil {
		// Attempt to re-add the entity to its original cell in case of failure
		rollbackErr := s.spatialMap.addEntity(currentX, currentY, entityId)
		if rollbackErr != nil {
			return fmt.Errorf("failed to move entity and failed to rollback: %v", rollbackErr)
		}
		return fmt.Errorf("failed to move entity and successfully rolled back")
	}

	// Update the entity's position
	s.entities[index].setPosition(newX, newY)

	return nil
}

func (s *simulation) SetEntityVector(entityId uuid.UUID, newVx, newVy int) error {
	// Lock the mutex to ensure thread safety when moving entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Find the index of the entity
	index, ok := s.entityMap[entityId]
	if !ok {
		return fmt.Errorf("entity %s not found", entityId)
	}

	// Update the entity's vector
	s.entities[index].setVector(newVx, newVy)

	return nil
}

// Clone method for simulation
func (s *simulation) clone() Simulation {
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Create a new simulation instance
	sClone := new(simulation)

	// Clone the spatial map
	sClone.spatialMap = s.spatialMap.Clone()

	// Clone entities
	sClone.entities = make([]Entity, len(s.entities))
	sClone.entityMap = make(map[uuid.UUID]int)
	for i, e := range s.entities {
		eClone := e.Clone()
		sClone.entities[i] = eClone
		sClone.entityMap[eClone.GetId()] = i
	}

	return sClone
}

// Implementing Clone() for simulation to satisfy Simulation interface
func (s *simulation) Clone() Simulation {
	return s.clone()
}
