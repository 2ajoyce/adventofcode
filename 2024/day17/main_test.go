package main

import (
	"day17/internal/day17"
	"fmt"
	"math/big"
	"os"
	"testing"
)

const INPUT_FILE = "test_input.txt"
const OUTPUT_FILE = "test_output.txt"

func validateOutput(t *testing.T, expectedContent string) bool {
	data, err := os.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
		return false
	}
	content := string(data)
	if content != expectedContent {
		t.Errorf("Expected output to contain '%s', but got: %v", expectedContent, content)
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
	os.Unsetenv("INPUT_FILE")
	os.Unsetenv("OUTPUT_FILE")
	os.Unsetenv("PARALLELISM")
	os.Unsetenv("DEBUG")

	// Exit with the same status as `go test`
	os.Exit(code)
}

func TestSolveComputer1(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(729))
	comp.SetOpcodes([]day17.Opcode{0, 1, 5, 4, 3, 0})
	expectedOutput := "4,6,3,5,6,3,5,2,1,0"

	output, err := SolveComputer(0, comp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s' but got: '%s'", expectedOutput, output)
	}
}

func TestSolveComputer2(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(10))
	comp.SetOpcodes([]day17.Opcode{5, 0, 5, 1, 5, 4})
	expectedOutput := "0,1,2"

	output, err := SolveComputer(0, comp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s', but got: %s", expectedOutput, output)
	}
}

func TestSolveComputer3(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(2024))
	comp.SetOpcodes([]day17.Opcode{0, 1, 5, 4, 3, 0})
	expectedOutput := "4,2,5,6,7,7,7,7,3,1,0"

	output, err := SolveComputer(0, comp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s', but got: %s", expectedOutput, output)
	}
}

func TestSolveCanSolvePart2(t *testing.T) {
	// This test is checking the known example from the description to prove that the program
	// works as expected
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(117440))
	comp.SetOpcodes([]day17.Opcode{0, 3, 5, 4, 3, 0})
	expectedOutput := "0,3,5,4,3,0"

	output, err := SolveComputer(0, comp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s', but got: %s", expectedOutput, output)
	}
}
func compareRanges(t *testing.T, r, expected Range, i int) bool {
	if r.Start.Cmp(expected.Start) != 0 {
		t.Errorf("Expected range %d to start at %v, but got %v", i, expected.Start, r.Start)
		return false
	}
	if r.End.Cmp(expected.End) != 0 {
		t.Errorf("Expected range %d to end at %v, but got %v", i, expected.End, r.End)
		return false
	}
	if r.Index != expected.Index {
		t.Errorf("Expected range %d to have index %d, but got %d", i, expected.Index, r.Index)
		return false
	}
	if r.Match != expected.Match {
		t.Errorf("Expected range %d to have match %s, but got %s", i, expected.Match, r.Match)
		return false
	}
	if r.OutputLength != expected.OutputLength {
		t.Errorf("Expected range %d to have match %d, but got %d", i, expected.OutputLength, r.OutputLength)
		return false
	}
	return true
}

func TestFindRangesIndexOne(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(117440))
	opCodes := []day17.Opcode{0, 3, 5, 4, 3, 0}
	comp.SetOpcodes(opCodes)

	index := len(opCodes)
	match := fmt.Sprintf("%d", opCodes[index-1])
	searchSpace := Range{Start: big.NewInt(0), End: big.NewInt(262144 - 8), Index: index, Match: match, OutputLength: len(opCodes)}
	ranges := findRanges(comp, searchSpace)

	expectedRanges := []Range{
		{Start: searchSpace.Start, End: searchSpace.End, Index: searchSpace.Index - 1, Match: searchSpace.Match, OutputLength: len(opCodes)},
	}

	if len(ranges) != len(expectedRanges) {
		t.Errorf("Expected %d ranges, but got %d", len(expectedRanges), len(ranges))
	}

	for i, r := range ranges {
		compareRanges(t, r, expectedRanges[i], i)
	}
}

