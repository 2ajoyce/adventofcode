package main

import (
	"day13/internal/aocUtils"
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

type ButtonA = [2]int64
type ButtonB = [2]int64
type Coordinate = [2]int64
type Machine struct {
	ButtonA         ButtonA
	ButtonB         ButtonB
	PrizeLocation   Coordinate
	currentPosition Coordinate
}

func NewMachine(buttonA ButtonA, buttonB ButtonB, prizeLocation Coordinate) *Machine {
	var m = new(Machine)
	m.ButtonA = buttonA
	m.ButtonB = buttonB
	m.PrizeLocation = prizeLocation
	m.currentPosition = [2]int64{0, 0}
	return m
}

// Move accepts two integers indicating how many times each button was pressed
// Move returns two boolean values indicating whether the machine moved past its X and Y axis
// Move will only update the current position if both pastX and pastY return false
func (m *Machine) Move(buttonA int, buttonB int) (pastX bool, pastY bool) {
	// Initialize the named return values
	newPosition := [2]int64{m.currentPosition[0], m.currentPosition[1]}
	pastX = false
	pastY = false
	// Increment the new location with the buttons
	for i := 0; i < buttonA; i++ {
		newPosition[0] += m.ButtonA[0]
	}
	for i := 0; i < buttonB; i++ {
		newPosition[1] += m.ButtonB[1]
	}
	// Check if the the new position has gone past the prize location
	if newPosition[0] > m.PrizeLocation[0] {
		pastX = true
	}
	if newPosition[1] > m.PrizeLocation[1] {
		pastY = true
	}
	return pastX, pastY
}

func (m *Machine) Reset() {
	m.currentPosition = Coordinate{0, 0} // reset to origin
}

func parseLines(lines []string) ([]*Machine, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	machines := []*Machine{}
	var buttonA ButtonA
	var buttonB ButtonB
	var prizeLocation Coordinate
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue // skip empty lines
		}
		switch line[7] {
		case 'A':
			buttonA = [2]int64{0, 0}
			parts := strings.Split(line[9:], ",")
			b, err := strconv.Atoi(strings.TrimPrefix(parts[0], " X+")) // This has a space in front
			if err != nil {
				return nil, fmt.Errorf("failed to parse buttonA[X]: %v", err)
			}
			buttonA[0] = int64(b)
			b, err = strconv.Atoi(strings.TrimPrefix(parts[1], " Y+")) // This has a space in front
			if err != nil {
				return nil, fmt.Errorf("failed to parse buttonA[Y]: %v", err)
			}
			buttonA[1] = int64(b)
		case 'B':
			buttonB = [2]int64{0, 0}
			parts := strings.Split(line[12:], ",")
			b, err := strconv.Atoi(strings.TrimPrefix(parts[0], " X+")) // This has a space in front
			if err != nil {
				return nil, fmt.Errorf("failed to parse buttonB[X]: %v", err)
			}
			buttonB[0] = int64(b)
			b, err = strconv.Atoi(strings.TrimPrefix(parts[1], " Y+")) // This has a space in front
			if err != nil {
				return nil, fmt.Errorf("failed to parse buttonB[Y]: %v", err)
			}
			buttonB[1] = int64(b)
		case 'X':
			prizeLocation = [2]int64{0, 0}
			parts := strings.Split(line[7:], ",")
			p, err := strconv.Atoi(strings.TrimPrefix(parts[0], "X=")) // This has no space in front
			if err != nil {
				return nil, fmt.Errorf("failed to parse prizeLocation[X]: %v", err)
			}
			prizeLocation[0] = int64(p) + int64(10000000000000)

			p, err = strconv.Atoi(strings.TrimPrefix(parts[1], " Y=")) // This has a space in front
			if err != nil {
				return nil, fmt.Errorf("failed to parse prizeLocation[Y]: %v", err)
			}
			prizeLocation[1] = int64(p) + int64(10000000000000)
			machines = append(machines, NewMachine(buttonA, buttonB, prizeLocation))
		}
	}
	fmt.Printf("Parsed %d machines\n", len(machines))
	return machines, nil
}

// Function to compute gcd and coefficients
func extendedGCD(a, b int64) (gcd, x, y int64) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, x1, y1 := extendedGCD(b, a%b)
	x = y1
	y = x1 - (a/b)*y1
	return
}

func solveDiophantine(ax, ay, bx, by, targetX, targetY int64) (nA, nB int64, winnable bool) {
	// Solve for nA using the first equation
	// nA * ax + nB * bx = targetX
	// nA * ay + nB * by = targetY

	// Let's express nB from the first equation:
	// nB = (targetX - nA * ax) / bx

	// Substitute into the second equation:
	// nA * ay + ((targetX - nA * ax) / bx) * by = targetY
	// Multiply through by bx to eliminate denominator:
	// nA * ay * bx + (targetX - nA * ax) * by = targetY * bx
	// nA * (ay * bx - ax * by) = targetY * bx - targetX * by

	denominator := ay*bx - ax*by
	numerator := targetY*bx - targetX*by

	if denominator == 0 {
		// Parallel lines or no solution
		return 0, 0, false
	}

	// Check if denominator divides numerator
	if numerator%denominator != 0 {
		return 0, 0, false
	}

	nA = numerator / denominator
	if nA < 0 {
		return 0, 0, false
	}

	// Now find nB from the first equation
	if (targetX-nA*ax) < 0 || (targetX-nA*ax)%bx != 0 {
		return 0, 0, false
	}
	nB = (targetX - nA*ax) / bx
	if nB < 0 {
		return 0, 0, false
	}

	return nA, nB, true
}

func solve1(machines []*Machine, parallelism int) ([]string, error) {
	output := []string{}
	results := make(chan int64, len(machines))
	semaphore := make(chan struct{}, parallelism)

	for i, m := range machines {
		go func(idx int64, machine *Machine) {
			var tokens = int64(0)
			semaphore <- struct{}{} // Limit the number of goroutines

			// Prize location coordinates
			targetX, targetY := machine.PrizeLocation[0], machine.PrizeLocation[1]

			// Movement vectors for Button A and Button B
			ax, ay := machine.ButtonA[0], machine.ButtonA[1]
			bx, by := machine.ButtonB[0], machine.ButtonB[1]

			nA, nB, winnable := solveDiophantine(ax, ay, bx, by, targetX, targetY)
			if winnable {
				tokens = nA*3 + nB*1
			}

			results <- tokens
			if winnable {
				fmt.Printf("Machine %d: Winnable with nA = %d, nB = %d, Tokens = %d\n", idx, nA, nB, tokens)
			} else {
				fmt.Printf("Machine %d: Not Winnable\n", idx)
			}
			<-semaphore // Release the semaphore
		}(int64(i), m)
	}

	// Collect results
	var totalTokens = int64(0)
	for i := 0; i < len(machines); i++ {
		t := <-results
		fmt.Printf("Result %d: Tokens = %d\n", i+1, t)
		totalTokens += t
	}
	output = append(output, fmt.Sprintf("Tokens: %d", totalTokens))
	return output, nil
}
