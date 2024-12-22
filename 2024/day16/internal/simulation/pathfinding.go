package simulation

import (
	"container/heap"
	"fmt"
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
	Cost        float64
	PriorNode   Coord
	CurrentNode Coord
	Path        Path
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
	heap.Push(pq, State{Cost: 0, CurrentNode: start, Path: Path{{Node: start, Cost: 0}}}) // TODO: Update this to save the prior node

	for pq.Len() > 0 {
		state := heap.Pop(pq).(State)
		currentNode := state.CurrentNode
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
				heap.Push(pq, State{Cost: newCost, CurrentNode: neighbor, Path: newPath}) // TODO: Update this to save the prior node
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
// - costFunc: a custom cost function that calculates the cost of transitioning between nodes.
//
// Returns:
// - A list of all optimal paths (each path as a Path struct).
// - The total cost of the optimal paths.
func ModifiedBFS(graph map[Coord]map[Coord]float64, start, target Coord, costFunc func(prior, current, next Coord) float64) ([]Path, float64) {
	dist := make(map[Coord]map[Coord]float64, len(graph)) // Map of current node to prior node and cost
	paths := make(map[Coord]map[Coord][]Path)
	pq := &PriorityQueue{}

	heap.Init(pq)
	for node := range graph {
		dist[node] = make(map[Coord]float64, len(graph[node]))
		paths[node] = make(map[Coord][]Path)
		for neighbor := range graph[node] {
			dist[node][neighbor] = math.MaxFloat64
			paths[node][neighbor] = []Path{}
		}
	}

	dist[start][start] = 0
	paths[start][start] = []Path{{{Node: start, Cost: 0}}}
	heap.Push(pq, State{Cost: 0, PriorNode: start, CurrentNode: start, Path: Path{{Node: start, Cost: 0}}})

	for pq.Len() > 0 {
		state := heap.Pop(pq).(State)
		priorNode := state.PriorNode
		currentNode := state.CurrentNode
		currentCost := state.Cost

		if currentNode == target {
			continue
		}

		for neighbor, weight := range graph[currentNode] {
			newCost := currentCost + weight + costFunc(priorNode, currentNode, neighbor)

			if newCost < dist[neighbor][currentNode] {
				dist[neighbor][currentNode] = newCost
				paths[neighbor][currentNode] = []Path{}
				for _, pathByDirection := range paths[currentNode] {
					for _, path := range pathByDirection {
						newPath := append(Path{}, path...)
						newPath = append(newPath, PathStep{Node: neighbor, Cost: newCost})
						paths[neighbor][currentNode] = append(paths[neighbor][currentNode], newPath)
					}
				}
				heap.Push(pq, State{Cost: newCost, PriorNode: currentNode, CurrentNode: neighbor})
			} else if newCost == dist[neighbor][currentNode] {
				for _, pathsByDirection := range paths[currentNode] {
					for _, path := range pathsByDirection {
						newPath := append(Path{}, path...)
						newPath = append(newPath, PathStep{Node: neighbor, Cost: newCost})
						paths[neighbor][currentNode] = append(paths[neighbor][currentNode], newPath)
					}
				}
			}
		}
	}

	allPathsToTarget := []Path{}
	for _, pathByDirection := range paths[target] {
		allPathsToTarget = append(allPathsToTarget, pathByDirection...)
	}
	fmt.Printf("Found %d paths to target\n", len(allPathsToTarget))

	cheapestCost := math.MaxFloat64
	for _, path := range allPathsToTarget {
		for _, step := range path {
			if step.Node == target && step.Cost < cheapestCost {
				cheapestCost = step.Cost
			}
		}
	}

	cheapestPaths := []Path{}
	for _, path := range allPathsToTarget {
		cost := 0.0
		priorStep := path[0]
		for i, step := range path {
			if i > 0 {
				priorStep = path[i-1]
			}
			if step.Node == target {
				break
			}
			cost += costFunc(priorStep.Node, step.Node, path[i+1].Node)
		}
		if cost == cheapestCost {
			cheapestPaths = append(cheapestPaths, path)
		}
	}
	if len(cheapestPaths) == 0 {
		return cheapestPaths, -1
	}
	return cheapestPaths, cheapestCost
}
