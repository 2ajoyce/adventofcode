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

	if err := originalCell.AddEntityId(entity1); err != nil {
		t.Fatalf("Failed to add entity1: %v", err)
	}
	if err := originalCell.AddEntityId(entity2); err != nil {
		t.Fatalf("Failed to add entity2: %v", err)
	}

	// Clone the cell
	clonedCell := originalCell.Clone()

	// Add a new entity to the cloned cell
	entity3, _ := uuid.NewV7()
	if err := clonedCell.AddEntityId(entity3); err != nil {
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

	sm.addEntity(5, 5, entityId1)
	sm.addEntity(5, 5, entityId2)

	smClone := sm.Clone()

	newEntityId, _ := uuid.NewV7()
	smClone.addEntity(6, 6, newEntityId)

	cell, _ := sm.GetCell(6, 6)
	if len(cell.GetEntityIds()) > 0 {
		t.Fatalf("Original map was modified after cloning!")
	}
	smCloneCell, _ := smClone.GetCell(5, 5)
	if len(smCloneCell.GetEntityIds()) != 2 {
		t.Fatalf("Clone lost entity IDs!")
	}
}
