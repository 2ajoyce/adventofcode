package main

import (
	"bufio"
	"day6/internal"
	"day6/internal/directions"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func ReadInput(INPUT_FILE string) ([]string, error) {
	inputFile, err := os.Open(INPUT_FILE)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %v", INPUT_FILE, err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func WriteOutput(OUTPUT_FILE string, results []string) error {
	outputFile, err := os.Create(OUTPUT_FILE)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", OUTPUT_FILE, err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	// Write the results to output.txt, one line per result
	for i, res := range results {
		_, err := writer.WriteString(res)
		if err != nil {
			return fmt.Errorf("error writing value to %s: %v", OUTPUT_FILE, err)
		}
		if i != len(results)-1 {
			_, err = writer.WriteString("\n")
			if err != nil {
				return fmt.Errorf("error writing newline to %s: %v", OUTPUT_FILE, err)
			}
		}
	}

	// Flush the writer to ensure all data is written to output.txt
	writer.Flush()
	return nil
}

type LoopError struct {
	Message string
}

func (e LoopError) Error() string {
	return e.Message
}

func main() {
	//os.Setenv("DEBUG", "true")
	OVERFLOW_LIMIT := 10000
	INPUT_FILE := os.Getenv("INPUT_FILE")
	OUTPUT_FILE := os.Getenv("OUTPUT_FILE")
	PARALLELISM, err := strconv.Atoi(os.Getenv("PARALLELISM"))
	if PARALLELISM < 1 || err != nil {
		PARALLELISM = 1
	}
	PARALLELISM = 10

	if INPUT_FILE == "" || OUTPUT_FILE == "" {
		fmt.Println("INPUT_FILE and OUTPUT_FILE environment variables not set")
		fmt.Println("Defaulting to input.txt and output.txt")
		INPUT_FILE = "input.txt"
		OUTPUT_FILE = "output.txt"
	}

	lines, err := ReadInput(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error reading from %s: %v", INPUT_FILE, err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// Start Solution Logic  ///////////////////////////////////////////
	////////////////////////////////////////////////////////////////////

	// Create an array of all coordinates containing the letter X
	gridMap, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve1(gridMap, OVERFLOW_LIMIT, false)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}
	gridMap, err = parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results2, err := solve2(gridMap, OVERFLOW_LIMIT, false, PARALLELISM)
	results = append(results, results2...)
	if err != nil {
		fmt.Println("Error solving 2:", err)
		return
	}
	////////////////////////////////////////////////////////////////////
	// End Solution Logic  /////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////

	err = WriteOutput(OUTPUT_FILE, results)
	if err != nil {
		fmt.Printf("Error writing to %s: %v", OUTPUT_FILE, err)
		return
	}

	fmt.Printf("Successfully processed %s and created %s", INPUT_FILE, OUTPUT_FILE)
}

func parseLines(lines []string) (internal.Gridmap, error) {
	DEBUG := os.Getenv("DEBUG")

	width := len(lines[0])
	height := len(lines)

	gridMap := internal.NewGridmap(width, height)

	guardCharacters := []rune{'^', 'v', '<', '>'}
	obstructionCharacter := rune('#')

	for y, line := range lines {
		// For every character in the line
		for x, char := range line {
			// if the character is a guard character
			if char == guardCharacters[0] || char == guardCharacters[1] || char == guardCharacters[2] || char == guardCharacters[3] {
				direction := directions.N // Default to north, but this should ALWAYS be set by the switch statement below.
				// Set the direction the guard is facing based on the char
				switch char {
				case '^':
					direction = directions.N
				case 'v':
					direction = directions.S
				case '<':
					direction = directions.W
				case '>':
					direction = directions.E
				default:
					return nil, fmt.Errorf("invalid guard character: %c", char)
				}
				newCoord := internal.NewCoord(x, y)
				newGuard, err := internal.NewGuard(newCoord, direction, char)
				if err != nil {
					return nil, fmt.Errorf("error creating guard at %s: %v", newCoord, err)
				}

				gridMap.SetGuards(append(gridMap.Guards(), newGuard))
			} else if char == obstructionCharacter {
				newObstruction, err := internal.NewObject(internal.NewCoord(x, y), obstructionCharacter)
				if err != nil {
					return nil, fmt.Errorf("error creating obstruction at (%d,%d): %v", x, y, err)
				}
				gridMap.SetObstructions(append(gridMap.Obstructions(), newObstruction))
			}
		}
	}
	if DEBUG == "true" {
		fmt.Printf("Parsed Input\n%s\n", gridMap.String())
	}
	return gridMap, nil
}

func simulate(gridMap internal.Gridmap, overflowLimit int, visualize bool) ([]internal.DirectedObject, error) {
	DEBUG := os.Getenv("DEBUG")

	// Simulate the guard moving
	// The guard moves in the direction it is facing until it reaches an obstruction or the edge of the grid.
	// Cycle through the guards, moving each one, till no more movement occurs.
	moved := false
	visitedDirectedObjects := []internal.DirectedObject{}
	pathMask := internal.NewGridmask([]internal.Coord{}, 'X')
	moves := 0

	for {
		moved = false

		for i := 0; i < len(gridMap.Guards()); i++ {
			guard := gridMap.Guards()[i]
			visitedDirectedObject, err := internal.NewGenericDirectedObject(guard.Location(), guard.FacingDirection(), 'X')
			if err != nil {
				return nil, fmt.Errorf("error creating visited directed object for guard at %s: %v", guard.Location(), err)
			}
			visitedDirectedObjects = append(visitedDirectedObjects, visitedDirectedObject)

			if DEBUG == "true" {
				fmt.Printf("Checking guard at %s facing %s\n", guard.Location(), guard.FacingDirection())
			}

			// Check if the guard can move in its current direction
			newCoords, err := guard.FacingCoord()
			if err != nil {
				return nil, err
			}
			if DEBUG == "true" {
				fmt.Printf("The guard would move to %s\n", newCoords)
			}
			// Check if the new coordinates are outside bounds
			if !gridMap.ValidateCoord(newCoords) {
				if DEBUG == "true" {
					fmt.Printf("Guard can not move from %s to %s as it is out of bounds\n", guard.Location(), newCoords)
				}
				continue
			}
			obstructed := false
			// Check if the new coordinates are an obstruction
			for _, obstruction := range gridMap.Obstructions() {
				if obstruction.Location() == newCoords {
					if DEBUG == "true" {
						fmt.Printf("Guard can not move from %s to %s as it is an obstruction\n", guard.Location(), newCoords)
					}
					obstructed = true
					continue
				}
			}
			if obstructed {
				guard = guard.TurnRight() // This assigns to the local variable, not to the object
				if DEBUG == "true" {
					fmt.Printf("Guard turned right to %s\n", guard.FacingDirection())
				}
			}
			if DEBUG == "true" {
				fmt.Printf("Moving guard from %s to %s\n", guard.Location(), newCoords)
			}
			// Move the guard forward
			guard, err = guard.Move()
			if err != nil {
				return nil, fmt.Errorf("error moving guard forward from %s: %v", guard.Location(), err)
			}

			moved = true
			moves++
			guards := gridMap.Guards()
			guards[i] = guard
			gridMap.SetGuards(guards)

			if visualize {
				// // Clear the screen and move cursor to the top-left corner
				fmt.Printf("\033[2J\033[H")
				// // Extract the visitedCoord locations to use as a mask
				visitedCoords := []internal.Coord{}
				for _, v := range visitedDirectedObjects {
					visitedCoords = append(visitedCoords, v.Location())
				}
				pathMask.SetLocations(visitedCoords)
				fmt.Printf("Move %d\n%s\n", moves, gridMap.String(pathMask))

				time.Sleep(5 * time.Millisecond)
			}

		}

		visited := make(map[internal.Coord]directions.Direction)
		for _, location := range visitedDirectedObjects {
			// Check to see if any locations have been visited more than once traveling in the same direction
			if _, ok := visited[location.Location()]; ok && visited[location.Location()] == location.FacingDirection() {
				return nil, LoopError{"loop detected"}
			}
			visited[location.Location()] = location.FacingDirection()
		}

		// Check if the guard has moved or if we have reached the overflow limit
		if !moved || moves > overflowLimit {
			break
		}
	}

	if DEBUG == "true" {
		fmt.Printf("Total moves: %d\n", moves)
	}

	return visitedDirectedObjects, nil
}

func solve1(gridMap internal.Gridmap, overflowLimit int, visualize bool) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG")

	results := []string{}

	visitedCoords, err := simulate(gridMap, overflowLimit, visualize)
	if err != nil {
		return nil, err
	}

	// Count the number of times each coordinate was visited
	coordCount := map[internal.Coord]int{}
	for _, coord := range visitedCoords {
		coordCount[coord.Location()]++
	}

	// Count the number of distinct coordinates visited by guards
	distinctCoordCount := len(coordCount)

	results = append(results, fmt.Sprintf("Distinct Moves: %s", strconv.Itoa(distinctCoordCount)))

	return results, nil
}

