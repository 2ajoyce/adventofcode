package internal

import (
	"day10/internal/simulation"
	"fmt"
	"os"
	"strconv"
)

type TopoEntity struct {
	simulation.Entity
	zHeight int
}

func NewTopoEntity(zHeight int) (*TopoEntity, error) {
	entity, err := simulation.NewEntity()
	if err != nil {
		return nil, err
	}
	var te = new(TopoEntity)
	te.Entity = entity
	te.zHeight = zHeight
	return te, nil
}

func (t TopoEntity) GetZHeight() int {
	return t.zHeight
}

func StringifySimulation(sim simulation.Simulation) (string, error) {
	simMap := sim.GetMap()
	y := simMap.GetHeight()
	x := simMap.GetWidth()
	output := ""
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			cell, err := simMap.GetCell(j, i)
			if err != nil {
				return "", err
			}
			entityIds, err := cell.GetEntityIds()
			if err != nil {
				return "", err
			}
			entity, err := sim.GetEntity(entityIds[0]) // Fragile, but we're only expecting one entity per cell for now.
			if err != nil {
				return "", err
			}
			topoEntity, ok := entity.(*TopoEntity)
			if !ok {
				return "", fmt.Errorf("entity is not of type TopoEntity")
			}
			output += strconv.Itoa(topoEntity.GetZHeight())
			if j < x-1 {
				output += " "
			}
		}
		if i < y-1 {
			output += "\n"
		}
	}
	return output, nil
}

type Coord struct {
	X, Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

func (c1 Coord) Equals(c2 Coord) bool {
	return c1.X == c2.X && c1.Y == c2.Y
}

// Return the coordinates of all trailheads in the simulation
// Trailheads are entities with the zHeight attribute set to 0
func FindTrailheads(sim simulation.Simulation) ([]Coord, error) {
	trailheads := []Coord{}
	sim.GetEntities()
	for _, entity := range sim.GetEntities() {
		topoEntity, ok := entity.(*TopoEntity)
		if !ok {
			continue
		}
		if topoEntity.GetZHeight() == 0 {
			x, y := topoEntity.GetPosition()
			trailheads = append(trailheads, Coord{X: x, Y: y})

		}
	}
	return trailheads, nil
}

// Starting at the trailhead coordinates, check every horizontally or vertically adjacent cell
// If the zHeight of the cell is one higher than the current cell, add it to a list of possible paths
// Continue this process until every path has been fully explored
// A path is fully explored
//   - when the current cell has a zHeight of 9
//   - if there are no more possible paths to explore
//
// Each cell can be part of multiple paths
// Return the number of paths which reached a cell with a zHeight of 9 as the score
func ScoreTrailhead(sim simulation.Simulation, trailhead Coord) (score int, err error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	score = 0

	if DEBUG {
		fmt.Printf("Starting ScoreTrailhead with trailhead at Coord%s\n", trailhead)
	}

	// Each trailhead -> zHeight = 9 path can only be scored once
	// This is a map of each trailhead to the height 9 tiles it can reach
	trailHeadToEnd := make(map[Coord][]Coord)

	// Create a cache of the zHeights to avoid redundant lookups
	zHeightCache := make(map[Coord]int)
	for _, entity := range sim.GetEntities() {
		topoEntity, ok := entity.(*TopoEntity)
		if !ok {
			if DEBUG {
				fmt.Printf("Skipping entity of type %T, not a TopoEntity\n", entity)
			}
			continue
		}
		zHeight := topoEntity.GetZHeight()
		x, y := topoEntity.GetPosition()
		coord := Coord{X: x, Y: y}
		zHeightCache[coord] = zHeight
		if DEBUG {
			fmt.Printf("Caching zHeight for Coord%s: %d\n", coord, zHeight)
		}

	}

	if DEBUG {
		fmt.Printf("Initial zHeightCache: %+v\n", zHeightCache)
	}

	// Create a queue of the cells to explore
	queue := []Coord{trailhead}
	if DEBUG {
		fmt.Printf("Initial queue: %+v\n", queue)
	}

	for len(queue) > 0 {
		current := queue[0] // Set the current cell to the first cell in the queue
		queue = queue[1:]   // Remove the first cell in the queue
		if DEBUG {
			fmt.Printf("Dequeued Coord%s, current score: %d\n", current, score)
		}

		currentZ, exists := zHeightCache[current]
		if !exists {
			if DEBUG {
				fmt.Printf("Coord%s not in zHeightCache, skipping\n", current)
			}
			continue
		}

		if DEBUG {
			fmt.Printf("Exploring neighbors for Coord%s with zHeight %d\n", current, currentZ)

		}

		if currentZ == 9 {
			priorTrailEnds := trailHeadToEnd[trailhead]
			alreadyVisited := false
			for _, end := range priorTrailEnds {
				// If we have already scored this trailend, skip it
				// Optionally, we could skip the rest of this loop
				if end == current {
					alreadyVisited = true
					break
				}
			}
			if alreadyVisited {
				if DEBUG {
					fmt.Printf("Already visited Coord%s, skipping\n", current)
				}
				continue
			}
			score++
			trailHeadToEnd[trailhead] = append(trailHeadToEnd[trailhead], current)
			if DEBUG {
				fmt.Println("Reached zHeight 9")
				fmt.Printf("Started from trailhead at Coord%s\n", trailhead)
				fmt.Printf("Ended at Coord%s\n", current)
				fmt.Printf("Score: %d\n", score)
				fmt.Printf("Trailhead %s to End: %+v\n", trailhead, trailHeadToEnd[trailhead])
			}
			continue
		}

		// Select the coordinates of the neighboring cells
		neighbors := []Coord{
			{current.X - 1, current.Y},
			{current.X + 1, current.Y},
			{current.X, current.Y - 1},
			{current.X, current.Y + 1},
		}

		if DEBUG {
			fmt.Printf("Neighbors of Coord%s: %+v\n", current, neighbors)
		}

		for _, neighbor := range neighbors {
			if !sim.GetMap().ValidateCoord(neighbor.X, neighbor.Y) {
				continue
			}

			neighborZ, neighborExists := zHeightCache[neighbor]
			if !neighborExists {
				if DEBUG {
					fmt.Printf("Neighbor Coord%s not in zHeightCache, skipping\n", neighbor)
				}
				continue
			}

			if neighborZ-currentZ == 1 {
				queue = append(queue, neighbor)
				if DEBUG {
					fmt.Printf("Neighbor Coord%s has zHeight %d (one higher than current %d). Added to queue.\n", neighbor, neighborZ, currentZ)
				}
			} else {
				if DEBUG {
					fmt.Printf("Neighbor Coord%s has zHeight %d (not one higher than current %d). Skipping.\n", neighbor, neighborZ, currentZ)
				}
			}
		}

		if DEBUG {
			fmt.Printf("Queue after processing Coord%s: %+v\n", current, queue)
		}
	}

	if DEBUG {
		fmt.Printf("Final score: %d\n", score)
	}

	return
}
