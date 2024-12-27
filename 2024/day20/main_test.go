package main

import (
	"day20/internal/aocUtils"
	"day20/internal/simulation"
	"os"
	"testing"
)

const INPUT_FILE = "test_input.txt"
const OUTPUT_FILE = "test_output.txt"

func validateOutput(t *testing.T, expectedOutput []string) bool {
	output, err := aocUtils.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
	}

	if len(output) == 0 {
		t.Errorf("Expected output to contain '%s', but got an empty string", expectedOutput)
		return false
	}

	if len(output) != len(expectedOutput) {
		t.Errorf("Expected output to be '%d' rows, but got '%d' rows", len(expectedOutput), len(output))
		return false
	}

	for i := range output {
		if output[i] != expectedOutput[i] {
			t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput[i], output[i])
		}
	}
	// If the validation fails, the input and output are retained for troubleshooting
	os.Remove(INPUT_FILE)
	os.Remove(OUTPUT_FILE)
	return true
}

func TestMain(m *testing.M) {
	// Set up environment variables here
	os.Setenv("INPUT_FILE", INPUT_FILE)
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	os.Setenv("PARALLELISM", "1")
	os.Setenv("DEBUG", "true")

	// Run all tests
	code := m.Run()

	// Clean up any resources if necessary
	// If the validation fails, the input and output are retained for troubleshooting
	os.Unsetenv("INPUT_FILE")
	os.Unsetenv("OUTPUT_FILE")
	os.Unsetenv("PARALLELISM")
	os.Unsetenv("DEBUG")

	// Exit with the same status as `go test`
	os.Exit(code)
}

func TestParseInputFullGrid(t *testing.T) {
	input := []string{
		"S##",
		"###",
		"##E",
	}

	path, err := parseLines(input)
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}

	if len(path) > 0 {
		t.Errorf("Expected path to have 0 elements, but got %d", len(path))
	}
}

func TestParseInputSimple(t *testing.T) {
	input := []string{
		"S.",
		"#E",
	}
	expectedPath := []simulation.Coord{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
	}

	path, err := parseLines(input)
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}

	if len(path) != len(expectedPath) {
		t.Errorf("Expected path to have %d elements, but got %d", len(expectedPath), len(path))
	}

	for i := range path {
		if path[i] != expectedPath[i] {
			t.Errorf("Expected path[%d] to be %v, but got %v", i, expectedPath[i], path[i])
		}
	}
}

func TestParseInputMedium(t *testing.T) {
	input := []string{
		"S...",
		"###.",
		"E.#.",
		"#...",
	}
	expectedPath := []simulation.Coord{
		{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 3, Y: 1}, {X: 3, Y: 2}, {X: 3, Y: 3}, {X: 2, Y: 3}, {X: 1, Y: 3}, {X: 1, Y: 2}, {X: 0, Y: 2},
	}

	path, err := parseLines(input)
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}

	if len(path) != len(expectedPath) {
		t.Errorf("Expected path to have 0 elements, but got %d", len(path))
	}

	for i := range path {
		if path[i] != expectedPath[i] {
			t.Errorf("Expected path[%d] to be %v, but got %v", i, expectedPath[i], path[i])
		}
	}
}

