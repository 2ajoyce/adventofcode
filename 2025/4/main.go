package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// First Problem
	input := make(chan [][]rune)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan [][]rune)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan [][]rune) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	result := [][]rune{}
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, ParseInput(line))
	}
	c <- result
	close(c)
}

// ParseInput parses the input into the necessary data structure.
// On more complex inputs, this allows us to use lines of text as input for tests
func ParseInput(input string) []rune {
	return []rune(input)
}

func Solve1(input chan [][]rune) (string, error) {
	total := 0
	// The use of a channel here is contrived, but most problems have been processed line by line
	grid := <-input
	PrintGrid(grid)
	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan [][]rune) (string, error) {
	total := 0
	// The use of a channel here is contrived, but most problems have been processed line by line
	grid := <-input
	PrintGrid(grid)
	return fmt.Sprintf("%d", total), nil
}

func PrintGrid(g [][]rune) {
	for _, row := range g {
		fmt.Print("|")
		for _, char := range row {
			fmt.Printf("%c|", char)
		}
		fmt.Println()
	}
}
