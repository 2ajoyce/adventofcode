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
		// Use a 5-bit input where bits 0,2,4 (0-based) are set -> 0b10101
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

func TestNewVoltageState(t *testing.T) {
	var testCases = []struct {
		name   string
		input  string
		output VoltageState
	}{
		{name: "Simple", input: "{1,2,3}", output: VoltageState{1, 2, 3}},
		{name: "Larger Values", input: "{217,234,214,41,203,41,236,197,221}", output: VoltageState{217, 234, 214, 41, 203, 41, 236, 197, 221}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			st := NewVoltageState(tc.input)
			if len(st) != len(tc.output) {
				t.Fatalf("NewVoltageState(%s) length = %d; want %d", tc.input, len(st), len(tc.output))
			}
			for i := range st {
				if st[i] != tc.output[i] {
					t.Fatalf("NewVoltageState(%s)[%d] = %d; want %d", tc.input, i, st[i], tc.output[i])
				}
			}
		})
	}
}

func TestVoltageStateString(t *testing.T) {
	var testCases = []struct {
		name   string
		input  VoltageState
		output string
	}{
		{name: "All Zero", input: VoltageState{0, 0, 0, 0, 0}, output: "{0,0,0,0,0}"},
		{name: "Mixed", input: VoltageState{1, 2, 3}, output: "{1,2,3}"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			s := tc.input.String()
			if s != tc.output {
				t.Fatalf("VoltageState.String produced %s; want %s", s, tc.output)
			}
		})
	}
}

func TestVoltageStatePressButton(t *testing.T) {
	var testCases = []struct {
		name        string
		initialVolt string
		button      string
		output      string
	}{
		{name: "Part 2 Base One Button", initialVolt: "{0,0,0,0,0}", button: "(0)", output: "{1,0,0,0,0}"},
		{name: "Multiple Indices", initialVolt: "{0,0,0,0,0}", button: "(0,2,4)", output: "{1,0,1,0,1}"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)

			// Parse initial voltage and button
			v := NewVoltageState(tc.initialVolt)
			b := NewButton(tc.button)

			// Apply button once
			v = v.PressButton(b)

			s := v.String()
			if s != tc.output {
				t.Fatalf("VoltageState.PressButton produced %s; want %s", s, tc.output)
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
		voltage      string // optional: only checked when non-empty
	}{
		{name: "Base - One Button", inputEq: "[.....] (0) {1,2,3,4,5}", initialState: "[.....]", output: "[#....]", voltage: "{1,2,3,4,5}"},
		{name: "Mixed State - One Button", inputEq: "[..##.] (0) {1,2,3,4,5}", initialState: "[..##.]", output: "[#.##.]", voltage: "{1,2,3,4,5}"},
		{name: "Invert", inputEq: "[..##.] (0,1,2,3,4) {1,2,3,4,5}", initialState: "[..##.]", output: "[##..#]", voltage: "{1,2,3,4,5}"},
		{name: "Partial Change", inputEq: "[..##.] (0) {1,2,3,4,5}", initialState: "[..##.]", output: "[#.##.]", voltage: "{1,2,3,4,5}"},
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

			// If we have an expected voltage string, check VoltageTarget too
			if tc.voltage != "" {
				vs := e.TargetVoltage.String()
				if vs != tc.voltage {
					t.Fatalf("NewEquation(%s) produced voltage %s; want %s", tc.inputEq, vs, tc.voltage)
				}
			}

			// Press each button once and check the state (Part 1 behavior)
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
