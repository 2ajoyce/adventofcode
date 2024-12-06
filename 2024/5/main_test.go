package main

import (
	"os"
	"strings"
	"testing"
)

func SetUpTestInput(t *testing.T, INPUT_FILE string) {
	// Write the input data to input.txt
	inputData := []string{}
	inputData = append(inputData, "47|53\n")
	inputData = append(inputData, "97|13\n")
	inputData = append(inputData, "97|61\n")
	inputData = append(inputData, "97|47\n")
	inputData = append(inputData, "75|29\n")
	inputData = append(inputData, "61|13\n")
	inputData = append(inputData, "75|53\n")
	inputData = append(inputData, "29|13\n")
	inputData = append(inputData, "97|29\n")
	inputData = append(inputData, "53|29\n")
	inputData = append(inputData, "61|53\n")
	inputData = append(inputData, "97|53\n")
	inputData = append(inputData, "61|29\n")
	inputData = append(inputData, "47|13\n")
	inputData = append(inputData, "75|47\n")
	inputData = append(inputData, "97|75\n")
	inputData = append(inputData, "47|61\n")
	inputData = append(inputData, "75|61\n")
	inputData = append(inputData, "47|29\n")
	inputData = append(inputData, "75|13\n")
	inputData = append(inputData, "53|13\n")
	inputData = append(inputData, "\n")
	inputData = append(inputData, "75,47,61,53,29\n")
	inputData = append(inputData, "97,61,53,29,13\n")
	inputData = append(inputData, "75,29,13\n")
	inputData = append(inputData, "75,97,47,61,53\n")
	inputData = append(inputData, "61,13,29\n")
	inputData = append(inputData, "97,13,75,29,47\n")

	err := os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)
	if err != nil {
		t.Errorf("Failed to write input data: %v", err)
	}
}

func TestMain(t *testing.T) {
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
	expectedContent := []byte("Result 1: 143\nResult 2: 123")
	if content != string(expectedContent) {
		t.Errorf("Expected \n%s\n but got \n%s\n", expectedContent, content)
	}
}
