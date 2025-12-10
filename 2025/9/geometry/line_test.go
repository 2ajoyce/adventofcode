package geometry

import (
	"testing"
)

// ai generated test file - only reviewed at a surface level
// This wasn't the point of the exercise, so minimal effort has
// been spent in this file.

func TestNewLine_Valid(t *testing.T) {
	a := NewPoint(0, 0)
	b := NewPoint(0, 2) // vertical, increasing
	l := NewLine(a, b)

	if l.A != a || l.B != b {
		t.Errorf("Expected A=%v B=%v, got A=%v B=%v", a, b, l.A, l.B)
	}
}

func TestNewLine_PanicsOnInvalid(t *testing.T) {
	tests := []struct {
		name string
		a, b *Point
	}{
		{name: "Not axis aligned", a: NewPoint(0, 0), b: NewPoint(1, 1)},
		{name: "Zero length", a: NewPoint(2, 2), b: NewPoint(2, 2)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for invalid line (%s), but none occurred", tc.name)
				}
			}()

			_ = NewLine(tc.a, tc.b)
		})
	}
}

func TestLineString(t *testing.T) {
	a := NewPoint(-1, 2)
	b := NewPoint(3, 2) // horizontal, increasing
	l := NewLine(a, b)

	want := "-1,2 -> 3,2"
	if got := l.String(); got != want {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestLength(t *testing.T) {
	var tests = []struct {
		name string
		a, b *Point
		want int
	}{
		{name: "Vertical up", a: NewPoint(0, 0), b: NewPoint(0, 2), want: 3},
		{name: "Horizontal right", a: NewPoint(1, 1), b: NewPoint(3, 1), want: 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLine(tc.a, tc.b)
			if got := l.Length(); got != tc.want {
				t.Errorf("%s: expected %d, got %d", tc.name, tc.want, got)
			}
		})
	}
}

func TestLengthPanicsForInvalidLine(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid line in Length, but none occurred")
		}
	}()

	// Construct an invalid line directly and call the public Length method.
	l := &Line{A: NewPoint(0, 0), B: NewPoint(1, 1)} // not axis aligned
	_ = l.Length()
}

func TestIsVerticalAndIsHorizontal(t *testing.T) {
	v := NewLine(NewPoint(0, 0), NewPoint(0, 5))
	if !v.IsVertical() || v.IsHorizontal() {
		t.Errorf("Expected vertical line to be vertical and not horizontal")
	}

	h := NewLine(NewPoint(1, 2), NewPoint(4, 2))
	if !h.IsHorizontal() || h.IsVertical() {
		t.Errorf("Expected horizontal line to be horizontal and not vertical")
	}
}

func TestAreCollinear(t *testing.T) {
	tests := []struct {
		name   string
		a1, b1 *Point
		a2, b2 *Point
		want   bool
	}{
		{
			name: "Vertical same X",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5),
			a2: NewPoint(0, 2), b2: NewPoint(0, 8),
			want: true,
		},
		{
			name: "Vertical different X",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5),
			a2: NewPoint(1, 2), b2: NewPoint(1, 8),
			want: false,
		},
		{
			name: "Horizontal same Y",
			a1:   NewPoint(0, 3), b1: NewPoint(5, 3),
			a2: NewPoint(2, 3), b2: NewPoint(8, 3),
			want: true,
		},
		{
			name: "Horizontal different Y",
			a1:   NewPoint(0, 3), b1: NewPoint(5, 3),
			a2: NewPoint(2, 4), b2: NewPoint(8, 4),
			want: false,
		},
		{
			name: "Perpendicular not collinear",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5),
			a2: NewPoint(0, 0), b2: NewPoint(5, 0),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l1 := NewLine(tc.a1, tc.b1)
			l2 := NewLine(tc.a2, tc.b2)
			if got := l1.AreCollinear(l2); got != tc.want {
				t.Errorf("AreCollinear(%s): expected %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestDoesOverlap(t *testing.T) {
	tests := []struct {
		name   string
		a1, b1 *Point
		a2, b2 *Point
		want   bool
	}{
		{
			name: "Vertical interior overlap",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 10),
			a2: NewPoint(0, 3), b2: NewPoint(0, 7),
			want: true,
		},
		{
			name: "Vertical endpoint touch only",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5),
			a2: NewPoint(0, 5), b2: NewPoint(0, 10),
			want: false,
		},
		{
			name: "Vertical disjoint",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5),
			a2: NewPoint(0, 6), b2: NewPoint(0, 10),
			want: false,
		},
		{
			name: "Parallel non-collinear",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5),
			a2: NewPoint(1, 2), b2: NewPoint(1, 4),
			want: false,
		},
		{
			name: "Horizontal interior overlap",
			a1:   NewPoint(0, 2), b1: NewPoint(10, 2),
			a2: NewPoint(3, 2), b2: NewPoint(7, 2),
			want: true,
		},
		{
			name: "Horizontal endpoint touch only",
			a1:   NewPoint(0, 2), b1: NewPoint(5, 2),
			a2: NewPoint(5, 2), b2: NewPoint(10, 2),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l1 := NewLine(tc.a1, tc.b1)
			l2 := NewLine(tc.a2, tc.b2)
			if got := l1.DoesOverlap(l2); got != tc.want {
				t.Errorf("DoesOverlap(%s): expected %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestDoesCross(t *testing.T) {
	tests := []struct {
		name   string
		a1, b1 *Point
		a2, b2 *Point
		want   bool
	}{
		{
			name: "Perpendicular proper cross",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5), // vertical
			a2: NewPoint(-2, 2), b2: NewPoint(2, 2), // horizontal
			want: true, // cross at (0,2)
		},
		{
			name: "Perpendicular endpoint touch only",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 2), // vertical
			a2: NewPoint(0, 2), b2: NewPoint(3, 2), // horizontal, touches at (0,2)
			want: false,
		},
		{
			name: "Perpendicular no intersection",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 1), // vertical
			a2: NewPoint(-2, 2), b2: NewPoint(2, 2), // horizontal above
			want: false,
		},
		{
			name: "Fully contained overlapping vertical lines do not count as cross",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 10),
			a2: NewPoint(0, 3), b2: NewPoint(0, 7),
			want: false,
		},
		{
			name: "Duplicate vertical lines do not count as cross",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 10),
			a2: NewPoint(0, 0), b2: NewPoint(0, 10),
			want: false,
		},
		{
			name: "Parallel vertical, disjoint",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5),
			a2: NewPoint(0, 6), b2: NewPoint(0, 10),
			want: false,
		},
		{
			name: "Parallel vertical, endpoint touch only",
			a1:   NewPoint(0, 0), b1: NewPoint(0, 5),
			a2: NewPoint(0, 5), b2: NewPoint(0, 10),
			want: false,
		},
		{
			name: "Fully contained overlapping horizontal lines do not count as cross",
			a1:   NewPoint(0, 2), b1: NewPoint(10, 2),
			a2: NewPoint(3, 2), b2: NewPoint(7, 2),
			want: false,
		},
		{
			name: "Duplicate horizontal lines do not count as cross",
			a1:   NewPoint(0, 2), b1: NewPoint(10, 2),
			a2: NewPoint(0, 2), b2: NewPoint(10, 2),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l1 := NewLine(tc.a1, tc.b1)
			l2 := NewLine(tc.a2, tc.b2)
			if got := l1.DoesCross(l2); got != tc.want {
				t.Errorf("DoesCross(%s): expected %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}
