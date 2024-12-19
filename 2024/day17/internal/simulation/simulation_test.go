package simulation

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestSpatialMapCellClone(t *testing.T) {
	originalCell := NewSpatialMapCell()

	// Add valid entities
	entity1, err := uuid.NewV7()
	if err != nil {
		t.Error("failed to create new uuid v7")
	}
	entity2, err := uuid.NewV7()
	if err != nil {
		t.Error("failed to create new uuid v7")
	}

	if err := originalCell.addEntityId(entity1); err != nil {
		t.Fatalf("Failed to add entity1: %v", err)
	}
	if err := originalCell.addEntityId(entity2); err != nil {
		t.Fatalf("Failed to add entity2: %v", err)
	}

	// Clone the cell
	clonedCell := originalCell.Clone()

	// Add a new entity to the cloned cell
	entity3, err := uuid.NewV7()
	if err != nil {
		t.Error("failed to create new uuid v7")
	}
	if err := clonedCell.addEntityId(entity3); err != nil {
		t.Fatalf("Failed to add entity3 to cloned cell: %v", err)
	}

	// Debug outputs for verification
	fmt.Printf("Original cell entity IDs: %v\n", originalCell.GetEntityIds())
	fmt.Printf("Cloned cell entity IDs: %v\n", clonedCell.GetEntityIds())

	// Validate that the original cell remains unchanged
	originalEntities := originalCell.GetEntityIds()
	clonedEntities := clonedCell.GetEntityIds()

	if len(originalEntities) != 2 {
		t.Fatalf("Original cell was modified: %v", originalEntities)
	}
	if originalEntities[0] != entity1 || originalEntities[1] != entity2 {
		t.Fatalf("Original cell content mismatch: %v", originalEntities)
	}
	if len(clonedEntities) != 3 {
		t.Fatalf("Cloned cell was not updated correctly: %v", clonedEntities)
	}
	if clonedEntities[0] != entity1 || clonedEntities[1] != entity2 || clonedEntities[2] != entity3 {
		t.Fatalf("Cloned cell content mismatch: %v", clonedEntities)
	}
}

func TestSpatialMapCloneIntegrity(t *testing.T) {
	sm := NewSpatialMap(10, 10)
	entityId1, err := uuid.NewV7()
	if err != nil {
		t.Error("failed to create new uuid v7")
	}
	entityId2, err := uuid.NewV7()
	if err != nil {
		t.Error("failed to create new uuid v7")
	}
	firstCoord := Coord{X: 5, Y: 5}

	if err := sm.addEntity(entityId1, firstCoord); err != nil {
		t.Fatalf("Failed to add entityId1: %v", err)
	}
	if err := sm.addEntity(entityId2, firstCoord); err != nil {
		t.Fatalf("Failed to add entityId2: %v", err)
	}

	smClone := sm.Clone()

	newEntityId, err := uuid.NewV7()
	if err != nil {
		t.Error("failed to create new uuid v7")
	}
	secondCoord := Coord{X: 6, Y: 6}
	if err := smClone.addEntity(newEntityId, secondCoord); err != nil {
		t.Fatalf("Failed to add newEntityId to cloned map: %v", err)
	}

	cell, err := sm.GetCell(secondCoord)
	if err != nil {
		t.Fatalf("Failed to get cell at secondCoord: %v", err)
	}
	if len(cell.GetEntityIds()) > 0 {
		t.Fatalf("Original map was modified after cloning!")
	}

	smCloneCell, err := smClone.GetCell(firstCoord)
	if err != nil {
		t.Fatalf("Failed to get cell at firstCoord in cloned map: %v", err)
	}
	if len(smCloneCell.GetEntityIds()) != 2 {
		t.Fatalf("Clone lost entity IDs! Expected 2, got %d", len(smCloneCell.GetEntityIds()))
	}
}

