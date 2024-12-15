package main

import (
	"day12/internal/aocUtils"
	"fmt"
	"os"
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

	lines, err := aocUtils.ReadInput(INPUT_FILE)
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
	results, err := solve1(input, PARALLELISM)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// WRITE OUTPUT FILE
	////////////////////////////////////////////////////////////////////

	err = aocUtils.WriteOutput(OUTPUT_FILE, results)
	if err != nil {
		fmt.Printf("Error writing to %s: %v", OUTPUT_FILE, err)
		return
	}

	fmt.Printf("Successfully processed %s and created %s", INPUT_FILE, OUTPUT_FILE)
}

func parseLines(lines []string) (RegionMap, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	rm := make(RegionMap, len(lines))
	for i, line := range lines {
		rm[i] = strings.TrimSpace(line)
	}

	fmt.Printf("Parsed map with height %d and width %d\n", len(rm), len(rm[0]))

	return rm, nil
}

// Coordinate represents a point in the grid
type Coordinate struct {
	X, Y int
}

// RegionMap represents the grid
type RegionMap []string

// UnionFind structure for region merging
type UnionFind struct {
	parent map[Coordinate]Coordinate
	rank   map[Coordinate]int
	size   map[Coordinate]int // To track the size (area) of each region
}

// Create a new UnionFind
func NewUnionFind() *UnionFind {
	return &UnionFind{
		parent: make(map[Coordinate]Coordinate),
		rank:   make(map[Coordinate]int),
		size:   make(map[Coordinate]int),
	}
}

// Find the root of a coordinate
func (uf *UnionFind) Find(coord Coordinate) Coordinate {
	if uf.parent[coord] == coord {
		return coord
	}
	uf.parent[coord] = uf.Find(uf.parent[coord]) // Path compression
	return uf.parent[coord]
}

// Union two coordinates
func (uf *UnionFind) Union(coord1, coord2 Coordinate) {
	root1 := uf.Find(coord1)
	root2 := uf.Find(coord2)

	if root1 != root2 {
		// Union by rank
		if uf.rank[root1] > uf.rank[root2] {
			uf.parent[root2] = root1
			uf.size[root1] += uf.size[root2]
		} else if uf.rank[root1] < uf.rank[root2] {
			uf.parent[root1] = root2
			uf.size[root2] += uf.size[root1]
		} else {
			uf.parent[root2] = root1
			uf.size[root1] += uf.size[root2]
			uf.rank[root1]++
		}
	}
}

// Initialize the UnionFind structure for the grid
func (uf *UnionFind) Initialize(grid RegionMap) {
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			coord := Coordinate{X: x, Y: y}
			uf.parent[coord] = coord
			uf.rank[coord] = 0
			uf.size[coord] = 1 // Initially, every cell is its own region
		}
	}
}

