package main

import (
	"2ajoyce/adventofcode/2025/10/equation"
	"fmt"
	"testing"
)

func TestSolve1(t *testing.T) {
	var testCases = []struct {
		name   string
		input  []string
		output string
	}{
		{name: "AOC Example 1", input: []string{
			"[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",                 // 2
			"[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",     // 3
			"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}", // 2
		}, output: "7"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan *equation.Equation)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- ParseInput(line)
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
		input  []string
		output string
	}{
		{name: "AOC Example 1", input: []string{
			"[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",                 // 2
			"[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",     // 3
			"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}", // 2
		}, output: "7"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Skip("Implement me!")
			fmt.Println(tc.name)
			inputChan := make(chan *equation.Equation)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- ParseInput(line)
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
