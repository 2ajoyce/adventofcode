package main

import (
	"fmt"
	"testing"
)

// I want to use strings in the tests matrix, not arrays of runes
type TestSpan struct {
	start string
	end   string
}

func (ts *TestSpan) Start() []rune {
	var result []rune
	for _, r := range ts.start {
		result = append(result, r)
	}
	return result
}

func (ts *TestSpan) End() []rune {
	var result []rune
	for _, r := range ts.end {
		result = append(result, r)
	}
	return result
}

func TestSolve1(t *testing.T) {
	var testCases = []struct {
		name   string
		input  []*TestSpan
		output string
	}{
		{name: "AOC Example 1", input: []*TestSpan{
			{start: "11", end: "22"},
			{start: "95", end: "115"},
			{start: "998", end: "1012"},
			{start: "1188511880", end: "1188511890"},
			{start: "222220", end: "222224"},
			{start: "1698522", end: "1698528"},
			{start: "446443", end: "446449"},
			{start: "38593856", end: "38593862"},
			{start: "565653", end: "565659"},
			{start: "824824821", end: "824824827"},
			{start: "2121212118", end: "2121212124"},
		}, output: "1227775554"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan *Span)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- &Span{start: line.Start(), end: line.End()}
				}
			}()
			result, err := Solve1(inputChan)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.output {
				t.Errorf("Expected %s, got %s", tc.output, result)
			}
		})
	}
}

func TestSolve2(t *testing.T) {
	var testCases = []struct {
		name   string
		input  []*TestSpan
		output string
	}{
		{name: "AOC Example 1", input: []*TestSpan{
			{start: "11", end: "22"},
			{start: "95", end: "115"},
			{start: "998", end: "1012"},
			{start: "1188511880", end: "1188511890"},
			{start: "222220", end: "222224"},
			{start: "1698522", end: "1698528"},
			{start: "446443", end: "446449"},
			{start: "38593856", end: "38593862"},
			{start: "565653", end: "565659"},
			{start: "824824821", end: "824824827"},
			{start: "2121212118", end: "2121212124"},
		}, output: "4174379265"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan *Span)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- &Span{start: line.Start(), end: line.End()}
				}
			}()
			result, err := Solve2(inputChan)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.output {
				t.Errorf("Expected %s, got %s", tc.output, result)
			}
		})
	}
}

func TestCheckSpan(t *testing.T) {
	var testCases = []struct {
		name   string
		input  *TestSpan
		output []string
	}{
		{name: "", input: &TestSpan{start: "11", end: "22"}, output: []string{"11", "22"}},
		{name: "", input: &TestSpan{start: "95", end: "115"}, output: []string{"99"}},
		{name: "", input: &TestSpan{start: "998", end: "1012"}, output: []string{"1010"}},
		{name: "", input: &TestSpan{start: "1188511880", end: "1188511890"}, output: []string{"1188511885"}},
		{name: "", input: &TestSpan{start: "222220", end: "222224"}, output: []string{"222222"}},
		{name: "", input: &TestSpan{start: "1698522", end: "1698528"}, output: []string{}},
		{name: "", input: &TestSpan{start: "446443", end: "446449"}, output: []string{"446446"}},
		{name: "", input: &TestSpan{start: "38593856", end: "38593862"}, output: []string{"38593859"}},
		{name: "", input: &TestSpan{start: "565653", end: "565659"}, output: []string{}},
		{name: "", input: &TestSpan{start: "824824821", end: "824824827"}, output: []string{}},
		{name: "", input: &TestSpan{start: "2121212118", end: "2121212124"}, output: []string{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckSpan(&Span{start: tc.input.Start(), end: tc.input.End()})
			if len(result) != len(tc.output) {
				t.Errorf("For %s | %s expected %v, got %q", tc.input.start, tc.input.end, tc.output, result)
			}
			for i, double := range result { // Assume outputs are sorted
				fmt.Printf("Result: %d | %d\n", StrToInt(tc.output[i]), ArrRuneToInt(double))
				if StrToInt(tc.output[i]) != ArrRuneToInt(double) {
					t.Errorf("For %s | %s expected %v, got %v", tc.input.start, tc.input.end, tc.output, result)
				}
			}
		})
	}
}
func TestCheckSpan2(t *testing.T) {
	var testCases = []struct {
		name   string
		input  *TestSpan
		output []string
	}{
		{name: "", input: &TestSpan{start: "11", end: "22"}, output: []string{"11", "22"}},
		{name: "", input: &TestSpan{start: "95", end: "115"}, output: []string{"99", "111"}},
		{name: "", input: &TestSpan{start: "998", end: "1012"}, output: []string{"999", "1010"}},
		{name: "", input: &TestSpan{start: "1188511880", end: "1188511890"}, output: []string{"1188511885"}},
		{name: "", input: &TestSpan{start: "222220", end: "222224"}, output: []string{"222222"}},
		{name: "", input: &TestSpan{start: "1698522", end: "1698528"}, output: []string{}},
		{name: "", input: &TestSpan{start: "446443", end: "446449"}, output: []string{"446446"}},
		{name: "", input: &TestSpan{start: "38593856", end: "38593862"}, output: []string{"38593859"}},
		{name: "", input: &TestSpan{start: "565653", end: "565659"}, output: []string{"565656"}},
		{name: "", input: &TestSpan{start: "824824821", end: "824824827"}, output: []string{"824824824"}},
		{name: "", input: &TestSpan{start: "2121212118", end: "2121212124"}, output: []string{"2121212121"}},
		{name: "", input: &TestSpan{start: "3", end: "16"}, output: []string{"11"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckSpan2(&Span{start: tc.input.Start(), end: tc.input.End()})
			if len(result) != len(tc.output) {
				t.Errorf("For %s | %s expected %v, got %q", tc.input.start, tc.input.end, tc.output, result)
			}
			for i, double := range result { // Assume outputs are sorted
				if StrToInt(tc.output[i]) != ArrRuneToInt(double) {
					t.Errorf("For %s | %s expected %v, got %q", tc.input.start, tc.input.end, tc.output, result)
				}
			}
		})
	}
}

