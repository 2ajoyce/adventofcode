package main

import (
	"testing"
)

func TestSolve1(t *testing.T) {
	var testCases = []struct {
		name   string
		input  string
		output string
	}{
		{name: "AOC Example 1", input: "1.txt", output: "2"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			problem := ReadInput("test/" + tc.input)
			result, err := Solve1(problem)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.output {
				t.Errorf("Expected %s, got %s", tc.output, result)
			}
		})
	}
}
