package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// First Problem
	input := make(chan []int)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan []int)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan []int) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c <- ParseLine(line)
	}
	close(c)
}

// Parse line converts a string into the type of the channel
// I'm not totally comfortable with this rune arithmatic approach, but it promises better performance
func ParseLine(s string) []int {
	out := make([]int, len(s))
	for i, r := range s {
		n := int(r - '0')
		if n < 0 || n > 9 {
			panic(fmt.Sprintf("invalid digit %q in %q", r, s))
		}
		out[i] = n
	}
	return out
}

func Solve1(input chan []int) (string, error) {
	total := 0
	for line := range input {
		digits := FindLargestPair(line)
		total += digitsToInt(digits)
	}
	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan []int) (string, error) {
	total := 0
	for line := range input {
		digits := FindLargestTwelve(line)
		total += digitsToInt(digits)
	}
	return fmt.Sprintf("%d", total), nil
}

func FindLargestPair(line []int) []int {
	if len(line) < 2 {
		return nil
	}

	// Initialize our two pointers
	p1, p2 := 0, 1

	// Iterate over the line
	for i := range line {
		if i == 0 { // Skip the first value
			continue
		}
		// If the current index is larger than existing p1 value, reassign
		// Skip if on the last value of the line
		if line[i] > line[p1] && i < len(line)-1 {
			p1 = i
			// Wiping p2 ensures that we don't hold onto an old reference
			p2 = -1
		} else {
			// If p2 is unassigned or we didn't reassign p1 we can consider reassigning p2
			if p2 == -1 || line[i] > line[p2] {
				p2 = i
			}
		}
	}

	if p2 == -1 { // This should never happen
		panic(fmt.Sprintf("we reached the end of the loop without a p2 value for line %v", line))
	}
	return []int{line[p1], line[p2]}
}

func FindLargestTwelve(line []int) []int {
	var NUMBER_OF_POINTERS = 12

	if len(line) < NUMBER_OF_POINTERS {
		return nil
	}

	// Initialize our pointers
	pointers := make([]int, NUMBER_OF_POINTERS)
	pointers[0] = 0
	for i := 1; i < len(pointers); i++ {
		pointers[i] = -1
	}

	// Iterate over the line
	for i := range line {
		if i == 0 { // Skip the first value
			continue
		}
		// If the current line[index] is larger than existing p0 value, reassign
		// Skip if there aren't enough values left to fit all pointers
		if line[i] > line[pointers[0]] && i <= len(line)-len(pointers) {
			pointers[0] = i
			wipePointers(pointers, 1)
		} else {
			// Iterate over the pointers to check if we can update any of them
			for j, p := range pointers {
				if j == 0 { // Skip the first pointer
					continue
				}

				// If the current pointer is unassigned
				if p == -1 {
					pointers[j] = i
					wipePointers(pointers, j+1)
					break
					// Else if the current pointer is assigned but smaller
					// AND we have enough digits left
				} else if line[i] > line[p] && len(line)-i >= len(pointers)-j {
					pointers[j] = i
					wipePointers(pointers, j+1)
					break
				}
			}
		}
	}

	// Initialize our result
	result := make([]int, NUMBER_OF_POINTERS)
	for i := 0; i < len(pointers); i++ {
		result[i] = line[pointers[i]]
	}
	return result
}

func wipePointers(p []int, start int) {
	for i := start; i < len(p); i++ {
		p[i] = -1
	}
}

func digitsToInt(digits []int) int {
	n := 0
	for _, d := range digits {
		n = n*10 + d
	}
	return n
}
