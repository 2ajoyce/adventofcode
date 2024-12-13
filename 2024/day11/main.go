package main

import (
	"day11/internal"
	"day11/internal/aocIo"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gosuri/uilive"
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

	lines, err := aocIo.ReadInput(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error reading from %s: %v", INPUT_FILE, err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// SOLUTION LOGIC
	////////////////////////////////////////////////////////////////////

	blink, stones, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve1(blink, stones, PARALLELISM)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// WRITE OUTPUT FILE
	////////////////////////////////////////////////////////////////////

	err = aocIo.WriteOutput(OUTPUT_FILE, results)
	if err != nil {
		fmt.Printf("Error writing to %s: %v", OUTPUT_FILE, err)
		return
	}

	fmt.Printf("Successfully processed %s and created %s", INPUT_FILE, OUTPUT_FILE)
}

func parseLines(lines []string) (int, []internal.Stone, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return -1, nil, errors.New("input is empty")
	}

	input := lines[0] // There is only one line of input
	// Input will be in the form "1:23 45"
	// 1 indicates the number of times to blink
	// There are two stones with values 23 and 45
	s := strings.Split(input, ":")
	if len(s) != 2 {
		return -1, nil, errors.New("input format invalid, expected 'blink:stone1 stone2 ...'")
	}
	b := s[0]
	blink, err := strconv.Atoi(b)
	if err != nil {
		return -1, nil, errors.New("Invalid blink value: " + b)
	}

	// Build Stones array
	s = strings.Split(s[1], " ") // Split the string after the colon
	stones := make([]internal.Stone, 0, len(s))
	for _, v := range s {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return -1, nil, fmt.Errorf("invalid stone value: %s", v)
		}
		stones = append(stones, internal.Stone(vInt))
	}

	return blink, stones, nil
}

// processChunk processes a single stone and returns the transformed stones.
func processStone(stone internal.Stone) ([]internal.Stone, error) {
	transformed := []internal.Stone{}

	if stone == 0 {
		transformed = append(transformed, 1)
	} else if internal.IsEven(int(stone)) {
		left, right, err := internal.Split(int(stone))
		if err != nil {
			return nil, err
		}
		transformed = append(transformed, internal.Stone(left), internal.Stone(right))
	} else {
		// Assuming that for odd stones, stone*2024 is the new stone value
		transformed = append(transformed, internal.Stone(int(stone)*2024))
	}

	return transformed, nil
}

func solve1(blink int, stones []internal.Stone, parallelism int) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"
	const maxChunkSize = 5_000_000

	// Initialize RunningTotal as a big.Int to handle large counts
	runningTotal := big.NewInt(0)

	// Initialize uilive writer
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	// Define how often to log status (e.g., every 5,000,000 stones)
	const statusInterval = 5_000_000
	nextStatusUpdate := int64(statusInterval)
	processedStones := int64(0)
	startTime := time.Now()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintf(writer, "\nInterrupt received, shutting down...\n")
		writer.Stop()
		os.Exit(1)
	}()

	// Initialize the memoization cache: map[stone][depth] = count
	cache := make(map[int]map[int]int)

	// Helper function to get count from cache or compute it
	var getStoneCount func(stone internal.Stone, depth int) (int, error)
	getStoneCount = func(stone internal.Stone, depth int) (int, error) {
		// Base case: if depth == blink, this stone counts as 1
		if depth == blink {
			return 1, nil
		}

		// Initialize inner map if not present
		if cache[int(stone)] == nil {
			cache[int(stone)] = make(map[int]int)
		}

		// Check if the count is already cached
		if count, exists := cache[int(stone)][depth]; exists {
			return count, nil
		}

		// Process the stone to get transformed stones
		transformedStones, err := processStone(stone)
		if err != nil {
			return 0, err
		}

		// Initialize count for this stone at this depth
		total := 0

		for _, newStone := range transformedStones {
			count, err := getStoneCount(newStone, depth+1)
			if err != nil {
				return 0, err
			}
			total += count
		}

		// Cache the computed count
		cache[int(stone)][depth] = total

		return total, nil
	}

	// Process each initial stone
	for _, stone := range stones {
		count, err := getStoneCount(stone, 0)
		if err != nil {
			fmt.Fprintf(writer, "\nError processing stone %d: %v\n", stone, err)
			writer.Stop()
			return nil, err
		}
		runningTotal.Add(runningTotal, big.NewInt(int64(count)))
		processedStones += 1

		// Periodic Status Updates
		if processedStones >= nextStatusUpdate {
			// Write to uilive writer
			status := fmt.Sprintf(
				"Status Update:\n  Total Processed Stones: %d\n  RunningTotal: %s\n  Time Elapsed: %s\n\n",
				processedStones, runningTotal.String(), time.Since(startTime),
			)
			fmt.Fprint(writer, status)

			// Update the nextStatusUpdate
			nextStatusUpdate += statusInterval
		}
	}

	// Final Status Update
	finalStatus := fmt.Sprintf(
		"\nFinal Status:\n  Total Processed Stones: %d\n  RunningTotal: %s\n  Time Elapsed: %s\n\n",
		processedStones, runningTotal.String(), time.Since(startTime),
	)
	fmt.Fprint(writer, finalStatus)

	// Prepare the result
	results := []string{fmt.Sprintf("Stones: %s", runningTotal.String())}
	return results, nil
}
