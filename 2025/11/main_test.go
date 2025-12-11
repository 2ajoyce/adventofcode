package main

import (
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
			"aaa: you hhh",
			"you: bbb ccc",
			"bbb: ddd eee",
			"ccc: ddd eee fff",
			"ddd: ggg",
			"eee: out",
			"fff: out",
			"ggg: out",
			"hhh: ccc fff iii",
			"iii: out",
		}, output: "5"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan *Graph)
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
			"svr: aaa bbb",
			"aaa: fft",
			"fft: ccc",
			"bbb: tty",
			"tty: ccc",
			"ccc: ddd eee",
			"ddd: hub",
			"hub: fff",
			"eee: dac",
			"dac: fff",
			"fff: ggg hhh",
			"ggg: out",
			"hhh: out",
		}, output: "2"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			inputChan := make(chan *Graph)
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
