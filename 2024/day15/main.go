package main

import (
	"day15/internal/aocUtils"
	"day15/internal/simulation"
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

	lines, err := aocUtils.ReadInput(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error reading from %s: %v", INPUT_FILE, err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// SOLUTION LOGIC
	////////////////////////////////////////////////////////////////////

	input, actions, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve1(input, actions, PARALLELISM)
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

func CalculateDirection(s string) (simulation.Direction, error) {
	switch s {
	case "^":
		return simulation.North, nil
	case ">":
		return simulation.East, nil
	case "v":
		return simulation.South, nil
	case "<":
		return simulation.West, nil
	default:
		return simulation.Direction{VX: 0, VY: 0}, fmt.Errorf("unknown direction: %s", s)
	}
}

const FishEntityType = "@"
const ObstacleEntityType = "#"
const BoxEntityType = "O"

func parseLines(lines []string) (simulation.Simulation, []simulation.Direction, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, nil, fmt.Errorf("input is empty")
	}

	var blankLineNum int
	// Iterate through the lines till a blank line is found
	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			blankLineNum = i
			break
		}
	}

	// Lines above the blank line are the map, Lines below the blank line are the series of actions to perform
	// Skip the first line as it represents the north wall of the map
	// Skip the first and last characters of each line as they represent the walls of the mam
	width := len(lines[1]) - 2
	height := blankLineNum - 2
	sim := simulation.NewSimulation(width, height)
	for y := 1; y <= height; y++ {
		for x := 1; x <= width; x++ {
			entityType := string(lines[y][x])
			if entityType == "." {
				continue
			}
			entity, err := simulation.NewEntity(entityType)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating fish entity: %v", err)
			}
			sim.AddEntity(entity, simulation.Coord{X: x - 1, Y: y - 1}, simulation.North)
		}
	}

	// Parse the actions
	actions := []simulation.Direction{}
	for _, line := range lines[blankLineNum+1:] {
		for _, char := range line {
			direction, err := CalculateDirection(string(char))
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing action: %v", err)
			}
			actions = append(actions, direction)
		}
	}

	return sim, actions, nil
}

