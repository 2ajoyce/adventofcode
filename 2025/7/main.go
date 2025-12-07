package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func main() {
	// First Problem
	input := make(chan string)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan string)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan string) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c <- ParseInput(line)
	}
	close(c)
}

// ParseInput parses the input into the necessary data structure.
// On more complex inputs, this allows us to use lines of text as input for tests
func ParseInput(input string) string {
	return input
}

func Solve1(input chan string) (string, error) {
	total := 0
	idx := [][]int{} // The locations of beams in each row
	rowNum := 0
	for line := range input {
		idx = append(idx, []int{}) // Add a new row for each line
		printableLine := []rune{}
		for i, r := range StrToArrRune(line) {

			// For the first row
			if rowNum == 0 {
				if r == 'S' {
					printableLine = append(printableLine, 'S')
					idx[rowNum] = append(idx[rowNum], i) // Add the starting index to the first row
				} else {
					printableLine = append(printableLine, '.')
				}
				continue
			}

			// For subsequent rows
			// If the prior row has a beam at this index, then
			// this row also has a beam at this index
			isBeamAbove := slices.Contains(idx[rowNum-1], i)
			isSplitter := r == '^'
			wasSplitter := i > 0 && StrToArrRune(line)[i-1] == '^'
			if isBeamAbove {
				if isSplitter {
					// Remove the last character
					// I've checked the input and "^^" does not occur
					// If it did this wouldn't work
					printableLine = printableLine[:len(printableLine)-1]
					printableLine = append(printableLine, '|')
					idx[rowNum] = append(idx[rowNum], i-1)
					printableLine = append(printableLine, '^')
					total++
				} else { // The beam continues straight
					printableLine = append(printableLine, '|')
					idx[rowNum] = append(idx[rowNum], i)
				}
			} else if wasSplitter {
				// Splitters are always followed by a beam
				printableLine = append(printableLine, '|')
				idx[rowNum] = append(idx[rowNum], i)
			} else {
				printableLine = append(printableLine, '.')
			}
		}
		fmt.Println(ArrRuneToStr(printableLine))
		idx[rowNum] = Dedupe(idx[rowNum])
		rowNum++
	}
	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan string) (string, error) {
	total := 0
	for line := range input {
		total += len(string(line)) // Increment the total by the number of characters in the line
	}
	return fmt.Sprintf("%d", total), nil
}
