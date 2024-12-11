package main

import (
	"bufio"
	"day10/internal"
	"day10/internal/simulation"
	"fmt"
	"os"
	"strconv"
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

func main() {
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

	lines, err := ReadInput(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error reading from %s: %v", INPUT_FILE, err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// Start Solution Logic  ///////////////////////////////////////////
	////////////////////////////////////////////////////////////////////

	// Create an array of all coordinates containing the letter X
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
	// End Solution Logic  /////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////

	err = WriteOutput(OUTPUT_FILE, results)
	if err != nil {
		fmt.Printf("Error writing to %s: %v", OUTPUT_FILE, err)
		return
	}

	fmt.Printf("Successfully processed %s and created %s", INPUT_FILE, OUTPUT_FILE)
}

func parseLines(lines []string) (simulation.Simulation, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")
	width := len(lines[0])
	height := len(lines)
	sim := simulation.NewSimulation(width, height)

	for y, line := range lines {
		for x, char := range line {
			zHeight := int(char - '0')
			e, err := internal.NewTopoEntity(zHeight)
			if err != nil {
				return nil, err
			}
			sim.AddEntity(e, x, y)
		}
	}
	return sim, nil
}

func solve1(sim simulation.Simulation, parallelism int) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning Solve 1")
	defer fmt.Println("Ending Solve 1")
	s, err := internal.StringifySimulation(sim)
	if err != nil {
		return nil, err
	}
	fmt.Println(s)
	trailheads, err := internal.FindTrailheads(sim)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Trailhead Count: %d\n", len(trailheads))
	mapScore := 0
	for t := range trailheads {
		fmt.Printf("Scoring trailhead at %d,%d\n", trailheads[t].X, trailheads[t].Y)
		score, err := internal.ScoreTrailhead(sim, trailheads[t])
		if err != nil {
			return nil, err
		}
		fmt.Printf("Score for trailhead at %d,%d is %d\n", trailheads[t].X, trailheads[t].Y, score)
		mapScore += score
	}
	fmt.Printf("Map Score: %d\n", mapScore)
	results := []string{}
	results = append(results, fmt.Sprintf("Score: %d", mapScore))
	return results, nil
}
