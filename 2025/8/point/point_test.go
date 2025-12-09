package point

import (
	"math"
	"testing"
)

func TestId(t *testing.T) {
	p := NewPoint(1, 2, 3)
	recieved := p.Id()
	expected := 1002000003
	if recieved != expected {
		t.Fatalf("Id() = %d, want %d", recieved, expected)
	}

	// Ensure Id() is stable when called multiple times
	if p.Id() != expected {
		t.Fatalf("Id changed across calls; got %d want %d", p.Id(), expected)
	}
}

func TestDistance(t *testing.T) {
	p1 := NewPoint(0, 0, 0)
	p2 := NewPoint(3, 4, 0)
	got := p1.Distance(p2)
	want := 5.0
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("Distance() = %v, want %v", got, want)
	}
}

func TestDistanceTo(t *testing.T) {
	p1 := NewPoint(1, 2, 3)
	p2 := NewPoint(4, 6, 3)

	// Build slice including both points (DistanceTo should ignore self)
	pts := []*Point{p1, p2}
	dmap := p1.DistanceTo(pts)

	// Expect exactly one entry (p2) because p1 should be skipped
	if len(dmap) != 1 {
		t.Fatalf("DistanceTo returned %d entries, want 1", len(dmap))
	}

	want := 5.0
	other, ok := dmap[want]
	if !ok {
		t.Fatalf("expected distance %v not found in map keys: %+v", want, dmap)
	}
	if other != p2 {
		t.Fatalf("DistanceTo map[%v] = %v, want %v", want, other, p2)
	}
}
