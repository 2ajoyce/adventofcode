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

func TestMainSmall(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"AAAA\n",
		"BBCD\n",
		"BBCC\n",
		"EEEC\n",
	}
	// Expense is count of sides * area of region
	// sides are the combined straight sides of a region, not the sides of a cell
	totalA := 4 * 4
	totalB := 4 * 4
	totalC := 4 * 8
	totalD := 1 * 4
	totalE := 3 * 4
	total := totalA + totalB + totalC + totalD + totalE
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Expense: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainMedium(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"OOOOO\n",
		"OXOXO\n",
		"OOOOO\n",
		"OXOXO\n",
		"OOOOO\n",
	}
	// Expense is count of sides * area of region
	// sides are the combined straight sides of a region, not the sides of a cell
	total := 436
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Expense: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainEShaped(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"EEEEE\n",
		"EXXXX\n",
		"EEEEE\n",
		"EXXXX\n",
		"EEEEE\n",
	}
	// Expense is count of sides * area of region
	// sides are the combined straight sides of a region, not the sides of a cell
	total := 236
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Expense: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainTouchingDiagonally(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"AAAAAA\n",
		"AAABBA\n",
		"AAABBA\n",
		"ABBAAA\n",
		"ABBAAA\n",
		"AAAAAA\n",
	}
	// Expense is count of sides * area of region
	// sides are the combined straight sides of a region, not the sides of a cell
	total := 368
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Expense: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainLarge(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"RRRRIICCFF\n",
		"RRRRIICCCF\n",
		"VVRRRCCFFF\n",
		"VVRCCCJFFF\n",
		"VVVVCJJCFE\n",
		"VVIVCCJJEE\n",
		"VVIIICJJEE\n",
		"MIIIIIJJEE\n",
		"MIIISIJEEE\n",
		"MMMISSJEEE\n",
	}
	// Expense is count of sides * area of region
	// sides are the combined straight sides of a region, not the sides of a cell	total1 := 12 * 18
	total := 1206
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Expense: %d", total)
	validateOutput(t, expectedContent)
}
