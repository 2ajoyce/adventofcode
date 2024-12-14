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

func TestMainSmallSuccess1(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"Button A: X+94, Y+34\n",
		"Button B: X+22, Y+67\n",
		"Prize: X=8400, Y=5400\n",
	}
	total := 280
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Tokens: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainSmallFailure1(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"Button A: X+26, Y+66\n",
		"Button B: X+67, Y+21\n",
		"Prize: X=12748, Y=12176\n",
	}
	total := 0
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Tokens: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainSmallSuccess2(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"Button A: X+17, Y+86\n",
		"Button B: X+84, Y+37\n",
		"Prize: X=7870, Y=6450\n",
	}
	total := 200
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Tokens: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainSmallFailure2(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"Button A: X+69, Y+23\n",
		"Button B: X+27, Y+71\n",
		"Prize: X=18641, Y=10279\n",
	}
	total := 0
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Tokens: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainMedium(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"Button A: X+94, Y+34\n",
		"Button B: X+22, Y+67\n",
		"Prize: X=8400, Y=5400\n",
		"\n",
		"Button A: X+26, Y+66\n",
		"Button B: X+67, Y+21\n",
		"Prize: X=12748, Y=12176\n",
		"\n",
		"Button A: X+17, Y+86\n",
		"Button B: X+84, Y+37\n",
		"Prize: X=7870, Y=6450\n",
		"\n",
		"Button A: X+69, Y+23\n",
		"Button B: X+27, Y+71\n",
		"Prize: X=18641, Y=10279\n",
	}
	total := 480
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Tokens: %d", total)
	validateOutput(t, expectedContent)
}
