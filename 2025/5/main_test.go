package main

import (
	"testing"
)

func TestSolve1(t *testing.T) {
	var testCases = []struct {
		name        string
		inputRanges []Range
		inputInts   []int
		output      string
	}{
		{name: "AOC Example 1",
			inputRanges: []Range{
				{start: 3, end: 5},
				{start: 10, end: 14},
				{start: 16, end: 20},
				{start: 12, end: 18},
			},
			inputInts: []int{
				1,
				5,
				8,
				11,
				17,
				32,
			}, output: "3"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cRange := make(chan Range)
			cInt := make(chan int)
			go func() {
				defer close(cRange)
				for _, r := range tc.inputRanges {
					cRange <- r
				}
			}()
			go func() {
				defer close(cInt)
				for _, i := range tc.inputInts {
					cInt <- i
				}
			}()
			result, err := Solve1(cRange, cInt)
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
		name        string
		inputRanges []Range
		inputInts   []int
		output      string
	}{
		{name: "AOC Example 1",
			inputRanges: []Range{
				{start: 3, end: 5},
				{start: 10, end: 14},
				{start: 16, end: 20},
				{start: 12, end: 18},
			},
			inputInts: []int{
				1,
				5,
				8,
				11,
				17,
				32,
			}, output: "3"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cRange := make(chan Range)
			cInt := make(chan int)
			go func() {
				defer close(cRange)
				for _, r := range tc.inputRanges {
					cRange <- r
				}
			}()
			go func() {
				defer close(cInt)
				for _, i := range tc.inputInts {
					cInt <- i
				}
			}()
			result, err := Solve2(cRange, cInt)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.output {
				t.Errorf("Expected %s, got %s", tc.output, result)
			}
		})
	}
}
