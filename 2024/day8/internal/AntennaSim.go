package internal

import (
	"day8/internal/simulation"
	"fmt"
	"os"
)

type AntennaSimulation struct {
	simulation.Simulation
}

func NewAntennaSimulation(width, height int) (*AntennaSimulation, error) {
	sim := simulation.NewSimulation(width, height)
	var as = new(AntennaSimulation)
	as.Simulation = sim
	return as, nil
}

func (as *AntennaSimulation) String() (result string) {
	height := as.Simulation.GetMap().GetHeight()
	width := as.Simulation.GetMap().GetWidth()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cell, err := as.Simulation.GetMap().GetCell(x, y)
			if err != nil {
				result += "?" // If there is an error fetching a cell, use "?" as a placeholder
				continue
			}

			entityIds, err := cell.GetEntityIds()
			if err != nil {
				if _, ok := err.(simulation.CellEmptyError); ok {
					result += "." // If the cell is not occupied, use "." as a placeholder
				} else {
					result += "?" // Use "?" as a placeholder if there is an error fetching the entity ID
				}
				continue
			}
			var entitySymbol = "?"
			var foundAntenna = false
			for _, e := range entityIds {
				entity, err := as.GetEntity(e)
				if err != nil {
					result += "?" // Use "?" as a placeholder if there is an error
				}
				switch e := entity.(type) {
				case *Antenna: // Order is important to prioritize showing antennas
					entitySymbol = e.String()
					foundAntenna = true
				case *Antinode:
					entitySymbol = e.String()
				}
				if foundAntenna {
					break
				}

			}
			result += entitySymbol
		}
		result += "\n"
	}
	return result
}

func (as *AntennaSimulation) AddAntenna(newAnt *Antenna, x int, y int) (*Antenna, error) {
	ant, err := as.Simulation.AddEntity(newAnt, x, y)
	if err != nil {
		return nil, err
	}
	as.updateAntinodes()
	return ant.(*Antenna), nil
}

// Helper function to place antinodes in a given direction until out of bounds
func (as *AntennaSimulation) placeAntinodesInDirection(startX, startY, ux, uy, step int) error {
	DEBUG := os.Getenv("DEBUG") == "true"

	x := startX + ux*step
	y := startY + uy*step

	for {
		if DEBUG {
			fmt.Printf("Checking position (%d, %d)\n", x, y)
		}

		if valid := as.GetMap().ValidateCoord(x, y); !valid {
			if DEBUG {
				fmt.Printf("Antinode position (%d, %d) is out of bounds. Stopping in this direction.\n", x, y)
			}
			break
		}

		// Place antinode
		if DEBUG {
			fmt.Printf("Placing antinode at (%d, %d)\n", x, y)
		}
		if err := as.placeAntinode(x, y); err != nil {
			return fmt.Errorf("failed to place antinode at (%d, %d): %v", x, y, err)
		}
		// Move to the next position in the direction
		x += ux * step
		y += uy * step
		if DEBUG {
			fmt.Printf("Moving to next position in direction (%d, %d)\n", ux, uy)
		}
	}

	return nil
}

