package main

import (
	"bufio"
	"day7/internal"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
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

type LoopError struct {
	Message string
}

func (e LoopError) Error() string {
	return e.Message
}

func main() {
	//os.Setenv("DEBUG", "true")
	INPUT_FILE := os.Getenv("INPUT_FILE")
	OUTPUT_FILE := os.Getenv("OUTPUT_FILE")
	PARALLELISM, err := strconv.Atoi(os.Getenv("PARALLELISM"))
	if PARALLELISM < 1 || err != nil {
		PARALLELISM = 1
	}
	fmt.Printf("PARALLELISM: %d\n", PARALLELISM)

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
	equations, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve1(equations, PARALLELISM)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}

	equations, err = parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results2, err := solve2(equations, PARALLELISM)
	results = append(results, results2...)
	if err != nil {
		fmt.Println("Error solving 2:", err)
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

func parseLines(lines []string) ([]internal.Equation, error) {
	//DEBUG := os.Getenv("DEBUG")
	equations := make([]internal.Equation, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ":")
		smallTotal, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("error converting total to integer: %v", err)
		}
		total := *big.NewInt(int64(smallTotal))
		combinedNumbers := strings.TrimSpace(parts[1])
		stringNumbers := strings.Split(combinedNumbers, " ")
		numbers := make([]int, len(stringNumbers))
		for j, strNumber := range stringNumbers {
			n, err := strconv.Atoi(strNumber)
			if err != nil {
				return nil, fmt.Errorf("error converting number to integer: %v", err)
			}
			numbers[j] = n
		}
		equations[i] = internal.NewEquation(total, numbers)
	}
	longestEquation := 0
	for _, e := range equations {
		eLen := len(e.Numbers())
		if eLen > longestEquation {
			longestEquation = eLen
		}
	}
	fmt.Printf("Longest equation: %d\n", longestEquation)

	return equations, nil
}

func solve1(equations []internal.Equation, parallelism int) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"

	results := []string{}

	// Open a channel to pass equations to the worker goroutines, a channel to signal completion, and a channel for errors
	// The goroutines will read equations from the channel, call their Solve() method, and end the routine when they receive a signal on doneChan.

	equationChan := make(chan internal.Equation)
	doneChan := make(chan bool)
	errorChan := make(chan error)
	bar := progressbar.Default(int64(len(equations)))

	for i := 0; i < parallelism; i++ {
		go func() {
			for equation := range equationChan {
				_, err := equation.Solve()
				if err != nil {
					errorChan <- err
					return
				}
				bar.Add(1)
			}
			doneChan <- true
		}()
	}
	for _, equation := range equations {
		equationChan <- equation
	}
	close(equationChan)
	for i := 0; i < parallelism; i++ {
		<-doneChan
	}
	close(errorChan)
	for err := range errorChan {
		return nil, err
	}

	solvedEquations := 0
	validEquations := 0
	rollingTotal := big.NewInt(0)
	for _, equation := range equations {
		if equation.IsSolved() {
			solvedEquations++
		}
		if equation.IsValid() {
			validEquations++
			total := equation.Total()
			rollingTotal.Add(rollingTotal, &total)
		}
	}
	fmt.Printf("Solved Equations: %d\n", solvedEquations)
	fmt.Printf("Valid Equations: %d\n", validEquations)
	fmt.Printf("Rolling Total: %d\n", rollingTotal)

	results = append(results, fmt.Sprintf("Calibration Result: %s", rollingTotal))

	return results, nil
}

func solve2(equations []internal.Equation, parallelism int) ([]string, error) {
	//DEBUG := os.Getenv("DEBUG") == "true"

	results := []string{}

	return results, nil
}
