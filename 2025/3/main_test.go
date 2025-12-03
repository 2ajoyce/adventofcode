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
			"987654321111111",
			"811111111111119",
			"234234234234278",
			"818181911112111",
		}, output: "357"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan []int)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- ParseLine(line)
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
			"987654321111111",
			"811111111111119",
			"234234234234278",
			"818181911112111",
		}, output: "357"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan []int)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- ParseLine(line)
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

func TestFindLargestPair(t *testing.T) {
	var testCases = []struct {
		name   string
		input  string
		output []int
	}{
		{name: "1273465", input: "1273465", output: []int{7, 6}},

		{name: "987654321111111", input: "987654321111111", output: []int{9, 8}},
		{name: "811111111111119", input: "811111111111119", output: []int{8, 9}},
		{name: "234234234234278", input: "234234234234278", output: []int{7, 8}},
		{name: "818181911112111", input: "818181911112111", output: []int{9, 2}},

		{name: "empty string", input: "", output: nil},
		{name: "length 1", input: "7", output: nil},

		{name: "two digits increasing", input: "12", output: []int{1, 2}},
		{name: "two digits decreasing", input: "21", output: []int{2, 1}},
		{name: "two digits equal", input: "11", output: []int{1, 1}},

		{name: "all digits same", input: "1111", output: []int{1, 1}},
		{name: "all digits same n=3", input: "999", output: []int{9, 9}},

		{name: "strictly increasing", input: "123456", output: []int{5, 6}},
		{name: "strictly increasing with 0", input: "012345", output: []int{4, 5}},

		{name: "strictly decreasing", input: "987654", output: []int{9, 8}},
		{name: "max at start zeros after", input: "9500", output: []int{9, 5}},

		{name: "max repeats at end", input: "9299", output: []int{9, 9}},
		{name: "alternating max", input: "9797", output: []int{9, 9}},
		{name: "max at start repeats later", input: "9191", output: []int{9, 9}},
		{name: "double max start", input: "9912", output: []int{9, 9}},

		{name: "global max only at last", input: "1239", output: []int{3, 9}},
		{name: "leading zeros max last", input: "009", output: []int{0, 9}},

		{name: "original example", input: "127456", output: []int{7, 6}},
		{name: "middle max complex", input: "539741", output: []int{9, 7}},

		{name: "tie first digit second decides", input: "7374", output: []int{7, 7}},
		{name: "tie first digit second decides 2", input: "8789", output: []int{8, 9}},
		{name: "tie on max many", input: "9498", output: []int{9, 9}},

		{name: "leading zeros pair", input: "0099", output: []int{9, 9}},
		{name: "zeros mixed", input: "0909", output: []int{9, 9}},
		{name: "all zeros", input: "0000", output: []int{0, 0}},
		{name: "max then zeros", input: "9000", output: []int{9, 0}},

		{name: "mixed structure 1", input: "506734", output: []int{7, 4}},
		{name: "mixed structure 2", input: "271936", output: []int{9, 6}},
		{name: "mixed structure 3", input: "864208", output: []int{8, 8}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FindLargestPair(ParseLine(tc.input))
			if len(result) != len(tc.output) {
				t.Errorf("Expected %d, got %d", tc.output, result)
			}
			for i, d := range result {
				if d != tc.output[i] {
					t.Errorf("Expected %d, got %d", tc.output, result)
				}
			}
		})
	}
}
