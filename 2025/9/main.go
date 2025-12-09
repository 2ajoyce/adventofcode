package main

import (
	"2ajoyce/adventofcode/2025/9/point"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// First Problem
	input := make(chan *point.Point)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan *point.Point)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan *point.Point) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c <- ParseInput(line)
	}
	close(c)
}

// ParseInput parses the input into the necessary data structure.
// On more complex inputs, this allows us to use lines of text as input for tests
func ParseInput(input string) *point.Point {
	// input is in the form "x,y"
	var x, y int
	fmt.Sscanf(input, "%d,%d", &x, &y)
	return point.NewPoint(x, y)
}

func Solve1(input chan *point.Point) (string, error) {
	total := 0
	points := []*point.Point{}
	for p := range input {
		points = append(points, p)
	}

	// Find the largest area between points
	// For each point, find the distance to every other point
	area := 0 // The largest area found so far
	for i := range points {
		for j := range points {
			if i == j {
				continue
			}
			a := Area(points[i], points[j])
			if a > area {
				area = a
			}
		}
	}
	total = area

	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan *point.Point) (string, error) {
	total := 0
	for p := range input {
		total += p.X + p.Y
	}
	return fmt.Sprintf("%d", total), nil
}

// Area calculates the area of the rectangle defined by two points on opposite corners.
// The area is INCLUSIVE of the points.
func Area(p1 *point.Point, p2 *point.Point) int {
	width := intAbs(p1.X-p2.X) + 1
	height := intAbs(p1.Y-p2.Y) + 1
	return width * height
}
