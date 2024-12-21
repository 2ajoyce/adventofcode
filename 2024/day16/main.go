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
				output += "O"
				continue
			}
			ids := cell.GetEntityIds()
			if len(ids) == 0 {
				output += " "
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

	path, cost, _, err := simulation.Dijkstra(sim, startLocation, endLocation, CostCustom)
	if err != nil {
		return output, err
	}

	if DEBUG {
		fmt.Println(PrintSim(sim, path))
	}
	fmt.Println("Total:", cost)
	output = append(output, fmt.Sprintf("Total: %.0f", cost))
	return output, nil
}

func CostCustom(dir simulation.Direction, from, to simulation.Coord) float64 {
	// If the direction is different from the previous direction, add 1000 to the cost
	if from.DirectionTo(to).VX != dir.VX && from.DirectionTo(to).VY != dir.VY {
		return 1000 + simulation.CostManhattan(dir, from, to)
	}

	return simulation.CostManhattan(dir, from, to)
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
