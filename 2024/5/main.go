package main

import (
	"bufio"
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
	rules, updates, err := ParseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := Solve1(rules, updates)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}

	results2, err := Solve2(rules, updates)
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

func ParseLines(lines []string) (map[int][]int, [][]int, error) {
	rules := map[int][]int{}
	updates := [][]int{}

	// Parse the input lines into rules and pages
	// The first section is all the rules in format 12|34
	// An empty line separates the first section from the second section
	secondSectionStart := 0
	for i, line := range lines {
		if line == "" {
			secondSectionStart = i + 1
			break
		}
		rule := strings.Split(line, "|")
		ruleKey, err := strconv.Atoi(strings.TrimSpace(rule[0]))
		if err != nil {
			return nil, nil, fmt.Errorf("invalid rule key: %s", rule[0])
		}
		ruleValue, err := strconv.Atoi(strings.TrimSpace(rule[1]))
		if err != nil {
			return nil, nil, fmt.Errorf("invalid rule value: %s", rule[1])
		}
		rules[ruleKey] = append(rules[ruleKey], ruleValue)
	}
	// The second section is all pages in format 56,78,90
	for i := secondSectionStart; i < len(lines); i++ {
		pages := strings.Split(lines[i], ",")
		update := []int{}
		for _, part := range pages {
			page, err := strconv.Atoi(part)
			if err != nil {
				fmt.Printf("Error converting %s to int: %v\n", part, err)
				continue
			}
			update = append(update, page)

		}
		updates = append(updates, update)
	}

	return rules, updates, nil
}

func Solve1(rules map[int][]int, updates [][]int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG")
	results := []string{}

	rollingSum := 0
	// For every update in the updates
	for _, update := range updates {
		seenPages := []int{}
		valid := true
		// For every page in the update
		for _, page := range update {
			// Add the page to the seenPages
			seenPages = append(seenPages, page)
			// Search the rules for the page number
			applicableRules := []int{}
			for ruleKey, ruleValues := range rules {
				if ruleKey == page {
					applicableRules = append(applicableRules, ruleValues...)
				}
			}
			if DEBUG == "true" {
				fmt.Printf("Page %d seen. Applicable rules: %v\n", page, applicableRules)
			}
			// For every applicable rule
			for _, rule := range applicableRules {
				if contains(seenPages, rule) {
					valid = false
				}
			}
		}
		if valid {
			// Find the middle value from the update
			middleValue := update[len(update)/2]
			rollingSum += middleValue
		}
	}
	results = append(results, fmt.Sprintf("Result 1: %s", strconv.Itoa(rollingSum)))
	return results, nil
}

func Solve2(rules map[int][]int, updates [][]int) ([]string, error) {
	DEBUG := os.Getenv("DEBUG")
	results := []string{}

	rollingSum := 0
	invalidUpdates := [][]int{}
	// Find the invalid updates
	for _, update := range updates {
		seenPages := []int{}
		valid := true
		// For every page in the update
		for _, page := range update {
			// Add the page to the seenPages
			seenPages = append(seenPages, page)
			// Search the rules for the page number
			applicableRules := []int{}
			for ruleKey, ruleValues := range rules {
				if ruleKey == page {
					applicableRules = append(applicableRules, ruleValues...)
				}
			}
			if DEBUG == "true" {
				fmt.Printf("Page %d seen. Applicable rules: %v\n", page, applicableRules)
			}
			// For every applicable rule
			for _, rule := range applicableRules {
				if contains(seenPages, rule) {
					valid = false
				}
			}
		}
		if !valid {
			sortedUpdate := sortUpdate(rules, update)
			invalidUpdates = append(invalidUpdates, sortedUpdate)
		}
	}
	// For every invalid update
	for _, update := range invalidUpdates {
		middleValue := update[len(update)/2]
		rollingSum += middleValue
	}

	results = append(results, fmt.Sprintf("Result 2: %s", strconv.Itoa(rollingSum)))
	return results, nil
}

func contains(seenPages []int, rule int) bool {
	for _, page := range seenPages {
		if page == rule {
			return true
		}
	}
	return false
}

func sortUpdate(rules map[int][]int, update []int) []int {
	DEBUG := os.Getenv("DEBUG")
	if DEBUG == "true" {
		fmt.Printf("Sorting update: %v\n", update)
	}
	// Some pages in the update are out of order
	// For every page in the update
	for i := 0; i < len(update); i++ {
		page := update[i]
		if DEBUG == "true" {
			fmt.Printf("Checking page %d\n", page)
		}
		// Find the applicable rules for the page number
		applicableRules := []int{}
		for ruleKey, ruleValues := range rules {
			if ruleKey == page {
				applicableRules = append(applicableRules, ruleValues...)
			}
		}
		if DEBUG == "true" {
			fmt.Printf("Applicable rules: %v\n", applicableRules)
		}
		// Sort the rule values by integer value
		//sort.Ints(applicableRules)
		// For every applicable rule check if a page matching that rule precedes the current page in the update
		for _, rule := range applicableRules {
			for j := 0; j < i; j++ {
				if DEBUG == "true" {
					fmt.Printf("Checking rule %v against update: %v with applicable rules: %v\n", rule, update, applicableRules)
				}
				// If an applicable rule is found
				if update[j] == rule {
					if DEBUG == "true" {
						fmt.Printf("Rule %d Applicable\n", rule)
						fmt.Printf("Moving page %d from index %d to index %d\n", page, i, j)
					}
					move(update, i, j)
					i = j // Adjust the index since we moved a page
					if DEBUG == "true" {
						fmt.Printf("Update after moving: %v\n", update)
					}
				}
			}
		}
	}
	if DEBUG == "true" {
		fmt.Printf("Sorted update: %v\n", update)
	}

	return update
}

func move(slice []int, startPosition int, endPosition int) {
	DEBUG := os.Getenv("DEBUG")
	if DEBUG == "true" {
		fmt.Printf("Moving %v from index %d to index %d\n", slice[startPosition], startPosition, endPosition)
	}
	if startPosition == endPosition {
		return // No movement needed if startPosition and endPosition are the same
	}

	if startPosition < endPosition {
		value := slice[startPosition]
		// Shift elements from startPosition+1 to endPosition one position to the left
		copy(slice[startPosition:endPosition], slice[startPosition+1:endPosition+1])
		slice[endPosition] = value
	} else {
		value := slice[startPosition]
		// Shift elements from endPosition to startPosition-1 one position to the right
		copy(slice[endPosition+1:startPosition+1], slice[endPosition:startPosition])
		slice[endPosition] = value
	}
}
