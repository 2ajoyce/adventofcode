package main

import (
	"2ajoyce/adventofcode/2025/8/dsu"
	"2ajoyce/adventofcode/2025/8/point"
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// First Problem
	input := make(chan *point.Point)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input, 1000)
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
	var x, y, z int
	_, err := fmt.Sscanf(input, "%d,%d,%d", &x, &y, &z)
	if err != nil {
		panic(fmt.Sprintf("failed to parse input line %q: %v", input, err))
	}
	return point.NewPoint(x, y, z)
}

func Solve1(input chan *point.Point, numConnections int) (string, error) {
	total := 0
	points := []*point.Point{}

	// Read all points
	for p := range input {
		points = append(points, p)
	}

	// Get the distance from each point to all other points
	pairs := len(points) * (len(points) - 1) / 2
	distances := make([]struct {
		a, b     int
		distance float64
	}, 0, pairs)
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			d := struct {
				a, b     int
				distance float64
			}{a: i, b: j, distance: points[i].Distance(points[j])}
			distances = append(distances, d)
		}
	}

	// Sort distances
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	// Initialize DSU with all points
	d := dsu.NewDSU(len(points))

	// Use DSU to connect the "numConnections" closest pairs
	for i := range numConnections {
		a := distances[i].a
		b := distances[i].b
		d.Union(a, b)
	}

	ds := d.SetSizes() // ds is a map of root -> size

	sizes := []int{}
	for _, size := range ds {
		sizes = append(sizes, size)
	}
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	// fmt.Println("Disjoint Sets Sizes:")
	// fmt.Printf(" - 1: %d\n", sizes[0])
	// fmt.Printf(" - 2: %d\n", sizes[1])
	// fmt.Printf(" - 3: %d\n", sizes[2])
	total = sizes[0] * sizes[1] * sizes[2]

	return fmt.Sprintf("%d", total), nil
}

func Solve2(input chan *point.Point) (string, error) {
	total := 0
	points := []*point.Point{}

	// Read all points
	for p := range input {
		points = append(points, p)
	}

	// Get the distance from each point to all other points
	pairs := len(points) * (len(points) - 1) / 2
	distances := make([]struct {
		a, b     int
		distance float64
	}, 0, pairs)
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			d := struct {
				a, b     int
				distance float64
			}{a: i, b: j, distance: points[i].Distance(points[j])}
			distances = append(distances, d)
		}
	}

	// Sort distances
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	// Initialize DSU with all points
	d := dsu.NewDSU(len(points))

	// Use DSU to connect pairs until all points are connected
	// We need to track the last two points connected for the result
	var p1, p2 *point.Point // The last two points connected
	for i := 0; d.Count() > 1; i++ {
		a := distances[i].a
		b := distances[i].b
		d.Union(a, b)
		p1 = points[a]
		p2 = points[b]
	}

	total = p1.X() * p2.X()

	return fmt.Sprintf("%d", total), nil
}