func TestStripPadding(t *testing.T) {
	var testCases = []struct {
		name   string
		input  string
		output string
	}{
		{name: "Zero", input: "0", output: "0"},
		{name: "One", input: "1", output: "1"},
		{name: "Zero Padded Odd (1)", input: "0101", output: "101"},
		{name: "Zero Padded Odd (2)", input: "00101", output: "101"},
		{name: "Zero Padded Even (1)", input: "01011", output: "1011"},
		{name: "Zero Padded Even (2)", input: "001011", output: "1011"},
		{name: "Zero End Odd (1)", input: "10110", output: "10110"},
		{name: "Zero End Odd (2)", input: "10100", output: "10100"},
		{name: "Zero End Even (1)", input: "1010", output: "1010"},
		{name: "Zero End Even (2)", input: "101100", output: "101100"},
		{name: "Zero Both Sides Odd (1)", input: "0110", output: "110"},
		{name: "Zero Both Sides Odd (2)", input: "00100", output: "100"},
		{name: "Zero Both Sides Even (1)", input: "010", output: "10"},
		{name: "Zero Both Sides Even (2)", input: "001100", output: "1100"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StripPadding(tc.input)
			// This check compares the len of []rune to len of string
			// That comparision is not long term durable, but works for this tes
			if len(result) != len(tc.output) {
				t.Errorf("Expected %s, got %c", tc.output, result)
			}
			for i, r := range tc.output {
				if result[i] != r {
					t.Errorf("Expected %s, got %c", tc.output, result)
				}
			}
		})
	}
}

func TestIsInvalidId(t *testing.T) {
	var testCases = []struct {
		name   string
		input  string
		output bool
	}{
		{name: "", input: "0", output: false},
		{name: "", input: "00", output: false},
		{name: "", input: "1", output: false},
		{name: "", input: "12", output: false},
		{name: "", input: "12341234", output: true},
		{name: "", input: "123123123", output: true},
		{name: "", input: "1212121212", output: true},
		{name: "", input: "1111111", output: true},
		{name: "", input: "11", output: true},
		{name: "", input: "22", output: true},
		{name: "", input: "99", output: true},
		{name: "", input: "111", output: true},
		{name: "", input: "999", output: true},
		{name: "", input: "1010", output: true},
		{name: "", input: "1188511885", output: true},
		{name: "", input: "446446", output: true},
		{name: "", input: "38593859", output: true},
		{name: "", input: "565656", output: true},
		{name: "", input: "824824824", output: true},
		{name: "", input: "2121212121", output: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsInvalidId(StrToArrRune(tc.input))
			if result != tc.output {
				t.Errorf("Expected %t, got %t", tc.output, result)
			}
		})
	}
}
