package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func validateOutput(t *testing.T, INPUT_FILE, OUTPUT_FILE, expectedContent string) bool {
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

func TestMainSmall(t *testing.T) {
	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"
	os.Setenv("INPUT_FILE", INPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	defer os.Unsetenv("OUTPUT_FILE")
	os.Setenv("PARALLELISM", "1")
	defer os.Unsetenv("PARALLELISM")
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	// Set up the input data
	inputData := []string{"1:0 1 10 99 999"}
	total := 7
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Stones: %d", total)
	validateOutput(t, INPUT_FILE, OUTPUT_FILE, expectedContent)
}

func TestMainMedium3(t *testing.T) {
	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"
	os.Setenv("INPUT_FILE", INPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	defer os.Unsetenv("OUTPUT_FILE")
	os.Setenv("PARALLELISM", "1")
	defer os.Unsetenv("PARALLELISM")
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	// Set up the input data
	inputData := []string{"3:125 17"}
	total := 5
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Stones: %d", total)
	validateOutput(t, INPUT_FILE, OUTPUT_FILE, expectedContent)
}

func TestMainMedium6(t *testing.T) {
	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"
	os.Setenv("INPUT_FILE", INPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	defer os.Unsetenv("OUTPUT_FILE")
	os.Setenv("PARALLELISM", "1")
	defer os.Unsetenv("PARALLELISM")
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	// Set up the input data
	inputData := []string{"6:125 17"}
	total := 22
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Stones: %d", total)
	validateOutput(t, INPUT_FILE, OUTPUT_FILE, expectedContent)
}

func TestMainMedium25(t *testing.T) {
	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"
	os.Setenv("INPUT_FILE", INPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	defer os.Unsetenv("OUTPUT_FILE")
	os.Setenv("PARALLELISM", "1")
	defer os.Unsetenv("PARALLELISM")
	os.Setenv("DEBUG", "false")
	defer os.Unsetenv("DEBUG")

	// Set up the input data
	inputData := []string{"25:125 17"}
	total := 55312
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Stones: %d", total)
	validateOutput(t, INPUT_FILE, OUTPUT_FILE, expectedContent)
}
