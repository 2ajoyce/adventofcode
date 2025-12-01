package main

import (
	"day21/internal/aocUtils"
	"day21/internal/day21"
	"fmt"
	"os"
	"slices"
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
	os.Setenv("DEPTH", "2")
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
	os.Unsetenv("DEPTH")
}

func TestMainExampleSmall(t *testing.T) {
	os.Setenv("DEPTH", "2")
	input := []string{
		"379A",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	main()
	expectedOutput := []string{
		"24256",
	}
	validateOutput(t, expectedOutput)
	os.Unsetenv("DEPTH")
}

func TestMainBaseCases(t *testing.T) {
	os.Setenv("DEPTH", "2")
	testCases := []struct {
		input          string
		expectedOutput string
	}{
		{input: "0A", expectedOutput: "0"},
		{input: "1A", expectedOutput: "48"},
		{input: "2A", expectedOutput: "76"},
		{input: "3A", expectedOutput: "84"},
		{input: "4A", expectedOutput: "200"},
		{input: "5A", expectedOutput: "200"},
		{input: "6A", expectedOutput: "180"},
		{input: "7A", expectedOutput: "364"},
		{input: "8A", expectedOutput: "336"},
		{input: "9A", expectedOutput: "288"},
	}
	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(string(tc.input), func(t *testing.T) {
			aocUtils.WriteToFile(INPUT_FILE, []string{tc.input})
			main()
			validateOutput(t, []string{tc.expectedOutput})
		})
	}
	os.Unsetenv("DEPTH")
}

func TestCalculateCost(t *testing.T) {
	testCases := []struct {
		code         string
		inputLen     int
		expectedCost int
	}{
		{code: "0A", inputLen: 5, expectedCost: 0},
		{code: "1A", inputLen: 5, expectedCost: 5},
		{code: "2A", inputLen: 5, expectedCost: 10},
		{code: "3A", inputLen: 5, expectedCost: 15},
		{code: "4A", inputLen: 5, expectedCost: 20},
		{code: "5A", inputLen: 5, expectedCost: 25},
		{code: "6A", inputLen: 5, expectedCost: 30},
		{code: "7A", inputLen: 5, expectedCost: 35},
		{code: "8A", inputLen: 5, expectedCost: 40},
		{code: "9A", inputLen: 5, expectedCost: 45},
		{code: "10A", inputLen: 5, expectedCost: 50},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.code, func(t *testing.T) {
			cost, _ := calculateCost(tc.code, tc.inputLen)
			if cost != tc.expectedCost {
				t.Errorf("Expected cost: %d, got: %d", tc.expectedCost, cost)
			}
		})
	}
}

func TestGenerateOptimalNumericValuesForCoordSimple(t *testing.T) {
	depth := 1
	input := '9'
	expectedOutput := "v<<A>>^AAAvA^A"
	optimalValueMap := generateOptimalDirectionalValues(depth) // This is non-deterministic, which is why we run it multiple times

	c := day21.Coord{X: 2, Y: 3} // Default starting position
	output := generateOptimalNumericValuesForCoord(optimalValueMap, c, input, depth)
	if output != expectedOutput {
		t.Errorf("Expected output to be %s, but got %s", expectedOutput, output)
	}
}