func TestFindRangesIndexTwo(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(117440))
	opCodes := []day17.Opcode{0, 3, 5, 4, 3, 0}
	comp.SetOpcodes(opCodes)

	index := len(opCodes) - 1
	// Convert the opcodes from the index to the end of the slice to a string
	match := ""
	for _, opcode := range opCodes[index-1:] {
		match += fmt.Sprintf("%d", opcode)
	}
	searchSpace := Range{Start: big.NewInt(0), End: big.NewInt(262136), Index: index, Match: match, OutputLength: len(opCodes)}
	ranges := findRanges(comp, searchSpace)

	expectedRanges := []Range{
		{Start: big.NewInt(98296), End: big.NewInt(163832), Index: searchSpace.Index - 1, Match: searchSpace.Match, OutputLength: len(opCodes)},
	}

	if len(ranges) != len(expectedRanges) {
		t.Errorf("Expected %d ranges, but got %d", len(expectedRanges), len(ranges))
	}

	for i, r := range ranges {
		compareRanges(t, r, expectedRanges[i], i)
	}
}

func TestFindRangesIndexThree(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(117440))
	opCodes := []day17.Opcode{0, 3, 5, 4, 3, 0}
	comp.SetOpcodes(opCodes)

	index := len(opCodes) - 2
	// Convert the opcodes from the index to the end of the slice to a string
	match := ""
	for _, opcode := range opCodes[index-1:] {
		match += fmt.Sprintf("%d", opcode)
	}
	searchSpace := Range{Start: big.NewInt(98296), End: big.NewInt(163832), Index: index, Match: match, OutputLength: len(opCodes)}
	ranges := findRanges(comp, searchSpace)

	expectedRanges := []Range{
		{Start: big.NewInt(114680), End: big.NewInt(122872), Index: searchSpace.Index - 1, Match: searchSpace.Match, OutputLength: len(opCodes)},
	}

	if len(ranges) != len(expectedRanges) {
		t.Errorf("Expected %d ranges, but got %d", len(expectedRanges), len(ranges))
	}

	for i, r := range ranges {
		compareRanges(t, r, expectedRanges[i], i)
	}
}

func TestFindRangesIndexFour(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(117440))
	opCodes := []day17.Opcode{0, 3, 5, 4, 3, 0}
	comp.SetOpcodes(opCodes)

	index := len(opCodes) - 3
	// Convert the opcodes from the index to the end of the slice to a string
	match := ""
	for _, opcode := range opCodes[index-1:] {
		match += fmt.Sprintf("%d", opcode)
	}
	searchSpace := Range{Start: big.NewInt(114680), End: big.NewInt(122872), Index: index, Match: match, OutputLength: len(opCodes)}
	ranges := findRanges(comp, searchSpace)

	expectedRanges := []Range{
		{Start: big.NewInt(117240), End: big.NewInt(118264), Index: searchSpace.Index - 1, Match: searchSpace.Match, OutputLength: len(opCodes)},
	}

	if len(ranges) != len(expectedRanges) {
		t.Errorf("Expected %d ranges, but got %d", len(expectedRanges), len(ranges))
	}

	for i, r := range ranges {
		compareRanges(t, r, expectedRanges[i], i)
	}
}

func TestFindRangesIndexFive(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(117440))
	opCodes := []day17.Opcode{0, 3, 5, 4, 3, 0}
	comp.SetOpcodes(opCodes)

	index := len(opCodes) - 4
	// Convert the opcodes from the index to the end of the slice to a string
	match := ""
	for _, opcode := range opCodes[index-1:] {
		match += fmt.Sprintf("%d", opcode)
	}
	searchSpace := Range{Start: big.NewInt(117240), End: big.NewInt(118264), Index: index, Match: match, OutputLength: len(opCodes)}
	ranges := findRanges(comp, searchSpace)

	expectedRanges := []Range{
		{Start: big.NewInt(117432), End: big.NewInt(117560), Index: searchSpace.Index - 1, Match: searchSpace.Match, OutputLength: len(opCodes)},
	}

	if len(ranges) != len(expectedRanges) {
		t.Errorf("Expected %d ranges, but got %d", len(expectedRanges), len(ranges))
	}

	for i, r := range ranges {
		compareRanges(t, r, expectedRanges[i], i)
	}
}

