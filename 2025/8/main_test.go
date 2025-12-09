package main

import (
	"2ajoyce/adventofcode/2025/8/point"
	"fmt"
	"testing"
)

func TestSolve1(t *testing.T) {
	var testCases = []struct {
		name           string
		input          []string
		numConnections int
		output         string
	}{
		{name: "AOC Example 1", input: []string{
			"162,817,812",
			"57,618,57",
			"906,360,560",
			"592,479,940",
			"352,342,300",
			"466,668,158",
			"542,29,236",
			"431,825,988",
			"739,650,466",
			"52,470,668",
			"216,146,977",
			"819,987,18",
			"117,168,530",
			"805,96,715",
			"346,949,466",
			"970,615,88",
			"941,993,340",
			"862,61,35",
			"984,92,344",
			"425,690,689",
		}, numConnections: 10, output: "40"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputChan := make(chan *point.Point)
			go func() {
				defer close(inputChan)
				for _, line := range tc.input {
					inputChan <- ParseInput(line)
				}
			}()
			result, err := Solve1(inputChan, tc.numConnections)
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
		name           string
		input          []string
		numConnections int
		output         string
	}{
		{name: "AOC Example 1", input: []string{
			"162,817,812",
			"57,618,57",
			"906,360,560",
			"592,479,940",
			"352,342,300",
			"466,668,158",
			"542,29,236",
			"431,825,988",
			"739,650,466",
			"52,470,668",
			"216,146,977",
			"819,987,18",
			"117,168,530",
			"805,96,715",
			"346,949,466",
			"970,615,88",
			"941,993,340",
			"862,61,35",
			"984,92,344",
			"425,690,689",
		}, output: "25272"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			inputChan := make(chan *point.Point)
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
