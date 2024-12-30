package day21

import (
	"slices"
	"testing"
)

func TestDirectionalCalculateMovementsBaseCases(t *testing.T) {
	testCases := []struct {
		name           string
		start          Coord
		input          rune
		expectedOutput []string
	}{
		{name: "Test A", start: Coord{2, 0}, input: 'A', expectedOutput: []string{}},
		{name: "Test ^", start: Coord{2, 0}, input: '^', expectedOutput: []string{"<A"}},
		{name: "Test <", start: Coord{2, 0}, input: '<', expectedOutput: []string{"v<<A", "<v<A"}},
		{name: "Test >", start: Coord{2, 0}, input: '>', expectedOutput: []string{"vA"}},
		{name: "Test v", start: Coord{2, 0}, input: 'v', expectedOutput: []string{"v<A", "<vA"}},
		{name: "Test special", start: Coord{0, 1}, input: 'A', expectedOutput: []string{">>^A", ">^>A"}},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			dk := NewDirectionalKeypad()
			dk.currentX = tc.start.X
			dk.currentY = tc.start.Y
			output := dk.CalculateMovements(tc.input)
			if len(output) != len(tc.expectedOutput) {
				t.Errorf("Expected %d outputs, but got %d, Outputs: %v", len(tc.expectedOutput), len(output), output)
				t.FailNow()
			}
			for _, o := range output {
				if !slices.Contains(tc.expectedOutput, o) {
					t.Errorf("Expected output to contain %s, but got %s", tc.expectedOutput, output)
					t.FailNow()
				}
			}
		})
	}
}

func TestDirectionalMoveBaseCases(t *testing.T) {
	testCases := []struct {
		name           string
		start          Coord
		end            Coord
		input          string
		expectedOutput string
	}{
		{name: "Test A", start: Coord{1, 1}, end: Coord{1, 1}, input: "A", expectedOutput: "A"},
		{name: "Test ^", start: Coord{1, 1}, end: Coord{1, 0}, input: "^", expectedOutput: "<A"},
		{name: "Test <", start: Coord{1, 1}, end: Coord{0, 1}, input: "<", expectedOutput: "v<<A"},
		{name: "Test >", start: Coord{1, 1}, end: Coord{2, 1}, input: ">", expectedOutput: "vA"},
		{name: "Test v", start: Coord{1, 1}, end: Coord{1, 2}, input: "v", expectedOutput: "v<A"},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			dk := NewDirectionalKeypad()
			dk.currentX = tc.start.X
			dk.currentY = tc.start.Y
			success := dk.Move(tc.input)
			if !success {
				t.Errorf("Failed to move to the target position")
				t.FailNow()
			}
			a := dk.GetCurrentPosition()
			if a != tc.end {
				t.Errorf("Expected %v, but got %v", tc.end, a)
			}
		})
	}
}
