package main

import (
	"day14/internal/aocUtils"
	"day14/internal/simulation"
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
	for y := 0; y < sim.GetMap().GetHeight(); y++ {
		for x := 0; x < sim.GetMap().GetWidth(); x++ {
			cell, err := sim.GetMap().GetCell(x, y)
			if err != nil {
				output += "?"
				continue
			}
			ids := cell.GetEntityIds()
			if len(ids) == 0 {
				output += " "
			} else {
				output += "#"
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

	for i := 1; i < len(lines); i++ {
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
		vx, err := strconv.Atoi(pVel[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing x velocity of entity in line %d: %s", i+1, p[0])
		}
		vy, err := strconv.Atoi(pVel[1])
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

	return sim, nil
}

func tick(sim simulation.Simulation, tickNumber int) (simulation.Simulation, error) {
	entities := sim.GetEntities()
	for _, e := range entities {
		// Move the entity based on its velocity
		x, y := e.GetPosition()
		velX, velY := e.GetVelocity()
		newX := x + velX
		newY := y + velY
		//fmt.Printf("Moving entity %s from (%d, %d) to (%d, %d)\n", e.GetId(), x, y, newX, newY)
		success, err := sim.MoveEntity(e.GetId(), newX, newY)
		if err != nil || !success {
			return sim, fmt.Errorf("error moving entity in tick %d: %s, %v", tickNumber, e.GetId(), err)
		}
	}
	return sim, nil
}

func calculateSafetyFactor(sim simulation.Simulation) int {
	// To determine the safest area, count the number of robots in each quadrant
	// Robots that are exactly in the middle (horizontally or vertically) don't count as being in any quadrant,

	//1. Divide the simulation space into four quadrants.
	height := sim.GetMap().GetHeight()
	width := sim.GetMap().GetWidth()
	quadrantCount := make(map[string]int)
	for _, e := range sim.GetEntities() {
		x, y := e.GetPosition()
		if x < width/2 && y < height/2 {
			quadrantCount["NW"]++
		} else if x < width/2 && y > height/2 {
			quadrantCount["NE"]++
		} else if x > width/2 && y < height/2 {
			quadrantCount["SW"]++
		} else if x > width/2 && y > height/2 {
			quadrantCount["SE"]++
		}
	}

	//2. Multiply the number of robots in each quadrant together.
	safetyFactor := quadrantCount["NW"] * quadrantCount["NE"] * quadrantCount["SW"] * quadrantCount["SE"]
	return safetyFactor
}

type OutputRecord struct {
	tick         int
	safetyFactor int
	visual       string
}

func solve1(sim simulation.Simulation, parallelism int) ([]string, error) {
	var output = []string{}
	const maxTick = 10000
	bar := progressbar.Default(maxTick)
	var safestTick = 0
	safetyMap := map[int]OutputRecord{}
	var safetyFactor = 0

	for tickNumber := 1; tickNumber <= maxTick; tickNumber++ {
		// fmt.Printf("Tick %d\n%s", tickNumber, PrintSim(sim))
		_, err := tick(sim, tickNumber)
		if err != nil {
			return nil, fmt.Errorf("error during tick %d: %s", tickNumber, err)
		}
		bar.Add(1)

		sf := calculateSafetyFactor(sim)
		if safetyFactor == 0 || sf <= safetyFactor {
			safetyFactor = sf
			visual := PrintSim(sim)
			safetyMap[tickNumber] = OutputRecord{tick: tickNumber, safetyFactor: sf, visual: visual}
		}
	}
	safetyFactor = calculateSafetyFactor(sim) // At the end, set the safety factor to the last frame

	// Print the output records
	safestTickText := []string{}
	for _, record := range safetyMap {
		safestTickText = append(output, fmt.Sprintf("Tick %d Safety Factor: %d\nVisual:\n%s\n", record.tick, record.safetyFactor, record.visual))
	}

	fmt.Printf("The safest tick is %d with a safety factor of %d\n", safestTick, safetyFactor)
	// Save a visual of the safest tick to a file
	aocUtils.WriteOutput("safest_tick.txt", safestTickText)
	output = append(output, fmt.Sprintf("Safety Factor: %d", safetyFactor))
	return output, nil
}