// Get neighbors of a coordinate
func getNeighbors(coord Coordinate, grid RegionMap) []Coordinate {
	dirs := []Coordinate{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	neighbors := []Coordinate{}

	for _, d := range dirs {
		nx, ny := coord.X+d.X, coord.Y+d.Y
		if nx >= 0 && ny >= 0 && ny < len(grid) && nx < len(grid[ny]) {
			neighbors = append(neighbors, Coordinate{X: nx, Y: ny})
		}
	}

	return neighbors
}

// Calculate the number of straight sides a region has
func calculateNumberOfSides(regionRoot Coordinate, grid RegionMap, uf *UnionFind) int {
	// Maps to store horizontal and vertical boundaries
	horizontalLines := make(map[int][]int)      // key: y, value: list of x where horizontal boundary starts
	verticalLines := make(map[int]map[int]bool) // key: x, value: set of y where vertical boundary starts

	// Identify all boundary edges
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			coord := Coordinate{X: x, Y: y}
			if uf.Find(coord) != regionRoot {
				continue
			}

			// Check four directions for boundary edges
			// Up
			if y == 0 || uf.Find(Coordinate{X: x, Y: y - 1}) != regionRoot {
				// Horizontal boundary at the top of the cell (y, x)
				horizontalLines[y] = append(horizontalLines[y], x)
			}

			// Down
			if y == len(grid)-1 || uf.Find(Coordinate{X: x, Y: y + 1}) != regionRoot {
				// Horizontal boundary at the bottom of the cell (y+1, x)
				horizontalLines[y+1] = append(horizontalLines[y+1], x)
			}

			// Left
			if x == 0 || uf.Find(Coordinate{X: x - 1, Y: y}) != regionRoot {
				// Vertical boundary on the left of the cell (x, y)
				if verticalLines[x] == nil {
					verticalLines[x] = make(map[int]bool)
				}
				verticalLines[x][y] = true
			}

			// Right
			if x == len(grid[y])-1 || uf.Find(Coordinate{X: x + 1, Y: y}) != regionRoot {
				// Vertical boundary on the right of the cell (x+1, y)
				if verticalLines[x+1] == nil {
					verticalLines[x+1] = make(map[int]bool)
				}
				verticalLines[x+1][y] = true
			}
		}
	}

	numberOfSides := 0

	// Function to count horizontal segments
	countHorizontalSegments := func(y int, xs []int) int {
		if len(xs) == 0 {
			return 0
		}
		count := 0
		sort.Ints(xs)
		prev := -2
		for i, x := range xs {
			if i == 0 {
				count++
			} else {
				// Check if current x is consecutive
				if x == prev+1 {
					// Check if there is a vertical boundary at x, y
					if verticalLines[x] != nil && verticalLines[x][y] {
						// There is a vertical boundary at x, y, cannot merge
						count++
					}
				} else {
					count++
				}
			}
			prev = x
		}
		return count
	}

	// Function to count vertical segments
	countVerticalSegments := func(x int, ys []int) int {
		if len(ys) == 0 {
			return 0
		}
		count := 0
		sortedYs := make([]int, len(ys))
		copy(sortedYs, ys)
		sort.Ints(sortedYs)
		prev := -2
		for i, y := range sortedYs {
			if i == 0 {
				count++
			} else {
				// Check if current y is consecutive
				if y == prev+1 {
					// Check if there is a horizontal boundary at y, x
					hasHorizontalBoundary := false
					if horizontalLines[y] != nil {
						for _, hX := range horizontalLines[y] {
							if hX == x {
								hasHorizontalBoundary = true
								break
							}
						}
					}
					if hasHorizontalBoundary {
						// There is a horizontal boundary at y, x, cannot merge
						count++
					}
				} else {
					count++
				}
			}
			prev = y
		}
		return count
	}

	// Count horizontal segments
	for y, xs := range horizontalLines {
		numberOfSides += countHorizontalSegments(y, xs)
	}

	// Count vertical segments
	for x, ysMap := range verticalLines {
		ys := []int{}
		for y := range ysMap {
			ys = append(ys, y)
		}
		numberOfSides += countVerticalSegments(x, ys)
	}

	return numberOfSides
}

// Solve function with expense calculation
func solve1(rm RegionMap, parallelism int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	var output = []string{}
	var totalExpense = 0 // The sum of the expense of every region

	uf := NewUnionFind()
	uf.Initialize(rm)

	// Merge regions based on adjacency and same value
	for y := 0; y < len(rm); y++ {
		for x := 0; x < len(rm[y]); x++ {
			coord := Coordinate{X: x, Y: y}
			for _, neighbor := range getNeighbors(coord, rm) {
				if rm[coord.Y][coord.X] == rm[neighbor.Y][neighbor.X] {
					uf.Union(coord, neighbor)
				}
			}
		}
	}

	// Calculate expenses for each region
	visited := make(map[Coordinate]bool)
	for y := 0; y < len(rm); y++ {
		for x := 0; x < len(rm[y]); x++ {
			coord := Coordinate{X: x, Y: y}
			root := uf.Find(coord)

			// Skip if already calculated
			if visited[root] {
				continue
			}
			visited[root] = true

			// Calculate area and number of sides
			area := uf.size[root]
			numberOfSides := calculateNumberOfSides(root, rm, uf)

			// Calculate expense
			expense := area * numberOfSides
			totalExpense += expense
			if DEBUG {
				fmt.Printf("Region rooted at (%d, %d): Area=%d, NumberOfSides=%d, Expense=%d\n", root.X, root.Y, area, numberOfSides, expense)
			}
		}
	}

	// Append total expense to output
	output = append(output, fmt.Sprintf("Expense: %d", totalExpense))
	return output, nil
}
