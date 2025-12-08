package graph

import (
	"fmt"
	"slices"
	"strconv"
)

type Graph struct {
	Nodes map[string][]string
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string][]string),
	}
}

// Add Edge adds an edge from a to b
func (g *Graph) AddEdge(a string, b string) {
	// Check that a exists
	if _, exists := g.Nodes[a]; !exists {
		g.AddNode(a)
	}
	if _, exists := g.Nodes[b]; !exists {
		g.AddNode(b)
	}

	// If the node has no edges, create the map
	if g.Nodes[a] == nil {
		g.Nodes[a] = []string{b}
		return
	}
	// If the node already contains this edge, return
	if slices.Contains(g.Nodes[a], b) {
		return
	}
	// Otherwise, add the edge to the map
	g.Nodes[a] = append(g.Nodes[a], b)
}

// Add Edge adds a new node
func (g *Graph) AddNode(n string) {
	if _, exists := g.Nodes[n]; exists {
		// already exists
		return
	}
	g.Nodes[n] = []string{}
}

func StrToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to convert string %q to integer: %v", s, err))
	}
	return num
}

func (g *Graph) CountPathsFrom(n string) int {
	totalEdges := 0
	for _, edges := range g.Nodes {
		totalEdges += len(edges)
	}
	memo := make(map[string]int)

	total := g.countPathsFrom(n, memo)

	return total
}

func (g *Graph) countPathsFrom(n string, memo map[string]int) int {
	// Return cache if possible
	if v, ok := memo[n]; ok {
		return v
	}

	// if n doesn't exist, return 0
	children, exists := g.Nodes[n]
	if !exists {
		memo[n] = 0
		return 0
	}

	// If the node has no children return 1
	if len(children) == 0 {
		memo[n] = 1
		return 1
	}

	total := 0
	for _, child := range children {
		total += g.countPathsFrom(child, memo)
	}

	memo[n] = total
	return total
}
