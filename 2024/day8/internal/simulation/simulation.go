package simulation

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

/////////////////////////////////////////////////////////////////////////////////////
// ERRORS
/////////////////////////////////////////////////////////////////////////////////////

type CellOccupiedError struct {
	entityId uuid.UUID
}

func (e CellOccupiedError) Error() string {
	return fmt.Sprintf("cell is already occupied by %s", e.entityId)
}

// ///////////////////////////////////////////////////////////////////////////////////
// SPATIAL MAP CELLS
// ///////////////////////////////////////////////////////////////////////////////////
type SpatialMapCell interface {
	GetEntityId() uuid.UUID
	IsAntinode() bool
	setEntityId(entityId uuid.UUID) (bool, error)
	setAntinode(antinode bool) (bool, error)
}

type spatialMapCell struct {
	entityId uuid.UUID
	antinode bool
}

func NewSpatialMapCell() SpatialMapCell {
	return &spatialMapCell{
		entityId: uuid.Nil,
		antinode: false,
	}
}

func (c *spatialMapCell) GetEntityId() uuid.UUID {
	return c.entityId
}

func (c *spatialMapCell) IsAntinode() bool {
	return c.antinode
}

func (c *spatialMapCell) setEntityId(entityId uuid.UUID) (bool, error) {
	// Check if the cell is already occupied by another entity before setting a new one.
	if c.entityId != uuid.Nil {
		return false, CellOccupiedError{entityId: c.GetEntityId()}
	}
	c.entityId = entityId
	return true, nil
}

func (c *spatialMapCell) setAntinode(status bool) (bool, error) {
	if c.antinode == status {
		return false, nil
	}
	c.antinode = status
	return true, nil
}

// ///////////////////////////////////////////////////////////////////////////////////
// SPATIAL MAP
// ///////////////////////////////////////////////////////////////////////////////////
type SpatialMap interface {
	GetWidth() int
	GetHeight() int
	GetIndex(x, y int) int
	GetCell(x, y int) (SpatialMapCell, error)
	ValidateCoord(x, y int) (bool, error)
	setEntity(x, y int, entityId uuid.UUID) (bool, error)
	removeEntity(x, y int) (bool, error)
	setAntinode(x, y int, status bool) (bool, error)
}

type spatialMap struct {
	width  int
	height int
	cells  []SpatialMapCell
}

func NewSpatialMap(width, height int) SpatialMap {
	cells := make([]SpatialMapCell, width*height)
	for i := range cells {
		cells[i] = &spatialMapCell{
			entityId: uuid.Nil,
			antinode: false, // default antinode status
		}
	}
	return &spatialMap{
		width:  width,
		height: height,
		cells:  cells,
	}
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
	if valid, _ := m.ValidateCoord(x, y); !valid {
		return nil, fmt.Errorf("coordinates (%d, %d) are out of bounds", x, y)
	}
	index := m.GetIndex(x, y)
	return m.cells[index], nil
}

func (m *spatialMap) ValidateCoord(x, y int) (bool, error) {
	// Check if the position is within bounds
	if x < 0 || x >= m.GetWidth() || y < 0 || y >= m.GetHeight() {
		return false, fmt.Errorf("position (%d, %d) out of bounds", x, y)
	}
	return true, nil
}

func (m *spatialMap) setEntity(x, y int, entityId uuid.UUID) (bool, error) {
	cell, err := m.GetCell(x, y)
	if err != nil {
		return false, err
	}
	return cell.setEntityId(entityId)
}

func (m *spatialMap) removeEntity(x, y int) (bool, error) {
	cell, err := m.GetCell(x, y)
	if err != nil {
		return false, err
	}

	return cell.setEntityId(uuid.Nil)
}

func (m *spatialMap) setAntinode(x, y int, status bool) (bool, error) {
	cell, err := m.GetCell(x, y)
	if err != nil {
		return false, err
	}
	return cell.setAntinode(status)
}

// ///////////////////////////////////////////////////////////////////////////////////
// ENTITY
// ///////////////////////////////////////////////////////////////////////////////////
type Entity interface {
	GetId() uuid.UUID
	GetPosition() (int, int)
	setPosition(x, y int) (bool, error)
}

type entity struct {
	id uuid.UUID
	x  int
	y  int
}

