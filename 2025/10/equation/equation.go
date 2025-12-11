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

// VoltageState represents the per-index voltages for part 2.
type VoltageState []uint16

// The initial state is all bits 0
type Equation struct {
	Buttons []Button
	// For Part 1
	Target State
	NumBits     int // Necessary for string conversion
	// For Part 2
	TargetVoltage VoltageState
}

// NewEquation builds an Equation from the string input
// s has form [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
//
//	[state] (button1) (button2) ... {batteries}
func NewEquation(s string) Equation {
	// Compilation errors are skipped as these have been hand tested already
	stateRegex, _ := regexp.Compile(`\[(.*?)\]`)
	buttonRegex, _ := regexp.Compile(`\((.*?)\)`)
	batteryRegex, _ := regexp.Compile(`\{(.*?)\}`)

	targetStr := stateRegex.FindString(s)
	buttonStrs := buttonRegex.FindAllString(s, -1)
	batteryStr := batteryRegex.FindString(s)

	buttons := make([]Button, len(buttonStrs))
	for i := range buttonStrs {
		buttons[i] = NewButton(buttonStrs[i])
	}

	state, numBits := NewState(targetStr)
	voltage := NewVoltageState(batteryStr)
	return Equation{
		Buttons:       buttons,
		Target:        state,
		NumBits:       numBits,
		TargetVoltage: voltage,
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

// NewVoltageState turns a string like "{3,5,4,7}" into a VoltageState
func NewVoltageState(s string) VoltageState {
	var st VoltageState
	// Remove the surrounding square brackets or braces
	s = s[1 : len(s)-1]

	// st starts with all bits 0
	// Set bits to the corresponding integer values
	for _, part := range strings.Split(s, ",") {
		var val uint16
		fmt.Sscanf(strings.TrimSpace(part), "%d", &val)
		st = append(st, val)
	}
	return st
}

// String converts a State back into a string of '.' and '#' or int
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

// String converts a VoltageState back into a string of "{v1,v2,...}"
func (st VoltageState) String() string {
	b := make([]string, len(st))
	for i := range st {
		b[i] = fmt.Sprintf("%d", st[i])
	}
	return fmt.Sprintf("{%s}", strings.Join(b, ","))
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

// PressButton increases the voltage state according to the button mask
func (st VoltageState) PressButton(b Button) VoltageState {
	for i := range st {
		if (b>>i)&1 == 1 {
			st[i]++
		}
	}
	return st
}
