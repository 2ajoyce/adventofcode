package main

import (
	"day14/internal/aocUtils"
	"day14/internal/simulation"
	"fmt"
	"os"
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

func PrintSim(sim simulation.Simulation) string {
	var output string
	// Map each entity id to a character for printing purposes
	entityChar := map[uuid.UUID]string{}
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
				// Check if id[0] is not in the map, assign it a letter
				if _, ok := entityChar[ids[0]]; !ok {
					entityChar[ids[0]] = string(rune(len(entityChar) + 'A'))
				}
				output += entityChar[ids[0]]
			}
		}
		output += "\n"
	}
	return output
}

func parseLines(lines []string) (simulation.Simulation, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	// Line 0 is the width and height
	p := strings.Split(lines[0], ":")
	width, err := strconv.Atoi(p[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing width: %s", p[0])
	}
	height, err := strconv.Atoi(p[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing height: %s", p[1])
	}
	sim := simulation.NewSimulation(width, height)

	for i := 1; i < len(lines[1:]); i++ {
		// Each line represents an entity with location and velocity
		p = strings.Split(lines[i], " ") //p[0]=location, p[1]=velocity
		pLoc := strings.Split(strings.TrimSpace(strings.TrimPrefix(p[0], "p=")), ",")
		pVel := strings.Split(strings.TrimSpace(strings.TrimPrefix(p[1], "v=")), ",")

		x, err := strconv.Atoi(pLoc[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing x coordinate of entity in line %d: %s", i+1, p[0])
		}
		y, err := strconv.Atoi(pLoc[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing y coordinate of entity in line %d: %s", i+1, p[0])
		}
		vx, err := strconv.ParseFloat(pVel[0], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing x velocity of entity in line %d: %s", i+1, p[0])
		}
		vy, err := strconv.ParseFloat(pVel[1], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing y velocity of entity in line %d: %s", i+1, p[0])
		}
		e, err := simulation.NewEntity()
		if err != nil {
			return nil, fmt.Errorf("error creating new entity in line %d: %s", i+1, p[0])
		}
		_, err = sim.AddEntity(e, x, y, vx, vy)
		if err != nil {
			return nil, fmt.Errorf("error adding entity to simulation in line %d: %s", i+1, p[0])
		}
	}

	fmt.Printf("Simulation\n%s", PrintSim(sim))
	return sim, nil
}

func solve1(sim simulation.Simulation, parallelism int) ([]string, error) {
	var output = []string{}
	var safetyFactor = 0

	output = append(output, fmt.Sprintf("Safety Factor: %d", safetyFactor))
	return output, nil
}
