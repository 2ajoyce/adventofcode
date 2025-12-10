package equation

import (
	"fmt"
	"testing"
)

func TestNewState(t *testing.T) {
	var testCases = []struct {
		name          string
		input         string
		outputState   State
		outputNumBits int
	}{
		{name: "Base", input: "[..##.]", outputState: 0b01100, outputNumBits: 5},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			st, nb := NewState(tc.input)
			if st != tc.outputState {
				t.Fatalf("NewState(%s) state = %016b; want %016b", tc.input, st, tc.outputState)
			}
			if nb != tc.outputNumBits {
				t.Fatalf("NewState(%s) numBits = %d; want %d", tc.input, nb, tc.outputNumBits)
			}

		})
	}
}

func TestStateString(t *testing.T) {
	var testCases = []struct {
		name   string
		input  string
		button string
		output string
	}{
		{name: "Base", input: "[.....]", button: "", output: "[.....]"},
		{name: "Mixed", input: "[..##.]", button: "", output: "[..##.]"},
		{name: "Invert", input: "[..##.]", button: "(0,1,2,3,4)", output: "[##..#]"},
		{name: "Partial Change", input: "[..##.]", button: "(0)", output: "[#.##.]"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			st, nb := NewState(tc.input)
			if tc.button != "" {
				b := NewButton(tc.button)
				st = st.PressButton(b)
			}
			s := st.String(nb)
			if s != tc.output {
				t.Fatalf("State.String produced %s; want %s", s, tc.output)
			}
		})
	}
}

func TestNewButton(t *testing.T) {
	var testCases = []struct {
		name        string
		input       string
		outputState State
	}{
		// Use a 5-bit input where bits 2 and 3 (0-based) are set -> 0b01100
		{name: "NewButton basic", input: "(0, 2, 4)", outputState: 0b10101},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			nb := NewButton(tc.input)
			if nb != Button(tc.outputState) {
				t.Fatalf("NewButton(%s) = %016b; want %016b", tc.input, nb, tc.outputState)
			}
		})
	}
}

func TestNewEquation(t *testing.T) {
	var testCases = []struct {
		name         string
		inputEq      string
		initialState string
		output       string
	}{
		{name: "Base - One Button", inputEq: "[.....] (0) {1,2,3}", initialState: "[.....]", output: "[#....]"},
		{name: "Mixed State - One Button", inputEq: "[..##.] (0) {}", initialState: "[..##.]", output: "[#.##.]"},
		{name: "Invert", inputEq: "[..##.] (0,1,2,3,4) {}", initialState: "[..##.]", output: "[##..#]"},
		{name: "Partial Change", inputEq: "[..##.] (0) {}", initialState: "[..##.]", output: "[#.##.]"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			e := NewEquation(tc.inputEq)

			// Check that the initial state is set
			s := e.Target.String(e.NumBits)
			if s != tc.initialState {
				t.Fatalf("NewEquation(%s) produced target %s; want %s", tc.inputEq, s, tc.initialState)
			}

			// Press each button once and check the state
			for _, b := range e.Buttons {
				e.Target = e.Target.PressButton(b)
			}
			s = e.Target.String(e.NumBits)
			if s != tc.output {
				t.Fatalf("After pressing buttons, state = %s; want %s", s, tc.output)
			}
		})
	}
}
