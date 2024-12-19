package main

import (
	"day15/internal/aocUtils"
	"day15/internal/simulation"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
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
	results, err := solve(input, actions)
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
const LeftBoxEntityType = "["
const RightBoxEntityType = "]"
const BoxEntityType = "O"

func parseLines(lines []string) (simulation.Simulation, []simulation.Direction, error) {
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

	var transformedLines []string
	for _, line := range lines {
		// If empty line, break
		if strings.TrimSpace(line) == "" {
			break
		}

		var transformedLine string
		for _, char := range line {
			switch char {
			case '#':
				transformedLine = transformedLine + "##"
			case 'O':
				transformedLine = transformedLine + "[]"
			case '.':
				transformedLine = transformedLine + ".."
			case '@':
				transformedLine = transformedLine + "@."
			}
		}
		transformedLine += "\n"
		transformedLines = append(transformedLines, transformedLine)
	}
	fmt.Println(transformedLines)

	// Lines above the blank line are the map, Lines below the blank line are the series of actions to perform
	// Skip the first line as it represents the north wall of the map
	// Skip the first 2 and last 2 characters of each line as they represent the walls of the mam
	width := len(transformedLines[0])
	height := blankLineNum
	sim := simulation.NewSimulation(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			entityType := string(transformedLines[y][x])
			if entityType == "." {
				continue
			}
			entity, err := simulation.NewEntity(entityType)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating fish entity: %v", err)
			}
			coords := []simulation.Coord{{X: x, Y: y}}
			if entityType == LeftBoxEntityType {
				continue // Skip the left box, add it when the right half is parsed
			}
			if entityType == RightBoxEntityType {
				entity, err = simulation.NewEntity(BoxEntityType)
				if err != nil {
					return nil, nil, fmt.Errorf("error creating box entity: %v", err)
				}
				coords = append(coords, simulation.Coord{X: x - 1, Y: y})
			}
			sim.AddEntity(entity, coords, simulation.North)
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
				if entity.GetEntityType() == BoxEntityType {
					if float64(x) == math.Max(float64(entity.GetPosition()[0].X), float64(entity.GetPosition()[1].X)) {
						output += string(RightBoxEntityType)
					} else {
						output += string(LeftBoxEntityType)
					}

				} else {
					output += string(entity.GetEntityType())
				}
			}
		}
		output += "\n"
	}
	return output
}

func solve(sim simulation.Simulation, actions []simulation.Direction) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	var output []string

	width := sim.GetMap().GetWidth()
	height := sim.GetMap().GetHeight()

	if DEBUG {
		fmt.Printf("Starting map with dimensions (%d, %d)\n", width, height)
		fmt.Printf("%s\n", PrintSim(sim))
		fmt.Printf("Actions: %v\n", actions)
	}

	fmt.Printf("Height: %d, Width: %d\n", height, width)
	fmt.Println(PrintSim(sim))

	// Find the coordinates of the fish
	fish := findFish(sim)
	if fish == nil {
		return nil, fmt.Errorf("fish entity not found in the simulation")
	}
	if DEBUG {
		fmt.Printf("Fish is at %s\n", fish.GetPosition()[0].String())
	}

	for _, direction := range actions {
		if DEBUG {
			fmt.Printf("Taking action: %v\n", direction)
		}

		canMove, err := canEntityMove(sim, fish, direction)
		if err != nil {
			return nil, fmt.Errorf("error checking if fish can move: %v", err)
		}
		if !canMove {
			if DEBUG {
				fmt.Printf("Cannot move fish. Skipping action %v\n", direction)
			}
			continue // If we can't move, continue to the next action
		}
		if DEBUG {
			fmt.Println("Fish can move")
		}

		moved, err := moveEntity(sim, fish, direction)
		if err != nil {
			return nil, fmt.Errorf("error moving fish: %v", err)
		}
		if !moved {
			if DEBUG {
				fmt.Printf("Did not move fish. Skipping action %v\n", direction)
			}
			return nil, fmt.Errorf("fish could move %v at %s, but did not move. This should never occur", direction, fish.GetPosition()[0].String())
		}

		fish, err = sim.GetEntity(fish.GetId())
		if err != nil {
			fmt.Printf("Error getting fish: %v\n", err)
		}
		if DEBUG {
			fmt.Printf("Moved fish to %v\n", fish.GetPosition())
		}
		if DEBUG {
			fmt.Printf("%s\n", PrintSim(sim))
		}
	}

	output = append(output, fmt.Sprintf("Total: %d", calculateTotal(sim)))
	return output, nil
}

func getEntityIds(sim simulation.Simulation, coords []simulation.Coord) ([]uuid.UUID, error) {
	var entityIds []uuid.UUID
	for _, coord := range coords {
		cell, err := sim.GetMap().GetCell(coord)
		if err != nil {
			return nil, fmt.Errorf("error getting cell at position %s: %v", coord.String(), err)
		}
		entityIds = append(entityIds, cell.GetEntityIds()...)
	}
	return entityIds, nil
}

