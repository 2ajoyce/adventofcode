package simulation

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestSpatialMapCellClone(t *testing.T) {
	originalCell := NewSpatialMapCell()

	// Add valid entities
	entity1, _ := uuid.NewV7()
	entity2, _ := uuid.NewV7()

	if err := originalCell.addEntityId(entity1); err != nil {
		t.Fatalf("Failed to add entity1: %v", err)
	}
	if err := originalCell.addEntityId(entity2); err != nil {
		t.Fatalf("Failed to add entity2: %v", err)
	}

	// Clone the cell
	clonedCell := originalCell.Clone()

	// Add a new entity to the cloned cell
	entity3, _ := uuid.NewV7()
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
	entityId1, _ := uuid.NewV7()
	entityId2, _ := uuid.NewV7()
	firstCoord := Coord{X: 5, Y: 5}
	sm.addEntity(firstCoord, entityId1)
	sm.addEntity(firstCoord, entityId2)

	smClone := sm.Clone()

	newEntityId, _ := uuid.NewV7()
	secondCoord := Coord{X: 6, Y: 6}
	smClone.addEntity(secondCoord, newEntityId)

	cell, _ := sm.GetCell(secondCoord)
	if len(cell.GetEntityIds()) > 0 {
		t.Fatalf("Original map was modified after cloning!")
	}
	smCloneCell, _ := smClone.GetCell(firstCoord)
	if len(smCloneCell.GetEntityIds()) != 2 {
		t.Fatalf("Clone lost entity IDs!")
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
	entityId, _ := uuid.NewV7()
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

	entityId, _ := uuid.NewV7()
	cell.addEntityId(entityId)

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

	entity, _ := NewEntity("test_entity")
	coord := Coord{X: 5, Y: 5}

	// Add entity
	_, err := sim.AddEntity(entity, coord, North)
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

	entity, _ := NewEntity("test_entity")
	sim.AddEntity(entity, Coord{X: 0, Y: 0}, Direction{VX: 1, VY: 1})

	err := sim.MoveEntity(entity.GetId(), Coord{X: -1, Y: -1}, true)
	if err != nil {
		t.Fatalf("Failed to move entity with wrapping: %v", err)
	}

	coord := entity.GetPosition()
	if coord.X != 9 || coord.Y != 9 {
		t.Fatalf("Entity did not wrap correctly: got position %s", coord.String())
	}
}

func TestSimulationClone(t *testing.T) {
	sim := NewSimulation(10, 10)

	entity, _ := NewEntity("test_entity")
	sim.AddEntity(entity, Coord{X: 5, Y: 5}, North)

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
	newEntity, _ := NewEntity("new_entity")
	simClone.AddEntity(newEntity, Coord{X: 6, Y: 6}, South)

	_, err = sim.GetEntity(newEntity.GetId())
	if err == nil {
		t.Fatalf("Original simulation was affected by changes to the clone")
	}
}
