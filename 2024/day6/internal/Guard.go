package internal

import (
	"day6/internal/directions"
	"fmt"
)

type Guard interface {
	DirectedObject
	String() string
}

type guard struct {
	*directedObject
}

func NewGuard(location Coord, facing directions.Direction, character rune) (Guard, error) {
	printMap := map[directions.Direction]rune{
		directions.N: '^',
		directions.S: 'v',
		directions.E: '>',
		directions.W: '<',
	}
	d, err := newDirectedObject(location, facing, printMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create guard: %w", err)
	}
	return &guard{directedObject: d}, nil
}