func TestAddAndRemoveEntityId(t *testing.T) {
	cell := NewSpatialMapCell()

	// Test adding a zero-value UUID
	zeroUUID := uuid.UUID{}
	err := cell.addEntityId(zeroUUID)
	if err == nil {
		t.Fatalf("Expected an error when adding a zero-value UUID, got none")
	}

	// Test adding and removing a valid UUID
	entityId, err := uuid.NewV7()
	if err != nil {
		t.Error("failed to create new uuid v7")
	}
	err = cell.addEntityId(entityId)
	if err != nil {
		t.Fatalf("Failed to add a valid UUID: %v", err)
	}

	err = cell.removeEntityId(entityId)
	if err != nil {
		t.Fatalf("Failed to remove a valid UUID: %v", err)
	}

	// Test removing a UUID not in the cell
	err = cell.removeEntityId(entityId)
	if err == nil {
		t.Fatalf("Expected an error when removing a non-existent UUID, got none")
	}
}

func TestIsEmptyAndGetEntityIds(t *testing.T) {
	cell := NewSpatialMapCell()

	if !cell.IsEmpty() {
		t.Fatalf("Expected cell to be empty, but it was not")
	}

	entityId, err := uuid.NewV7()
	if err != nil {
		t.Error("failed to create new uuid v7")
	}
	err = cell.addEntityId(entityId)
	if err != nil {
		t.Fatalf("Failed to add entityId: %v", err)
	}

	if cell.IsEmpty() {
		t.Fatalf("Expected cell to not be empty after adding an entity, but it was")
	}

	ids := cell.GetEntityIds()
	if len(ids) != 1 || ids[0] != entityId {
		t.Fatalf("GetEntityIds returned incorrect data: %v", ids)
	}
}

func TestSpatialMapValidateCoord(t *testing.T) {
	sm := NewSpatialMap(10, 10)

	tests := []struct {
		coord   Coord
		isValid bool
	}{
		{Coord{-1, 5}, false},
		{Coord{5, -1}, false},
		{Coord{10, 5}, false},
		{Coord{5, 10}, false},
		{Coord{0, 0}, true},
		{Coord{9, 9}, true},
	}

	for _, test := range tests {
		if sm.ValidateCoord(test.coord) != test.isValid {
			t.Fatalf("Validation failed for coordinates %s: expected %v", test.coord.String(), test.isValid)
		}
	}
}

func TestSimulationAddAndRemoveEntity(t *testing.T) {
	sim := NewSimulation(10, 10)

	entity, err := NewEntity("test_entity")
	if err != nil {
		t.Fatalf("Failed to create new entity: %v", err)
	}
	coord := Coord{X: 5, Y: 5}

	// Add entity
	_, err = sim.AddEntity(entity, []Coord{coord}, North)
	if err != nil {
		t.Fatalf("Failed to add entity: %v", err)
	}

	// Verify entity exists
	retrievedEntity, err := sim.GetEntity(entity.GetId())
	if err != nil || retrievedEntity == nil {
		t.Fatalf("Failed to retrieve added entity: %v", err)
	}

	// Remove entity
	err = sim.RemoveEntity(entity.GetId())
	if err != nil {
		t.Fatalf("Failed to remove entity: %v", err)
	}

	// Verify entity is removed
	_, err = sim.GetEntity(entity.GetId())
	if err == nil {
		t.Fatalf("Expected error when retrieving removed entity, but got none")
	}
}

