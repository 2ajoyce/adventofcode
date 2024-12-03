package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	leftList, rightList, err := ParseInput(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}

	results, err := Solve1(leftList, rightList)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}
	results2, err := Solve2(leftList, rightList)
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

func ParseInput(lines []string) ([]int, []int, error) {
	DEBUG := os.Getenv("DEBUG")

	leftList := []int{}
	rightList := []int{}

	// Each line in the input file is a string of two space separated integers
	// Iterate over the lines, splitting them on spaces and converting to integers.
	for _, line := range lines {
		if DEBUG == "true" {
			fmt.Printf("Processing line: %s\n", line)
		}
		parts := strings.Split(line, "   ")
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid line format: %s", line)
		}
		num1, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, nil, fmt.Errorf("error converting number in line %s to integer: %v", line, err)
		}
		num2, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, nil, fmt.Errorf("error converting number in line %s to integer: %v", line, err)
		}
		leftList = append(leftList, num1)
		rightList = append(rightList, num2)
		if DEBUG == "true" {
			fmt.Printf("Added to leftList: %d, Added to rightList: %d\n", num1, num2)
		}
	}
	return leftList, rightList, nil
}

func Solve1(leftList []int, rightList []int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG")

	results := []string{}

	// Sort each array by value from smallest to largest
	sort.Ints(leftList)
	sort.Ints(rightList)
	if DEBUG == "true" {
		fmt.Printf("Sorted leftList: %v\n", leftList)
		fmt.Printf("Sorted rightList: %v\n", rightList)
	}
	// Calculate the absolute difference between each pair of integers from leftList and rightList
	// Keep a running total of these differences
	rollingSum := 0
	for i := range leftList {
		diff := abs(leftList[i] - rightList[i])
		rollingSum += diff
		if DEBUG == "true" {
			fmt.Printf("Difference at index %d: %d, Rolling Sum: %d\n", i, diff, rollingSum)
		}
	}

	// Return the final result as a string
	results = append(results, fmt.Sprintf("Total Difference: %s", strconv.Itoa(rollingSum)))
	if DEBUG == "true" {
		fmt.Printf("Final Result: %s\n", results[0])
	}
	return results, nil
}

func Solve2(leftList []int, rightList []int) ([]string, error) {
	//os.getEnv("DEBUG")
	results := []string{}

	rollingSum := 0
	// Iterate over each integer in the left list
	for i := range leftList {
		// Check how many times this integer appears in the right list
		count := 0
		for _, num := range rightList {
			if num == leftList[i] {
				count++
			}
		}
		rollingSum += leftList[i] * count
	}
	results = append(results, fmt.Sprintf("Similarity Score: %s", strconv.Itoa(rollingSum)))
	return results, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
