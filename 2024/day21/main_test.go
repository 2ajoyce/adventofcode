package main

import (
	"day21/internal/aocUtils"
	"day21/internal/day21"
	"os"
	"strings"
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
		t.Errorf("Expected output to contain '%s', but got an empty string", strings.Join(expectedOutput, "\n"))
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

func TestMainExample(t *testing.T) {
	input := []string{
		"029A",
		"980A",
		"179A",
		"456A",
		"379A",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	main()
	expectedOutput := []string{
		"126384",
	}
	validateOutput(t, expectedOutput)
}

func TestMainExampleSmall(t *testing.T) {
	input := []string{
		"379A",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	main()
	expectedOutput := []string{
		"24256",
	}
	validateOutput(t, expectedOutput)
}

func TestComplexMachine(t *testing.T) {
	// 7 8 9
	// 4 5 6
	// 1 2 3
	// _ 0 A

	// _ ^ A
	// < v >
    
	//            3  
	//        ^   A         <
	//    <   A > A  v <<   A
	// <v<A>>^AvA^A<vA<AA>>^A
	// v<<A>>^AvA^A

	testCases := []struct {
		input          rune
		expectedOutput string
	}{
		{input: '3', expectedOutput: "<v<A>>^AvA^A<vA<AA>"},
		{input: '0', expectedOutput: "<vA<AA>>^AvAA<^A>A"},
		{input: '9', expectedOutput: "<v<A>>^AAAvA^A<vA<AA>>^AvAA<^A>A<"},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(string(tc.input), func(t *testing.T) {
			c := day21.Coord{X: 2, Y: 3} // Default starting position
			output := complexMachine1(c, tc.input)
			if output != tc.expectedOutput {
				t.Errorf("Expected output to be '%s', but got '%s'", tc.expectedOutput, output)
			}
		})
	}
}
