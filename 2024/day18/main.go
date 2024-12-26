package main

import (
	"day18/internal/aocUtils"
	"day18/internal/simulation"
	"fmt"
	"os"
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

const ObstacleEntityType = "#"

func parseLines(lines []string) (simulation.Simulation, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	// The first line is the simulation width and height
	// The rest of the lines are the simulation layout
	dimension := lines[0]
	dimensions := strings.Split(dimension, ":")
	if len(dimensions) != 2 {
		return nil, fmt.Errorf("invalid dimension format: %s", dimension)
	}
	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, fmt.Errorf("invalid width: %s", dimensions[0])
	}
	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, fmt.Errorf("invalid height: %s", dimensions[1])
	}

	// Initialize the simulation
	sim := simulation.NewSimulation(width, height)

	// Remove the dimensions from the lines now that it is parsed
	lines = lines[1:]

	if DEBUG {
		for i, line := range lines {
			fmt.Printf("Line %4d :   %s\n", i, line)
		}
		fmt.Printf("\n")
	}

	// Parse the rest of the lines
	for i, line := range lines {
		if i == 1024 {
			// Part 1 only reads in the first 1024 lines (2025 if you include the dimensions)
			expectedValue := "52,51"
			if line != expectedValue {
				return nil, fmt.Errorf("line 1024 is %s, expected %s", line, expectedValue)
			}
			break
		}
		if DEBUG {
			fmt.Printf("Parsing line %d: %s\n", i, line)
		}
		coords := strings.Split(line, ",")
		if len(coords) != 2 {
			return nil, fmt.Errorf("invalid coordinate format on line %d: %s", i, line)
		}
		x, err := strconv.Atoi(coords[0])
		if err != nil {
			return nil, fmt.Errorf("invalid x coordinate: %s", coords[0])
		}
		y, err := strconv.Atoi(coords[1])
		if err != nil {
			return nil, fmt.Errorf("invalid y coordinate: %s", coords[1])
		}
		coord := simulation.Coord{X: x, Y: y}
		obst, err := simulation.NewEntity(ObstacleEntityType)
		if err != nil {
			return nil, fmt.Errorf("error creating entity %s at %v: %v", ObstacleEntityType, coord, err)
		}
		sim.AddEntity(obst, []simulation.Coord{coord}, simulation.North)
		if DEBUG {
			fmt.Printf("    Added entity %s at %v\n", ObstacleEntityType, coord)
		}
	}

	// Verify the simulation was created correctly
	entities := sim.GetEntities()
	if len(entities) > 100 && len(entities) != 1024 { // This is hacky, but it should avoid running validation on the tests
		return nil, fmt.Errorf("expected 1024 entities, got %d", len(entities))
	}

	if DEBUG {
		fmt.Printf("Parsing complete\n\n")
	}

	return sim, nil
}

type printMask map[simulation.Coord]string

func stringifySimulation(sim simulation.Simulation, masks []printMask) string {
	output := ""
	m := sim.GetMap()
	for y := 0; y < m.GetHeight(); y++ {
		for x := 0; x < m.GetWidth(); x++ {
			sym := "%"
			cell, err := m.GetCell(simulation.Coord{X: x, Y: y})
			if err != nil {
				return fmt.Sprintf("error getting cell at %v: %v", simulation.Coord{X: x, Y: y}, err)
			}
			entityIds := cell.GetEntityIds()
			if len(entityIds) == 0 { // Empty cell
				sym = "."
			} else { // Entity in cell
				entity, err := sim.GetEntity(entityIds[0])
				if err != nil {
					return fmt.Sprintf("error getting entity %s: %v", entityIds[0], err)
				}
				sym = entity.GetEntityType()
			}
			// Apply masks
			for _, mask := range masks {
				if mask != nil {
					if val, ok := mask[simulation.Coord{X: x, Y: y}]; ok {
						sym = val
					}
				}
			}
			output += sym
		}
		output += "\n"
	}
	return output
}

func solve(sim simulation.Simulation) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning single-threaded solve")

	// Print the initial state
	fmt.Printf("Initial State:\n%s\n", stringifySimulation(sim, nil))

	// Make the graph
	fmt.Printf("Making graph...\n")
	g, err := makeGraph(sim)
	if err != nil {
		return nil, fmt.Errorf("error making graph: %v", err)
	}

	// Solve the graph using Dijkstra's algorithm
	start := simulation.Coord{X: 0, Y: 0}
	target := simulation.Coord{X: sim.GetMap().GetWidth() - 1, Y: sim.GetMap().GetHeight() - 1}
	fmt.Printf("Solving graph from %v to %v...\n", start, target)
	path, cost := simulation.Dijkstra(g, start, target, simulation.CostManhattan)
	fmt.Printf("Path found with cost %.1f:\n", cost)
	if DEBUG {
		fmt.Printf("Path: %v\n", path)
	}
	result := []string{fmt.Sprintf("%.0f", cost)}
	return result, nil
}

type graph = map[simulation.Coord]map[simulation.Coord]float64

func makeGraph(sim simulation.Simulation) (graph, error) {
	// Turn the simulation into a generic graph for pathfinding
	var g graph = make(map[simulation.Coord]map[simulation.Coord]float64)
	m := sim.GetMap()
	for y := 0; y < m.GetHeight(); y++ {
		for x := 0; x < m.GetWidth(); x++ {
			coord := simulation.Coord{X: x, Y: y}
			cell, err := m.GetCell(coord)
			if err != nil {
				return g, fmt.Errorf("error getting cell at %v: %v", coord, err)
			}
			entityIds := cell.GetEntityIds()
			if len(entityIds) == 0 { // Empty cell
				g[coord] = make(map[simulation.Coord]float64)
				sim.GetMap().GetNeighbors(coord)
				for _, neighbor := range sim.GetMap().GetNeighbors(coord) {
					g[coord][neighbor] = 0
				}
			}
		}
	}
	return g, nil
}
