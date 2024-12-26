package main

import (
	"day18/internal/aocUtils"
	"day18/internal/simulation"
	"os"
	"testing"
)

const INPUT_FILE = "test_input.txt"
const OUTPUT_FILE = "test_output.txt"

func validateOutput(t *testing.T, expectedOutput string) bool {
	output, err := aocUtils.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
	}

	if len(output) == 0 {
		t.Errorf("Expected output to contain '%s', but got an empty string", expectedOutput)
		return false
	}

	if len(output) > 1 {
		t.Errorf("Expected output to contain '%s', but got multiple lines", expectedOutput)
		return false
	}

	if output[0] != expectedOutput {
		t.Errorf("Expected output to contain '%s', but got: %s", expectedOutput, output[0])
		return false
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

func TestMainParseEmptyGrid(t *testing.T) {
	// *012
	// 0...
	// 1...
	// 2...
	input := []string{
		"3:3",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "No obstacle fully blocks the path"

	main()

	validateOutput(t, expectedOutput)
}

func TestMainParseFullGrid(t *testing.T) {
	// *012
	// 0###
	// 1###
	// 2###
	input := []string{
		"3:3",
		"2,2",
		"2,1",
		"1,2",
		"2,0",
		"0,2",
		"1,1",
		"1,0",
		"0,1",
		"0,0",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "(2, 2)" // No path

	main()

	validateOutput(t, expectedOutput)
}

func TestMainParseBaseCaseCenteredGrid(t *testing.T) {
	// *012
	// 0...
	// 1.#.
	// 2...
	input := []string{
		"3:3",
		"1,1",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "No obstacle fully blocks the path"

	main()

	validateOutput(t, expectedOutput)
}

func TestMainParseBaseCaseRightGrid(t *testing.T) {
	// *012
	// 0...
	// 1..#
	// 2...
	input := []string{
		"3:3",
		"2,1",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "No obstacle fully blocks the path"

	main()

	validateOutput(t, expectedOutput)
}

func TestMainParseBaseCaseLeftGrid(t *testing.T) {
	// *012
	// 0...
	// 1#..
	// 2...
	input := []string{
		"3:3",
		"0,1",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "No obstacle fully blocks the path"

	main()

	validateOutput(t, expectedOutput)
}

func TestMainParseBaseCaseOneWayGrid(t *testing.T) {
	// *012
	// 0...
	// 1.##
	// 2...
	input := []string{
		"3:3",
		"1,1",
		"2,1",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "No obstacle fully blocks the path"

	main()

	validateOutput(t, expectedOutput)
}

func TestMainPart1Example(t *testing.T) {
	input := []string{
		"7:7",
		"5,4",
		"4,2",
		"4,5",
		"3,0",
		"2,1",
		"6,3",
		"2,4",
		"1,5",
		"0,6",
		"3,3",
		"2,6",
		"5,1",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "No obstacle fully blocks the path"

	main()

	validateOutput(t, expectedOutput)
}

func TestMainPart2Example(t *testing.T) {
	input := []string{
		"7:7",
		"5,4",
		"4,2",
		"4,5",
		"3,0",
		"2,1",
		"6,3",
		"2,4",
		"1,5",
		"0,6",
		"3,3",
		"2,6",
		"5,1",
		"1,2",
		"5,5",
		"2,5",
		"6,5",
		"1,4",
		"0,4",
		"6,4",
		"1,1",
		"6,1",
		"1,0",
		"0,5",
		"1,6",
		"2,0",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "(6, 1)"

	main()

	validateOutput(t, expectedOutput)
}

func Test_makeGraphBaseCase(t *testing.T) {
	// Simple 2x2 grid with no obstacles
	// *01
	// 0..
	// 1..
	sim := simulation.NewSimulation(2, 2)
	graph, err := makeGraph(sim)
	if err != nil {
		t.Errorf("Failed to make graph: %v", err)
	}

	expectedResults := map[simulation.Coord]map[simulation.Coord]float64{
		{X: 0, Y: 0}: {
			simulation.Coord{X: 0, Y: 1}: 0,
			simulation.Coord{X: 1, Y: 0}: 0,
		},
		{X: 0, Y: 1}: {
			simulation.Coord{X: 0, Y: 0}: 0,
			simulation.Coord{X: 1, Y: 1}: 0,
		},
		{X: 1, Y: 0}: {
			simulation.Coord{X: 0, Y: 0}: 0,
			simulation.Coord{X: 1, Y: 1}: 0,
		},
		{X: 1, Y: 1}: {
			simulation.Coord{X: 0, Y: 1}: 0,
			simulation.Coord{X: 1, Y: 0}: 0,
		},
	}

	if len(graph) != len(expectedResults) {
		t.Errorf("Expected graph to have %d keys, but got %d", len(expectedResults), len(graph))
	}

	for key, value := range expectedResults {
		if _, ok := graph[key]; !ok {
			t.Errorf("Expected graph to contain key %v, but it was not found", key)
		}

		for neighbor, weight := range value {
			if _, ok := graph[key][neighbor]; !ok {
				t.Errorf("Expected graph to contain neighbor %v for key %v, but it was not found", neighbor, key)
			}

			if graph[key][neighbor] != weight {
				t.Errorf("Expected graph to contain weight %f for neighbor %v of key %v, but got %f", weight, neighbor, key, graph[key][neighbor])
			}
		}
	}
}

func Test_makeGraphWithEntity(t *testing.T) {
	// Simple 2x2 grid with one obstacle
	// *01
	// 0.#
	// 1..
	sim := simulation.NewSimulation(2, 2)
	entity, err := simulation.NewEntity("#")
	if err != nil {
		t.Errorf("Failed to create entity: %v", err)
	}
	sim.AddEntity(entity, []simulation.Coord{{X: 1, Y: 0}}, simulation.North)
	graph, err := makeGraph(sim)
	if err != nil {
		t.Errorf("Failed to make graph: %v", err)
	}

	expectedResults := map[simulation.Coord]map[simulation.Coord]float64{
		{X: 0, Y: 0}: {
			simulation.Coord{X: 0, Y: 1}: 0,
		},
		{X: 0, Y: 1}: {
			simulation.Coord{X: 0, Y: 0}: 0,
			simulation.Coord{X: 1, Y: 1}: 0,
		},
		{X: 1, Y: 1}: {
			simulation.Coord{X: 0, Y: 1}: 0,
		},
	}

	if len(graph) != len(expectedResults) {
		t.Errorf("Expected graph to have %d keys, but got %d", len(expectedResults), len(graph))
	}

	for key, value := range expectedResults {
		if _, ok := graph[key]; !ok {
			t.Errorf("Expected graph to contain key %v, but it was not found", key)
		}

		for neighbor, weight := range value {
			if _, ok := graph[key][neighbor]; !ok {
				t.Errorf("Expected graph to contain neighbor %v for key %v, but it was not found", neighbor, key)
			}

			if graph[key][neighbor] != weight {
				t.Errorf("Expected graph to contain weight %f for neighbor %v of key %v, but got %f", weight, neighbor, key, graph[key][neighbor])
			}
		}
	}
}
