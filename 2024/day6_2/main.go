package main

import (
	"day6/internal/aocUtils"
	"day6/internal/simulation"
	"fmt"
	"os"
	"strconv"

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
			if lines[y][x] == '.' { // Skip empty cells
				continue
			}
			unitLetter := string(lines[y][x])
			entity, err := simulation.NewEntity(unitLetter)
			if err != nil {
				return nil, fmt.Errorf("error creating entity for (%d,%d): %v", x, y, err)
			}
			_, err = sim.AddEntity(entity, x, y, North.vx, North.vy)
			if err != nil {
				return nil, fmt.Errorf("error adding entity to sim at (%d,%d): %v", x, y, err)
			}

		}

	}
	return sim, nil
}

func PrintSim(sim simulation.Simulation) string {
	var output string
	for y := 0; y < sim.GetMap().GetHeight(); y++ {
		for x := 0; x < sim.GetMap().GetWidth(); x++ {
			cell, err := sim.GetMap().GetCell(x, y)
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

type Coord struct {
	x, y int
}

const GuardEntityType = "^"
const ObstacleEntityType = "#"

type Direction struct {
	vx, vy int
}

// Define the four cardinal directions
var (
	North = Direction{vx: 0, vy: -1}
	East  = Direction{vx: 1, vy: 0}
	South = Direction{vx: 0, vy: 1}
	West  = Direction{vx: -1, vy: 0}
)

// TurnRight rotates the direction 90 degrees clockwise
func TurnRight(dir Direction) Direction {
	return Direction{vx: dir.vy, vy: -dir.vx}
}
func solve1(sim simulation.Simulation, workerCount int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	var output []string
	positions := make(map[Coord]int) // Record the count of times each position is visited

	if DEBUG {
		fmt.Printf("Starting simulation with %d workers...\n", workerCount)
		fmt.Println(PrintSim(sim))
	}

	// Helper function to turn right
	turnRight := func(current Direction) Direction {
		switch current {
		case North:
			return East
		case East:
			return South
		case South:
			return West
		case West:
			return North
		default:
			return current
		}
	}

	// Identify the guard entity
	var guardID uuid.UUID
	var guardLocation Coord
	var guardDirection Direction
	entities := sim.GetEntities()
	for _, entity := range entities {
		if entity.GetEntityType() == GuardEntityType {
			guardID = entity.GetId()
			x, y := entity.GetPosition()
			guardLocation = Coord{x: x, y: y}
			vx, vy := entity.GetVector() // Returns vx, vy as integers
			guardDirection = Direction{vx: vx, vy: vy}
			break
		}
	}

	// Check if guard was found
	if guardID == uuid.Nil {
		return nil, fmt.Errorf("no guard entity found in the simulation")
	}

	if DEBUG {
		fmt.Println("Guard Location:", guardLocation)
		fmt.Println("Guard Direction:", guardDirection)
	}

	positions[guardLocation]++ // Starting position of the guard is counted as a visited position

	// Simulation loop
	for {
		// Calculate the position in front of the guard based on current direction
		newX := guardLocation.x + guardDirection.vx
		newY := guardLocation.y + guardDirection.vy

		// Check if the new position is outside the map boundaries
		if newX < 0 || newX >= sim.GetMap().GetWidth() || newY < 0 || newY >= sim.GetMap().GetHeight() {
			if DEBUG {
				fmt.Println("Guard has left the map.")
			}
			break // Guard has left the map; end simulation
		}

		// Check if the space in front is empty
		cell, err := sim.GetMap().GetCell(newX, newY)
		if err != nil {
			return nil, fmt.Errorf("error accessing cell (%d,%d): %v", newX, newY, err)
		}

		isEmpty := len(cell.GetEntityIds()) == 0

		if isEmpty {
			// Move the guard forward
			success, err := sim.MoveEntity(guardID, newX, newY)
			if err != nil {
				return nil, fmt.Errorf("error moving guard: %v", err)
			}
			if !success {
				return nil, fmt.Errorf("failed to move guard to (%d,%d)", newX, newY)
			}

			// Update guard's current location
			guardLocation = Coord{x: newX, y: newY}

			// Increment the visit count for the new position
			positions[guardLocation]++

			if DEBUG {
				fmt.Printf("Guard moved to (%d, %d)\n", newX, newY)
				fmt.Println(PrintSim(sim))
			}
		} else {
			// Turn the guard to the right
			guardDirection = turnRight(guardDirection)

			// Update the guard's direction vector in the simulation
			err := sim.SetEntityVector(guardID, guardDirection.vx, guardDirection.vy)
			if err != nil {
				return nil, fmt.Errorf("error updating guard direction: %v", err)
			}

			if DEBUG {
				fmt.Printf("Guard turned right. New direction: (%d, %d)\n", guardDirection.vx, guardDirection.vy)
				fmt.Println(PrintSim(sim))
			}

			// Do not move the guard on the same iteration it turned right
		}
	}

	// Calculate distinct and total positions visited
	distinctPositions := len(positions)
	totalPositions := 0
	for _, count := range positions {
		totalPositions += count
	}

	if DEBUG {
		fmt.Println("Total positions visited by the guard:", totalPositions)
		fmt.Println("Distinct positions visited by the guard:", distinctPositions)
	}

	output = append(output, fmt.Sprintf("Distinct Positions: %d", distinctPositions))
	return output, nil
}
