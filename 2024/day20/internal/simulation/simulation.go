package simulation

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

/////////////////////////////////////////////////////////////////////////////////////
// COORD & DIRECTIONS
/////////////////////////////////////////////////////////////////////////////////////

type Coord struct {
	X, Y int
}

func (c *Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

func (c *Coord) Move(d Direction) Coord {
	return Coord{X: c.X + d.VX, Y: c.Y + d.VY}
}

// GetNeighbors returns the neighbors of a coordinate
// Since this operation happens outside the context of a map, no validation
// is done to check if the neighbors are within expected bounds
func (c *Coord) GetNeighbors() []Coord {
	neighbors := []Coord{
		c.Move(North),
		c.Move(East),
		c.Move(South),
		c.Move(West),
	}

	return neighbors
}

// DirectionTo returns the direction from one coordinate to another
func (c *Coord) DirectionTo(coord Coord) Direction {
	dx := coord.X - c.X
	dy := coord.Y - c.Y
	if dx > 0 {
		return East
	} else if dx < 0 {
		return West
	} else if dy > 0 {
		return South
	} else if dy < 0 {
		return North
	}
	return Direction{VX: 0, VY: 0}
}

type Direction struct {
	VX, VY int
}

func (d Direction) String() string {
	// Turn the direction into a string
	switch d {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	}
	return "Unknown Direction" // Default case should not be reached if directions are correct
}

func (d Direction) TurnLeft() Direction {
	switch d {
	case North:
		return West
	case East:
		return North
	case South:
		return East
	case West:
		return South
	default:
		return d
	}
}

func (d Direction) TurnRight() Direction {
	switch d {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	default:
		return d
	}
}

// Define the four cardinal directions
var (
	North = Direction{VX: 0, VY: -1}
	East  = Direction{VX: 1, VY: 0}
	South = Direction{VX: 0, VY: 1}
	West  = Direction{VX: -1, VY: 0}
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
	addEntityId(entityId uuid.UUID) error
	removeEntityId(entityId uuid.UUID) error
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

func (c *spatialMapCell) addEntityId(entityId uuid.UUID) error {
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

func (c *spatialMapCell) removeEntityId(entityId uuid.UUID) error {
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
	GetCell(coord Coord) (SpatialMapCell, error)
	GetHeight() int
	GetWidth() int
	GetIndex(coord Coord) int
	GetNeighbors(coord Coord) []Coord
	removeEntity(entityId uuid.UUID, coords ...Coord) error // Mutation functions are private
	addEntity(entityId uuid.UUID, coords ...Coord) error    // Renamed from setEntity
	ValidateCoord(coord Coord) bool
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

func (m *spatialMap) GetIndex(coord Coord) int {
	return coord.Y*m.width + coord.X
}

func (m *spatialMap) GetCell(coord Coord) (SpatialMapCell, error) {
	if valid := m.ValidateCoord(coord); !valid {
		return nil, fmt.Errorf("coordinates %s are out of bounds", coord.String())
	}
	index := m.GetIndex(coord)
	return m.cells[index], nil
}

func (m *spatialMap) ValidateCoord(coord Coord) bool {
	// Check if the position is within bounds
	if coord.X < 0 || coord.X >= m.GetWidth() || coord.Y < 0 || coord.Y >= m.GetHeight() {
		return false
	}
	return true
}

func (m *spatialMap) addEntity(entityId uuid.UUID, coords ...Coord) error {
	var modifiedCells []SpatialMapCell
	var cellErrs []error = nil
	var rollbackErrs []error = nil
	for _, coord := range coords {
		cell, cellErr := m.GetCell(coord)
		if cellErr != nil {
			cellErrs = append(cellErrs, fmt.Errorf("error accessing cell at coordinates %s: %w", coord.String(), cellErr))
			break
		}
		cellErr = cell.addEntityId(entityId)
		if cellErr != nil {
			cellErrs = append(cellErrs, fmt.Errorf("error adding entity %s to cell %s: %w", entityId, coord.String(), cellErr))
			break
		}
		modifiedCells = append(modifiedCells, cell)
	}
	if cellErrs != nil {
		for _, cell := range modifiedCells {
			rollbackErr := cell.removeEntityId(entityId)
			if rollbackErr != nil {
				rollbackErrs = append(rollbackErrs, rollbackErr)
			}
		}
	}
	cellErr := errors.Join(cellErrs...)
	rollbackErr := errors.Join(rollbackErrs...)
	return errors.Join(cellErr, rollbackErr)
}

func (m *spatialMap) removeEntity(entityId uuid.UUID, coords ...Coord) error {
	var modifiedCells []SpatialMapCell
	var cellErrs []error = nil
	var rollbackErrs []error = nil
	for _, coord := range coords {
		cell, cellErr := m.GetCell(coord)
		if cellErr != nil {
			cellErrs = append(cellErrs, fmt.Errorf("error accessing cell at coordinates %s: %w", coord.String(), cellErr))
			break
		}
		cellErr = cell.removeEntityId(entityId)
		if cellErr != nil {
			cellErrs = append(cellErrs, fmt.Errorf("error removing entity %s from cell %s: %w", entityId, coord.String(), cellErr))
			break
		}
		modifiedCells = append(modifiedCells, cell)
	}
	if cellErrs != nil { // Rollback
		for _, cell := range modifiedCells {
			rollbackErr := cell.addEntityId(entityId)
			if rollbackErr != nil {
				rollbackErrs = append(rollbackErrs, rollbackErr)
			}
		}
	}
	cellErr := errors.Join(cellErrs...)
	rollbackErr := errors.Join(rollbackErrs...)
	return errors.Join(cellErr, rollbackErr)
}

func (m *spatialMap) GetNeighbors(coord Coord) []Coord {
	neighbors := []Coord{
		coord.Move(North),
		coord.Move(East),
		coord.Move(South),
		coord.Move(West),
	}

	validNeighbors := []Coord{}
	for _, neighbor := range neighbors {
		if m.ValidateCoord(neighbor) {
			validNeighbors = append(validNeighbors, neighbor)
		}
	}

	return validNeighbors
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
	GetPosition() (coords []Coord)
	GetDirection() (direction Direction)
	setPosition(coords ...Coord)      // Mutation functions are private
	setDirection(direction Direction) // Mutation functions are private
	Clone() Entity                    // Added Clone() method
}

type entity struct {
	id         uuid.UUID
	entityType string
	coords     []Coord
	direction  Direction
}

func NewEntity(entityType string) (Entity, error) {
	e := new(entity)
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	e.id = id
	e.entityType = entityType
	e.setPosition(Coord{X: 0, Y: 0})        // Default position is the origin
	e.setDirection(Direction{VX: 0, VY: 0}) // Default vector is zero
	return e, nil
}

func (e *entity) GetId() uuid.UUID {
	return e.id
}

func (e *entity) GetPosition() (coords []Coord) {
	coords = make([]Coord, len(e.coords))
	for i, coord := range e.coords {
		coords[i].X = coord.X
		coords[i].Y = coord.Y
	}
	return
}

func (e *entity) setPosition(coords ...Coord) {
	e.coords = make([]Coord, len(coords))
	for i, coord := range coords {
		e.coords[i].X = coord.X
		e.coords[i].Y = coord.Y
	}
}

func (e *entity) GetDirection() (direction Direction) {
	direction.VX = e.direction.VX
	direction.VY = e.direction.VY
	return
}

func (e *entity) setDirection(direction Direction) {
	e.direction.VX = direction.VX
	e.direction.VY = direction.VY
}

func (e *entity) GetEntityType() string {
	return e.entityType
}

// Clone method for entity
func (e *entity) clone() Entity {
	clone := new(entity)
	clone.id = e.id
	clone.entityType = e.entityType
	clone.setPosition(e.coords...)
	clone.setDirection(e.direction)
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
	AddEntity(e Entity, coords []Coord, direction Direction) (Entity, error)
	GetEntity(entityId uuid.UUID) (Entity, error)
	GetEntities() []Entity
	MoveEntity(entityId uuid.UUID, wrapping bool) error
	SetEntityDirection(entityId uuid.UUID, newDirection Direction) error
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

func (s *simulation) AddEntity(e Entity, coords []Coord, direction Direction) (Entity, error) {
	// Lock the mutex to ensure thread safety when adding entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Validate the coordinates before adding the entity
	for _, coord := range coords {
		if valid := s.spatialMap.ValidateCoord(coord); !valid {
			return nil, fmt.Errorf("invalid coordinates %s for new entity", coord.String())
		}
	}

	// Add the entity to the spatial map at the specified coordinates
	err := s.spatialMap.addEntity(e.GetId(), coords...)
	if err != nil {
		var coordStr string
		for _, coord := range coords {
			coordStr += coord.String()
		}
		return nil, fmt.Errorf("failed to add entity at coordinates %s", coordStr)
	}

	// Update the entity's position
	e.setPosition(coords...)

	// Update the entity's vector
	e.setDirection(direction)

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
	coords := entityToRemove.GetPosition()

	// Remove the entity from the spatial map
	err := s.spatialMap.removeEntity(entityId, coords...)
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

// Move the entity based on the entity location and direction
func (s *simulation) MoveEntity(entityId uuid.UUID, wrapping bool) error {
	// Lock the mutex to ensure thread safety when moving entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Get the map's dimensions
	width := s.spatialMap.GetWidth()
	height := s.spatialMap.GetHeight()

	entity, err := s.GetEntity(entityId)
	if err != nil {
		return fmt.Errorf("error accessing entity %s", entityId.String())
	}

	var newCoords []Coord
	for _, coord := range entity.GetPosition() {
		newCoord := Coord{coord.X, coord.Y}
		newCoord = newCoord.Move(entity.GetDirection())

		if wrapping {
			// Wrap the new coordinates so that when entities leave the map they re-enter on the other side.
			newCoord.X = ((newCoord.X % width) + width) % width
			newCoord.Y = ((newCoord.Y % height) + height) % height
		}
		newCoords = append(newCoords, newCoord)
	}

	// Validate the new coordinates
	for _, newCoord := range newCoords {
		if success := s.spatialMap.ValidateCoord(newCoord); !success {
			return fmt.Errorf("can not move entity to invalid coordinates %s", newCoord.String())
		}
	}

	// Find the index of the entity
	index, ok := s.entityMap[entityId]
	if !ok {
		fmt.Println(s.entityMap)
		return fmt.Errorf("entity %s not found", entityId)
	}

	// Get the current position of the entity
	currentCoords := s.entities[index].GetPosition()

	// Remove the entity from its current cell
	err = s.spatialMap.removeEntity(entityId, currentCoords...)
	if err != nil {
		return fmt.Errorf("failed to remove entity from current location: %v", err)
	}

	// Add the entity to the new cell
	err = s.spatialMap.addEntity(entityId, newCoords...)
	if err != nil {
		// Attempt to re-add the entity to its original cell in case of failure
		rollbackErr := s.spatialMap.addEntity(entityId, currentCoords...)
		if rollbackErr != nil {
			return fmt.Errorf("failed to move entity and failed to rollback: %v", rollbackErr)
		}
		return fmt.Errorf("failed to move entity and successfully rolled back")
	}

	// Update the entity's position
	s.entities[index].setPosition(newCoords...)

	return nil
}

func (s *simulation) SetEntityDirection(entityId uuid.UUID, newDirection Direction) error {
	// Lock the mutex to ensure thread safety when moving entities
	s.updateMutex.Lock()
	defer s.updateMutex.Unlock()

	// Find the index of the entity
	index, ok := s.entityMap[entityId]
	if !ok {
		return fmt.Errorf("entity %s not found", entityId)
	}

	// Update the entity's vector
	s.entities[index].setDirection(newDirection)

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