func TestSimulationMoveEntityWithWrapping(t *testing.T) {
	sim := NewSimulation(10, 10)

	entity, err := NewEntity("test_entity")
	if err != nil {
		t.Fatalf("Failed to create new entity: %v", err)
	}
	_, err = sim.AddEntity(entity, []Coord{{X: 0, Y: 0}}, North)
	if err != nil {
		t.Fatalf("Failed to add entity: %v", err)
	}

	err = sim.MoveEntity(entity.GetId(), true)
	if err != nil {
		t.Fatalf("Failed to move entity with wrapping: %v", err)
	}

	coords := entity.GetPosition()
	if len(coords) == 0 {
		t.Fatalf("Entity has no position after moving")
	}
	coord := coords[0]
	if coord.X != 0 || coord.Y != 9 {
		t.Fatalf("Entity did not move correctly with wrapping: got position %s", coord.String())
	}

	// Move entity to trigger wrapping
	for i := 0; i < 10; i++ {
		err = sim.MoveEntity(entity.GetId(), true)
		if err != nil {
			t.Fatalf("Failed to move entity with wrapping: %v", err)
		}
	}

	coords = entity.GetPosition()
	if len(coords) == 0 {
		t.Fatalf("Entity has no position after moving multiple times")
	}
	coord = coords[0]
	expectedX := 0
	expectedY := 9
	if coord.X != expectedX || coord.Y != expectedY {
		t.Fatalf("Entity did not wrap correctly: expected (%d, %d), got (%d, %d)", expectedX, expectedY, coord.X, coord.Y)
	}
}

func TestSimulationClone(t *testing.T) {
	sim := NewSimulation(10, 10)

	entity, err := NewEntity("test_entity")
	if err != nil {
		t.Fatalf("Failed to create new entity: %v", err)
	}
	_, err = sim.AddEntity(entity, []Coord{Coord{X: 5, Y: 5}}, North)
	if err != nil {
		t.Fatalf("Failed to add entity: %v", err)
	}

	simClone := sim.Clone()

	// Verify cloned simulation has the entity
	cloneEntity, err := simClone.GetEntity(entity.GetId())
	if err != nil {
		t.Fatalf("Cloned simulation does not have the original entity: %v", err)
	}

	if cloneEntity.GetId() != entity.GetId() {
		t.Fatalf("Cloned entity ID mismatch: expected %v, got %v", entity.GetId(), cloneEntity.GetId())
	}

	// Verify original simulation is unaffected by changes to the clone
	newEntity, err := NewEntity("new_entity")
	if err != nil {
		t.Fatalf("Failed to create new entity: %v", err)
	}
	_, err = simClone.AddEntity(newEntity, []Coord{Coord{X: 6, Y: 6}}, South)
	if err != nil {
		t.Fatalf("Failed to add new entity to cloned simulation: %v", err)
	}

	_, err = sim.GetEntity(newEntity.GetId())
	if err == nil {
		t.Fatalf("Original simulation was affected by changes to the clone")
	}
}

func TestAddMultiCellEntity(t *testing.T) {
	sim := NewSimulation(10, 10)

	entity, err := NewEntity("multi_cell_entity")
	if err != nil {
		t.Fatalf("Failed to create new entity: %v", err)
	}
	coords := []Coord{
		{X: 2, Y: 2},
		{X: 3, Y: 2},
		{X: 2, Y: 3},
	}

	// Add entity
	_, err = sim.AddEntity(entity, coords, North)
	if err != nil {
		t.Fatalf("Failed to add multi-cell entity: %v", err)
	}

	// Verify entity is added to all specified cells
	for _, coord := range coords {
		cell, err := sim.GetMap().GetCell(coord)
		if err != nil {
			t.Fatalf("Failed to retrieve cell at %s: %v", coord.String(), err)
		}
		entityIds := cell.GetEntityIds()
		if len(entityIds) != 1 || entityIds[0] != entity.GetId() {
			t.Fatalf("Cell at %s does not contain the expected entity ID: %v", coord.String(), entityIds)
		}
	}
}

