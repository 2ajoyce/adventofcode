package simulation

import (
	"container/heap"
	"math"
)

func CostManhattan(from, to Coord) float64 {
	return math.Abs(float64(from.X-to.X)) + math.Abs(float64(from.Y-to.Y))
}

// PathStep represents a step in a path, including the node and the cost to reach that node.
type PathStep struct {
	Node Coord
	Cost float64
}

// Path represents a list of PathSteps.
type Path []PathStep

// State represents a state in the priority queue for pathfinding.
type State struct {
	Cost float64
	Node Coord
	Path Path
}

// PriorityQueue is a min-heap of States.
type PriorityQueue []State

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Cost < pq[j].Cost }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(State))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// Dijkstra finds the single shortest path in a graph.
//
// Arguments:
// - graph: a map where keys are nodes (Coord) and values are maps of neighbors (Coord) and their weights.
// - start: the starting node as a Coord.
// - target: the target node as a Coord.
// - costFn: a custom cost function that calculates the cost of transitioning between nodes.
//
// Returns:
// - The shortest path as a Path struct.
// - The total cost of the shortest path.
func Dijkstra(graph map[Coord]map[Coord]float64, start, target Coord, costFn func(prior, current, next Coord) float64) (Path, float64) {
	dist := make(map[Coord]float64)
	prev := make(map[Coord]Coord)
	pq := &PriorityQueue{}

	heap.Init(pq)

	for node := range graph {
		dist[node] = math.MaxInt
	}
	dist[start] = 0
	heap.Push(pq, State{Cost: 0, Node: start, Path: Path{{Node: start, Cost: 0}}})

	for pq.Len() > 0 {
		state := heap.Pop(pq).(State)
		currentNode := state.Node
		currentCost := state.Cost

		if currentNode == target {
			return state.Path, currentCost
		}

		for neighbor, weight := range graph[currentNode] {
			newCost := currentCost + costFn(Coord{}, currentNode, neighbor) + weight
			if newCost < dist[neighbor] {
				dist[neighbor] = newCost
				prev[neighbor] = currentNode
				newPath := append(Path{}, state.Path...)
				newPath = append(newPath, PathStep{Node: neighbor, Cost: newCost})
				heap.Push(pq, State{Cost: newCost, Node: neighbor, Path: newPath})
			}
		}
	}

	return nil, -1 // No path found
}

// ModifiedBFS finds all optimal paths in a graph.
//
// Arguments:
// - graph: a map where keys are nodes (Coord) and values are maps of neighbors (Coord) and their weights.
// - start: the starting node as a Coord.
// - target: the target node as a Coord.
//
// Returns:
// - A list of all optimal paths (each path as a Path struct).
// - The total cost of the optimal paths.
func ModifiedBFS(graph map[Coord]map[Coord]float64, start, target Coord) ([]Path, float64) {
	dist := make(map[Coord]float64)
	paths := make(map[Coord][]Path)
	pq := &PriorityQueue{}

	heap.Init(pq)

	for node := range graph {
		dist[node] = math.MaxInt
		paths[node] = []Path{}
	}
	dist[start] = 0
	paths[start] = []Path{{{Node: start, Cost: 0}}}
	heap.Push(pq, State{Cost: 0, Node: start, Path: Path{{Node: start, Cost: 0}}})

	for pq.Len() > 0 {
		state := heap.Pop(pq).(State)
		currentNode := state.Node
		currentCost := state.Cost

		for neighbor, weight := range graph[currentNode] {
			newCost := currentCost + weight

			if newCost < dist[neighbor] {
				dist[neighbor] = newCost
				paths[neighbor] = []Path{}
				for _, path := range paths[currentNode] {
					newPath := append(Path{}, path...)
					newPath = append(newPath, PathStep{Node: neighbor, Cost: newCost})
					paths[neighbor] = append(paths[neighbor], newPath)
				}
				heap.Push(pq, State{Cost: newCost, Node: neighbor})
			} else if newCost == dist[neighbor] {
				for _, path := range paths[currentNode] {
					newPath := append(Path{}, path...)
					newPath = append(newPath, PathStep{Node: neighbor, Cost: newCost})
					paths[neighbor] = append(paths[neighbor], newPath)
				}
			}
		}
	}

	return paths[target], dist[target]
}