func TestGenerateOptimalNumericValuesForCoord(t *testing.T) {
	depth := 1
	// This test case is slow, but it fully captures the inputs returning multiple outputs when run from the initial position
	testCases := []struct {
		input           rune
		possibleOutputs []string
	}{
		{input: '0', possibleOutputs: []string{"<vA<AA>>^AvAA<^A>A", "<vA<AA>>^AvAA^<A>A", "v<A<AA>>^AvAA<^A>A", "v<A<AA>>^AvAA^<A>A"}},
		{input: '1', possibleOutputs: []string{"v<<A>>^A<vA<A>>^AAvAA<^A>A", "v<<A>>^Av<A<A>>^AAvAA<^A>A", "v<<A>>^A<vA<A>>^AAvAA^<A>A", "v<<A>>^Av<A<A>>^AAvAA^<A>A"}},
		{input: '2', possibleOutputs: []string{"v<A<AA>>^AvA<^A>AvA^A", "<vA<AA>>^AvA^<A>AvA^A", "<vA<AA>>^AvA<^A>AvA^A", "v<A<AA>>^AvA^<A>AvA^A"}},
		{input: '3', possibleOutputs: []string{"v<<A>>^AvA^A"}},
		{input: '4', possibleOutputs: []string{"v<<A>>^AAv<A<A>>^AAvAA<^A>A", "v<<A>>^AA<vA<A>>^AAvAA<^A>A", "v<<A>>^AA<vA<A>>^AAvAA^<A>A", "v<<A>>^AAv<A<A>>^AAvAA^<A>A"}},
		{input: '5', possibleOutputs: []string{"<vA<AA>>^AvA^<A>AAvA^A", "<vA<AA>>^AvA<^A>AAvA^A", "v<A<AA>>^AvA<^A>AAvA^A", "v<A<AA>>^AvA^<A>AAvA^A"}},
		{input: '6', possibleOutputs: []string{"v<<A>>^AAvA^A"}},
		{input: '7', possibleOutputs: []string{"v<<A>>^AAA<vA<A>>^AAvAA<^A>A", "v<<A>>^AAAv<A<A>>^AAvAA<^A>A", "v<<A>>^AAA<vA<A>>^AAvAA^<A>A", "v<<A>>^AAAv<A<A>>^AAvAA^<A>A"}},
		{input: '8', possibleOutputs: []string{"v<A<AA>>^AvA<^A>AAAvA^A", "<vA<AA>>^AvA<^A>AAAvA^A", "<vA<AA>>^AvA^<A>AAAvA^A", "v<A<AA>>^AvA^<A>AAAvA^A"}},
		{input: '9', possibleOutputs: []string{"v<<A>>^AAAvA^A"}},
		{input: 'A', possibleOutputs: []string{""}},
	}
	outputsSeen := make(map[rune]map[string]int)

	for i := range 100 {
		optimalValueMap := generateOptimalDirectionalValues(depth) // This is non-deterministic, which is why we run it multiple times
		for _, tc := range testCases {
			tc := tc // capture range variable

			// Define the test function
			testName := fmt.Sprintf("Input: %s, Iteration: %d", string(tc.input), i)
			t.Run(testName, func(t *testing.T) {
				c := day21.Coord{X: 2, Y: 3} // Default starting position
				output := generateOptimalNumericValuesForCoord(optimalValueMap, c, tc.input, depth)
				if !slices.Contains(tc.possibleOutputs, output) {
					t.Errorf("Expected output to contain %s, but got %s", tc.possibleOutputs, output)
				}
				if _, ok := outputsSeen[tc.input]; !ok {
					outputsSeen[tc.input] = make(map[string]int)
				}
				outputsSeen[tc.input][output]++
			})
		}
	}

	// Check that the count of outputs seen is the same as the count of expected outputs for each test case
	for _, tc := range testCases {
		tc := tc // capture range variable
		outputsSeenCount := len(outputsSeen[tc.input])
		if outputsSeenCount != len(tc.possibleOutputs) {
			e := fmt.Sprintf("\nTest Case: %s\n", string(tc.input))
			e += fmt.Sprintf("    Expected %d unique outputs, but got %d\n", len(tc.possibleOutputs), outputsSeenCount)
			for k, v := range outputsSeen[tc.input] {
				e += fmt.Sprintf("        %s: %d\n", k, v)
			}
			t.Error(e)
		}

	}
}

