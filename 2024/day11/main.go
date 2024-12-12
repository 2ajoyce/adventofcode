package main

import (
	"day11/internal"
	"day11/internal/io"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

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

	lines, err := io.ReadInput(INPUT_FILE)
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

	err = io.WriteOutput(OUTPUT_FILE, results)
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
	localNewStones := make([]internal.Stone, 0, len(stones)*2)
	for _, stone := range stones {
		if stone == 0 {
			localNewStones = append(localNewStones, 1)
		} else if internal.IsEven(stone) {
			left, right, err := internal.Split(stone)
			if err != nil {
				return nil, err
			}
			localNewStones = append(localNewStones, left, right)
		} else {
			localNewStones = append(localNewStones, stone*2024)
		}
	}
	return localNewStones, nil
}

func solve1(blink int, stones []internal.Stone, parallelism int) ([]string, error) {
    DEBUG := os.Getenv("DEBUG") == "true"
    const maxChunkSize = 100_000

    for blinkIteration := 0; blinkIteration < blink; blinkIteration++ {
        // Print statement indicating progression of blinks
        fmt.Printf("Processing Blink %d out of %d\n", blinkIteration+1, blink)

        // Create a new progress bar for this blink iteration
        bar := progressbar.Default(int64(len(stones)))

        if DEBUG {
            fmt.Printf("Blink Iteration %d\n", blinkIteration+1)
        }

        if parallelism <= 1 {
            // Serial processing
            newStones, err := processChunk(stones)
            if err != nil {
                return nil, err
            }
            bar.Add(len(stones))
            stones = newStones
        } else {
            // Parallel processing
            chunkSize := (len(stones) + parallelism - 1) / parallelism
            if chunkSize > maxChunkSize {
                chunkSize = maxChunkSize
            }

            // Determine the total number of chunks
            numChunks := (len(stones) + chunkSize - 1) / chunkSize

            // Channels now have enough capacity to hold all results/errors
            newStonesCh := make(chan []internal.Stone, numChunks)
            errCh := make(chan error, numChunks)

            var wg sync.WaitGroup

            // Split stones into chunks and process them in parallel
            for start := 0; start < len(stones); start += chunkSize {
                end := start + chunkSize
                if end > len(stones) {
                    end = len(stones)
                }
                chunk := stones[start:end]

                wg.Add(1)
                go func(ch []internal.Stone) {
                    defer wg.Done()
                    res, err := processChunk(ch)
                    if err != nil {
                        // Try to send the error, if no room, just return
                        select {
                        case errCh <- err:
                        default: // In case buffer is full (unlikely)
                        }
                        return
                    }

                    // Attempt to send results
                    // This will not block because the channel is large enough
                    // to hold all chunks' results.
                    newStonesCh <- res
                    // Update progress bar for this chunk
                    bar.Add(len(ch))
                }(chunk)
            }

            // Wait for all goroutines to finish
            wg.Wait()
            close(newStonesCh)
            close(errCh)

            // Check if there were any errors
            select {
            case err := <-errCh:
                if err != nil {
                    return nil, err
                }
            default:
                // No error
            }

            // Combine results from the channel
            newStones := make([]internal.Stone, 0, len(stones)*2)
            for res := range newStonesCh {
                newStones = append(newStones, res...)
            }

            stones = newStones
        }

        if DEBUG {
            internal.PrintStones(stones)
        }
    }

    results := []string{fmt.Sprintf("Stones: %d", len(stones))}
    return results, nil
}
