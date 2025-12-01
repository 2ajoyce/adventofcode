package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const DIAL_SIZE = 100

func main() {
	// First Problem
	input := make(chan string)
	go ReadInput("input1.txt", input)
	result, err := Solve(input)
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

// ReadInput reads the input from the input.txt file and sends each line to the provided channel.
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
		c <- line
	}
	close(c)
}

func Solve(input chan string) (string, error) {
	number := 50 // The instructions specify that the dial starts on 50
	total := 0   // This is our output, indicating the number of times the dial landed on zero
	for line := range input {
		// Read in the input
		direction := string(line[0])
		increment, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}
		increment = increment % DIAL_SIZE // Adding a full turn does not change the result

		// Rotate the dial
		if string(direction) == "L" {
			number = MoveLeft(number, increment)
		}
		if string(direction) == "R" {
			number = MoveRight(number, increment)
		}

		// If the dial ends on zero, increment the total.
		if number == 0 {
			total++
		}
	}
	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan string) (string, error) {
	number := 50 // The instructions specify that the dial starts on 50
	total := 0   // This is our output, indicating the number of times the dial passed zero
	for line := range input {
		// Read in the input
		direction := string(line[0])
		increment, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}
		total += increment / DIAL_SIZE    // Count each full spin as passing zero
		increment = increment % DIAL_SIZE // Adding a full turn does not change the result

		pre := number
		// Rotate the dial
		if string(direction) == "L" {
			number = MoveLeft(number, increment)
			if number != 0 && pre != 0 && number > pre {
				total++
			}
		}
		if string(direction) == "R" {
			number = MoveRight(number, increment)
			if number != 0 && pre != 0 && number < pre {
				total++
			}
		}

		// If the dial ends on zero, increment the total.
		if number == 0 {
			total++
		}
	}
	return fmt.Sprintf("%d", total), nil
}

func MoveRight(current int, increment int) int {
	if current > DIAL_SIZE || increment > DIAL_SIZE {
		panic("MoveRight recieved an input larger than DIAL_SIZE")
	}
	return (current + increment) % DIAL_SIZE
}

func MoveLeft(current int, increment int) int {
	if current > DIAL_SIZE || increment > DIAL_SIZE {
		panic("MoveLeft recieved an input larger than DIAL_SIZE")
	}
	return (current - increment + DIAL_SIZE) % DIAL_SIZE
}