// updateAntinodes updates the antinodes in the simulation based on the current antennas.
func (as *AntennaSimulation) updateAntinodes() (*AntennaSimulation, error) {
	// Antinodes occur when two antennas with the same signal are in line with each other
	// Each pair of aligned antennas will create a line of antinodes
	// The spacing between each antinode is the same as the distance between the antenna and the first antinode

	DEBUG := os.Getenv("DEBUG") == "true"

	// First, remove all existing antinodes
	entities := as.GetEntities()
	for _, e := range entities {
		if _, ok := e.(*Antinode); ok {
			success, err := as.RemoveEntity(e.GetId())
			if !success || err != nil {
				return as, fmt.Errorf("failed to remove existing antinode: %v", err)
			}
			if DEBUG {
				fmt.Printf("Removed existing antinode with ID %s\n", e.GetId())
			}
		}
	}

	// Gather all antennas by their signals
	antennaMap := make(map[string][]*Antenna)
	for _, e := range as.GetEntities() {
		if ant, ok := e.(*Antenna); ok {
			antennaMap[ant.GetSignal()] = append(antennaMap[ant.GetSignal()], ant)
			if DEBUG {
				fmt.Printf("Grouped antenna ID %s with signal '%s'\n", ant.GetId(), ant.GetSignal())
			}
		}
	}

	if DEBUG {
		antennaCount := 0
		for _, ants := range antennaMap {
			antennaCount += len(ants)
		}
		fmt.Printf("The map contains %d antennas\n", antennaCount)
		fmt.Printf("The map contains %d signals\n", len(antennaMap))
	}

	// For each signal group, examine every pair of antennas
	for signal, ants := range antennaMap {
		if len(ants) < 2 {
			if DEBUG {
				fmt.Printf("Signal '%s' has less than two antennas. Skipping.\n", signal)
			}
			continue
		}

		if DEBUG {
			fmt.Printf("Processing signal '%s' with %d antennas\n", signal, len(ants))
		}

		// Check all unique pairs
		for i := 0; i < len(ants); i++ {
			for j := i + 1; j < len(ants); j++ {
				a1 := ants[i]
				a2 := ants[j]

				x1, y1 := a1.GetPosition()
				x2, y2 := a2.GetPosition()

				dx := x2 - x1
				dy := y2 - y1

				// If they occupy the same spot, ignore
				if dx == 0 && dy == 0 {
					if DEBUG {
						fmt.Printf("Antennas at (%d,%d) and (%d,%d) occupy the same position. Skipping.\n", x1, y1, x2, y2)
					}
					continue
				}

				// Calculate GCD to determine the minimal step
				stepGCD := gcd(abs(dx), abs(dy))
				ux := dx / stepGCD
				uy := dy / stepGCD

				if DEBUG {
					fmt.Printf("Processing pair (%d,%d) and (%d,%d): dx=%d, dy=%d, gcd=%d, ux=%d, uy=%d\n",
						x1, y1, x2, y2, dx, dy, stepGCD, ux, uy)
				}

				// Place an antinode at the location of both antennas
				if DEBUG {
					fmt.Printf("Placing antinode at (%d,%d) and (%d,%d)\n", x1, y1, x2, y2)
				}
				if err := as.placeAntinode(x1, y1); err != nil {
					return as, fmt.Errorf("error placing antinode at (%d,%d): %v", x1, y1, err)
				}
				if err := as.placeAntinode(x2, y2); err != nil {
					return as, fmt.Errorf("error placing antinode at (%d,%d): %v", x2, y2, err)
				}

				// Calculate spacing (distance between antenna and first antinode)
				spacing := stepGCD

				// Place antinodes in both directions using the helper function
				// Starting from a1, stepping backward
				err := as.placeAntinodesInDirection(x1, y1, -ux, -uy, spacing)
				if err != nil {
					return as, err
				}

				// Starting from a2, stepping forward
				err = as.placeAntinodesInDirection(x2, y2, ux, uy, spacing)
				if err != nil {
					return as, err
				}
			}
		}
	}

	return as, nil
}

// placeAntinode places an antinode at the specified (x, y) position.
// It removes any existing antinodes in the cell to avoid duplicates.
func (as *AntennaSimulation) placeAntinode(x, y int) error {
	// Create a new antinode using the new() syntax
	a, err := NewAntinode()
	if err != nil {
		return fmt.Errorf("failed to create new antinode: %v", err)
	}

	// Get the cell at the specified position
	cell, err := as.GetMap().GetCell(x, y)
	if err != nil {
		return fmt.Errorf("failed to get cell at (%d, %d): %v", x, y, err)
	}

	// Retrieve all entity IDs in the cell
	entityIds, err := cell.GetEntityIds()
	if err != nil && err.Error() != "cell is empty" {
		return fmt.Errorf("failed to get entities in cell (%d, %d): %v", x, y, err)
	}

	// Iterate over all entities and remove existing antinodes
	for _, eid := range entityIds {
		entity, err := as.GetEntity(eid)
		if err != nil {
			// Log the error and continue
			fmt.Printf("warning: failed to get entity with ID %s: %v\n", eid, err)
			continue
		}

		if _, ok := entity.(*Antinode); ok {
			success, err := as.RemoveEntity(eid)
			if !success || err != nil {
				// Log the error and continue
				fmt.Printf("warning: failed to remove existing antinode with ID %s: %v\n", eid, err)
			}
		}
	}

	// Add the new antinode to the simulation
	_, err = as.AddEntity(a, x, y)
	if err != nil {
		return fmt.Errorf("failed to add new antinode at (%d, %d): %v", x, y, err)
	}

	return nil
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

type Antinode struct {
	simulation.Entity
}

func NewAntinode() (*Antinode, error) {
	var a = new(Antinode)
	entity, err := simulation.NewEntity()
	if err != nil {
		return nil, err
	}
	a.Entity = entity
	return a, nil
}

func (a *Antinode) String() string {
	// Default character for an antinode is "#"
	return "#"
}

type Antenna struct {
	simulation.Entity
	signal string
}

func NewAntenna(signal string) (*Antenna, error) {
	var a = new(Antenna)
	entity, err := simulation.NewEntity()
	if err != nil {
		return nil, err
	}
	a.Entity = entity
	a.signal = signal
	return a, nil
}

func (a *Antenna) GetSignal() string {
	return a.signal
}

func (a *Antenna) String() string {
	// If the signal is empty, return a default string
	if a.signal == "" {
		return "@"
	}
	// If the signal is longer than 1 character, truncate it to 1 character
	if len(a.signal) > 1 {
		return a.signal[:1]
	}
	// Else return the signal as is
	return a.signal
}