func findLoops(gridmap internal.Gridmap, visitedCoords []internal.DirectedObject) ([]internal.Coord, error) {
	// I believe that there is a more efficient way to do this, but these are daily tasks and I'm out of time
	// findLoops clears the tests, but it fails to work on the larger grid. I've solved several issues, but apparently
	// there are still bugs in my logic.

	DEBUG := os.Getenv("DEBUG")

	// Count the number of times each coordinate was visited
	coordCount := map[internal.Coord]int{}
	for _, coord := range visitedCoords {
		coordCount[coord.Location()]++
	}

	// Count the number of coordinates visited more than one time
	duplicateCoordCount := 0
	for _, count := range coordCount {
		if count > 1 {
			duplicateCoordCount++
		}
	}

	turningPoints := map[internal.Coord]bool{}

	for i := 0; i < len(visitedCoords); i++ {
		for j := 0; j < i; j++ {
			currentLocation := visitedCoords[i].Location()
			currentDirection := visitedCoords[i].FacingDirection()
			previousLocation := visitedCoords[j].Location()
			previousDirection := visitedCoords[j].FacingDirection()
			if DEBUG == "true" {
				fmt.Printf("Checking if %s, %s and %s, %s are on the same line\n", currentLocation, currentDirection, previousLocation, previousDirection)
			}
			// If the visitedCoords[i].facing is N or S then the y axis has to be the matching coordinate
			// If the visitedCoords[i].facing is E or W then the x axis has to be the matching coordinate
			if (currentDirection == directions.N || currentDirection == directions.S) && currentLocation.Y() != previousLocation.Y() {
				continue
			}
			if (currentDirection == directions.E || currentDirection == directions.W) && currentLocation.X() != previousLocation.X() {
				continue
			}
			if i != j && (currentLocation.X() == previousLocation.X() || currentLocation.Y() == previousLocation.Y()) {
				if DEBUG == "true" {
					fmt.Printf("Checking if %s, %s and %s, %s are different directions\n", currentLocation, currentDirection, previousLocation, previousDirection)
				}

				// Get the coordinates of the space one move from the previous location in the direction of a right turn from the current direction
				// If that space contains an obstacle, then it will cause a loop
				possibleObstacle, err := previousLocation.Move(currentDirection.TurnRight())
				if err != nil {
					return nil, err
				}
				hasObstacle := false
				for _, obstacle := range gridmap.ObstructionLocations() {
					if possibleObstacle == obstacle {
						hasObstacle = true
						break
					}
				}
				if !hasObstacle {
					continue
				}

				// Check if the space between visitedCoords[i] and visitedCoords[j] is clear of obstacles
				clear := true
				for _, obstruction := range gridmap.Obstructions() {
					obstructionLocation := obstruction.Location()
					if currentLocation.X() == previousLocation.X() && currentLocation.X() == obstructionLocation.X() {
						if obstructionLocation.BetweenY(currentLocation, previousLocation) {
							clear = false
							break
						}
					}
					if currentLocation.Y() == previousLocation.Y() && currentLocation.Y() == obstructionLocation.Y() {
						if obstructionLocation.BetweenX(currentLocation, previousLocation) {
							clear = false
							break
						}
					}
				}

				if clear {
					turningPoint := internal.NewCoord(currentLocation.X(), currentLocation.Y())

					// Move the turning point one step in the direction of the guard's movement to correctly place it
					turningPoint, err := turningPoint.Move(currentDirection)
					if err != nil {
						return nil, err
					}
					turningPoints[turningPoint] = true

					if DEBUG == "true" {
						fmt.Printf("Found a TurningPoint at %s\n", turningPoint)
					}
				}
			}
		}
	}

	points := []internal.Coord{}
	for point := range turningPoints {
		points = append(points, point)
	}

	if DEBUG == "true" {
		obsMask := internal.NewGridmask(points, 'O')
		fmt.Println(gridmap.String(obsMask))
		fmt.Printf("Number of turning points identified: %d\n", len(turningPoints))
		fmt.Printf("Turning points:\n")
		for k, v := range turningPoints {
			if v {
				fmt.Printf("%s\n", k)
			}
		}
	}

	return points, nil
}

