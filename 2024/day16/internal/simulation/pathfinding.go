package simulation

import (
	hp "container/heap"
	"fmt"
	"math"
)

type path struct {
	value float64
	nodes []Coord
}

type minPath []path

func (h minPath) Len() int           { return len(h) }
func (h minPath) Less(i, j int) bool { return h[i].value < h[j].value }
func (h minPath) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minPath) Push(x interface{}) {
	*h = append(*h, x.(path))
}

func (h *minPath) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type pathHeap struct {
	values *minPath
}

func newPathHeap() *pathHeap {
	return &pathHeap{values: &minPath{}}
}

func (h *pathHeap) push(p path) {
	hp.Push(h.values, p)
}

func (h *pathHeap) pop() path {
	i := hp.Pop(h.values)
	return i.(path)
}

func CostManhattan(dir Direction, from, to Coord) float64 {
	return math.Abs(float64(from.X-to.X)) + math.Abs(float64(from.Y-to.Y))
}

type CoordDirection struct {
	Coord     Coord
	Direction Direction
}

// Dijkstra finds the shortest path from the start coordinate to the goal coordinate
// using Dijkstra's algorithm. It returns the path as a slice of coordinates, the total
// cost of the path, the number of visited nodes, and an error if the path is not found.
//
// Parameters:
// - sim: The simulation containing the map and other relevant data.
// - start: The starting coordinate.
// - goal: The goal coordinate.
// - costFunction: A function that calculates the cost between two coordinates.
//
// Returns:
// - []Coord: The shortest path from start to goal.
// - int: The total cost of the path.
// - int: The number of visited nodes.
// - error: An error if the path is not found.
func Dijkstra(sim Simulation, start, goal Coord, costFunction func(Direction, Coord, Coord) float64) (shortestPath []Coord, cost float64, visitedNodes int, err error) {
	h := newPathHeap()
	h.push(path{value: 0, nodes: []Coord{start}})
	visited := make(map[CoordDirection]bool)
	solutions := make([]path, 0)

	// Find the other entities on the map
	entities := sim.GetEntities()
	impassable := make(map[Coord]bool, len(entities)-2) // Minus 2 for the start and goal
	for _, entity := range entities {
		coords := entity.GetPosition()
		for _, coord := range coords {
			impassable[coord] = true
		}
	}
	impassable[start] = false
	impassable[goal] = false

	for len(*h.values) > 0 {
		p := h.pop()
		node := p.nodes[len(p.nodes)-1]
		priorCord := node
		dir := East
		if len(p.nodes) > 1 {
			priorCord = p.nodes[len(p.nodes)-2]
			dir = priorCord.DirectionTo(node)
		}
		cordDir := CoordDirection{Coord: node, Direction: dir}
		if visited[cordDir] {
			continue
		}

		if node == goal {
			solutions = append(solutions, path{p.value, p.nodes})
		}

		for _, neighbor := range sim.GetMap().GetNeighbors(node) {
			if len(p.nodes) < 1 {
				dir = East
			}
			if !visited[CoordDirection{Coord: neighbor, Direction: dir}] {
				if impassable[neighbor] {
					continue
				}
				cost := costFunction(dir, node, neighbor)
				h.push(path{value: p.value + cost, nodes: append([]Coord{}, append(p.nodes, neighbor)...)})
			}
		}
		visited[cordDir] = true
	}

	fmt.Printf("Found %d solutions\n", len(solutions))
	fmt.Printf("Visited %d nodes\n", len(visited))
	if len(solutions) == 0 {
		return nil, 0, 0, fmt.Errorf("path not found")
	}
	CheapestSolution := solutions[0]
	for _, solution := range solutions {
		if solution.value < CheapestSolution.value {
			CheapestSolution = solution
		}
	}
	fmt.Printf("Cheapest solution costs %v\n", CheapestSolution.value)
	return CheapestSolution.nodes, CheapestSolution.value, len(visited), nil

}
