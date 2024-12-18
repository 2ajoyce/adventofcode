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
		"########\n",
		"#..O.O.#\n",
		"##@.O..#\n",
		"#...O..#\n",
		"#.#.O..#\n",
		"#...O..#\n",
		"#......#\n",
		"########\n",
		"\n",
		"<^^>>>vv<v>>v<<\n",
	}
	// End State
	// ########
	// #....OO#
	// ##.....#
	// #.....O#
	// #.#O@..#
	// #...O..#
	// #...O..#
	// ########
	const total = 1751
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Total: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainDoubleBox(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"#######\n",
		"#...#.#\n",
		"#.....#\n",
		"#..OO@#\n",
		"#..O..#\n",
		"#.....#\n",
		"#######\n",
		"\n",
		"<vv<<^^<<^^\n",
	}
	// 	100 * 1 + 5 = 105
	// 100 * 2 + 7 = 207
	// 100 * 3 + 6 = 306
	const total = 618
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Total: %d", total)
	validateOutput(t, expectedContent)
}

func TestMainLarge(t *testing.T) {
	// Set up the input data
	inputData := []string{
		"##########\n",
		"#..O..O.O#\n",
		"#......O.#\n",
		"#.OO..O.O#\n",
		"#..O@..O.#\n",
		"#O#..O...#\n",
		"#O..O..O.#\n",
		"#.OO.O.OO#\n",
		"#....O...#\n",
		"##########\n",
		"\n",
		"<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^\n",
		"vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v\n",
		"><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<\n",
		"<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^\n",
		"^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><\n",
		"^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^\n",
		">^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^\n",
		"<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>\n",
		"^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>\n",
		"v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^\n",
	}
	// End State
	// ##########
	// #.O.O.OOO#
	// #........#
	// #OO......#
	// #OO@.....#
	// #O#.....O#
	// #O.....OO#
	// #O.....OO#
	// #OO....OO#
	// ##########
	const total = 9021
	os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)

	// Run the main function
	main()

	expectedContent := fmt.Sprintf("Total: %d", total)
	validateOutput(t, expectedContent)
}
