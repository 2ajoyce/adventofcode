package main

import (
	"fmt"
	"os"
	"strings"
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

func TestMain1(t *testing.T) {
	inputData := []string{
		"Register A: 729\n",
		"Register B: 0\n",
		"Register C: 0\n",
		"\n",
		"Program: 0,1,5,4,3,0\n",
	}
	const output = "4,6,3,5,6,3,5,2,1,0"
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	main()

	expectedContent := fmt.Sprintf("Output: %s", output)
	validateOutput(t, expectedContent)
}

func TestMain2(t *testing.T) {
	inputData := []string{
		"Register A: 10\n",
		"Register B: 0\n",
		"Register C: 0\n",
		"\n",
		"Program: 5,0,5,1,5,4\n",
	}
	const output = "0,1,2"
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	main()

	expectedContent := fmt.Sprintf("Output: %s", output)
	validateOutput(t, expectedContent)
}

func TestMain3(t *testing.T) {
	inputData := []string{
		"Register A: 2024\n",
		"Register B: 0\n",
		"Register C: 0\n",
		"\n",
		"Program: 0,1,5,4,3,0\n",
	}
	const output = "4,2,5,6,7,7,7,7,3,1,0"
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	main()

	expectedContent := fmt.Sprintf("Output: %s", output)
	validateOutput(t, expectedContent)
}

func TestMainCanSolvePart2(t *testing.T) {
	// This test is checking the known example from the description to prove that the program
	// works as expected
	inputData := []string{
		"Register A: 117440\n",
		"Register B: 0\n",
		"Register C: 0\n",
		"\n",
		"Program: 0,3,5,4,3,0\n",
	}
	const output = "0,3,5,4,3,0"
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	main()

	expectedContent := fmt.Sprintf("Output: %s", output)
	validateOutput(t, expectedContent)
}