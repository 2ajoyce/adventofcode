package main

import (
	"bufio"
	"fmt"
	"maps"
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

	total = graph.countPathsDfs("you", "out", nil)

	return fmt.Sprintf("%d", total), nil
}

type Graph map[string][]string

func (g *Graph) countPathsDfs(current, target string, mustVisit []string) int {
	visited := make(map[string]bool)

	// Track required nodes using a map[string]bool
	req := make(map[string]bool)
	for _, m := range mustVisit {
		req[m] = false
	}

	return dfsHelper(*g, current, target, visited, req)
}

func dfsHelper(g Graph, current, target string, visited map[string]bool, req map[string]bool) int {
	// Mark if current is a required node
	if _, ok := req[current]; ok {
		req[current] = true
	}

	// If we reached the target, all required nodes must be visited
	if current == target {
		for _, seen := range req {
			if !seen {
				return 0
			}
		}
		return 1
	}

	visited[current] = true
	defer func() { visited[current] = false }()

	total := 0
	for _, next := range g[current] {
		if visited[next] {
			continue
		}

		// Clone req map for the child call (important!)
		nextReq := make(map[string]bool, len(req))
		maps.Copy(nextReq, req)

		total += dfsHelper(g, next, target, visited, nextReq)
	}

	return total
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

	//Find every path from "svr" to "out"
	// The paths must all also visit both "dac" and "fft" (in any order).
	// Assume the graph is acyclic for this DP to be correct
	svrToDac := graph.CountPaths("svr", "dac")
	dacToFft := graph.CountPaths("dac", "fft")
	fftToOut := graph.CountPaths("fft", "out")

	svrToFft := graph.CountPaths("svr", "fft")
	fftToDac := graph.CountPaths("fft", "dac")
	dacToOut := graph.CountPaths("dac", "out")

	total = svrToDac*dacToFft*fftToOut + svrToFft*fftToDac*dacToOut

	return fmt.Sprintf("%d", total), nil
}

func (g Graph) countPathsMemo(start, target string, memo map[string]int) int {
	if val, ok := memo[start]; ok {
		return val
	}
	if start == target {
		return 1
	}

	total := 0
	for _, next := range g[start] {
		total += g.countPathsMemo(next, target, memo)
	}
	memo[start] = total
	return total
}

func (g Graph) CountPaths(start, target string) int {
	memo := make(map[string]int)
	return g.countPathsMemo(start, target, memo)
}