func (e entity) GetId() uuid.UUID {
	return e.id
}

func (e entity) GetPosition() (int, int) {
	return e.x, e.y
}

func (e *entity) setPosition(x, y int) (bool, error) {
	e.x = x
	e.y = y
	return true, nil
}

// ///////////////////////////////////////////////////////////////////////////////////
// SIMULATION
// ///////////////////////////////////////////////////////////////////////////////////
type Simulation interface {
	GetMap() SpatialMap
	GetEntities() []Entity
	MoveEntity(entityId uuid.UUID, newX, newY int) (bool, error)
	AddEntity(x, y int, antinode bool) (Entity, error)
	RemoveEntity(entityId uuid.UUID) (bool, error)
}

type simulation struct {
	updateMutex sync.Mutex
	spatialMap  SpatialMap
	entities    []Entity
	entityMap   map[uuid.UUID]int // EntityID -> index in Entities slice
}

func NewSimulation(width, height int) Simulation {
	return &simulation{
		updateMutex: sync.Mutex{},
		spatialMap:  NewSpatialMap(width, height),
		entities:    make([]Entity, 0),
		entityMap:   make(map[uuid.UUID]int),
	}
}

func (s *simulation) GetMap() SpatialMap {
	return s.spatialMap
}

// GetEntities returns a copy of the current list of entities.
func (s *simulation) GetEntities() []Entity {
	return s.entities
}

func (s *simulation) AddEntity(x, y int, antinode bool) (Entity, error) {
	// Lock the mutex to ensure thread safety when adding entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Validate the coordinates before adding the entity
	valid, err := s.spatialMap.ValidateCoord(x, y)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid coordinates for new entity: %v", err)
	}

	// Create a new entity and add it to the simulation
	newId, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %v", err)
	}

	newEntity := &entity{id: newId, x: x, y: y}
	s.entities = append(s.entities, newEntity)

	// Add the entity to the spatial map at the specified coordinates
	success, err := s.spatialMap.setEntity(x, y, newEntity.id)
	if err != nil || !success {
		return nil, err
	}

	// Set the antinode status if specified
	success, err = s.spatialMap.setAntinode(x, y, antinode)
	if err != nil || !success {
		return nil, err
	}

	// Add the entity to the map for quick lookup
	s.entityMap[newEntity.id] = len(s.entities) - 1

	return newEntity, nil
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

	// Get the cell at the entity's position in the spatial map
	cell, err := s.spatialMap.GetCell(x, y)
	if err != nil {
		return false, err
	}

	// Check if the entity at those coordinates is the one we want to remove
	if entityId != entityToRemove.GetId() || entityToRemove.GetId() != cell.GetEntityId() {
		return false, fmt.Errorf("entity with ID %v not found at coordinates (%d, %d)", entityId, x, y)
	}

	// Remove the entity from the slice and update the map
	lastIndex := len(s.entities) - 1
	s.entities[index] = s.entities[lastIndex]
	delete(s.entityMap, s.entities[lastIndex].GetId())

	// Update the spatial map to remove the entity
	success, err := s.spatialMap.removeEntity(x, y)
	if !success || err != nil {
		return false, err
	}

	// Remove the last element from the slice
	s.entities = s.entities[:lastIndex]

	return true, nil
}

func (s *simulation) MoveEntity(entityId uuid.UUID, newX, newY int) (bool, error) {
	// Lock the mutex to ensure thread safety when moving entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Check if the new position is within bounds
	success, err := s.spatialMap.ValidateCoord(newX, newY)
	if err != nil || !success {
		return false, fmt.Errorf("invalid coordinates for moving entity: %v", err)
	}

	// Find the index of the entity
	index, ok := s.entityMap[entityId]
	if !ok {
		return false, fmt.Errorf("entity %d not found", entityId)
	}

	// Using the index, get the coords
	x, y := s.entities[index].GetPosition()

	// Remove the entity from the map
	success, err = s.spatialMap.removeEntity(x, y)
	if err != nil {
		return success, err
	}
	// Set the entity in the new position on the map
	success, err = s.spatialMap.setEntity(newX, newY, entityId)
	if err != nil {
		return success, err
	}
	s.entities[index].setPosition(newX, newY)
	return true, nil
}
