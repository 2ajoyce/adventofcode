package day21

import (
	"fmt"
	"sync"
)

// 7 8 9
// 4 5 6
// 1 2 3
// _ 0 A

// The keypad is a 3x2 grid with the following layout:
// _ ^ A
// < v >

//            3                          7          9                 A
//        ^   A       ^^        <<       A     >>   A        vvv      A
//    <   A > A   <   AA  v <   AA >>  ^ A  v  AA ^ A  v <   AAA >  ^ A
// v<<A>>^AvA^Av<<A>>^AAv<A<A>>^AAvAA^<A>Av<A>^AA<A>Av<A<A>>^AAAvA^<A>A
//            3                      7          9                 A
//        ^   A         <<      ^^   A     >>   A        vvv      A
//    <   A > A  v <<   AA >  ^ AA > A  v  AA ^ A   < v  AAA >  ^ A
// <v<A>>^AvA^A<vA<AA>>^AAvA<^A>AAvA^A<vA>^AA<A>A<v<A>A>^AAAvA<^A>A

type DirectionalKeypad struct {
	currentX int
	currentY int
	mutex    sync.Mutex
}

// newDirectionalKeypad creates a new DirectionalKeypad
// The starting position is 2, 0 (0-indexed) over the A
func NewDirectionalKeypad() *DirectionalKeypad {
	return &DirectionalKeypad{
		currentX: 2,
		currentY: 0,
	}
}

func (d *DirectionalKeypad) GetPosition(c rune) Coord {
	var targetX, targetY int
	switch c {
	case 'A':
		targetX = 2
		targetY = 0
	case '^':
		targetX = 1
		targetY = 0
	case '<':
		targetX = 0
		targetY = 1
	case 'v':
		targetX = 1
		targetY = 1
	case '>':
		targetX = 2
		targetY = 1
	}
	return Coord{X: targetX, Y: targetY}
}

// CalculateMovements calculates the movements necessary to input a string on a directional keypad
// The input rune is in the range [<^v>] and A
// The output is a series of movements to input the string
// The movements are: ^ (up), v (down), < (left), > (right), A (press)
func (d *DirectionalKeypad) CalculateMovements(input rune) []string {
	// Lock the mutex
	d.mutex.Lock()
	defer d.mutex.Unlock()

	originalX := d.currentX
	originalY := d.currentY

	output := d.calculateMovement(input)

	d.currentX = originalX
	d.currentY = originalY

	if len(output) == 0 {
		return []string{"A"}
	}

	permutations := permutateSubstring(output)
	validPermutations := []string{}
	shortestPermutationLength := len(permutations[0])
	for i := range permutations {
		if d.validateMove(permutations[i]) {
			validPermutations = append(validPermutations, permutations[i])
		}
		if len(permutations[i]) < shortestPermutationLength {
			shortestPermutationLength = len(permutations[i])
		}
	}
	shortestPermutations := []string{}
	for i := range validPermutations {
		if len(validPermutations[i]) == shortestPermutationLength {
			shortestPermutations = append(shortestPermutations, validPermutations[i])
		}
	}
	return shortestPermutations
}

// calculateMovement calculates the movements to move from the current position to the target position
// This function MUTATES the currentX and currentY values and should be called within a mutex lock
// where the currentX and currentY values are reset to their original values after the function is called
func (d *DirectionalKeypad) calculateMovement(c rune) string {
	target := d.GetPosition(c)
	output := ""
	for d.currentX < target.X {
		output += ">"
		d.currentX++
	}
	for d.currentX > target.X {
		output += "<"
		d.currentX--

		// special case to avoid the empty space
		if d.currentX == 0 && d.currentY == 0 {
			// Undo the last step
			output = output[:len(output)-1]
			d.currentX++

			// Move down and over instead
			output += "v"
			d.currentY++
			output += "<"
			d.currentX--
		}
	}
	for d.currentY < target.Y {
		output += "v"
		d.currentY++
	}
	for d.currentY > target.Y {
		output += "^"
		d.currentY--

		// special case to avoid the empty space
		if d.currentX == 0 && d.currentY == 0 {
			// Undo the last step
			output = output[:len(output)-1]
			d.currentY++

			// Move over and down instead
			output += ">"
			d.currentX++
			output += "^"
			d.currentY--
		}
	}
	if len(output) == 0 {
		return output
	}
	output += "A"
	return output
}

func (d *DirectionalKeypad) GetCurrentPosition() Coord {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return Coord{X: d.currentX, Y: d.currentY}
}

func (d *DirectionalKeypad) SetCurrentPosition(x, y int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.currentX = x
	d.currentY = y
}

func (d *DirectionalKeypad) ResetPosition() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.currentX = 2
	d.currentY = 0
}

func (d *DirectionalKeypad) Move(input string) bool {
	// Lock the mutex
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, c := range input {
		if !d.moveTo(c) {
			return false
		}
	}
	return true
}

func (d *DirectionalKeypad) moveTo(c rune) bool {
	switch c {
	case '^':
		d.currentY--
	case 'v':
		d.currentY++
	case '<':
		d.currentX--
	case '>':
		d.currentX++
	case 'A':
		return true
	}
	if d.currentX < 0 || d.currentX > 3 || d.currentY < 0 || d.currentY > 2 {
		msg := fmt.Sprintf("Invalid position: Moved out of bounds: %d, %d", d.currentX, d.currentY)
		panic(msg)
	}
	if d.currentX == 0 && d.currentY == 0 {
		panic("Invalid position: Moved to empty space")
	}

	return true
}

func (d *DirectionalKeypad) validateMove(m string) bool {
	cloneX := d.currentX
	cloneY := d.currentY
	for i, move := range m {
		switch move {
		case '^':
			cloneY--
		case 'v':
			cloneY++
		case '<':
			cloneX--
		case '>':
			cloneX++
		case 'A':
			// check that this is the final move
			if i != len(m)-1 {
				fmt.Println("'A' must be the final character in the string")
				return false
			}
			return true
		}
		if cloneX < 0 || cloneX > 3 || cloneY < 0 || cloneY > 2 {
			//fmt.Println("Invalid position: Moved out of bounds")
			return false
		}
		if cloneX == 0 && cloneY == 0 {
			//fmt.Println("Invalid position: Moved to empty space")
			return false
		}
	}

	return true
}
