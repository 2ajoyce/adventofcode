package main

import (
	"2ajoyce/adventofcode/2025/9/point"
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
			"7,1",
			"11,1",
			"11,7",
			"9,7",
			"9,5",
			"2,5",
			"2,3",
			"7,3",
		}, output: "50"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan *point.Point)
			go func() {
				defer close(inputChan)
				for _, point := range tc.input {
					inputChan <- ParseInput(point)
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
			"7,1",
			"11,1",
			"11,7",
			"9,7",
			"9,5",
			"2,5",
			"2,3",
			"7,3",
		}, output: "50"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Skip("Implement me!")
			fmt.Println(tc.name)
			inputChan := make(chan *point.Point)
			go func() {
				defer close(inputChan)
				for _, point := range tc.input {
					inputChan <- ParseInput(point)
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

func TestArea(t *testing.T) {
	var testCases = []struct {
		p1, p2 *point.Point
		output int
	}{
		{p1: point.NewPoint(2, 5), p2: point.NewPoint(9, 7), output: 24},
		{p1: point.NewPoint(7, 1), p2: point.NewPoint(11, 7), output: 35},
		{p1: point.NewPoint(7, 3), p2: point.NewPoint(2, 3), output: 6},
		{p1: point.NewPoint(2, 5), p2: point.NewPoint(11, 1), output: 50},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s-%s", tc.p1, tc.p2), func(t *testing.T) {
			result := Area(tc.p1, tc.p2)
			if result != tc.output {
				t.Errorf("Expected %d, got %d", tc.output, result)
			}
		})
	}
}
