package main

import (
	"fmt"
	"testing"
)

func TestSolve(t *testing.T) {
	var testCases = []struct {
		name   string
		input  []string
		output string
	}{
		{name: "AOC Input 1", input: []string{
			"L68",
			"L30",
			"R48",
			"L5",
			"R60",
			"L55",
			"L1",
			"L99",
			"R14",
			"L82",
		}, output: "3"},
		{name: "Right Increment larger than DIAL_SIZE", input: []string{
			"R160",
			"L10", // lands on zero
			"R10",
		}, output: "1"},
		{name: "Left Increment larger than DIAL_SIZE", input: []string{
			"L160",
			"R10", // lands on zero
			"L10",
		}, output: "1"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan string)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- line
				}
			}()
			result, err := Solve(inputChan)
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
		{name: "AOC Input 1", input: []string{
			"L68",
			"L30",
			"R48",
			"L5",
			"R60",
			"L55",
			"L1",
			"L99",
			"R14",
			"L82",
		}, output: "6"},
		{name: "Right Increment smaller than DIAL_SIZE", input: []string{
			"R10",
			"L10",
			"R10",
		}, output: "0"},
		{name: "Right Increment larger than DIAL_SIZE", input: []string{
			"R160", // 10: Passes zero twice
			"L10",  // 0: lands on zero
			"R10",  // 10
		}, output: "3"},
		{name: "Left Increment larger than DIAL_SIZE", input: []string{
			"L160", // 90: Passes zero twice
			"R10",  // 0: lands on zero
			"L10",  // 90
		}, output: "3"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			inputChan := make(chan string)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- line
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

func TestRight(t *testing.T) {
	var testCases = []struct {
		name      string
		current   int
		increment int
		output    int
	}{
		{name: "No Spin from 0", current: 0, increment: 0, output: 0},
		{name: "No Spin from 99", current: 99, increment: 0, output: 99},
		{name: "Partial Spin from 0", current: 0, increment: 55, output: 55},
		{name: "Partial Spin from 55", current: 55, increment: 10, output: 65},
		{name: "Full Spin", current: 0, increment: 100, output: 0},
		{name: "Full Spin from 55", current: 55, increment: 100, output: 55},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MoveRight(tc.current, tc.increment)
			if result != tc.output {
				t.Errorf("Expected %d, got %d", tc.output, result)
			}
		})
	}
}
func TestLeft(t *testing.T) {
	var testCases = []struct {
		name      string
		current   int
		increment int
		output    int
	}{
		{name: "No Spin from 0", current: 0, increment: 0, output: 0},
		{name: "No Spin from 99", current: 99, increment: 0, output: 99},
		{name: "Partial Spin from 0", current: 0, increment: 55, output: 45},
		{name: "Partial Spin from 45", current: 45, increment: 10, output: 35},
		{name: "Full Spin", current: 0, increment: 100, output: 0},
		{name: "Full Spin from 55", current: 55, increment: 100, output: 55},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MoveLeft(tc.current, tc.increment)
			if result != tc.output {
				t.Errorf("Expected %d, got %d", tc.output, result)
			}
		})
	}
}
