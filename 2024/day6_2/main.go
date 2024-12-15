package main

import (
	"day6/internal/aocUtils"
	"day6/internal/simulation"
	"fmt"
	"os"
	"strconv"
	"sync"

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
	switch dir {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	default:
		return dir // Return unchanged if not one of the four cardinal directions
	}
}
func solve1(sim simulation.Simulation, workerCount int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	var output []string
	var obstructionPositions []Coord

	width := sim.GetMap().GetWidth()
	height := sim.GetMap().GetHeight()

	// Collect all empty cells where an obstacle can be placed
	var emptyCells []Coord
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cell, err := sim.GetMap().GetCell(x, y)
			if err != nil {
				if DEBUG {
					fmt.Printf("Skipping cell (%d, %d) due to error: %v\n", x, y, err)
				}
				continue
			}
			if len(cell.GetEntityIds()) == 0 {
				emptyCells = append(emptyCells, Coord{x: x, y: y})
			}
		}
	}

	if DEBUG {
		fmt.Printf("Total empty cells to evaluate: %d\n", len(emptyCells))
	}

	// Set up concurrency
	type result struct {
		coord Coord
		loop  bool
	}

	jobs := make(chan Coord, len(emptyCells))
	resultsChan := make(chan result, len(emptyCells))
	errorsChan := make(chan error, len(emptyCells))

	var wg sync.WaitGroup

	// Worker function
	worker := func(workerId int) {
		defer wg.Done()
		for coord := range jobs {
			if DEBUG {
				fmt.Printf("Worker %d: Processing coordinate (%d,%d)\n", workerId, coord.x, coord.y)
			}
			// Clone the original simulation
			clonedSim := sim.Clone()

			// Create an obstacle entity
			obstacle, err := simulation.NewEntity(ObstacleEntityType)
			if err != nil {
				if DEBUG {
					fmt.Printf("Worker %d: error creating obstacle entity: %v\n", workerId, err)
				}
				// Send error to errorsChan
				errorsChan <- fmt.Errorf("worker %d: error creating obstacle entity at (%d,%d): %v", workerId, coord.x, coord.y, err)

				// Continue to next job
				continue
			}

			// Place the obstacle at the specified coordinate
			_, err = clonedSim.AddEntity(obstacle, coord.x, coord.y, 0, 0) // Obstacles don't need movement vectors
			if err != nil {
				errorsChan <- fmt.Errorf("worker %d: error adding obstacle at (%d,%d): %v", workerId, coord.x, coord.y, err)
				continue
			}

			// Run the simulation and check for loops
			loopDetected, err := detectLoop(workerId, clonedSim)
			if err != nil {
				errorsChan <- fmt.Errorf("worker %d: error running simulation: %v", workerId, err)
				continue
			}
			resultsChan <- result{coord: coord, loop: loopDetected}
		}
	}

	// Start worker goroutines
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(i)
		if DEBUG {
			fmt.Printf("Worker %d: Started.\n", i)
		}
	}

	// Send jobs
	for _, coord := range emptyCells {
		jobs <- coord
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(resultsChan)
	close(errorsChan)

	// Collect errors
	errSlice := []error{}
	for err := range errorsChan {
		errSlice = append(errSlice, err)
	}

	// If any errors were encountered, return the first one
	if len(errSlice) > 0 {
		return nil, errSlice[0]
	}

	// Collect results
	for res := range resultsChan {
		if res.loop {
			obstructionPositions = append(obstructionPositions, res.coord)
		}
	}

	if DEBUG {
		fmt.Printf("Total obstruction positions causing loops: %d\n", len(obstructionPositions))
		fmt.Printf("Obstruction Positions: %v\n", obstructionPositions)
	}

	output = append(output, fmt.Sprintf("Obstruction Positions: %d", len(obstructionPositions)))
	return output, nil
}

// detectLoop runs the simulation and determines if the guard enters a loop
func detectLoop(workerId int, sim simulation.Simulation) (bool, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	if DEBUG {
		fmt.Printf("Worker %d: Running detectLoop\n", workerId)
		fmt.Printf("%s\n", PrintSim(sim))
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
			vx, vy := entity.GetVector()
			guardDirection = Direction{vx: vx, vy: vy}
			break
		}
	}

	// If no guard found, cannot detect loop
	if guardID == uuid.Nil {
		return false, fmt.Errorf("no guard found on map")
	}

	if DEBUG {
		fmt.Printf("Worker %d: Guard Location: (%d, %d), Direction: (%d, %d)\n", workerId, guardLocation.x, guardLocation.y, guardDirection.vx, guardDirection.vy)
	}

	// State tracking: position + direction
	type State struct {
		x, y   int
		vx, vy int
	}
	visitedStates := make(map[State]bool)
	visitedStates[State{guardLocation.x, guardLocation.y, guardDirection.vx, guardDirection.vy}] = true

	// Simulation loop with loop detection
	for {
		// Calculate the position in front of the guard based on current direction
		newX := guardLocation.x + guardDirection.vx
		newY := guardLocation.y + guardDirection.vy

		if DEBUG {
			fmt.Printf("Worker %d: Moving Guard to (%d, %d)\n", workerId, newX, newY)
		}

		// Check if the new position is outside the map boundaries
		if !sim.GetMap().ValidateCoord(newX, newY) {
			// Guard has left the map; no loop
			if DEBUG {
				fmt.Printf("worker %d: Guard moved (%d, %d) outside the map boundaries (width: %d, height: %d)\n",
					workerId, newX, newY, sim.GetMap().GetWidth(), sim.GetMap().GetHeight())
			}
			return false, nil
		}

		// Check if the space in front is empty
		cell, err := sim.GetMap().GetCell(newX, newY)
		if err != nil {
			// Error accessing cell; assume no loop
			return false, fmt.Errorf("worker %d: Error accessing cell at (%d, %d): %v", workerId, newX, newY, err)
		}

		if cell.IsEmpty() {
			// Move the guard forward
			err := sim.MoveEntity(guardID, newX, newY, false)
			if err != nil {
				// Failed to move; assume no loop
				return false, fmt.Errorf("worker %d: Failed to move guard to (%d, %d): %v", workerId, newX, newY, err)
			}

			// Update guard's current location
			guardLocation = Coord{x: newX, y: newY}

			if DEBUG {
				fmt.Printf("Guard moved to (%d, %d)\n", newX, newY)
			}

			// Check for loop
			currentState := State{guardLocation.x, guardLocation.y, guardDirection.vx, guardDirection.vy}
			if visitedStates[currentState] {
				// Loop detected
				return true, nil
			}
			visitedStates[currentState] = true
		} else {
			// Turn the guard to the right
			guardDirection = TurnRight(guardDirection)

			// Update the guard's direction vector in the simulation
			err := sim.SetEntityVector(guardID, guardDirection.vx, guardDirection.vy)
			if err != nil {
				// Error updating direction; assume no loop
				return false, fmt.Errorf("worker %d: Error setting entity vector for guard at (%d, %d): %v", workerId, newX, newY, err)
			}

			// Check for loop after turning
			currentState := State{guardLocation.x, guardLocation.y, guardDirection.vx, guardDirection.vy}
			if visitedStates[currentState] {
				// Loop detected
				return true, nil
			}
			visitedStates[currentState] = true
		}
	}
}
