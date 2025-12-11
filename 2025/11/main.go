package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	// First Problem
	input := make(chan *Graph)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan *Graph)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan *Graph) {
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
var (
	nameRegex  = regexp.MustCompile(`^([a-z]{3}):`)
	childRegex = regexp.MustCompile(`(?::\s*|\s)([a-z]{3})`)
)

func ParseInput(input string) *Graph {
	name := nameRegex.FindStringSubmatch(input)[1]

	childMatches := childRegex.FindAllStringSubmatch(input, -1)
	children := make([]string, 0, len(childMatches))
	for _, match := range childMatches {
		children = append(children, match[1])
	}

	g := make(Graph)
	g[name] = children
	return &g
}

func Solve1(input chan *Graph) (string, error) {
	total := 0
	graph := Graph{}
	for n := range input {
		// Merge the parsed graph into the main graph
		for k, v := range *n {
			// if key already exists, panic
			if _, ok := graph[k]; ok {
				panic(fmt.Sprintf("Input contains duplicate key: %s", k))
			}
			graph[k] = v
		}
	}

	total = graph.countPathsDfs("you", "out")

	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan *Graph) (string, error) {
	total := 0
	graph := Graph{}
	for n := range input {
		// Merge the parsed graph into the main graph
		for k, v := range *n {
			// if key already exists, panic
			if _, ok := graph[k]; ok {
				panic(fmt.Sprintf("Input contains duplicate key: %s", k))
			}
			graph[k] = v
		}
	}

	total = graph.countPathsDfs("you", "out")

	return fmt.Sprintf("%d", total), nil
}

type Graph map[string][]string

func (g *Graph) countPathsDfs(current, target string) int {
	return dfsHelper(*g, current, target, make(map[string]bool))
}

func dfsHelper(g Graph, current, target string, visited map[string]bool) int {
	if current == target {
		return 1
	}

	visited[current] = true
	defer func() { visited[current] = false }() // backtrack

	total := 0
	for _, next := range g[current] {
		if visited[next] {
			// avoid cycles
			continue
		}
		total += dfsHelper(g, next, target, visited)
	}
	return total
}
