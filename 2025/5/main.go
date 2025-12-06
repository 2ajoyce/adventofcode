package main

import (
	"2ajoyce/adventofcode/2025/5/interval"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	start, end int
}

func main() {
	// First Problem
	cRange := make(chan Range)
	cInt := make(chan int)
	go ReadInput("input1.txt", cRange, cInt)
	result, err := Solve1(cRange, cInt)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	cRange = make(chan Range)
	cInt = make(chan int)
	go ReadInput("input2.txt", cRange, cInt)
	result, err = Solve2(cRange)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, cRange chan Range, cInt chan int) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)

	// The input will have a top section and a bottom section, separated by a newline
	closedRange := false
	for scanner.Scan() {
		line := scanner.Text()
		// Skip the newline when we get there
		if line == "\n" || len(line) < 1 {
			continue
		}
		if strings.Contains(line, "-") {
			cRange <- ParseRange(line)
		} else {
			// When we get the first search, we need to close the range channel
			if !closedRange {
				close(cRange)
				closedRange = true
			}
			cInt <- StrToInt(line)
		}
	}
	close(cInt)
}

func ParseRange(input string) Range {
	arr := strings.Split(input, "-")
	if len(arr) != 2 {
		panic(fmt.Sprintf("error parsing range range %s", input))
	}
	r := Range{start: StrToInt(arr[0]), end: StrToInt(arr[1])}
	return r
}

func Solve1(cRange chan Range, cInt chan int) (string, error) {
	total := 0

	tree := interval.NewIntervalTree()

	for r := range cRange {
		tree.Insert(r.start, r.end)
	}

	for i := range cInt {
		nodes := tree.Search(i)
		if len(nodes) > 0 {
			total++
		}
	}
	return fmt.Sprintf("%d", total), nil
}

func Solve2(cRange chan Range) (string, error) {
	total := int64(0)

	tree := interval.NewIntervalTree()

	for r := range cRange {
		tree.InsertWithoutOverlap(r.start, r.end)
	}

	// DFS to count all ints in all ranges
	total = Dive(tree.Root)

	return fmt.Sprintf("%d", total), nil
}

func Dive(n *interval.Node) int64 {
	total := int64(0)

	total = total + int64(n.End-n.Start+1) // Add one to include end

	if n.Left != nil {
		total = total + Dive(n.Left)
	}

	if n.Right != nil {
		total = total + Dive(n.Right)
	}
	return total
}

func StrToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to convert string %q to integer: %v", s, err))
	}
	return num
}
