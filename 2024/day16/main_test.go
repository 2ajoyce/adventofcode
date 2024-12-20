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

func TestMainSmall1(t *testing.T) {
	inputData := []string{
		"###############\n",
		"#.......#....E#\n",
		"#.#.###.#.###.#\n",
		"#.....#.#...#.#\n",
		"#.###.#####.#.#\n",
		"#.#.#.......#.#\n",
		"#.#.#####.###.#\n",
		"#...........#.#\n",
		"###.#.#####.#.#\n",
		"#...#.....#.#.#\n",
		"#.#.#.###.#.#.#\n",
		"#.....#...#.#.#\n",
		"#.###.#.#.#.#.#\n",
		"#S..#.....#...#\n",
		"###############\n",
	}
	const total = 7036
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Total: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainSmall2(t *testing.T) {
	inputData := []string{
		"#################\n",
		"#...#...#...#..E#\n",
		"#.#.#.#.#.#.#.#.#\n",
		"#.#.#.#...#...#.#\n",
		"#.#.#.#.###.#.#.#\n",
		"#...#.#.#.....#.#\n",
		"#.#.#.#.#.#####.#\n",
		"#.#...#.#.#.....#\n",
		"#.#.#####.#.###.#\n",
		"#.#.#.......#...#\n",
		"#.#.###.#####.###\n",
		"#.#.#...#.....#.#\n",
		"#.#.#.#####.###.#\n",
		"#.#.#.........#.#\n",
		"#.#.#.#########.#\n",
		"#S#.............#\n",
		"#################\n",
	}
	const total = 11048
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Total: %d", total)
	validateOutput(t, expectedContent)
}
