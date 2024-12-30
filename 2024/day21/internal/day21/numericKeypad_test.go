package day21

import (
	"slices"
	"testing"
)

func TestNumericCalculateMovementBaseCases(t *testing.T) {
	testCases := []struct {
		name           string
		startingX      int
		startingY      int
		input          rune
		expectedOutput []string
	}{
		{name: "Test A", startingX: 2, startingY: 3, input: 'A', expectedOutput: []string{}},
		{name: "Test 0", startingX: 2, startingY: 3, input: '0', expectedOutput: []string{"<A"}},
		{name: "Test 1", startingX: 2, startingY: 3, input: '1', expectedOutput: []string{"^<<A", "<^<A"}},
		{name: "Test 2", startingX: 2, startingY: 3, input: '2', expectedOutput: []string{"^<A", "<^A"}},
		{name: "Test 3", startingX: 2, startingY: 3, input: '3', expectedOutput: []string{"^A"}},
		{name: "Test 4", startingX: 2, startingY: 3, input: '4', expectedOutput: []string{"^<<^A", "^<^<A", "^^<<A", "<^<^A", "<^^<A"}},
		{name: "Test 5", startingX: 2, startingY: 3, input: '5', expectedOutput: []string{"^^<A", "^<^A", "<^^A"}},
		{name: "Test 6", startingX: 2, startingY: 3, input: '6', expectedOutput: []string{"^^A"}},
		{name: "Test 7", startingX: 2, startingY: 3, input: '7', expectedOutput: []string{"^<^^<A", "^^<^<A", "^^^<<A", "<^^^<A", "^<<^^A", "^<^<^A", "^^<<^A", "<^<^^A", "<^^<^A"}},
		{name: "Test 8", startingX: 2, startingY: 3, input: '8', expectedOutput: []string{"<^^^A", "^<^^A", "^^<^A", "^^^<A"}},
		{name: "Test 9", startingX: 2, startingY: 3, input: '9', expectedOutput: []string{"^^^A"}},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			nk := NewNumericKeypad()
			nk.currentX = tc.startingX
			nk.currentY = tc.startingY
			output := nk.CalculateMovements(tc.input)
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

func TestNumericMoveBaseCases(t *testing.T) {
	testCases := []struct {
		name  string
		start Coord
		end   Coord
		input string
	}{
		{name: "Test ^", start: Coord{2, 2}, end: Coord{2, 1}, input: "^"},
		{name: "Test v", start: Coord{2, 2}, end: Coord{2, 3}, input: "v"},
		{name: "Test <", start: Coord{2, 2}, end: Coord{1, 2}, input: "<"},
		{name: "Test >", start: Coord{1, 2}, end: Coord{2, 2}, input: ">"},
		{name: "Test A", start: Coord{2, 2}, end: Coord{2, 2}, input: "A"},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			nk := NewNumericKeypad()
			nk.currentX = tc.start.X
			nk.currentY = tc.start.Y
			success := nk.Move(tc.input)
			if !success {
				t.Errorf("Failed to move to the target position")
				t.FailNow()
			}
			a := nk.GetCurrentPosition()
			if a != tc.end {
				t.Errorf("Expected %v, but got %v", tc.end, a)
			}
		})
	}
}

func TestPermutateSubstring(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{input: "^<A", expected: []string{"^<A", "<^A"}},
		{input: "><A", expected: []string{"><A", "<>A"}},
		{input: "^^A", expected: []string{"^^A"}},
		{input: "A", expected: []string{"A"}},
		{input: "", expected: []string{}},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.input, func(t *testing.T) {
			output := permutateSubstring(tc.input)
			if len(output) != len(tc.expected) {
				t.Errorf("Expected %d outputs, but got %d, Outputs: %v", len(tc.expected), len(output), output)
				t.FailNow()
			}
			for _, o := range output {
				if !slices.Contains(tc.expected, o) {
					t.Errorf("Expected output to contain %s, but got %s", tc.expected, output)
					t.FailNow()
				}
			}
		})
	}
}
