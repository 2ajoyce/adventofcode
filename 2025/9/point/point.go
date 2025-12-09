package point

import (
	"fmt"
	"math"
)

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y}
}

func (p *Point) String() string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

// The manhattan distance between two points
func (p *Point) DistanceTo(other *Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Abs(float64(dx)) + math.Abs(float64(dy))
}
