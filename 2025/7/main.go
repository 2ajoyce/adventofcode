package main

import (
	"2ajoyce/adventofcode/2025/7/graph"
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
		// printableLine := []rune{}
		for i, r := range StrToArrRune(line) {

			// For the first row
			if rowNum == 0 {
				if r == 'S' {
					// printableLine = append(printableLine, 'S')
					idx[rowNum] = append(idx[rowNum], i) // Add the starting index to the first row
				} else {
					// printableLine = append(printableLine, '.')
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
					// printableLine = printableLine[:len(printableLine)-1]
					// printableLine = append(printableLine, '|')
					idx[rowNum] = append(idx[rowNum], i-1)
					// printableLine = append(printableLine, '^')
					total++
				} else { // The beam continues straight
					// printableLine = append(printableLine, '|')
					idx[rowNum] = append(idx[rowNum], i)
				}
			} else if wasSplitter {
				// Splitters are always followed by a beam
				// printableLine = append(printableLine, '|')
				idx[rowNum] = append(idx[rowNum], i)
			} else {
				// printableLine = append(printableLine, '.')
			}
		}
		// fmt.Println(ArrRuneToStr(printableLine))
		idx[rowNum] = Dedupe(idx[rowNum])
		rowNum++
	}
	return fmt.Sprintf("%d", total), nil
}

type Path struct {
	StartRow int
	StartCol int // The column where the path originated
	PathCol  int // The column the path travels down
}

func Solve2(input chan string) (string, error) {
	total := 0
	idx := [][]int{} // The locations of beams in each row
	rowNum := 0
	g := graph.NewGraph()
	p := map[int][]Path{} // Maps from column index to where path started
	start := ""
	for line := range input {
		idx = append(idx, []int{}) // Add a new row for each line
		for i, r := range StrToArrRune(line) {

			// For the first row
			if rowNum == 0 {
				if r == 'S' {
					idx[rowNum] = append(idx[rowNum], i) // Add the starting index to the first row
					start = fmt.Sprintf("%d-%d", rowNum, i)
					g.AddNode(start)
					p[i] = append(p[i], Path{StartRow: rowNum, StartCol: i, PathCol: i})
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
					g.AddNode(fmt.Sprintf("%d-%d", rowNum, i))
					for _, path := range p[i] {
						g.AddEdge(fmt.Sprintf("%d-%d", path.StartRow, path.StartCol), fmt.Sprintf("%d-%d", rowNum, i))
					}
					delete(p, i) // Remove the completed path

					// Remove the last character
					// I've checked the input and "^^" does not occur
					// If it did this wouldn't work
					idx[rowNum] = append(idx[rowNum], i-1)
					// New path starts here
					p[i-1] = append(p[i-1], Path{StartRow: rowNum, StartCol: i, PathCol: i - 1})
				} else { // The beam continues straight
					idx[rowNum] = append(idx[rowNum], i)
				}
			}
			if wasSplitter {
				// Splitters are always followed by a beam
				idx[rowNum] = append(idx[rowNum], i)
				// new path starts here
				p[i] = append(p[i], Path{StartRow: rowNum, StartCol: i - 1, PathCol: i})
			}
		}
		idx[rowNum] = Dedupe(idx[rowNum])
		rowNum++
	}
	// If there are any open paths, create nodes and terminate them
	for _, paths := range p {
		for _, path := range paths {
			g.AddNode(fmt.Sprintf("%d-%d", rowNum, path.PathCol))
			g.AddEdge(fmt.Sprintf("%d-%d", path.StartRow, path.StartCol), fmt.Sprintf("%d-%d", rowNum, path.PathCol))
		}
	}
	total = g.CountPathsFrom(start)
	return fmt.Sprintf("%d", total), nil
}
