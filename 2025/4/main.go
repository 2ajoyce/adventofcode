package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// First Problem
	input := make(chan [][]rune)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan [][]rune)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan [][]rune) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	result := [][]rune{}
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, ParseInput(line))
	}
	c <- result
	close(c)
}

// ParseInput parses the input into the necessary data structure.
// On more complex inputs, this allows us to use lines of text as input for tests
func ParseInput(input string) []rune {
	return []rune(input)
}

func Solve1(input chan [][]rune) (string, error) {
	total := 0
	// The use of a channel here is contrived, but most problems have been processed line by line
	grid := <-input
	// PrintGrid(grid)
	h := CalculateHeatmap(grid)
	// PrintHeatmap(h)

	for y, row := range grid {
		for x, c := range row {
			if IsPaper(c) && h[y][x] < 4 {
				total++
			}
		}
	}

	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan [][]rune) (string, error) {
	total := 0
	// The use of a channel here is contrived, but most problems have been processed line by line
	grid := <-input
	// PrintGrid(grid)
	h := CalculateHeatmap(grid)
	// PrintHeatmap(h)

	for {
		paperRemoved := false
		// Iterate over the grid, removing all paper
		for y, row := range grid {
			for x, c := range row {
				if IsPaper(c) && h[y][x] < 4 {
					grid[y][x] = '.' // Remove the paper from the grid
					total++
					paperRemoved = true
				}
			}
		}
		if paperRemoved == false {
			break // No more paper can be removed
		}
		// Recalculate the new heatmap
		h = CalculateHeatmap(grid)
	}

	return fmt.Sprintf("%d", total), nil
}

func PrintGrid(g [][]rune) {
	for _, row := range g {
		fmt.Print("|")
		for _, char := range row {
			fmt.Printf("%c|", char)
		}
		fmt.Println()
	}
}

func PrintHeatmap(g [][]int) {
	for _, row := range g {
		fmt.Print("|")
		for _, char := range row {
			fmt.Printf("%d|", char)
		}
		fmt.Println()
	}
}

func IsEmpty(r rune) bool {
	if r == '.' {
		return true
	}
	return false
}

func IsPaper(r rune) bool {
	if r == '@' {
		return true
	}
	return false
}

type Coord struct {
	x int
	y int
}

// Nearby returns the 8 coordinates surrounding c.
// It removes coordinates that fall outside the range [0, max].
func (c *Coord) Nearby(maxX, maxY int) []Coord {
	var result []Coord

	// Offsets for the 8 surrounding cells
	deltas := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for _, d := range deltas {
		nx := c.x + d.dx
		ny := c.y + d.dy

		// Check bounds (≥0 and ≤max)
		if nx >= 0 && ny >= 0 && nx <= maxX && ny <= maxY {
			result = append(result, Coord{nx, ny})
		}
	}

	return result
}

func CalculateHeatmap(g [][]rune) [][]int {
	// Initialize the heatmap
	h := make([][]int, len(g))
	for i, _ := range h {
		h[i] = make([]int, len(g[0]))
	}

	// Loop over every square in the input
	for y, row := range g {
		for x, char := range row {
			// If the square is paper, increase the surrounding heatmap values by one
			if IsPaper(char) {
				current := Coord{x, y}
				surrounding := current.Nearby(len(g[0])-1, len(g)-1)
				for _, c := range surrounding {
					h[c.y][c.x]++
				}
			}
		}
	}
	return h
}