func PrintSim(sim simulation.Simulation) string {
	var output string
	for y := 0; y < sim.GetMap().GetHeight(); y++ {
		for x := 0; x < sim.GetMap().GetWidth(); x++ {
			cell, err := sim.GetMap().GetCell(simulation.Coord{X: x, Y: y})
			if err != nil {
				output += "?"
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

// solve1 attempts to execute a series of actions, moving the fish and any boxes as needed.
func solve1(sim simulation.Simulation, actions []simulation.Direction, workerCount int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	var output []string

	width := sim.GetMap().GetWidth()
	height := sim.GetMap().GetHeight()

	if DEBUG {
		fmt.Printf("Starting map with dimensions (%d, %d)\n", width, height)
		fmt.Printf("%s\n", PrintSim(sim))
		fmt.Printf("Actions: %v\n", actions)
	}

	// Find the coordinates of the fish
	fish := findFish(sim)
	if fish == nil {
		return nil, fmt.Errorf("fish entity not found in the simulation")
	}
	fishPosition := fish.GetPosition()

	for _, direction := range actions {
		if DEBUG {
			fmt.Printf("Taking action: %v\n", direction)
		}
		// Calculate the next position the fish wants to move to
		nextPosition := fishPosition.Move(direction)
		if !sim.GetMap().ValidateCoord(nextPosition) {
			if DEBUG {
				fmt.Printf("Skipping action due to leaving the map at position %s\n", nextPosition.String())
			}
			continue
		}

		// Attempt to move boxes recursively if necessary
		canMove, err := canMoveBoxesRecursive(sim, nextPosition, direction, DEBUG)
		if err != nil {
			return nil, err
		}
		if !canMove {
			if DEBUG {
				fmt.Printf("Cannot move in direction %v due to obstacles or map boundaries\n", direction)
			}
			continue
		}

		// If canMove is true, proceed to move all boxes
		if err := moveBoxesRecursive(sim, nextPosition, direction, DEBUG); err != nil {
			return nil, err
		}

		// Now, move the fish
		err = sim.MoveEntity(fish.GetId(), nextPosition, false)
		if err != nil {
			return nil, fmt.Errorf("error moving fish to position %s: %v", nextPosition.String(), err)
		}
		fishPosition = nextPosition

		if DEBUG {
			fmt.Printf("Moved fish to %s\n", nextPosition.String())
			fmt.Printf("%s\n", PrintSim(sim))
		}
	}

	output = append(output, fmt.Sprintf("Total: %d", calculateTotal(sim)))
	return output, nil
}

// canMoveBoxesRecursive checks if all boxes in the specified direction can be moved.
func canMoveBoxesRecursive(sim simulation.Simulation, position simulation.Coord, direction simulation.Direction, DEBUG bool) (bool, error) {
	if !sim.GetMap().ValidateCoord(position) {
		if DEBUG {
			fmt.Printf("Cannot move: position %s is out of bounds\n", position.String())
		}
		return false, nil
	}

	cell, err := sim.GetMap().GetCell(position)
	if err != nil {
		return false, fmt.Errorf("error retrieving cell at position %s: %v", position.String(), err)
	}

	if cell.IsEmpty() {
		// Base case: no box to move, so movement is possible
		return true, nil
	}

	entityID := cell.GetEntityIds()[0]
	entity, err := sim.GetEntity(entityID)
	if err != nil {
		return false, fmt.Errorf("error retrieving entity at position %s: %v", position.String(), err)
	}

	switch entity.GetEntityType() {
	case ObstacleEntityType:
		if DEBUG {
			fmt.Printf("Cannot move: obstacle at position %s\n", position.String())
		}
		return false, nil
	case BoxEntityType:
		// Calculate the next position in the same direction
		nextPosition := position.Move(direction)
		// Recursively check if the next box can be moved
		return canMoveBoxesRecursive(sim, nextPosition, direction, DEBUG)
	default:
		// Unknown entity type; treat as non-movable
		if DEBUG {
			fmt.Printf("Cannot move: unknown entity type '%s' at position %s\n", entity.GetEntityType(), position.String())
		}
		return false, nil
	}
}

// moveBoxesRecursive moves all boxes in the specified direction.
func moveBoxesRecursive(sim simulation.Simulation, position simulation.Coord, direction simulation.Direction, DEBUG bool) error {
	cell, err := sim.GetMap().GetCell(position)
	if err != nil {
		return fmt.Errorf("error retrieving cell at position %s: %v", position.String(), err)
	}

	if cell.IsEmpty() {
		// Base case: no box to move
		return nil
	}

	entityID := cell.GetEntityIds()[0]
	entity, err := sim.GetEntity(entityID)
	if err != nil {
		return fmt.Errorf("error retrieving entity at position %s: %v", position.String(), err)
	}

	if entity.GetEntityType() != BoxEntityType {
		return fmt.Errorf("entity at position %s is not a box", position.String())
	}

	// Calculate the next position in the same direction
	nextPosition := position.Move(direction)

	// Recursively move the next boxes first
	if err := moveBoxesRecursive(sim, nextPosition, direction, DEBUG); err != nil {
		return err
	}

	// Now, move the current box
	err = sim.MoveEntity(entity.GetId(), nextPosition, false)
	if err != nil {
		return fmt.Errorf("error moving box from %s to %s: %v", position.String(), nextPosition.String(), err)
	}
	if DEBUG {
		fmt.Printf("Moved box from %s to %s\n", position.String(), nextPosition.String())
	}
	return nil
}

func findFish(sim simulation.Simulation) simulation.Entity {
	for _, entity := range sim.GetEntities() {
		if entity.GetEntityType() == FishEntityType {
			return entity
		}
	}
	return nil
}

func calculateTotal(sim simulation.Simulation) int {
	// The total of a box is equal to 100 times its distance from the top edge of the map plus its distance from the left edge of the map
	total := 0
	for _, entity := range sim.GetEntities() {
		if entity.GetEntityType() == BoxEntityType {
			topEdge := entity.GetPosition().Y + 1
			leftEdge := entity.GetPosition().X + 1
			total += (100 * topEdge) + leftEdge
		}
	}
	return total
}