func getEntities(sim simulation.Simulation, coords []simulation.Coord) ([]simulation.Entity, error) {
	var entities []simulation.Entity
	entityIds, err := getEntityIds(sim, coords)
	if err != nil {
		return nil, fmt.Errorf("error getting entity ids for coords %v", coords)
	}
	for _, entityId := range entityIds {
		entity, err := sim.GetEntity(entityId)
		if err != nil {
			return nil, fmt.Errorf("error getting entity for id %s", entityId)
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func findFish(sim simulation.Simulation) simulation.Entity {
	for _, entity := range sim.GetEntities() {
		if entity.GetEntityType() == FishEntityType {
			return entity
		}
	}
	return nil
}

// This recursive function will check if an object can move
func canEntityMove(sim simulation.Simulation, entity simulation.Entity, direction simulation.Direction) (bool, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	canMove := true
	switch entity.GetEntityType() {
	case ObstacleEntityType:
		if DEBUG {
			fmt.Printf("Found Obstacle EntityType at position: %v\n", entity.GetPosition())
		}
		canMove = false
	case FishEntityType, BoxEntityType:
		if DEBUG {
			fmt.Printf("Found %s EntityType at position: %v\n", entity.GetEntityType(), entity.GetPosition())
		}
		coords := entity.GetPosition()
		if DEBUG {
			fmt.Printf("Checking coordinates %v\n", coords)
		}
		nextCoords := []simulation.Coord{}
		for _, coord := range coords {
			if !sim.GetMap().ValidateCoord(coord) {
				return false, nil
			}
			nextCoord := simulation.Coord{X: coord.X, Y: coord.Y}
			nextCoord = nextCoord.Move(direction)
			if !sim.GetMap().ValidateCoord(nextCoord) {
				return false, nil

			}
			if slices.Contains(coords, nextCoord) {
				continue // Don't check your own coords
			}
			nextCoords = append(nextCoords, nextCoord)
		}
		entities, err := getEntities(sim, nextCoords)
		if err != nil {
			return false, fmt.Errorf("error getting entities for coords %v", nextCoords)
		}
		dedupedEntities := []simulation.Entity{}
		for _, e := range entities {
			if !slices.Contains(dedupedEntities, e) {
				dedupedEntities = append(dedupedEntities, e)
			}
		}
		if DEBUG {
			fmt.Printf("Found %d entities at coords %v\n", len(dedupedEntities), nextCoords)
		}
		for _, e := range dedupedEntities {
			success, err := canEntityMove(sim, e, direction)
			if err != nil {
				return false, fmt.Errorf("error checking if entity can move for direction %v", direction)
			}
			if !success {
				canMove = false
				break
			}
		}
	}
	return canMove, nil
}

func moveEntity(sim simulation.Simulation, entity simulation.Entity, direction simulation.Direction) (bool, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	switch entity.GetEntityType() {
	case ObstacleEntityType:
		if DEBUG {
			fmt.Printf("Tried to move Obstacle EntityType at position: %v\n", entity.GetPosition())
		}
		return false, nil // Obstacles can not move
	case FishEntityType, BoxEntityType:
		if DEBUG {
			fmt.Printf("Tried to move %s EntityType at position: %v\n", entity.GetEntityType(), entity.GetPosition())
		}
		coords := entity.GetPosition()
		nextCoords := []simulation.Coord{}
		for _, coord := range coords {
			if !sim.GetMap().ValidateCoord(coord) {
				return false, nil // Can not move to invalid coords
			}
			nextCoord := simulation.Coord{X: coord.X, Y: coord.Y}
			nextCoords = append(nextCoords, nextCoord.Move(direction))
		}
		entities, err := getEntities(sim, nextCoords)
		if err != nil {
			return false, fmt.Errorf("error getting entities for coords %v", coords)
		}
		dedupedEntities := []simulation.Entity{}
		for _, e := range entities {
			if !slices.Contains(dedupedEntities, e) {
				dedupedEntities = append(dedupedEntities, e)
			}
		}
		for _, e := range dedupedEntities {
			if entity.GetId() == e.GetId() {
				continue // Don't try to move yourself
			}
			success, err := moveEntity(sim, e, direction)
			if err != nil {
				return false, fmt.Errorf("error in moveEntity for entity type %s at position: %v", e.GetEntityType(), e.GetPosition())

			}
			if !success {
				return false, nil // If any entity can not move, return false
			}
		}
		// If every location is valid and every entity in those locations can move
		err = sim.SetEntityDirection(entity.GetId(), direction)
		if err != nil {
			return false, fmt.Errorf("error setting direction for entity %s at position %v", entity.GetEntityType(), entity.GetPosition())
		}
		err = sim.MoveEntity(entity.GetId(), false)
		if err != nil {
			return false, fmt.Errorf("error in MoveEntity for entity %s at position %v", entity.GetEntityType(), entity.GetPosition())

		}
	}
	return true, nil
}

func calculateTotal(sim simulation.Simulation) int {
	DEBUG := os.Getenv("DEBUG") == "true"
	total := 0
	fmt.Println("Calculating Total")
	for _, entity := range sim.GetEntities() {
		if entity.GetEntityType() == BoxEntityType {
			coords := entity.GetPosition()
			var topEdge = 100_000_000_000
			var leftEdge = 100_000_000_000
			var rightEdge = 0
			for _, coord := range coords {
				topEdge = int(math.Min(float64(topEdge), float64(coord.Y)))
				leftEdge = int(math.Min(float64(leftEdge), float64(coord.X)))
				rightEdge = int(math.Max(float64(rightEdge), float64(coord.X)))
			}

			subtotal := (100 * topEdge) + leftEdge
			total += subtotal

			if DEBUG {
				rightEdge = sim.GetMap().GetWidth() - rightEdge // Not needed, but keeping because it helps with debugging
				fmt.Printf("Entity ID: %s, Top Edge: %d, Left Edge: %d, Right Edge: %d\n", entity.GetId().String(), topEdge, leftEdge, rightEdge)
				fmt.Printf("Entity ID: %s, Position: %v, Subtotal: %d\n", entity.GetId().String(), entity.GetPosition(), subtotal)
				fmt.Printf("Entity ID: %s, Position: %v, Total: %d\n", entity.GetId().String(), entity.GetPosition(), total)
				fmt.Println()
			}
		}
	}
	return total
}