func findLoopsBruteForce(gridmap internal.Gridmap, visitedCoords []internal.DirectedObject, parallelism int) ([]internal.Coord, error) {
	DEBUG := os.Getenv("DEBUG") == "true"

	// Channel to collect turning points
	turningPointsCh := make(chan internal.Coord, len(visitedCoords))
	// Channel to collect errors
	errCh := make(chan error, 1)
	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup
	// Mutex to protect the turningPoints map
	var mu sync.Mutex
	// Map to store unique turning points
	turningPointsMap := make(map[internal.Coord]bool)

	// Initialize progress bar
	bar := progressbar.Default(int64(len(visitedCoords)))

	// Create a buffered channel for tasks
	tasksCh := make(chan internal.DirectedObject, len(visitedCoords))

	// Start worker goroutines
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for possibleLocation := range tasksCh {
				// Add obstacle to the gridmap
				if DEBUG {
					fmt.Printf("[Worker %d] Adding obstacle at %s\n", workerID, possibleLocation.Location())
				}
				updatedGridmap := gridmap.Clone()
				updatedGridmap.SetObstructions(append(updatedGridmap.Obstructions(), possibleLocation))

				// Run the simulation
				_, err := simulate(updatedGridmap, 10000, false)
				if err != nil {
					if DEBUG {
						fmt.Printf("[Worker %d] Error running simulation with obstacle at %s: %v\n", workerID, possibleLocation.Location(), err)
					}
					// Check if the error is of type LoopError
					if _, ok := err.(LoopError); ok {
						if DEBUG {
							fmt.Printf("[Worker %d] Found a TurningPoint at %s\n", workerID, possibleLocation.Location())
						}
						// Send the turning point to the channel
						turningPointsCh <- possibleLocation.Location()
					} else {
						// Send the first non-LoopError encountered
						select {
						case errCh <- err:
						default:
						}
						return
					}
				}
				// Update the progress bar
				bar.Add(1)
			}
		}(i + 1)
	}

	// Send tasks to the workers
	go func() {
		// Skip the first element as it's the starting point
		for _, possibleLocation := range visitedCoords {
			tasksCh <- possibleLocation
		}
		if DEBUG {
			fmt.Println("[Main] All tasks sent to workers")
		}
		close(tasksCh)
	}()

	// Goroutine to wait for all workers and then close the turningPointsCh
	go func() {
		wg.Wait()
		if DEBUG {
			fmt.Println("[Main] All workers have finished")
		}
		close(turningPointsCh)
	}()

	// Collect turning points and handle errors
	for {
		select {
		case point, ok := <-turningPointsCh:
			if !ok {
				fmt.Println("[Main] All turning points collected")
				fmt.Println("[Main] Closing turningPointsCh")

				// Channel closed, all turning points received
				turningPointsCh = nil
				errCh = nil

			} else {
				if DEBUG {
					fmt.Printf("[Main] Received TurningPoint at %s\n", point)
					fmt.Println("[Main] Locking Mutex")
				}
				mu.Lock()
				turningPointsMap[point] = true
				if DEBUG {
					fmt.Println("[Main] Unlocking Mutex")
				}
				mu.Unlock()
			}
		case err, ok := <-errCh:
			if ok {
				// An error occurred, return immediately
				return nil, err
			}
			errCh = nil
		}
		if turningPointsCh == nil && errCh == nil {
			break
		}
	}

	// Convert the map to a slice
	points := make([]internal.Coord, 0, len(turningPointsMap))
	for point := range turningPointsMap {
		points = append(points, point)
	}

	return points, nil
}

func solve2(gridMap internal.Gridmap, overflowLimit int, visualize bool, parallelism int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"

	results := []string{}

	simulationGridmap := gridMap.Clone()
	visitedCoords, err := simulate(simulationGridmap, overflowLimit, visualize)
	if err != nil {
		return nil, err
	}

	turningPointsGridmap := gridMap.Clone()
	//turningPoints, err := findLoops(turningPointsGridmap, visitedCoords)
	turningPoints, err := findLoopsBruteForce(turningPointsGridmap, visitedCoords, parallelism)
	if err != nil {
		return nil, err
	}
	if DEBUG {
		fmt.Printf("Turning Points: %v\n", turningPoints)
	}

	results = append(results, fmt.Sprintf("Turning Points: %s", strconv.Itoa(len(turningPoints))))

	return results, nil
}
