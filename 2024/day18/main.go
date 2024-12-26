package main

import (
	"day18/internal/aocUtils"
	"day18/internal/simulation"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
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

	sim, obstacles, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve(sim, obstacles)
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

func parseLines(lines []string) (simulation.Simulation, []simulation.Coord, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, nil, fmt.Errorf("input is empty")
	}

	// The first line is the simulation width and height
	// The rest of the lines are the simulation layout
	dimension := lines[0]
	dimensions := strings.Split(dimension, ":")
	if len(dimensions) != 2 {
		return nil, nil, fmt.Errorf("invalid dimension format: %s", dimension)
	}
	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid width: %s", dimensions[0])
	}
	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid height: %s", dimensions[1])
	}

	// Initialize the simulation
	sim := simulation.NewSimulation(width, height)

	// Remove the dimensions from the lines now that it is parsed
	lines = lines[1:]

	obstacles := []simulation.Coord{}
	// Parse the rest of the lines
	for i, line := range lines {
		coords := strings.Split(line, ",")
		if len(coords) != 2 {
			return nil, nil, fmt.Errorf("invalid coordinate format on line %d: %s", i, line)
		}
		x, err := strconv.Atoi(coords[0])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid x coordinate: %s", coords[0])
		}
		y, err := strconv.Atoi(coords[1])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid y coordinate: %s", coords[1])
		}
		coord := simulation.Coord{X: x, Y: y}
		obstacles = append(obstacles, coord)
	}

	// Verify the simulation was created correctly
	if len(obstacles) > 100 && len(obstacles) != 3450 { // This is hacky, but it should avoid running validation on the tests
		return nil, nil, fmt.Errorf("expected 3450 obstacles, got %d", len(obstacles))
	}

	if DEBUG {
		fmt.Printf("Parsing complete\n\n")
	}

	return sim, obstacles, nil
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

func solve(sim simulation.Simulation, obstacles []simulation.Coord) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning single-threaded solve")

	if DEBUG {
		// Print the initial state
		fmt.Printf("Initial State:\n%s\n", stringifySimulation(sim, nil))
	}

	//////////////////////////////////////////////////////////////////////////////////////
	// This giant loop isn't optimal, but it worked for the scope of the problem
	//////////////////////////////////////////////////////////////////////////////////////
	var alwaysValid bool = true
	var finalObstacle simulation.Coord // The obstacle that when added, prevents any path from being found
	bar := progressbar.Default(int64(len(obstacles)))
	for i, obstacle := range obstacles {
		// Clone the simulation
		cloneSim := sim.Clone()

		// Add all obstacles [0, i] to the simulation
		for j := 0; j <= i; j++ {
			entity, err := simulation.NewEntity(ObstacleEntityType)
			if err != nil {
				return nil, fmt.Errorf("error creating obstacle entity: %v", err)
			}
			_, err = cloneSim.AddEntity(entity, []simulation.Coord{obstacles[j]}, simulation.North)
			if err != nil {
				return nil, fmt.Errorf("error adding obstacle %d: %v", j, err)
			}
		}

		// Make the graph
		g, err := makeGraph(cloneSim)
		if err != nil {
			return nil, fmt.Errorf("error making graph: %v", err)
		}

		// Solve the graph using Dijkstra's algorithm
		start := simulation.Coord{X: 0, Y: 0}
		target := simulation.Coord{X: cloneSim.GetMap().GetWidth() - 1, Y: cloneSim.GetMap().GetHeight() - 1}
		_, cost := simulation.Dijkstra(g, start, target, simulation.CostManhattan)
		if cost < 0 {
			alwaysValid = false
			finalObstacle = obstacle

			if DEBUG {
				fmt.Printf("Path:\n%s\n", stringifySimulation(cloneSim, nil))
			}

			break
		}
		bar.Add(1)
	}

	if alwaysValid { // This should only happen in test cases
		return []string{"No obstacle fully blocks the path"}, nil
	}

	result := []string{fmt.Sprintf("%s", finalObstacle.String())}
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
