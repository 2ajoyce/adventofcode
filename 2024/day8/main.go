package main

import (
	"bufio"
	"day8/internal"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	// input, err = parseLines(lines)
	// if err != nil {
	// 	fmt.Println("Error parsing input:", err)
	// 	return
	// }
	// results2, err := solve2(input, PARALLELISM)
	// results = append(results, results2...)
	// if err != nil {
	// 	fmt.Println("Error solving 2:", err)
	// 	return
	//}
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

func parseLines(lines []string) (*internal.AntennaSimulation, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")
	sim, err := internal.NewAntennaSimulation(len(lines[0]), len(lines))
	if err != nil {
		return nil, err
	}
	for y, line := range lines {
		c := strings.Split(line, "")
		for x, char := range c {
			switch char {
			case ".": // These indicate empty cells
				continue
			default: // A cell with any other character is considered an entity
				if DEBUG {
					fmt.Printf("Creating antenna at (%d,%d)\n", x, y)
				}
				a, err := internal.NewAntenna(char)
				if err != nil {
					return nil, err
				}

				insertedAnt, err := sim.AddAntenna(a, x, y)
				if err != nil {
					return nil, err
				}
				if DEBUG {
					fmt.Printf("Inserted antenna at (%d,%d) with signal: %s\n", x, y, insertedAnt)
				}
			}
		}
	}

	return sim, nil
}

func solve1(sim *internal.AntennaSimulation, parallelism int) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning Solve 1")
	fmt.Println(sim)
	// Count the number of antinodes in the simulation
	entities := sim.GetEntities()
	rollingTotal := 0
	for _, entity := range entities {
		switch entity.(type) {
		case *internal.Antinode:
			rollingTotal++
		}
	}
	fmt.Printf("Found %d antinodes\n", rollingTotal)
	results := []string{}
	results = append(results, fmt.Sprintf("Unique Antinode Locations: %d", rollingTotal))
	return results, nil
}
