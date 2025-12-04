package main

import (
	"testing"
)

func TestSolve1(t *testing.T) {
	var testCases = []struct {
		name   string
		input  []string
		output string
	}{
		{name: "AOC Example 1", input: []string{
			"..@@.@@@@.", // ..xx.xx@x.
			"@@@.@.@.@@", // x@@.@.@.@@
			"@@@@@.@.@@", // @@@@@.x.@@
			"@.@@@@..@.", // @.@@@@..@.
			"@@.@@@@.@@", // x@.@@@@.@x
			".@@@@@@@.@", // .@@@@@@@.@
			".@.@.@.@@@", // .@.@.@.@@@
			"@.@@@.@@@@", // x.@@@.@@@@
			".@@@@@@@@.", // .@@@@@@@@.
			"@.@.@@@.@.", // x.x.@@@.x.
		}, output: "13"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan [][]rune)
			go func() {
				defer close(inputChan)
				input := [][]rune{}
				for _, line := range tc.input {
					input = append(input, ParseInput(line))
				}
				inputChan <- input
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
		input  []string
		output string
	}{
		{name: "AOC Example 1", input: []string{
			"..@@.@@@@.", // ..xx.xx@x.
			"@@@.@.@.@@", // x@@.@.@.@@
			"@@@@@.@.@@", // @@@@@.x.@@
			"@.@@@@..@.", // @.@@@@..@.
			"@@.@@@@.@@", // x@.@@@@.@x
			".@@@@@@@.@", // .@@@@@@@.@
			".@.@.@.@@@", // .@.@.@.@@@
			"@.@@@.@@@@", // x.@@@.@@@@
			".@@@@@@@@.", // .@@@@@@@@.
			"@.@.@@@.@.", // x.x.@@@.x.
		}, output: "13"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan [][]rune)
			go func() {
				defer close(inputChan)
				input := [][]rune{}
				for _, line := range tc.input {
					input = append(input, ParseInput(line))
				}
				inputChan <- input
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
