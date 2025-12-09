package point

import (
	"fmt"
	"math"
	"strconv"
)

type Point struct {
	_id int
	x   int
	y   int
	z   int
}

// NewPoint constructs a new immutable Point.
func NewPoint(x, y, z int) *Point {
	return &Point{x: x, y: y, z: z}
}

// Getters for coordinates
func (p *Point) X() int { return p.x }
func (p *Point) Y() int { return p.y }
func (p *Point) Z() int { return p.z }

// Support unique integer IDs for points up to
// (1,000,000, 1,000,000, 1,000,000)
func (p *Point) Id() int {
	if p._id == 0 {
		p._id = p.x*1000000000 + p.y*1000000 + p.z
	}
	return p._id
}

func (p *Point) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p.x, p.y, p.z)
}

// Calculates the Euclidean distance between two points
func (p1 *Point) Distance(p2 *Point) float64 {
	dx := float64(p1.x - p2.x)
	dy := float64(p1.y - p2.y)
	dz := float64(p1.z - p2.z)

	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// DistanceTo returns an map of distance -> *Point for all points in the input slice
func (p *Point) DistanceTo(points []*Point) map[float64]*Point {
	distances := make(map[float64]*Point)
	for _, other := range points {
		if other.Id() == p.Id() {
			continue
		}
		dist := p.Distance(other)
		distances[dist] = other
	}
	return distances
}

func StrToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to convert string %q to integer: %v", s, err))
	}
	return num
}
