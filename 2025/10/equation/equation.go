package equation

import (
	"fmt"
	"regexp"
	"strings"
)

// Up to 16 cells can fit in a uint16
// The current max is 10 in the input.txt
type State uint16
type Button uint16

// The initial state is all bits 0
type Equation struct {
	Buttons []Button
	Target  State
	NumBits int // Necessary for string conversion
}

// NewEquation builds an Equation from the string input
// s has form [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
//
//	[state] (button1) (button2) ... {batteries}
func NewEquation(s string) Equation {
	// Compilation errors are skipped as these have been hand tested already
	stateRegex, _ := regexp.Compile(`\[(.*?)\]`)
	buttonRegex, _ := regexp.Compile(`\((.*?)\)`)
	// batteryRegex, _ := regexp.Compile(`\{(.*?)\}`) // Batteries are ignored for now

	targetStr := stateRegex.FindString(s)
	buttonStrs := buttonRegex.FindAllString(s, -1)
	// batteryStrs := batteryRegex.FindAllString(s, -1) // Batteries are ignored for now

	buttons := make([]Button, len(buttonStrs))
	for i := range buttonStrs {
		buttons[i] = NewButton(buttonStrs[i])
	}

	state, numBits := NewState(targetStr)
	return Equation{
		Buttons: buttons,
		Target:  state,
		NumBits: numBits,
	}
}

// NewState turns a string like "[.##.###.##]" into a State bitmask
// '.' -> 0, '#' -> 1
// The number of bits is necessary to convert back to string
func NewState(s string) (State, int) {
	var st State
	// Remove the surrounding square brackets
	s = s[1 : len(s)-1]

	// st starts with all bits 0
	// Set bits to 1 where there is a '#'
	for i, ch := range s {
		if ch == '#' {
			st |= 1 << i
		}
	}
	return st, len(s)
}

// stateToString converts a State back into a string of '.' and '#'
// numBits is necessary to determine output length
func (st State) String(numBits int) string {
	b := make([]byte, numBits)
	for i := range numBits {
		if (st>>i)&1 == 1 {
			b[i] = '#'
		} else {
			b[i] = '.'
		}
	}
	return fmt.Sprintf("[%s]", b)
}

// NewButton builds a button maskfrom a string by toggling bits to 1
// Example: "(0, 2, 4)" -> 0b10101
func NewButton(s string) Button {
	var b Button
	// Remove surrounding parentheses
	s = s[1 : len(s)-1]

	// Since we don't know how many indices there are, split by comma first
	var indices []int
	for part := range strings.SplitSeq(s, ",") {
		var idx int
		fmt.Sscanf(strings.TrimSpace(part), "%d", &idx)
		indices = append(indices, idx)
	}

	for i := range indices {
		b |= 1 << indices[i]
	}

	return b
}

// PressButton XORs the button mask with the state and returns the new state
func (st State) PressButton(b Button) State {
	st ^= State(b)
	return st
}
