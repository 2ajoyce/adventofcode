package main

import (
	"day20/internal/aocUtils"
	"day20/internal/simulation"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {

	////////////////////////////////////////////////////////////////////
	// ENVIRONMENT SETUP
	////////////////////////////////////////////////////////////////////

	//os.Setenv("DEBUG", "true")
	INPUT_FILE := os.Getenv("INPUT_FILE")
	OUTPUT_FILE := os.Getenv("OUTPUT_FILE")
	PARALLELISM, err := strconv.Atoi(os.Getenv("PARALLELISM"))
	if PARALLELISM < 1 || err != nil {
		PARALLELISM = 1
	}
	fmt.Printf("PARALLELISM: %d\n\n", PARALLELISM)

	if INPUT_FILE == "" || OUTPUT_FILE == "" {
		fmt.Println("INPUT_FILE and OUTPUT_FILE environment variables not set")
		fmt.Println("Defaulting to input.txt and output.txt")
		INPUT_FILE = "input.txt"
		OUTPUT_FILE = "output.txt"
	}

	////////////////////////////////////////////////////////////////////
	// READ INPUT FILE
	////////////////////////////////////////////////////////////////////

	lines, err := aocUtils.ReadFile(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error reading from %s: %v", INPUT_FILE, err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// SOLUTION LOGIC
	////////////////////////////////////////////////////////////////////

	input, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve(input)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// WRITE OUTPUT FILE
	////////////////////////////////////////////////////////////////////

	err = aocUtils.WriteToFile(OUTPUT_FILE, results)
	if err != nil {
		fmt.Printf("Error writing to %s: %v", OUTPUT_FILE, err)
		return
	}

	fmt.Printf("Successfully processed %s and created %s\n", INPUT_FILE, OUTPUT_FILE)
}

func parseLines(lines []string) ([]simulation.Coord, error) {
	// DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	startCoord := simulation.Coord{X: -1, Y: -1}
	endCoord := simulation.Coord{X: -1, Y: -1}
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
		for j, r := range line {
			if r == 'S' {
				startCoord = simulation.Coord{X: j, Y: i}
			} else if r == 'E' {
				endCoord = simulation.Coord{X: j, Y: i}
			}
		}
	}

	fmt.Printf("Parsed Grid\n    Start:%s, End:%s\n%s\n", startCoord.String(), endCoord.String(), stringifyGrid(grid, nil, 4))
	fmt.Println()

	// Find the path from start to end
	path := findPath(grid, startCoord, endCoord)

	return path, nil
}

type GridMask struct {
	coords []simulation.Coord
	mask   rune
}

func stringifyGrid(grid [][]rune, mask []GridMask, indent int) string {
	result := ""
	newChar := ' '
	for y, row := range grid {
		result += strings.Repeat(" ", indent)
		for x, r := range row {
			newChar = r
			for _, m := range mask {
				for _, coord := range m.coords {
					if coord.X == x && coord.Y == y {
						newChar = m.mask
					}
				}
			}
			result += string(newChar)
		}
		result += "\n"
	}
	return result
}

func findPath(grid [][]rune, startCoord, endCoord simulation.Coord) []simulation.Coord {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Finding Path...")

	path := []simulation.Coord{startCoord}
	priorCoords := []simulation.Coord{startCoord}
	currentCoord := startCoord

	for currentCoord != endCoord {
		priorCoords = append(priorCoords, currentCoord)
		neighbors := currentCoord.GetNeighbors()
		if DEBUG {
			gridString := stringifyGrid(grid, []GridMask{{coords: path, mask: '@'}, {coords: neighbors, mask: 'N'}}, 4)
			fmt.Println(gridString)
		}
		for _, neighbor := range neighbors {
			if neighbor.X < 0 || neighbor.Y < 0 || neighbor.Y >= len(grid) || neighbor.X >= len(grid[0]) {
				continue
			}
			if slices.Contains(priorCoords, neighbor) {
				continue
			}
			// Only one neighbor is valid
			if grid[neighbor.Y][neighbor.X] == '.' || grid[neighbor.Y][neighbor.X] == 'E' {
				// Move to the neighbor
				currentCoord = neighbor
				path = append(path, currentCoord)
				break
			}
		}
		if slices.Contains(priorCoords, currentCoord) {
			// No valid neighbors
			break
		}
	}

	// Validate the path
	if path[0] != startCoord || path[len(path)-1] != endCoord {
		fmt.Println("No valid path found")
		return []simulation.Coord{}
	}

	if DEBUG {
		fmt.Printf("Path Found: %v\n", path)
		gridString := stringifyGrid(grid, []GridMask{{coords: path, mask: '@'}}, 4)
		fmt.Println(gridString)
	}
	return path
}

func solve(path []simulation.Coord) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning solve...")

	return nil, nil
}
