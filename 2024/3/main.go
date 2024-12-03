package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	//os.Setenv("DEBUG", "true")
	INPUT_FILE := os.Getenv("INPUT_FILE")
	OUTPUT_FILE := os.Getenv("OUTPUT_FILE")

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
	equations, err := ParseInput(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	results, err := Solve1(equations)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}
	results2, err := Solve2(equations)
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

func ParseInput(lines []string) ([]string, error) {
	DEBUG := os.Getenv("DEBUG")

	equations := []string{}

	// Compile the regex pattern to find multiplication expressions of the form mul(xxx,yyy)
	re := regexp.MustCompile(`(?<do>do\(\))|(?<dont>don't\(\))|(?<mul>mul\(\d{1,3},\d{1,3}\))`)

	for _, line := range lines {
		if DEBUG == "true" {
			fmt.Printf("Processing line: %s\n", line)
		}

		// Find all matches of the regex in the current line
		matches := re.FindAllString(line, -1)

		if DEBUG == "true" {
			fmt.Printf("Matches found: %v\n", matches)
		}

		// Append each match to the equations slice
		equations = append(equations, matches...)
	}
	return equations, nil
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

func Solve1(equations []string) ([]string, error) {
	DEBUG := os.Getenv("DEBUG")
	results := []string{}

	rollingSum := 0
	for _, equation := range equations {
		// if equation does not begin with mul, skip it
		if !strings.HasPrefix(equation, "mul(") {
			if DEBUG == "true" {
				fmt.Printf("Skipping equation: %s\n", equation)
			}
			continue
		}
		if DEBUG == "true" {
			fmt.Printf("Processing equation: %s\n", equation)
		}
		// Each equation is of the form "(mul\(\d{1,3},\d{1,3}\)+)"
		// Extract the numbers from the equation and multiply them together
		numbers := strings.Split(equation[4:len(equation)-1], ",")
		x, _ := strconv.Atoi(numbers[0])
		y, _ := strconv.Atoi(numbers[1])

		rollingSum += x * y
		if DEBUG == "true" {
			fmt.Printf("Result of %s: %d\n", equation, x*y)
		}
	}
	results = append(results, fmt.Sprintf("Result 1: %s", strconv.Itoa(rollingSum)))
	return results, nil
}

func Solve2(equations []string) ([]string, error) {
	results := []string{}

	rollingSum := 0
	enabled := true
	for _, equation := range equations {
		if strings.HasPrefix(equation, "do(") {
			enabled = true
		} else if strings.HasPrefix(equation, "don't(") {
			enabled = false
		}
		if strings.HasPrefix(equation, "mul(") && enabled {
			// Each equation is of the form "(mul\(\d{1,3},\d{1,3}\)+)"
			// Extract the numbers from the equation and multiply them together
			numbers := strings.Split(equation[4:len(equation)-1], ",")
			x, _ := strconv.Atoi(numbers[0])
			y, _ := strconv.Atoi(numbers[1])

			rollingSum += x * y
		}
	}
	results = append(results, fmt.Sprintf("Result 2: %s", strconv.Itoa(rollingSum)))
	return results, nil
}
