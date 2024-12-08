package directions

import "fmt"

type Direction string

// Conventionally in English the terms Westnorth(WN), Eastnorth(EN), Westsouth(WS), and Eastsouth(ES) are not used.
// Instead, the terms Northwest(NW), Northeast(NE), Southeast(SE), and Southwest(SW) are used.
// That convention has been maintained here for now.
const (
	N  Direction = "N"
	S  Direction = "S"
	E  Direction = "E"
	W  Direction = "W"
	NE Direction = "NE"
	NW Direction = "NW"
	SE Direction = "SE"
	SW Direction = "SW"
)

func (d Direction) TurnRight() Direction {
	switch d {
	case N:
		return E
	case S:
		return W
	case E:
		return S
	case W:
		return N
	default:
		fmt.Println("Invalid direction")
	}
	return d
}

func (d Direction) TurnLeft() Direction {
	switch d {
	case N:
		return W
	case S:
		return E
	case E:
		return N
	case W:
		return S
	default:
		fmt.Println("Invalid direction")
	}
	return d
}
