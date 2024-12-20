package main

import (
	"day16/internal/aocUtils"
	"day16/internal/simulation"
	"errors"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/exp/slices"
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
	results, err := solve(input, PARALLELISM)
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

	fmt.Printf("Successfully processed %s and created %s\n", INPUT_FILE, OUTPUT_FILE)
}

const ReindeerEntityType = "@"
const ObstacleEntityType = "#"
const StartTileEntityType = "S"
const EndTileEntityType = "E"

func parseLines(lines []string) (simulation.Simulation, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	width := len(lines[0])
	height := len(lines)
	sim := simulation.NewSimulation(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			entityType := string(lines[y][x])
			if entityType == "." {
				continue
			}
			entity, err := simulation.NewEntity(entityType)
			if err != nil {
				return nil, fmt.Errorf("error creating entity of type %s: %v", entityType, err)
			}
			coords := []simulation.Coord{{X: x, Y: y}}
			sim.AddEntity(entity, coords, simulation.North)
		}
	}

	return sim, nil
}

func PrintSim(sim simulation.Simulation, path []simulation.Coord) string {
	var output string
	for y := 0; y < sim.GetMap().GetHeight(); y++ {
		for x := 0; x < sim.GetMap().GetWidth(); x++ {
			cell, err := sim.GetMap().GetCell(simulation.Coord{X: x, Y: y})
			if err != nil {
				output += "?"
				continue
			}
			if slices.Contains(path, simulation.Coord{X: x, Y: y}) {
				output += "X"
				continue
			}
			ids := cell.GetEntityIds()
			if len(ids) == 0 {
				output += "."
			} else {
				entity, err := sim.GetEntity(ids[0])
				if err != nil {
					output += "?"
				}
				output += string(entity.GetEntityType())
			}
		}
		output += "\n"
	}
	return output
}

func solve(sim simulation.Simulation, WORKER_COUNT int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	var output []string

	width := sim.GetMap().GetWidth()
	height := sim.GetMap().GetHeight()

	if DEBUG {
		fmt.Printf("Map Dimensions - Width: %d Height: %d\n", width, height)
	}
	startTile := findStart(sim)
	if startTile == nil {
		return output, errors.New("no Start Tile found")
	}
	startLocation := startTile.GetPosition()[0]
	if DEBUG {
		fmt.Printf("Starting Location: %s\n", startLocation.String())
	}
	endTile := findEnd(sim)
	if endTile == nil {
		return output, errors.New("no End Tile found")
	}
	endLocation := endTile.GetPosition()[0]
	if DEBUG {
		fmt.Printf("Ending Location: %s\n", endLocation.String())
	}

	path, steps, turns, err := bfs(sim, startLocation, endLocation)
	if err != nil {
		return output, err
	}

	total := calculateTotal(steps, turns)
	if DEBUG {
		fmt.Println(PrintSim(sim, path))
	}
	fmt.Println("Total:", total)
	output = append(output, fmt.Sprintf("Total: %d", total))
	return output, nil
}

func findStart(sim simulation.Simulation) simulation.Entity {
	for _, entity := range sim.GetEntities() {
		if entity.GetEntityType() == StartTileEntityType {
			return entity
		}
	}
	return nil
}

func findEnd(sim simulation.Simulation) simulation.Entity {
	for _, entity := range sim.GetEntities() {
		if entity.GetEntityType() == EndTileEntityType {
			return entity
		}
	}
	return nil
}

func calculateTotal(steps int, turns int) int {
	DEBUG := os.Getenv("DEBUG") == "true"
	total := 0
	if DEBUG {
		fmt.Println("Calculating Total")
	}
	total += steps
	total += turns * 1000
	return total
}

func bfs(sim simulation.Simulation, start, end simulation.Coord) ([]simulation.Coord, int, int, error) {
	directions := []simulation.Direction{simulation.North, simulation.East, simulation.South, simulation.West}
	queue := []simulation.Coord{start}
	visited := make(map[simulation.Coord]bool)
	parent := make(map[simulation.Coord]*simulation.Coord)
	directionMap := make(map[simulation.Coord]simulation.Direction)

	visited[start] = true
	directionMap[start] = simulation.North // Assume starting direction is North

	steps := 0
	turns := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			path := reconstructPath(parent, start, end)
			return path, steps, turns, nil
		}

		for _, dir := range directions {
			neighbor := current.Move(dir)
			if isValid(sim, neighbor) && !visited[neighbor] {
				queue = append(queue, neighbor)
				visited[neighbor] = true
				parent[neighbor] = &current
				directionMap[neighbor] = dir

				steps++
				if directionMap[current] != dir {
					turns++
				}
			}
		}
	}

	return nil, 0, 0, errors.New("no path found")
}

func isValid(sim simulation.Simulation, pos simulation.Coord) bool {
	width := sim.GetMap().GetWidth()
	height := sim.GetMap().GetHeight()
	if pos.X < 0 || pos.X >= width || pos.Y < 0 || pos.Y >= height {
		return false
	}
	cell, err := sim.GetMap().GetCell(pos)
	if err != nil {
		return false
	}
	ids := cell.GetEntityIds()
	for _, id := range ids {
		entity, _ := sim.GetEntity(id)
		if entity.GetEntityType() == ObstacleEntityType {
			return false
		}
	}
	return true
}

func reconstructPath(parent map[simulation.Coord]*simulation.Coord, start, end simulation.Coord) []simulation.Coord {
	var path []simulation.Coord
	for at := &end; at != nil; at = parent[*at] {
		path = append([]simulation.Coord{*at}, path...)
	}
	return path
}
