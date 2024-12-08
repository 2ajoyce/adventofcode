package internal

import (
	"day6/internal/directions"
	"fmt"
)

type DirectedObject interface {
	Location() Coord
	String() string
	Move() (DirectedObject, error)
	FacingDirection() directions.Direction
	FacingCoord() (Coord, error)
	TurnLeft() DirectedObject
	TurnRight() DirectedObject
}

type DirectedObjectPrintMap map[directions.Direction]rune

type directedObject struct {
	*object
	facing   directions.Direction
	printMap DirectedObjectPrintMap
}

func (d *directedObject) FacingDirection() directions.Direction {
	return d.facing
}

func (d *directedObject) Move() (DirectedObject, error) {
	newLocation, err := d.coord.Move(d.facing)
	if err != nil {
		return nil, err
	}
	newDirectedObject := *d
	newCoord, ok := newLocation.(coord)
	if !ok {
		return nil, fmt.Errorf("can not move to invalid coordinate")
	}
	newDirectedObject.object.coord = newCoord
	return &newDirectedObject, nil
}

// FacingCoord returns the coordinates that the directedObject is facing based on its current direction.
func (d *directedObject) FacingCoord() (Coord, error) {
	newLocation, err := d.coord.Move(d.facing)
	if err != nil {
		return nil, err
	}
	newCoord, ok := newLocation.(coord)
	if !ok {
		return nil, fmt.Errorf("can not move to invalid coordinate")
	}
	return newCoord, nil
}

func (d *directedObject) TurnLeft() DirectedObject {
	newFacing := d.facing.TurnLeft()
	newDirectedObject, err := newDirectedObject(d.coord, newFacing, d.printMap)
	if err != nil {
		panic(fmt.Errorf("error creating new directed object while turning left: %v", err.Error()))
	}
	return newDirectedObject
}

func (d *directedObject) TurnRight() DirectedObject {
	newFacing := d.facing.TurnRight()
	newDirectedObject, err := newDirectedObject(d.coord, newFacing, d.printMap)
	if err != nil {
		panic(fmt.Errorf("error creating new directed object while turning right: %v", err.Error()))
	}
	return newDirectedObject
}

func (d *directedObject) String() string {
	return string(d.printMap[d.facing])

}

func newDirectedObject(location Coord, facing directions.Direction, printMap DirectedObjectPrintMap) (*directedObject, error) {
	directedObj, err := NewObject(location, printMap[directions.N])
	if err != nil {
		return nil, err
	}
	return &directedObject{object: directedObj.(*object), facing: facing, printMap: printMap}, nil
}
