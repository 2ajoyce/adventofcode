package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
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
	problems := [][]int{}  // An array of the problems. Each problem is an array of numbers.
	operands := []string{} // Each problem has an operand
	firstLine := true
	for line := range input {
		// The last line is operands
		if strings.Contains(line, "+") || strings.Contains(line, "*") {
			for v := range strings.FieldsSeq(line) {
				operands = append(operands, v)
			}
			continue
		}

		for i, v := range strings.Fields(line) {
			if firstLine {
				problems = append(problems, []int{})
			}
			problems[i] = append(problems[i], StrToInt(v))
			i++
		}
		firstLine = false
	}

	// For every problem
	for i, p := range problems {
		solution := p[0] // Start with the first number
		for j, v := range p {
			if j == 0 { // Skip the first number
				continue
			}

			switch operands[i] {
			case "+":
				solution = solution + v
			case "*":
				solution = solution * v
			}
		}
		total = total + solution
	}

	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan string) (string, error) {
	total := 0
	var problems [][]int   // An array of the problems. Each problem is an array of numbers.
	operands := []string{} // Each problem has an operand
	lines := [][]rune{}
	for line := range input {
		lines = append(lines, StrToArrRune(line))
	}

	// Remove the last line (operands)
	operandLine := lines[len(lines)-1]
	lines = lines[:len(lines)-1]
	for v := range strings.FieldsSeq(ArrRuneToStr(operandLine)) {
		operands = append(operands, v)
	}
	// Reverse opperands
	slices.Reverse(operands)

	// Find the longest line
	longest := 0
	problemCount := 0
	for _, line := range lines {
		if longest < len(line) {
			longest = len(line)
			problemCount = len(strings.Fields(ArrRuneToStr(line)))
		}
	}

	// Initialize the problems array
	problems = make([][]int, problemCount)

	// The longest line determines how many columns we will read in
	currentProblem := 0
	for i := longest; i >= 0; i-- {
		emptyDigitCount := 0
		problemString := []rune{}
		for _, line := range lines {
			if len(line) <= i {
				continue
			}
			if line[i] == ' ' {
				emptyDigitCount++
				continue
			}
			problemString = append(problemString, line[i])
		}
		if emptyDigitCount == len(lines) {
			currentProblem++
		}
		if len(problemString) > 0 {
			problems[currentProblem] = append(problems[currentProblem], ArrRuneToInt(problemString))
		}
	}

	// For every problem
	for i, p := range problems {
		solution := p[0] // Start with the first number
		for j, v := range p {
			if j == 0 { // Skip the first number
				continue
			}

			switch operands[i] {
			case "+":
				solution = solution + v
			case "*":
				solution = solution * v
			}
		}
		total = total + solution
	}

	return fmt.Sprintf("%d", total), nil
}
