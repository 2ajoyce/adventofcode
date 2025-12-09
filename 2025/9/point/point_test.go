package point

import (
	"testing"
)

func TestNewPoint(t *testing.T) {
	var testCases = []struct {
		name   string
		x, y   int
		output *Point
	}{
		{name: "Origin", x: 0, y: 0, output: &Point{0, 0}},
		{name: "TopLeft", x: -1, y: 1, output: &Point{-1, 1}},
		{name: "TopRight", x: 1, y: 1, output: &Point{1, 1}},
		{name: "BottomRight", x: 1, y: -1, output: &Point{1, -1}},
		{name: "BottomLeft", x: -1, y: -1, output: &Point{-1, -1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewPoint(tc.x, tc.y)
			if result.X != tc.output.X || result.Y != tc.output.Y {
				t.Errorf("Expected %v, got %v", tc.output, result)
			}
		})
	}
}

func TestString(t *testing.T) {
	var testCases = []struct {
		name   string
		x, y   int
		output string
	}{
		{name: "Origin", x: 0, y: 0, output: "0,0"},
		{name: "TopLeft", x: -1, y: 1, output: "-1,1"},
		{name: "TopRight", x: 1, y: 1, output: "1,1"},
		{name: "BottomRight", x: 1, y: -1, output: "1,-1"},
		{name: "BottomLeft", x: -1, y: -1, output: "-1,-1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPoint(tc.x, tc.y)
			result := p.String()
			if result != tc.output {
				t.Errorf("Expected %v, got %v", tc.output, result)
			}
		})
	}
}

func TestDistanceTo(t *testing.T) {
	var testCases = []struct {
		name   string
		x, y   int
		output float64
	}{
		{name: "Distance zero", x: 0, y: 0, output: 0},
		{name: "Distance 3,4", x: 3, y: 4, output: 7},
		{name: "Distance negatives", x: -1, y: -2, output: 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p1 := NewPoint(0, 0)
			p2 := NewPoint(tc.x, tc.y)
			result := p1.DistanceTo(p2)
			if result != tc.output {
				t.Errorf("Expected %v, got %v", tc.output, result)
			}
		})
	}
}