func TestGenerateOptimalDirectionalValuesBaseCase(t *testing.T) {
	depth := 4
	symbols := []rune{'<', '>', '^', 'v', 'A'}
	// Keypad layout:
	// _ ^ A
	// < v >

	aCoord := day21.Coord{X: 2, Y: 0} // This test will only validate movements from A to symbols

	output := generateOptimalDirectionalValues(depth)
	for _, o := range output[aCoord] {
		fmt.Printf("Output: %v\n", o)
	}
	// Starting with depth 1, insert the output into the first keypad
	for _, symbol := range symbols {
		for d := 1; d <= depth; d++ {
			subOutput := output[aCoord][symbol][d]
			fmt.Printf("Depth %d, Symbol %s: SubOutput: %s\n", d, string(symbol), subOutput)
			dk := day21.NewDirectionalKeypad()
			for d2 := d; d2 > 0; d2-- {
				subResult := ""
				for _, c := range subOutput {
					if c == 'A' {
						p := dk.GetCurrentPosition()
						for _, s := range symbols {
							if p == dk.GetPosition(s) {
								subResult += string(s)
								break
							}
						}
					}
					f := dk.GetCurrentPosition()
					dk.Move(string(c)) // This doesn't move to the rune. It moves BY the rune.
					t := dk.GetCurrentPosition()
					fmt.Printf("    Depth %d, Symbol %s: From %v to %v: %c\n", d, string(symbol), f, t, c)
					// fmt.Printf("Depth %d, Symbol %s: SubResult: %s\n", d, string(symbol), subResult)
				}
				if d2 == 1 {
					fmt.Printf("Depth %d:%d, Symbol %s: Result: %v\n", d, d2, string(symbol), subResult)
					if subResult != string(symbol) {
						t.Errorf("Depth %d:%d, Symbol %s Expected output to be %s, but got %s", d, d2, string(symbol), string(symbol), subResult)
					}
				}
				subOutput = subResult
			}
		}
	}
}

func TestGenerateDirectionalValuesForCoordBaseCase(t *testing.T) {
	depth := 0
	testCases := []struct {
		coord           day21.Coord
		input           rune
		possibleOutputs []string
	}{
		{coord: day21.Coord{X: 2, Y: 0}, input: '^', possibleOutputs: []string{"<A"}},
		{coord: day21.Coord{X: 2, Y: 0}, input: '<', possibleOutputs: []string{"<v<A", "v<<A"}},
		{coord: day21.Coord{X: 2, Y: 0}, input: 'v', possibleOutputs: []string{"<vA", "v<A"}},
		{coord: day21.Coord{X: 2, Y: 0}, input: '>', possibleOutputs: []string{"vA"}},
		{coord: day21.Coord{X: 2, Y: 0}, input: 'A', possibleOutputs: []string{"A"}},
	}

	outputsSeen := make(map[rune]map[string]int)
	for i := range 20 {
		t.Logf("Test %d\n", i)
		for _, tc := range testCases {
			tc := tc // capture range variable
			t.Run(fmt.Sprintf("Coord: (%d, %d), Input: %c", tc.coord.X, tc.coord.Y, tc.input), func(t *testing.T) {
				output := generateDirectionalValuesForCoord(tc.coord, tc.input)[depth]
				if !slices.Contains(tc.possibleOutputs, output) {
					t.Errorf("Expected output to contain %s, but got %s", tc.possibleOutputs, output)
				}
				if _, ok := outputsSeen[tc.input]; !ok {
					outputsSeen[tc.input] = make(map[string]int)
				}
				outputsSeen[tc.input][output]++
			})
		}
	}
	// Check that the count of outputs seen is the same as the count of expected outputs for each test case
	for _, tc := range testCases {
		tc := tc // capture range variable
		outputsSeenCount := len(outputsSeen[tc.input])
		if outputsSeenCount != len(tc.possibleOutputs) {
			e := fmt.Sprintf("\nTest Case: %s\n", string(tc.input))
			e += fmt.Sprintf("    Expected %d unique outputs, but got %d\n", len(tc.possibleOutputs), outputsSeenCount)
			for k, v := range outputsSeen[tc.input] {
				e += fmt.Sprintf("        %s: %d\n", k, v)
			}
			t.Error(e)
		}

	}
}

