package internal

import (
	"day6/internal/directions"
	"fmt"
	"os"
)

type Coord interface {
	X() int
	Y() int
	String() string
	Move(direction directions.Direction) (Coord, error)
	Equals(other Coord) bool
	BetweenX(other1 Coord, other2 Coord) bool
	BetweenY(other1 Coord, other2 Coord) bool
}

type coord struct {
	x int
	y int
}

// Move updates the coordinates of a coord instance based on the specified moveDirection.
// It returns a new Coord with the updated position or an error if an invalid direction is provided.
//
// The function supports moving in the following directions:
// - Directions.N  (North)
// - Directions.S  (South)
// - Directions.E  (East)
// - Directions.W  (West)
// - Directions.NE (Northeast)
// - Directions.NW (Northwest)
// - Directions.SE (Southeast)
// - Directions.SW (Southwest)
//
// Example usage:
// coord := NewCoord(0, 0)
// newCoord, err := coord.Move(Directions.N)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// fmt.Println(newCoord) // Output: (0, -1)
func (c coord) Move(moveDirection directions.Direction) (Coord, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	if DEBUG {
		fmt.Printf("Moving %s from %s\n", moveDirection, c)
	}
	x := c.x
	y := c.y
	switch moveDirection {
	case directions.N:
		y -= 1
	case directions.S:
		y += 1
	case directions.E:
		x += 1
	case directions.W:
		x -= 1
	case directions.NE:
		x += 1
		y -= 1
	case directions.NW:
		x -= 1
		y -= 1
	case directions.SE:
		x += 1
		y += 1
	case directions.SW:
		x -= 1
		y += 1
	default:
		return nil, fmt.Errorf("invalid direction: %s", moveDirection)
	}
	newCoord := NewCoord(x, y)
	if DEBUG {
		fmt.Printf("Moved to %s\n", newCoord)
	}
	return newCoord, nil
}

func (c coord) X() int {
	return c.x
}

func (c coord) Y() int {
	return c.y
}

// String returns a string representation of the coord in the format "(x, y)".
func (c coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.x, c.y)
}

func (c coord) Equals(other Coord) bool {
	otherC, ok := other.(coord)
	if !ok {
		return false
	}
	return c.x == otherC.x && c.y == otherC.y
}

// This function checks if the current coordinate is between two other coordinates on the X axis
// It returns true if it is between them, false otherwise.
func (c coord) BetweenX(other1 Coord, other2 Coord) bool {
	return checkBetween(c.X(), other1.X(), other2.X())
}

// This function checks if the current coordinate is between two other coordinates on the Y axis
// It returns true if it is between them, false otherwise.
func (c coord) BetweenY(other1 Coord, other2 Coord) bool {
	return checkBetween(c.Y(), other1.Y(), other2.Y())
}

func checkBetween(middle int, end1 int, end2 int) bool {
	smaller := 0
	larger := 0
	if end1 < end2 {
		smaller = end1
		larger = end2
	} else {
		smaller = end2
		larger = end1
	}
	if middle >= smaller && middle <= larger {
		return true
	} else {
		return false
	}
}

// NewCoord creates a new Coord with the given x and y values.
func NewCoord(x int, y int) Coord {
	return coord{x, y}
}
