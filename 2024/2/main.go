package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	INPUT_FILE := os.Getenv("INPUT_FILE")
	OUTPUT_FILE := os.Getenv("OUTPUT_FILE")
	if INPUT_FILE == "" || OUTPUT_FILE == "" {
		fmt.Println("INPUT_FILE and OUTPUT_FILE environment variables not set")
		fmt.Println("Defaulting to input.txt and output.txt")
		INPUT_FILE = "input.txt"
		OUTPUT_FILE = "output.txt"
	}

	inputFile, err := os.Open(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error opening %s: %v", INPUT_FILE, err)
		return
	}
	defer inputFile.Close()

	outputFile, err := os.Create(OUTPUT_FILE)
	if err != nil {
		fmt.Printf("Error creating %s: %v", OUTPUT_FILE, err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	////////////////////////////////////////////////////////////////////
	// Start Solution Logic  ///////////////////////////////////////////
	////////////////////////////////////////////////////////////////////
	reports, err := ParseInput(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	results, err := Solve1(reports)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}
	results2, err := Solve2(reports)
	results = append(results, results2...)
	if err != nil {
		fmt.Println("Error solving 2:", err)
		return
	}
	////////////////////////////////////////////////////////////////////
	// End Solution Logic  /////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////

	// Write the results to output.txt, one line per result
	for i, res := range results {
		_, err := writer.WriteString(res)
		if err != nil {
			fmt.Printf("error writing value to %s: %v", OUTPUT_FILE, err)
			return
		}
		if i != len(results)-1 {
			_, err = writer.WriteString("\n")
			if err != nil {
				fmt.Printf("error writing newline to %s: %v", OUTPUT_FILE, err)
				return
			}
		}
	}

	// Flush the writer to ensure all data is written to output.txt
	writer.Flush()

	fmt.Printf("Successfully processed %s and created %s", INPUT_FILE, OUTPUT_FILE)
}

func ParseInput(lines []string) ([][]int, error) {
	DEBUG := os.Getenv("DEBUG")

	reports := [][]int{}

	// Each line in the input file is a string of space separated integers representing a report
	// Iterate over the lines, splitting them on spaces and converting to integers
	for _, line := range lines {
		if DEBUG == "true" {
			fmt.Printf("Processing line: %s\n", line)
		}
		parts := strings.Split(line, " ")
		report := []int{}
		for part := range parts {
			num, err := strconv.Atoi(parts[part])
			if err != nil {
				return nil, fmt.Errorf("error converting %s to int: %v", parts[part], err)
			}
			report = append(report, num)
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func Solve1(reports [][]int) ([]string, error) {
	results := []string{}

	rollingSum := 0
	for _, report := range reports {
		validity, err := ValidateReport(report)
		if err != nil {
			return nil, fmt.Errorf("error validating report: %v", err)
		}
		if validity {
			rollingSum += 1
		}

	}
	results = append(results, fmt.Sprintf("Safe Reports: %s", strconv.Itoa(rollingSum)))
	return results, nil
}

func Solve2(reports [][]int) ([]string, error) {
	results := []string{}

	rollingSum := 0
	for _, report := range reports {
		validity, err := ValidateReport(report)
		if err != nil {
			return nil, fmt.Errorf("error validating report: %v", err)
		}
		if !validity {
			// Loop through the report removing one item each loop to see if the report becomes valid
			for i := range report {
				newReport := append([]int(nil), report[:i]...)
				newReport = append(newReport, report[i+1:]...)
				newValidity, err := ValidateReport(newReport)
				if err != nil {
					return nil, fmt.Errorf("error validating new report: %v", err)
				}
				if newValidity {
					validity = true
					break
				}
			}
		}
		if validity {
			rollingSum += 1
		}
	}
	results = append(results, fmt.Sprintf("Actually Safe Reports: %s", strconv.Itoa(rollingSum)))
	return results, nil
}

func ValidateReport(report []int) (bool, error) {
	validity := true

	// Rule 1: The values in a report are either all increasing or all decreasing.
	increasing := false
	decreasing := false
	for i := range report {
		// Exit one iteration early since we are comparing against the next value.
		if i == len(report)-1 {
			break
		}
		if report[i] < report[i+1] {
			increasing = true
		}
		if report[i] > report[i+1] {
			decreasing = true
		}
		if increasing && decreasing {
			validity = false
			break
		}
	}

	// Rule 2:  Any two adjacent values in a report differ by at least one and at most three.
	for i := range report {
		// Exit one iteration early since we are comparing against the next value.
		if i == len(report)-1 {
			break
		}
		diff := abs(report[i] - report[i+1])
		if diff < 1 || diff > 3 {
			validity = false
			break
		}
	}
	return validity, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
