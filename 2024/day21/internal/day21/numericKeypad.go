package day21

import (
	"fmt"
	"sync"
)

type Coord struct {
	X int
	Y int
}

// The keypad is a 3x4 grid of numbers 0-9, with an A in the bottom right corner
// 7 8 9
// 4 5 6
// 1 2 3
// _ 0 A
type NumericKeypad struct {
	currentX int
	currentY int
	mutex    sync.Mutex
}

// NewNumericKeypad creates a new NumericKeypad
// The starting position is 2, 3 (0-indexed) over the A
func NewNumericKeypad() *NumericKeypad {
	return &NumericKeypad{
		currentX: 2,
		currentY: 3,
	}
}

func (n *NumericKeypad) GetPosition(c rune) Coord {
	var targetX, targetY int
	switch c {
	case 'A':
		targetX = 2
		targetY = 3
	case '0':
		targetX = 1
		targetY = 3
	case '1':
		targetX = 0
		targetY = 2
	case '2':
		targetX = 1
		targetY = 2
	case '3':
		targetX = 2
		targetY = 2
	case '4':
		targetX = 0
		targetY = 1
	case '5':
		targetX = 1
		targetY = 1
	case '6':
		targetX = 2
		targetY = 1
	case '7':
		targetX = 0
		targetY = 0
	case '8':
		targetX = 1
		targetY = 0
	case '9':
		targetX = 2
		targetY = 0
	}
	return Coord{X: targetX, Y: targetY}
}

// CalculateMovements calculates the movements necessary to input a value on a numeric keypad
// The input rune is in the range of [0-9] or 'A'
// The output is a series of movements to input the rune
// The movements are: ^ (up), v (down), < (left), > (right), A (press)
func (n *NumericKeypad) CalculateMovements(input rune) []string {
	// Lock the mutex
	n.mutex.Lock()
	defer n.mutex.Unlock()

	originalX := n.currentX
	originalY := n.currentY

	output := n.calculateMovement(input)

	n.currentX = originalX
	n.currentY = originalY

	permutations := permutateSubstring(output)
	validPermutations := []string{}
	for i := range permutations {
		if n.validateMove(permutations[i]) {
			validPermutations = append(validPermutations, permutations[i])
		}
	}

	return validPermutations
}

// calculateMovement calculates the movements necessary to input a single character on a numeric keypad
// This function MUTATES the currentX and currentY values and should only be called within a mutex lock
// with the currentX and currentY values reset to their original values after use
func (n *NumericKeypad) calculateMovement(c rune) string {
	target := n.GetPosition(c)
	output := ""
	for n.currentX < target.X {
		output += ">"
		n.currentX++
	}
	for n.currentX > target.X {
		output += "<"
		n.currentX--

		// special case to avoid the empty space
		if n.currentX == 0 && n.currentY == 3 {
			// Undo the last step
			output = output[:len(output)-1]
			n.currentX++

			// Move up and over instead
			output += "^"
			n.currentY--
			output += "<"
			n.currentX--
		}
	}
	for n.currentY < target.Y {
		output += "v"
		n.currentY++

		// special case to avoid the empty space
		if n.currentX == 0 && n.currentY == 3 {
			// Undo the last step
			output = output[:len(output)-1]
			n.currentY--

			// Move over and down instead
			output += ">"
			n.currentX++
			output += "v"
			n.currentY++
		}
	}
	for n.currentY > target.Y {
		output += "^"
		n.currentY--
	}
	if len(output) == 0 {
		return output
	}
	output += "A"
	return output
}

func (n *NumericKeypad) GetCurrentPosition() Coord {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return Coord{X: n.currentX, Y: n.currentY}
}

func (n *NumericKeypad) SetCurrentPosition(x, y int) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.currentX = x
	n.currentY = y
}

func (n *NumericKeypad) ResetPosition() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.currentX = 2
	n.currentY = 3
}

func (n *NumericKeypad) Move(input string) bool {
	// Lock the mutex
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for _, c := range input {
		if !n.moveTo(c) {
			return false
		}
	}
	return true
}

func (n *NumericKeypad) moveTo(c rune) bool {
	switch c {
	case '^':
		n.currentY--
	case 'v':
		n.currentY++
	case '<':
		n.currentX--
	case '>':
		n.currentX++
	case 'A':
		return true
	}
	if n.currentX < 0 || n.currentX > 2 || n.currentY < 0 || n.currentY > 3 {
		panic("Invalid position: Moved out of bounds")
	}
	if n.currentX == 0 && n.currentY == 3 {
		panic("Invalid position: Moved to empty space")
	}

	return true
}

func (n *NumericKeypad) validateMove(c string) bool {
	cloneX := n.currentX
	cloneY := n.currentY
	for i, move := range c {
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
			if i != len(c)-1 {
				fmt.Println("'A' must be the final character in the string")
				return false
			}

			return true
		}
		if cloneX < 0 || cloneX > 2 || cloneY < 0 || cloneY > 3 {
			//fmt.Println("Invalid position: Moved out of bounds")
			return false
		}
		if cloneX == 0 && cloneY == 3 {
			//fmt.Println("Invalid position: Moved to empty space")
			return false
		}
	}

	return true
}

// This function will accept a string and return an array of permutations
// For example, the input "^<A" will return ["^<A", "<^A"]
// Each string in the output will be an enumeration of the input string
// The order of the strings in the output array does not matter
func permutateSubstring(input string) []string {
	if len(input) == 0 {
		return []string{}
	}
	if len(input) == 1 {
		if input == "A" {
			return []string{"A"}
		}
	}
	// Remove the A character from the input string
	input = input[:len(input)-1]
	var result []string
	permute([]rune(input), 0, &result)

	// Deduplicate the results
	uniqueResults := make(map[string]struct{})
	for _, r := range result {
		uniqueResults[r] = struct{}{}
	}
	result = []string{}
	for k := range uniqueResults {
		result = append(result, k)
	}

	// Add the A character back to the end of each string
	for i := range result {
		result[i] += "A"
	}
	return result
}

func permute(runes []rune, start int, result *[]string) {
	if start == len(runes)-1 {
		*result = append(*result, string(runes))
		return
	}
	for i := start; i < len(runes); i++ {
		runes[start], runes[i] = runes[i], runes[start]
		permute(runes, start+1, result)
		runes[start], runes[i] = runes[i], runes[start]
	}
}
