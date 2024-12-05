package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ReadInput(INPUT_FILE string) ([]string, error) {
	inputFile, err := os.Open(INPUT_FILE)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %v", INPUT_FILE, err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func WriteOutput(OUTPUT_FILE string, results []string) error {
	outputFile, err := os.Create(OUTPUT_FILE)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", OUTPUT_FILE, err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	// Write the results to output.txt, one line per result
	for i, res := range results {
		_, err := writer.WriteString(res)
		if err != nil {
			return fmt.Errorf("error writing value to %s: %v", OUTPUT_FILE, err)
		}
		if i != len(results)-1 {
			_, err = writer.WriteString("\n")
			if err != nil {
				return fmt.Errorf("error writing newline to %s: %v", OUTPUT_FILE, err)
			}
		}
	}

	// Flush the writer to ensure all data is written to output.txt
	writer.Flush()
	return nil
}

func main() {
	//os.Setenv("DEBUG", "true")
	INPUT_FILE := os.Getenv("INPUT_FILE")
	OUTPUT_FILE := os.Getenv("OUTPUT_FILE")

	if INPUT_FILE == "" || OUTPUT_FILE == "" {
		fmt.Println("INPUT_FILE and OUTPUT_FILE environment variables not set")
		fmt.Println("Defaulting to input.txt and output.txt")
		INPUT_FILE = "input.txt"
		OUTPUT_FILE = "output.txt"
	}

	lines, err := ReadInput(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error reading from %s: %v", INPUT_FILE, err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// Start Solution Logic  ///////////////////////////////////////////
	////////////////////////////////////////////////////////////////////

	// Create an array of all coordinates containing the letter X
	startingCoords, err := FindLetter(lines, 'X')
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	// Given starting coordinates, determine how many possible matching words exist
	results, err := Solve1(lines, startingCoords)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}

	// Create an array of all coordinates containing the letter A
	startingCoords, err = FindLetter(lines, 'A')
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results2, err := Solve2(lines, startingCoords)
	results = append(results, results2...)
	if err != nil {
		fmt.Println("Error solving 2:", err)
		return
	}
	////////////////////////////////////////////////////////////////////
	// End Solution Logic  /////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////

	err = WriteOutput(OUTPUT_FILE, results)
	if err != nil {
		fmt.Printf("Error writing to %s: %v", OUTPUT_FILE, err)
		return
	}

	fmt.Printf("Successfully processed %s and created %s", INPUT_FILE, OUTPUT_FILE)
}

func FindLetter(lines []string, letter byte) ([][2]int, error) {
	DEBUG := os.Getenv("DEBUG")

	// Create a slice of tuples (x, y) for each match found in the lines
	startingCoords := [][2]int{}

	for i, line := range lines {
		// Find all occurrences of the letter in the line
		re := regexp.MustCompile(fmt.Sprintf(`%c`, letter))
		matches := re.FindAllStringIndex(line, -1)

		// Append each match to the slice
		for _, match := range matches {
			startingCoords = append(startingCoords, [2]int{match[0], i})

		}
	}
	if DEBUG == "true" {
		PrintGrid(fmt.Sprintf("Find Letter %s", string(letter)), lines, startingCoords)
	}
	return startingCoords, nil
}

func Solve1(grid []string, startingCoords [][2]int) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG")
	results := []string{}

	rollingSum := 0
	for _, startingCoord := range startingCoords {
		gridSize := len(grid)
		directions := [][]int{
			{-1, 0},  // up
			{1, 0},   // down
			{0, -1},  // left
			{0, 1},   // right
			{-1, -1}, // up-left
			{-1, 1},  // up-right
			{1, -1},  // down-left
			{1, 1},   // down-right
		}

		for _, direction := range directions {
			x, y := startingCoord[0], startingCoord[1]
			wordFound := true

			for j := 0; j < len("XMAS"); j++ {
				if x < 0 || x >= gridSize || y < 0 || y >= gridSize {
					wordFound = false
					break
				}
				if grid[y][x] != "XMAS"[j] {
					wordFound = false
					break
				}
				x += direction[0]
				y += direction[1]
			}

			if wordFound {
				rollingSum += 1
			}
		}
	}
	results = append(results, fmt.Sprintf("Result 1: %s", strconv.Itoa(rollingSum)))
	return results, nil
}

func Solve2(grid []string, startingCoords [][2]int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG")
	results := []string{}

	rollingSum := 0
	masFound := [][2]int{}
	for _, startingCoord := range startingCoords {
		diagonals := [][][2]int{
			{{-1, -1}, {1, 1}}, // top-left to bottom-right diagonal
			{{1, 1}, {-1, -1}}, // bottom-right to top-left diagonal
			{{1, -1}, {-1, 1}}, // top-right to bottom-left diagonal
			{{-1, 1}, {1, -1}}, // bottom-left to top-right diagonal
		}

		validDiagonals := 0
		for _, diagonal := range diagonals {
			x1, y1 := startingCoord[0]+diagonal[0][0], startingCoord[1]+diagonal[0][1]
			x2, y2 := startingCoord[0]+diagonal[1][0], startingCoord[1]+diagonal[1][1]

			if x1 >= 0 && x2 >= 0 && y1 >= 0 && y2 >= 0 && x1 < len(grid) && x2 < len(grid) && y1 < len(grid[x1]) && y2 < len(grid[x2]) {
				if string(grid[y1][x1]) == "M" && string(grid[y2][x2]) == "S" {
					validDiagonals += 1
				}
			} else {
				validDiagonals -= 1
			}

			if validDiagonals == 2 {
				rollingSum += 1
				masFound = append(masFound, startingCoord)
				break
			} else if validDiagonals < 0 {
				break
			}
		}
	}

	if DEBUG == "true" {
		PrintGrid("Result 2:", grid, masFound)
	}
	results = append(results, fmt.Sprintf("Result 2: %s", strconv.Itoa(rollingSum)))
	return results, nil
}

func PrintGrid(label string, grid []string, selectedCoords [][2]int) {
	gridString := ""
	for i := range grid {
		for j := range grid[i] {
			found := false
			for k := range selectedCoords {
				if selectedCoords[k][0] == j && selectedCoords[k][1] == i {
					found = true
				}
			}
			if found {
				gridString += "*"
			} else {
				gridString += strings.ToLower(string(grid[i][j]))
			}
		}
		gridString += "\n"
	}
	fmt.Printf("%s\n%s\n", label, gridString)
}