func TestParseInputExample(t *testing.T) {
	input := []string{
		"###############",
		"#...#...#.....#",
		"#.#.#.#.#.###.#",
		"#S#...#.#.#...#",
		"#######.#.#.###",
		"#######.#.#...#",
		"#######.#.###.#",
		"###..E#...#...#",
		"###.#######.###",
		"#...###...#...#",
		"#.#####.#.###.#",
		"#.#...#.#.#...#",
		"#.#.#.#.#.#.###",
		"#...#...#...###",
		"###############",
	}
	expectedPath := []simulation.Coord{
		{X: 1, Y: 3}, {X: 1, Y: 2}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 3, Y: 2}, {X: 3, Y: 3}, {X: 4, Y: 3}, {X: 5, Y: 3}, {X: 5, Y: 2}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1}, {X: 7, Y: 2}, {X: 7, Y: 3}, {X: 7, Y: 4}, {X: 7, Y: 5}, {X: 7, Y: 6}, {X: 7, Y: 7}, {X: 8, Y: 7}, {X: 9, Y: 7}, {X: 9, Y: 6}, {X: 9, Y: 5}, {X: 9, Y: 4}, {X: 9, Y: 3}, {X: 9, Y: 2}, {X: 9, Y: 1}, {X: 10, Y: 1}, {X: 11, Y: 1}, {X: 12, Y: 1}, {X: 13, Y: 1}, {X: 13, Y: 2}, {X: 13, Y: 3}, {X: 12, Y: 3}, {X: 11, Y: 3}, {X: 11, Y: 4}, {X: 11, Y: 5}, {X: 12, Y: 5}, {X: 13, Y: 5}, {X: 13, Y: 6}, {X: 13, Y: 7}, {X: 12, Y: 7}, {X: 11, Y: 7}, {X: 11, Y: 8}, {X: 11, Y: 9}, {X: 12, Y: 9}, {X: 13, Y: 9}, {X: 13, Y: 10}, {X: 13, Y: 11}, {X: 12, Y: 11}, {X: 11, Y: 11}, {X: 11, Y: 12}, {X: 11, Y: 13}, {X: 10, Y: 13}, {X: 9, Y: 13}, {X: 9, Y: 12}, {X: 9, Y: 11}, {X: 9, Y: 10}, {X: 9, Y: 9}, {X: 8, Y: 9}, {X: 7, Y: 9}, {X: 7, Y: 10}, {X: 7, Y: 11}, {X: 7, Y: 12}, {X: 7, Y: 13}, {X: 6, Y: 13}, {X: 5, Y: 13}, {X: 5, Y: 12}, {X: 5, Y: 11}, {X: 4, Y: 11}, {X: 3, Y: 11}, {X: 3, Y: 12}, {X: 3, Y: 13}, {X: 2, Y: 13}, {X: 1, Y: 13}, {X: 1, Y: 12}, {X: 1, Y: 11}, {X: 1, Y: 10}, {X: 1, Y: 9}, {X: 2, Y: 9}, {X: 3, Y: 9}, {X: 3, Y: 8}, {X: 3, Y: 7}, {X: 4, Y: 7}, {X: 5, Y: 7},
	}

	path, err := parseLines(input)
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}

	if len(path) != len(expectedPath) {
		t.Errorf("Expected path to have %d elements, but got %d", len(expectedPath), len(path))
	}

	for i := range path {
		if path[i] != expectedPath[i] {
			t.Errorf("Expected path[%d] to be %v, but got %v", i, expectedPath[i], path[i])
		}
	}
}

func TestMainSimple(t *testing.T) {
	input := []string{
		"S.",
		"#.",
		"E.",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	main()
	expectedOutput := []string{
		"Steps Saved, Count of Cheats",
		"2,1",
	}
	validateOutput(t, expectedOutput)
}

func TestMainExample(t *testing.T) {
	input := []string{
		"###############",
		"#...#...#.....#",
		"#.#.#.#.#.###.#",
		"#S#...#.#.#...#",
		"#######.#.#.###",
		"#######.#.#...#",
		"#######.#.###.#",
		"###..E#...#...#",
		"###.#######.###",
		"#...###...#...#",
		"#.#####.#.###.#",
		"#.#...#.#.#...#",
		"#.#.#.#.#.#.###",
		"#...#...#...###",
		"###############",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	main()
	expectedOutput := []string{
		"Steps Saved, Count of Cheats",
		"2,14",
		"4,14",
		"6,2",
		"8,4",
		"10,2",
		"12,3",
		"20,1",
		"36,1",
		"38,1",
		"40,1",
		"64,1",
	}
	validateOutput(t, expectedOutput)
}
