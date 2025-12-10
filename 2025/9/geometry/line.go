package geometry

import "fmt"

type Line struct {
	A *Point
	B *Point
}

// NewLine creates a new line from point a to point b
// Panics if the line is invalid
func NewLine(a, b *Point) *Line {
	line := &Line{A: a, B: b}
	if err := line.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", a.String(), b.String(), err))
	}
	return line
}

// Validate checks that the line is valid
//   - the line is axis-aligned
//   - the points are not the same
func (l1 *Line) Validate() error {
	// Verify A and B are non nil
	if l1.A == nil || l1.B == nil {
		return fmt.Errorf("line points cannot be nil")
	}

	// If not axis aligned
	if l1.A.X != l1.B.X && l1.A.Y != l1.B.Y {
		return fmt.Errorf("line is not axis aligned")
	}

	// If same point
	if l1.A.X == l1.B.X && l1.A.Y == l1.B.Y {
		return fmt.Errorf("line cannot have zero length")
	}

	return nil
}

// Unsafe string method that does not validate the line
func (l1 *Line) string() string {
	return fmt.Sprintf("%s -> %s", l1.A.String(), l1.B.String())
}
func (l1 *Line) String() string {
	if err := l1.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l1.A.String(), l1.B.String(), err))
	}
	return l1.string()
}

// unsafe Length method that does not validate the line
func (l1 *Line) length() int {

	if l1.IsVertical() {
		return intAbs(l1.A.Y-l1.B.Y) + 1
	} else { //horizontal line
		return intAbs(l1.A.X-l1.B.X) + 1
	}
}

// Length returns the length of an axis-aligned line
// This length is inclusive of both endpoints
func (l1 *Line) Length() int {
	if err := l1.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l1.A.String(), l1.B.String(), err))
	}
	return l1.length()
}

// unsafe IsVertical method that does not validate the line
func (l1 *Line) isVertical() bool {
	return l1.A.X == l1.B.X
}

func (l1 *Line) IsVertical() bool {
	if err := l1.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l1.A.String(), l1.B.String(), err))
	}
	return l1.isVertical()
}

// unsafe IsHorizontal method that does not validate the line
func (l1 *Line) isHorizontal() bool {
	return l1.A.Y == l1.B.Y
}

func (l1 *Line) IsHorizontal() bool {
	if err := l1.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l1.A.String(), l1.B.String(), err))
	}
	return l1.isHorizontal()
}

// unsafe areCollinear method that does not validate the line
func (l1 *Line) areCollinear(l2 *Line) bool {
	// If not the same orientation, cannot be colinear
	if l1.isVertical() != l2.isVertical() {
		return false
	}

	if l1.isVertical() {
		return l1.A.X == l2.A.X
	} else { // horizontal
		return l1.A.Y == l2.A.Y
	}
}

func (l1 *Line) AreCollinear(l2 *Line) bool {
	if err := l1.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l1.A.String(), l1.B.String(), err))
	}
	if err := l2.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l2.A.String(), l2.B.String(), err))
	}

	return l1.areCollinear(l2)
}

// unsafe doesOverlap method that does not validate the line
func (l1 *Line) doesOverlap(l2 *Line) bool {
	var lStart, lEnd, oStart, oEnd int

	// Must be colinear to overlap
	if !l1.areCollinear(l2) {
		return false
	}

	if l1.isVertical() {
		lStart, lEnd = l1.A.Y, l1.B.Y
		oStart, oEnd = l2.A.Y, l2.B.Y
	} else {
		lStart, lEnd = l1.A.X, l1.B.X
		oStart, oEnd = l2.A.X, l2.B.X
	}

	// Normalize to [start <= end]
	if lStart > lEnd {
		lStart, lEnd = lEnd, lStart
	}
	if oStart > oEnd {
		oStart, oEnd = oEnd, oStart
	}

	// Make sure lStart is the smaller start
	if lStart > oStart {
		lStart, lEnd, oStart, oEnd = oStart, oEnd, lStart, lEnd
	}

	// oStart < lEnd => interior overlap (no shared-endpoint-only)
	return oStart < lEnd
}

// DoesOverlap determines if two axis-aligned in the same orientation overlap
// Shared start or end points are NOT considered overlapping
func (l1 *Line) DoesOverlap(l2 *Line) bool {
	if err := l1.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l1.A.String(), l1.B.String(), err))
	}
	if err := l2.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l2.A.String(), l2.B.String(), err))
	}
	return l1.doesOverlap(l2)
}

// unsafe doesCross method that does not validate the line
func (l1 *Line) doesCross(l2 *Line) bool {
	// if l1.isVertical() && l2.isVertical() {
	// 	return l1.doesOverlap(l2) // Overlaps count as a cross
	// }

	// if l1.isHorizontal() && l2.isHorizontal() {
	// 	return l1.doesOverlap(l2) // Overlaps count as a cross
	// }

	// Ensure l1 is vertical and l2 is horizontal
	if l1.isHorizontal() && l2.isVertical() {
		l1, l2 = l2, l1
	}

	// Fixed endpoints
	x := l1.A.X
	y := l2.A.Y

	// X-range of horizontal line (normalize)
	hx1, hx2 := l2.A.X, l2.B.X
	if hx1 > hx2 {
		hx1, hx2 = hx2, hx1
	}

	// Y-range of vertical line (normalize)
	vy1, vy2 := l1.A.Y, l1.B.Y
	if vy1 > vy2 {
		vy1, vy2 = vy2, vy1
	}

	// Strict inequalities: endpoint-only touching does NOT count
	return hx1 < x && x < hx2 &&
		vy1 < y && y < vy2
}

// DoesCross determines if two axis-aligned lines cross each other
// A start point or end point touching is not considered crossing
func (l1 *Line) DoesCross(l2 *Line) bool {
	if err := l1.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l1.A.String(), l1.B.String(), err))
	}
	if err := l2.Validate(); err != nil {
		panic(fmt.Sprintf("invalid line from %s to %s: %v", l2.A.String(), l2.B.String(), err))
	}
	return l1.doesCross(l2)
}

func intAbs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// Unsafe versions of methods that do not validate the line
// Added to enable higher performance
// Todo: Come back and clean up the methods in this file to reduce duplication
func (l1 *Line) IsVerticalUnsafe() bool   { return l1.A.X == l1.B.X }
func (l1 *Line) IsHorizontalUnsafe() bool { return l1.A.Y == l1.B.Y }
func (l1 *Line) DoesCrossUnsafe(l2 *Line) bool {
	return l1.doesCross(l2)
}
func (l1 *Line) DoesOverlapUnsafe(l2 *Line) bool {
	return l1.doesOverlap(l2)
}
