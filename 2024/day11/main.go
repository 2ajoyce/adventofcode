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

	input := lines[0] // There is only one line of input
	// Input will be in the form "1:23 45"
	// 1 indicates the number of times to blink
	// There are two stones with values 23 and 45
	s := strings.Split(input, ":")
	b := s[0]
	blink, err := strconv.Atoi(b)
	if err != nil {
		return -1, nil, errors.New("Invalid blink value: " + b)
	}

	// Build Stones array
	s = strings.Split(s[1], " ") // Corrected this line to split the string after the colon
	stones := make([]internal.Stone, len(s))
	for i, v := range s {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return -1, nil, fmt.Errorf("invalid stone value: %s", v)
		}
		stones[i] = vInt
	}

	return blink, stones, nil
}

// processChunk processes a slice of stones and returns the transformed slice.
func processChunk(stones []internal.Stone) ([]internal.Stone, error) {
	// Initialize a slice to hold all transformed stones.
	// Assuming each stone can produce up to 2 new stones, preallocate for efficiency.
	transformed := make([]internal.Stone, 0, len(stones)*2)

	for _, stone := range stones {
		if stone == 0 {
			transformed = append(transformed, 1)
		} else if internal.IsEven(stone) {
			left, right, err := internal.Split(stone)
			if err != nil {
				return nil, err
			}
			transformed = append(transformed, left, right)
		} else {
			transformed = append(transformed, stone*2024)
		}
	}

	return transformed, nil
}
func solve1(blink int, stones []internal.Stone, parallelism int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	const maxChunkSize = 5_000_000

	// Define how many stones to display from each queue
	const displayStones = 5
	const maxPruneThreshold = 500_000_000

	// Initialize queues for each depth (0 to blink)
	queues := make([][]internal.Stone, blink+1)
	queues[0] = stones

	// Initialize indices to track the next stone to process in each queue
	indices := make([]int, blink+1) // Corrected length

	// Initialize RunningTotal as a big.Int to handle large counts
	runningTotal := big.NewInt(0)

	// Initialize uilive writer
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	// Define how often to log status (e.g., every 1,000,000 stones)
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

	// Helper function to ensure queues have enough capacity
	ensureCapacity := func(queues [][]internal.Stone, index int) [][]internal.Stone {
		if index < len(queues) {
			return queues
		}
		newLength := len(queues) * 2
		if newLength <= index {
			newLength = index + 1
		}
		return append(queues, make([][]internal.Stone, newLength-len(queues))...)
	}

	// Processing loop
	for {
		// Find the deepest non-empty queue
		deepest := -1
		for d := blink; d >= 0; d-- {
			if indices[d] < len(queues[d]) {
				deepest = d
				break
			}
		}

		// If no queues have stones left to process, break the loop
		if deepest == -1 {
			break
		}

		if DEBUG {
			fmt.Fprintf(writer, "Comparing deepest=%d with blink=%d\n", deepest, blink)
		}

		// Determine how many stones to process in this chunk
		remainingStones := len(queues[deepest]) - indices[deepest]
		chunkSize := maxChunkSize
		if remainingStones < maxChunkSize {
			chunkSize = remainingStones
		}

		// Get the chunk of stones to process
		stonesToProcess := queues[deepest][indices[deepest] : indices[deepest]+chunkSize]
		indices[deepest] += chunkSize // Mark these stones as processed

		if deepest == blink {
			// **COUNTING WITHOUT PROCESSING**
			// Add the number of stones in this chunk directly to runningTotal
			runningTotal.Add(runningTotal, big.NewInt(int64(chunkSize)))
			processedStones += int64(chunkSize)

			if DEBUG {
				fmt.Fprintf(writer, "Counting %d stones from queue[%d] without processing\n", chunkSize, deepest)
			}
		} else {
			if DEBUG {
				fmt.Fprintf(writer, "Processing %d stones from queue[%d]\n", chunkSize, deepest)
			}

			// Transform the stones using the updated processChunk
			transformedStones, err := processChunk(stonesToProcess)
			if err != nil {
				fmt.Fprintf(writer, "\nError processing stones at depth %d: %v\n", deepest, err)
				writer.Stop()
				return nil, err
			}

			// Determine the next depth
			nextDepth := deepest + 1

			// Ensure queues have enough capacity
			queues = ensureCapacity(queues, nextDepth)

			// Add the transformed stones to the next depth queue
			queues[nextDepth] = append(queues[nextDepth], transformedStones...)
			processedStones += int64(chunkSize)
		}

		// Periodic Status Updates
		if processedStones >= nextStatusUpdate {
			// Write to uilive writer
			status := fmt.Sprintf(
				"Status Update:\n  Total Processed Stones: %d\n  RunningTotal: %s\n  Stones in Queues:\n",
				processedStones, runningTotal.String(),
			)
			for d := 0; d <= blink; d++ {
				remaining := len(queues[d]) - indices[d]
				if remaining <= 0 {
					status += fmt.Sprintf("    Depth %d: 0 stones remaining\n", d)
					continue
				}

				// Determine the number of stones to display (min(displayStones, remaining))
				displayCount := displayStones
				if remaining < displayStones {
					displayCount = remaining
				}

				// Slice to get the first N stones
				firstNStones := queues[d][indices[d] : indices[d]+displayCount]

				// Format stones for display
				stoneStrs := make([]string, len(firstNStones))
				for i, stone := range firstNStones {
					stoneStrs[i] = strconv.Itoa(int(stone))
				}

				status += fmt.Sprintf(
					"    Depth %d: %d stones remaining | First %d stones: [%s]\n",
					d, remaining, displayCount, strings.Join(stoneStrs, ", "),
				)
			}
			status += fmt.Sprintf("  Time Elapsed: %s\n\n", time.Since(startTime))
			fmt.Fprint(writer, status)

			// Update the nextStatusUpdate
			nextStatusUpdate += statusInterval
		}

		// Pruning Logic
		if indices[deepest] >= maxPruneThreshold {
			// Prune the processed stones from the queue
			queues[deepest] = queues[deepest][indices[deepest]:]
			// Reset the index for this queue
			indices[deepest] = 0
		}
	}

	// Final Status Update
	finalStatus := fmt.Sprintf(
		"\nFinal Status:\n  Total Processed Stones: %d\n  RunningTotal: %s\n  Time Elapsed: %s\n\n",
		processedStones, runningTotal.String(), time.Since(startTime),
	)
	fmt.Fprint(writer, finalStatus)

	// After processing all stones, prepare the result
	results := []string{fmt.Sprintf("Stones: %s", runningTotal.String())}
	return results, nil
}
