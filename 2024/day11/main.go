package main

import (
	"day11/internal"
	"day11/internal/io"
	"errors"
	"fmt"
	"math/big"
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
		vInt, success := new(big.Int).SetString(v, 10)
		if !success {
			return -1, nil, errors.New("Invalid stone value: " + v)
		}
		stones[i] = *internal.NewStone(vInt)
	}

	return blink, stones, nil
}

func solve1(blink int, stones []internal.Stone, parallelism int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	bar := progressbar.Default(int64(blink))
	fmt.Println("Beginning Solve 1")
	defer fmt.Println("Ending Solve 1")
	zero := big.NewInt(0) // Predefined for use in Rule 1
	one := big.NewInt(1)
	twentyTwentyFour := big.NewInt(2024) // Used in every iteration of Rule 3

	// For however many blinks were specified
	for blinkIteration := 0; blinkIteration < blink; blinkIteration++ {
		bar.Add(1)
		if DEBUG {
			fmt.Printf("Blink Iteration %d\n", blinkIteration+1)
		}
		// Check each rule against each stone
		for stoneIndex := 0; stoneIndex < len(stones); stoneIndex++ {
			if DEBUG {
				internal.PrintStones(stones)
			}
			// Rule 1: Turn 0's into 1's
			if stones[stoneIndex].Value.Cmp(zero) == 0 {
				stones[stoneIndex].ChangeValue(*one) // Todo: This may cause a bug, flagging it for inspection later
				continue
			}
			// Rule 2: If the stone number has an even number of digits
			if stones[stoneIndex].IsEven() {
				left, right, err := stones[stoneIndex].Split()
				if err != nil {
					return nil, err
				}
				stones = append(stones[:stoneIndex], append([]internal.Stone{*left, *right}, stones[stoneIndex+1:]...)...)
				stoneIndex++ // Skip the next stone as we've added two new ones in its place
				continue
			}
			// Rule 3: Multiply the value by 2024
			newValue := new(big.Int).Mul(&stones[stoneIndex].Value, twentyTwentyFour)
			stones[stoneIndex].ChangeValue(*newValue)
		}
		if DEBUG {
			internal.PrintStones(stones)
		}
	}

	results := []string{}
	results = append(results, fmt.Sprintf("Stones: %d", len(stones)))
	return results, nil
}
