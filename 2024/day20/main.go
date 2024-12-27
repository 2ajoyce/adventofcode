package main

import (
	"day20/internal/aocUtils"
	"day20/internal/simulation"
	"fmt"
	"math"
	"os"
	"slices"
	"sort"
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
		if len(path) < 10 {
			fmt.Printf("    Path Found: %v\n", path)
		} else {
			fmt.Println("    Path Found")
			fmt.Printf("    First 5 Elements: %v\n", path[:5])
			fmt.Printf("    Last 5 Elements: %v\n", path[len(path)-5:])
		}
		gridString := stringifyGrid(grid, []GridMask{{coords: path, mask: '@'}}, 4)
		fmt.Println(gridString)
	}
	return path
}

func solve(path []simulation.Coord) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning solve...")

	type Cheat struct {
		start simulation.Coord
		end   simulation.Coord
	}

	// The number of steps saved mapped to the cheatsByStepsSaved enabling the savings
	var cheatsByStepsSaved map[int][]Cheat = make(map[int][]Cheat)

	// For every coordinate in the path
	// "steps" are number of steps to get to that coordinate
	for steps, coord := range path {
		// Compare every coord to subsequent coords
		for i := steps + 1; i < len(path); i++ {
			if DEBUG {
				fmt.Printf("Comparing %v to %v\n", coord, path[i])
			}

			if !canReach(coord, path[i]) {
				continue
			}
			// The coord is within 2 steps of the current coord
			if DEBUG {
				fmt.Println("    Found a reachable location")
			}

			stepsToCurrent := steps
			stepsToNew := i
			stepsSaved := stepsToNew - stepsToCurrent - 2
			if DEBUG {
				fmt.Printf("    Steps to Current: %d, Steps to New: %d, Steps Saved: %d\n", stepsToCurrent, stepsToNew, stepsSaved)
			}

			if stepsSaved > 0 { // It takes 2 steps to move to the new location, less than 2 steps saved is not worth it
				cheatsByStepsSaved[stepsSaved] = append(cheatsByStepsSaved[stepsSaved], Cheat{start: coord, end: path[i]})
			}
		}
	}
	if DEBUG {
		fmt.Println()
	}

	if DEBUG {
		fmt.Println("Cheats:")
		for stepsSaved, cheat := range cheatsByStepsSaved {
			fmt.Printf("    Steps Saved: %d\n", stepsSaved)
			for _, c := range cheat {
				fmt.Printf("        %v -> %v\n", c.start, c.end)
			}
		}
	}

	totalCheatsGreaterThanOrEqualTo100 := 0
	result := []string{"Steps Saved, Count of Cheats"}

	// Extract and sort the keys
	sortedStepsSaved := make([]int, 0, len(cheatsByStepsSaved))
	for stepsSaved := range cheatsByStepsSaved {
		sortedStepsSaved = append(sortedStepsSaved, stepsSaved)
	}
	sort.Ints(sortedStepsSaved)

	// Iterate over the sorted keys
	for _, stepsSaved := range sortedStepsSaved {
		cheats := cheatsByStepsSaved[stepsSaved]
		if stepsSaved >= 100 {
			totalCheatsGreaterThanOrEqualTo100 += len(cheats)
		}
		result = append(result, fmt.Sprintf("%d,%d", stepsSaved, len(cheats)))
	}
	fmt.Printf("Total Cheats Over 100: %d\n", totalCheatsGreaterThanOrEqualTo100)
	return result, nil
}

// canReach returns true if coord1 can reach coord2 in 2 steps
func canReach(coord1, coord2 simulation.Coord) bool {
	// If the coords are the same, then they can reach each other
	if coord1.X == coord2.X && coord1.Y == coord2.Y {
		return true
	}

	// If the coords are within 2 steps of each other, then they can reach each other
	dx := math.Abs(float64(coord1.X - coord2.X))
	dy := math.Abs(float64(coord1.Y - coord2.Y))

	if (dx == 0 && dy == 2) || (dx == 2 && dy == 0) {
		return true
	}

	return false
}
