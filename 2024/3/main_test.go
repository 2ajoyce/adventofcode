package main

import (
	"os"
	"strings"
	"testing"
)

func SetUpTestInput(t *testing.T, INPUT_FILE string) {
	// Write the input data to input.txt
	inputData := []string{}
	inputData = append(inputData, "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))")

	err := os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)
	if err != nil {
		t.Errorf("Failed to write input data: %v", err)
	}
}

func TestMainWritesToFile(t *testing.T) {
	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"

	os.Setenv("INPUT_FILE", INPUT_FILE)
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	os.Setenv("DEBUG", "true")

	// Don't forget to clean up! :D
	defer os.Remove(INPUT_FILE)
	defer os.Remove(OUTPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	defer os.Unsetenv("OUTPUT_FILE")
	defer os.Unsetenv("DEBUG")

	// Set up the input data
	SetUpTestInput(t, INPUT_FILE)

	// Run the main function
	main()

	// Read the content of output.txt
	data, err := os.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
	}

	ValidateOutput(t, string(data))
}

func ValidateOutput(t *testing.T, content string) {
	expectedContent := []byte("Result 1: 161\nResult 2: 48")
	if content != string(expectedContent) {
		t.Errorf("Expected \n%s\n but got \n%s\n", expectedContent, content)
	}
}