func TestFindRangesIndexSix(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(117440))
	opCodes := []day17.Opcode{0, 3, 5, 4, 3, 0}
	comp.SetOpcodes(opCodes)

	index := len(opCodes) - 5
	// Convert the opcodes from the index to the end of the slice to a string
	match := ""
	for _, opcode := range opCodes[index-1:] {
		match += fmt.Sprintf("%d", opcode)
	}
	searchSpace := Range{Start: big.NewInt(117432), End: big.NewInt(117560), Index: index, Match: match, OutputLength: len(opCodes)}
	ranges := findRanges(comp, searchSpace)

	expectedRanges := []Range{
		{Start: big.NewInt(117440), End: big.NewInt(117440), Index: searchSpace.Index - 1, Match: searchSpace.Match, OutputLength: len(opCodes)},
	}

	if len(ranges) != len(expectedRanges) {
		t.Errorf("Expected %d ranges, but got %d", len(expectedRanges), len(ranges))
	}

	for i, r := range ranges {
		compareRanges(t, r, expectedRanges[i], i)
	}
}

func TestFindRangesWithBFSTest(t *testing.T) {
	comp := day17.NewComputer()
	opCodes := []day17.Opcode{0, 3, 5, 4, 3, 0}
	comp.SetOpcodes(opCodes)
	index := len(opCodes)
	initialRange := Range{Start: big.NewInt(0), End: big.NewInt(262144 - 8), Index: index - 0, Match: "0", OutputLength: len(opCodes)}
	// Define expected results at each depth level
	expectedResults := [][]Range{
		{{Start: big.NewInt(0), End: big.NewInt(262136), Index: index - 1, Match: "0", OutputLength: len(opCodes)}},
		{{Start: big.NewInt(98296), End: big.NewInt(163832), Index: index - 2, Match: "30", OutputLength: len(opCodes)}},
		{{Start: big.NewInt(114680), End: big.NewInt(122872), Index: index - 3, Match: "430", OutputLength: len(opCodes)}},
		{{Start: big.NewInt(117240), End: big.NewInt(118264), Index: index - 4, Match: "5430", OutputLength: len(opCodes)}},
		{{Start: big.NewInt(117432), End: big.NewInt(117560), Index: index - 5, Match: "35430", OutputLength: len(opCodes)}},
		{{Start: big.NewInt(117440), End: big.NewInt(117440), Index: index - 6, Match: "035430", OutputLength: len(opCodes)}},
	}

	// Initialize the heap
	heap := []Range{initialRange}

	for _, expected := range expectedResults {
		// Pop the range from the heap
		r := heap[0]
		heap = heap[1:]

		r.Match = ""
		for _, opcode := range opCodes[r.Index-1:] {
			r.Match += fmt.Sprintf("%d", opcode)
		}

		t.Logf("Index: %d: Processing range %s", r.Index, r)

		// Find the ranges for the current range
		ranges := findRanges(comp, r)
		if len(ranges) != len(expected) {
			t.Errorf("Expected %d ranges, but got %d", len(expected), len(ranges))
			break
		}
		t.Logf("Index: %d: Found the correct number of ranges: %d", r.Index, len(ranges))

		success := true
		for j, e := range expected {
			success = compareRanges(t, ranges[j], e, j)
			if !success {
				success = false
			}
		}
		if !success {
			break
		}
		t.Logf("Index: %d: Found the correct ranges", r.Index)

		// Add the ranges to the heap
		heap = append(heap, ranges...)
	}
}

func TestFindRangesBFS(t *testing.T) {
	opCodes := []day17.Opcode{0, 3, 5, 4, 3, 0}

	results := findRangesBFS(opCodes)
	expectedResults := []*big.Int{big.NewInt(117440)}

	if len(results) != len(expectedResults) {
		t.Errorf("Expected %d results, but got %d", len(expectedResults), len(results))
	}

	for i, r := range results {
		if r.Cmp(expectedResults[i]) != 0 {
			t.Errorf("Expected result %d to be %v, but got %v", i, expectedResults[i], r)
		}
	}
}
