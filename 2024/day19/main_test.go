package main

import (
	"day19/internal/aocUtils"
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

func TestDecomposeSentence(t *testing.T) {
	tests := []struct {
		name     string
		terms    []Term
		sentence Sentence
		expected []Term
	}{
		{
			name:     "Empty terms",
			terms:    []Term{},
			sentence: Sentence("r"),
			expected: nil,
		},
		{
			name:     "Empty sentence",
			terms:    []Term{"r"},
			sentence: Sentence(""),
			expected: []Term{},
		},
		{
			name:     "One term sentence",
			terms:    []Term{"r"},
			sentence: Sentence("r"),
			expected: []Term{"r"},
		},
		{
			name:     "Two char term",
			terms:    []Term{"rr"},
			sentence: Sentence("rr"),
			expected: []Term{"rr"},
		},
		{
			name:     "Double term sentence",
			terms:    []Term{"r"},
			sentence: Sentence("rr"),
			expected: []Term{"r", "r"},
		},
		{
			name:     "Two term sentence",
			terms:    []Term{"r", "w"},
			sentence: Sentence("rw"),
			expected: []Term{"r", "w"},
		},
		{
			name:     "Three term sentence",
			terms:    []Term{"r", "rw", "b"},
			sentence: Sentence("rrwb"),
			expected: []Term{"r", "rw", "b"},
		},
		{
			name:     "Sentence with 2 solutions",
			terms:    []Term{"r", "w", "rw"},
			sentence: Sentence("rrw"),
			expected: []Term{"r", "rw"}, // We will prefer less terms when possible
		},
		{
			name:     "Unsolvable sentence",
			terms:    []Term{"r", "w"},
			sentence: Sentence("rbw"),
			expected: nil,
		},
		{
			name:     "Large Term is a trap",
			terms:    []Term{"rr", "w", "wr"},
			sentence: Sentence("wrr"),
			expected: []Term{"w", "rr"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := decomposeSentence(tt.terms, tt.sentence)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expected != nil && result == nil {
				t.Errorf("Expected result to be non-nil")
			}
			if len(result) != len(tt.expected) {
				t.Errorf("Expected result to have %d elements, but got: %v", len(tt.expected), result)
			}
			for i, term := range result {
				if term != tt.expected[i] {
					t.Errorf("Expected result to contain '%s', but got: %v", tt.expected[i], term)
				}
			}
		})
	}
}

func TestMainExample(t *testing.T) {
	input := []string{
		"r, wr, b, g, bwu, rb, gb, br",
		"",
		"brwrr",
		"bggr",
		"gbbr",
		"rrbgbr",
		"ubwu",
		"bwurrg",
		"brgr",
		"bbrgwb",
		"",
	}
	aocUtils.WriteToFile(INPUT_FILE, input)
	expectedOutput := "6"

	main()

	validateOutput(t, expectedOutput)
}
