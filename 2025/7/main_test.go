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
		{name: "Simple Beam", input: []string{
			".......S.......",
			"...............",
			"...............",
		}, output: "0"},
		{name: "One Splitter", input: []string{
			".......S.......",
			".......^.......",
			"...............",
		}, output: "1"},
		{name: "Three Splitters", input: []string{
			".......S.......",
			".......^.......",
			"......^.^......",
		}, output: "3"},
		{name: "AOC Example 1", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".....^.^.^.....",
			"...............",
			"....^.^...^....",
			"...............",
			"...^.^...^.^...",
			"...............",
			"..^...^.....^..",
			"...............",
			".^.^.^.^.^...^.",
			"...............",
		}, output: "21"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan string)
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
		{name: "Simple Beam", input: []string{
			".......S.......",
			"...............",
		}, output: "1"},
		{name: "One Splitter", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
		}, output: "2"},
		{name: "Three Splitters", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
		}, output: "4"},
		{name: "Four Splitters", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".........^.....",
			"...............",
		}, output: "5"},
		{name: "Merging Beams", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".......^.......",
			"...............",
		}, output: "6"},
		{name: "Skip Level", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^........",
			"...............",
		}, output: "3"},
		{name: "Overlapping Terminals", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".....^...^.....",
			"...............",
			"....^.^.^.^....",
			"...............",
		}, output: "10"},
		{name: "AOC Example 1", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".....^.^.^.....",
			"...............",
			"....^.^...^....",
			"...............",
			"...^.^...^.^...",
			"...............",
			"..^...^.....^..",
			"...............",
			".^.^.^.^.^...^.",
			"...............",
		}, output: "40"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			inputChan := make(chan string)
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

func TestSolve3(t *testing.T) {
	var testCases = []struct {
		name   string
		input  []string
		output string
	}{
		{name: "Simple Beam", input: []string{
			".......S.......",
			"...............",
		}, output: "1"},
		{name: "One Splitter", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
		}, output: "2"},
		{name: "Three Splitters", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
		}, output: "4"},
		{name: "Four Splitters", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".........^.....",
			"...............",
		}, output: "5"},
		{name: "Merging Beams", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".......^.......",
			"...............",
		}, output: "6"},
		{name: "Skip Level", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^........",
			"...............",
		}, output: "3"},
		{name: "Overlapping Terminals", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".....^...^.....",
			"...............",
			"....^.^.^.^....",
			"...............",
		}, output: "10"},
		{name: "AOC Example 1", input: []string{
			".......S.......",
			"...............",
			".......^.......",
			"...............",
			"......^.^......",
			"...............",
			".....^.^.^.....",
			"...............",
			"....^.^...^....",
			"...............",
			"...^.^...^.^...",
			"...............",
			"..^...^.....^..",
			"...............",
			".^.^.^.^.^...^.",
			"...............",
		}, output: "40"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			inputChan := make(chan string)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- ParseInput(line)
				}
			}()
			result, err := Solve3(inputChan)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tc.output {
				t.Errorf("Expected %s, got %s", tc.output, result)
			}
		})
	}
}
