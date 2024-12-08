package internal

import (
	"fmt"
)

type Object interface {
	Location() Coord
	String() string
}

// Base objects are directionless
type object struct {
	coord     coord
	character rune
}

func (o *object) Location() Coord {
	return o.coord
}

func (o *object) String() string {
	return string(o.character)
}

func NewObject(location Coord, character rune) (Object, error) {
	newLocation, ok := location.(coord)
	if !ok {
		return nil, fmt.Errorf("invalid coordinate")
	}
	return &object{newLocation, character}, nil
}