func TestGenerateDirectionalValuesForCoordSecondLocation(t *testing.T) {
	testCases := []struct {
		coord          day21.Coord
		input          rune
		expectedOutput []string
	}{
		{coord: day21.Coord{X: 0, Y: 1}, input: '^', expectedOutput: []string{">^A"}},
		{coord: day21.Coord{X: 0, Y: 1}, input: '<', expectedOutput: []string{"A"}},
		{coord: day21.Coord{X: 0, Y: 1}, input: 'v', expectedOutput: []string{">A"}},
		{coord: day21.Coord{X: 0, Y: 1}, input: '>', expectedOutput: []string{">>A"}},
		{coord: day21.Coord{X: 0, Y: 1}, input: 'A', expectedOutput: []string{">>^A", ">^>A"}},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(fmt.Sprintf("Coord: (%d, %d), Input: %c", tc.coord.X, tc.coord.Y, tc.input), func(t *testing.T) {
			output := generateDirectionalValuesForCoord(tc.coord, tc.input)
			if len(output) != len(tc.expectedOutput) {
				t.Errorf("Expected output to be %d rows, but got %d rows", len(tc.expectedOutput), len(output))
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

func TestGenerateOptimalDirectionalValuesForCoordBaseCase(t *testing.T) {
	depth := 1
	testCases := []struct {
		coord           day21.Coord
		input           rune
		possibleOutputs []string
	}{
		{coord: day21.Coord{X: 2, Y: 0}, input: '^', possibleOutputs: []string{"<A"}},
		{coord: day21.Coord{X: 2, Y: 0}, input: '<', possibleOutputs: []string{"<v<A"}}, // This value is ONLY valid at depth 1
		{coord: day21.Coord{X: 2, Y: 0}, input: 'v', possibleOutputs: []string{"<vA", "v<A"}},
		{coord: day21.Coord{X: 2, Y: 0}, input: '>', possibleOutputs: []string{"vA"}},
		{coord: day21.Coord{X: 2, Y: 0}, input: 'A', possibleOutputs: []string{"A"}},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(fmt.Sprintf("Coord: (%d, %d), Input: %c", tc.coord.X, tc.coord.Y, tc.input), func(t *testing.T) {
			output := generateOptimalDirectionalValuesForCoord(tc.coord, tc.input, depth)[depth]
			if !slices.Contains(tc.possibleOutputs, output) {
				t.Errorf("Expected output to contain %s, but got %s", tc.possibleOutputs, output)
			}
		})
	}

}

func TestGenerateOptimalDirectionalValuesForCoordSecondLocation(t *testing.T) {
	depth := 1
	testCases := []struct {
		coord           day21.Coord
		input           rune
		possibleOutputs []string
	}{
		{coord: day21.Coord{X: 0, Y: 1}, input: '^', possibleOutputs: []string{">^A"}},
		{coord: day21.Coord{X: 0, Y: 1}, input: '<', possibleOutputs: []string{"A"}},
		{coord: day21.Coord{X: 0, Y: 1}, input: 'v', possibleOutputs: []string{">A"}},
		{coord: day21.Coord{X: 0, Y: 1}, input: '>', possibleOutputs: []string{">>A"}},
		{coord: day21.Coord{X: 0, Y: 1}, input: 'A', possibleOutputs: []string{">>^A"}},
	}

	outputsSeen := make(map[rune]map[string]int)
	for i := range 20 {
		t.Logf("Test %d\n", i)
		for _, tc := range testCases {
			tc := tc // capture range variable
			t.Run(fmt.Sprintf("Coord: (%d, %d), Input: %c", tc.coord.X, tc.coord.Y, tc.input), func(t *testing.T) {
				output := generateOptimalDirectionalValuesForCoord(tc.coord, tc.input, depth)[depth]
				if !slices.Contains(tc.possibleOutputs, output) {
					t.Errorf("Expected output to contain %s, but got %s", tc.possibleOutputs, output)
				}
				if _, ok := outputsSeen[tc.input]; !ok {
					outputsSeen[tc.input] = make(map[string]int)
				}
				outputsSeen[tc.input][output]++
			})
		}
	}
	// Check that the count of outputs seen is the same as the count of expected outputs for each test case
	for _, tc := range testCases {
		tc := tc // capture range variable
		outputsSeenCount := len(outputsSeen[tc.input])
		if outputsSeenCount != len(tc.possibleOutputs) {
			e := fmt.Sprintf("\nTest Case: %s\n", string(tc.input))
			e += fmt.Sprintf("    Expected %d unique outputs, but got %d\n", len(tc.possibleOutputs), outputsSeenCount)
			for k, v := range outputsSeen[tc.input] {
				e += fmt.Sprintf("        %s: %d\n", k, v)
			}
			t.Error(e)
		}

	}
}
