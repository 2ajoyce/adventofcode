package main

import (
	"2ajoyce/adventofcode/2025/10/equation"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// First Problem
	input := make(chan *equation.Equation)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan *equation.Equation)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan *equation.Equation) {
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

// Split the input line into
func ParseInput(input string) *equation.Equation {
	// [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
	eq := equation.NewEquation(strings.TrimSpace(input))
	return &eq
}

func Solve1(input chan *equation.Equation) (string, error) {
	total := 0

	// Collect all equations
	equations := []*equation.Equation{}
	for eq := range input {
		equations = append(equations, eq)
	}

	// Solve each equation
	for _, eq := range equations {
		minButtonPushes, err := minButtonPressesBFS(eq)
		if err != nil {
			return "", err
		}
		total += minButtonPushes
	}

	return fmt.Sprintf("%d", total), nil
}

// minButtonPressesBFS returns the minimum number of button presses to
// get from the start state (all zeros) to eq.Target using BFS or an
// error if the target is unreachable.
func minButtonPressesBFS(eq *equation.Equation) (int, error) {
	// The  initial state is all zeros
	// Since State is a uint16, all zeros is just 0
	var start equation.State = 0

	// If start is already the target return 0
	if start == eq.Target {
		return 0, nil
	}

	// The queue holds states to explore
	queue := make([]equation.State, 0)
	// visited tracks states which have been seen
	visited := make(map[equation.State]bool)
	// dist tracks the number of button presses to reach each state
	dist := make(map[equation.State]int)

	// Add the start state to the queue, mark it visited, and set its distance to 0
	queue = append(queue, start)
	visited[start] = true
	dist[start] = 0

	head := 0
	for head < len(queue) {
		// Dequeue the next state
		cur := queue[head]
		head++
		// Try pressing each button
		for _, btn := range eq.Buttons {
			next := cur.PressButton(btn)
			// If we reached the target, return the distance
			if next == eq.Target {
				return dist[cur] + 1, nil
			}
			// If already visited, skip
			if !visited[next] {
				// Mark as visisted and record distance
				visited[next] = true
				dist[next] = dist[cur] + 1
				// Enqueue the new state
				queue = append(queue, next)
			}
		}
	}

	return 0, fmt.Errorf("no solution for equation: %v", eq)
}

func Solve2(input chan *equation.Equation) (string, error) {
	total := 0

	// Collect all equations
	equations := []*equation.Equation{}
	for eq := range input {
		equations = append(equations, eq)
	}

	// Solve each equation
	for _, eq := range equations {
		minButtonPushes := 0
		eq.Target.PressButton(eq.Buttons[0]) // placeholder, impliment bfs
		total += minButtonPushes
	}

	return fmt.Sprintf("%d", total), nil
}
