package internal

import (
	"day6/internal/directions"
	"fmt"
)

type GenericDirectedObject interface {
	DirectedObject
	String() string
}

type genericDirectedObject struct {
	*directedObject
}

func (g *genericDirectedObject) String() string {
	return string(g.character)
}

func NewGenericDirectedObject(location Coord, facing directions.Direction, character rune) (GenericDirectedObject, error) {
	printMap := map[directions.Direction]rune{
		directions.N: character,
		directions.S: character,
		directions.E: character,
		directions.W: character,
	}

	g, err := newDirectedObject(location, facing, printMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create GenericDirectedObject: %w", err)
	}
	return &genericDirectedObject{directedObject: g}, nil
}
