package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func writeInputToFile(INPUT_FILE string, inputData []string, t *testing.T) {
	err := os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)
	if err != nil {
		t.Errorf("Failed to write input data: %v", err)
	}
}

func validateOutput(t *testing.T, content string, expectedContent string) bool {
	if content != expectedContent {
		t.Errorf("Expected \n%s\n but got \n%s\n", expectedContent, content)
		return false
	}
	return true
}

func SetUpFullTestInput(t *testing.T, INPUT_FILE string) (total int, inputData []string) {
	inputData = append(inputData, "89010123\n")
	inputData = append(inputData, "78121874\n")
	inputData = append(inputData, "87430965\n")
	inputData = append(inputData, "96549874\n")
	inputData = append(inputData, "45678903\n")
	inputData = append(inputData, "32019012\n")
	inputData = append(inputData, "01329801\n")
	inputData = append(inputData, "10456732\n")
	total = 36
	return total, inputData
}

func TestMain(t *testing.T) {
	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"

	// Don't forget to clean up!
	os.Setenv("INPUT_FILE", INPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	defer os.Unsetenv("OUTPUT_FILE")
	os.Setenv("PARALLELISM", "1")
	defer os.Unsetenv("PARALLELISM")
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

	// Set up the input data
	total, inputData := SetUpFullTestInput(t, INPUT_FILE)
	writeInputToFile(INPUT_FILE, inputData, t)

	// Run the main function
	main()

	// Read the content of output.txt
	data, err := os.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
	}
	expectedContent := fmt.Sprintf("Score: %d", total)
	valid := validateOutput(t, string(data), expectedContent)
	if valid {
		os.Remove(INPUT_FILE)
		os.Remove(OUTPUT_FILE)
	}
}