func TestMoveMultiCellEntity(t *testing.T) {
	sim := NewSimulation(10, 10)

	entity, err := NewEntity("multi_cell_entity")
	if err != nil {
		t.Fatalf("Failed to create new entity: %v", err)
	}
	initialCoords := []Coord{
		{X: 4, Y: 4},
		{X: 5, Y: 4},
		{X: 4, Y: 5},
	}

	// Add entity
	_, err = sim.AddEntity(entity, initialCoords, West)
	if err != nil {
		t.Fatalf("Failed to add multi-cell entity: %v", err)
	}

	// Move entity
	err = sim.MoveEntity(entity.GetId(), true)
	if err != nil {
		t.Fatalf("Failed to move multi-cell entity: %v", err)
	}
	err = sim.MoveEntity(entity.GetId(), true)
	if err != nil {
		t.Fatalf("Failed to move multi-cell entity: %v", err)
	}

	// Verify new positions
	newCoords := []Coord{
		{X: 2, Y: 4},
		{X: 3, Y: 4},
		{X: 2, Y: 5},
	}
	for _, coord := range newCoords {
		cell, err := sim.GetMap().GetCell(coord)
		if err != nil {
			t.Fatalf("Failed to retrieve cell at %s after moving: %v", coord.String(), err)
		}
		entityIds := cell.GetEntityIds()
		if len(entityIds) != 1 || entityIds[0] != entity.GetId() {
			t.Fatalf("Cell at %s does not contain the expected entity ID: %v", coord.String(), entityIds)
		}
	}

	// Verify old positions are cleared
	for _, coord := range initialCoords {
		cell, err := sim.GetMap().GetCell(coord)
		if err != nil {
			t.Fatalf("Failed to retrieve cell at %s: %v", coord.String(), err)
		}
		if len(cell.GetEntityIds()) != 0 {
			t.Fatalf("Cell at %s was not cleared after moving entity: %v", coord.String(), cell.GetEntityIds())
		}
	}
}

func TestRemoveMultiCellEntity(t *testing.T) {
	sim := NewSimulation(10, 10)

	entity, err := NewEntity("multi_cell_entity")
	if err != nil {
		t.Fatalf("Failed to create new entity: %v", err)
	}
	coords := []Coord{
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 1, Y: 2},
	}

	// Add entity
	_, err = sim.AddEntity(entity, coords, South)
	if err != nil {
		t.Fatalf("Failed to add multi-cell entity: %v", err)
	}

	// Remove entity
	err = sim.RemoveEntity(entity.GetId())
	if err != nil {
		t.Fatalf("Failed to remove multi-cell entity: %v", err)
	}

	// Verify all positions are cleared
	for _, coord := range coords {
		cell, err := sim.GetMap().GetCell(coord)
		if err != nil {
			t.Fatalf("Failed to retrieve cell at %s: %v", coord.String(), err)
		}
		if len(cell.GetEntityIds()) != 0 {
			t.Fatalf("Cell at %s was not cleared after removing entity: %v", coord.String(), cell.GetEntityIds())
		}
	}
}

func TestAddMultiCellEntityWithInvalidCoords(t *testing.T) {
	sim := NewSimulation(10, 10)

	entity, err := NewEntity("multi_cell_entity")
	if err != nil {
		t.Fatalf("Failed to create new entity: %v", err)
	}
	coords := []Coord{
		{X: 5, Y: 5},
		{X: -1, Y: 5}, // Invalid coordinate
		{X: 5, Y: -1}, // Invalid coordinate
	}

	// Attempt to add entity with invalid coordinates
	_, err = sim.AddEntity(entity, coords, North)
	if err == nil {
		t.Fatalf("Expected an error when adding entity with invalid coordinates, got none")
	}

	// Verify no cells contain the entity
	for _, coord := range coords {
		if sim.GetMap().ValidateCoord(coord) {
			cell, err := sim.GetMap().GetCell(coord)
			if err != nil {
				t.Fatalf("Failed to retrieve cell at %s: %v", coord.String(), err)
			}
			if len(cell.GetEntityIds()) != 0 {
				t.Fatalf("Cell at %s contains entity despite invalid addition: %v", coord.String(), cell.GetEntityIds())
			}
		}
	}
}
