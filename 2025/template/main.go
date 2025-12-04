package main

import (
	"bufio"
	"fmt"
	"os"
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
	for line := range input {
		total += len(string(line)) // Increment the total by the number of characters in the line
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
