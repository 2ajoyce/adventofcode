package main

import (
	"bufio"
	"day9/internal"
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

func parseLines(lines []string) (*internal.DiskMap, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")
	diskmap := internal.NewDiskMap()
	input := lines[0] // Input is a single line

	maxBlockId := 0
	currentFileId := 0
	for i, char := range input {
		if DEBUG {
			fmt.Printf("Parsing char %c at index %d\n", char, i)
			fmt.Printf("MaxBlockId: %d\n", maxBlockId)
		}
		if i%2 == 0 { // if i is even it is a file
			blockCount, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, fmt.Errorf("error parsing block count: %v", err)
			}
			if DEBUG {
				fmt.Printf("Adding file %d with %d blocks\n", currentFileId, blockCount)
			}
			blocks := make([]int, blockCount)
			for i := range blockCount {
				blocks[i] = maxBlockId
				maxBlockId++
			}
			diskmap.AddFile(currentFileId, blocks)
			currentFileId++
		} else { // it is empty space
			blockCount, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, fmt.Errorf("error parsing block count: %v", err)
			}
			if DEBUG {
				fmt.Printf("Skipping %d empty blocks\n", blockCount)
			}
			maxBlockId += blockCount
		}
		if DEBUG {
			fmt.Println()
		}
	}

	return diskmap, nil
}

func solve1(diskmap *internal.DiskMap, parallelism int) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning Solve 1")
	defer fmt.Println("Ending Solve 1")
	fmt.Println(diskmap)
	diskmap.Compact()
	fmt.Println(diskmap)
	results := []string{}
	results = append(results, fmt.Sprintf("Checksum: %s", diskmap.GetChecksum()))
	return results, nil
}
