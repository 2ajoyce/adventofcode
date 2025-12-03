package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// First Problem
	input := make(chan *Span)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan *Span)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

type Span struct {
	start []rune
	end   []rune
}

// ReadInput reads the input from the filepath and sends each span to the provided channel.
func ReadInput(filepath string, c chan *Span) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// Line has format 11-22,95-115
		// Split the line into distinct spans
		spans := strings.Split(line, ",")
		var start, end []rune
		for _, span := range spans {
			splitSpan := strings.Split(span, "-")
			start = StripPadding(splitSpan[0])
			end = StripPadding(splitSpan[1])
			c <- &Span{start, end}
		}
	}
	close(c)
}

func StripPadding(s string) []rune {
	var output []rune
	// No numbers should be zero padded, but we're going to convert to/from int just in case
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert '%s' to integer", s))
	}
	for _, r := range fmt.Sprintf("%d", i) {
		output = append(output, r)
	}
	return output
}

func Solve1(input chan *Span) (string, error) {
	result := 0
	for span := range input {
		invalidIds := CheckSpan(span)
		for _, id := range invalidIds {
			result += ArrRuneToInt(id)
		}
	}
	return fmt.Sprintf("%d", result), nil
}

func Solve2(input chan *Span) (string, error) {
	result := 0
	for span := range input {
		invalidIds := CheckSpan2(span)
		for _, id := range invalidIds {
			result += ArrRuneToInt(id)
		}
	}
	return fmt.Sprintf("%d", result), nil
}

// CheckSpan will assess all numbers in a span, returning any that are "doubles"
// This version takes the easy approach iterating over every number in the span.
func CheckSpan(s *Span) [][]rune {
	// Doubled numbers can't be odd so we can rule out any span where the start and end have N digits and N is odd
	if len(s.start) == len(s.end) && len(s.start)%2 != 0 {
		return [][]rune{}
	}

	sNum := ArrRuneToInt(s.start)
	// If the start length is odd, we need to increment it till it is even
	for len(IntToArrRune(sNum))%2 != 0 {
		sNum++
	}
	start := IntToArrRune(sNum)

	eNum := ArrRuneToInt(s.end)
	// If the end length is odd we need to decrement it till it is even
	for len(IntToArrRune(eNum))%2 != 0 {
		eNum--
	}
	end := IntToArrRune(eNum)

	// Split the number in half to dummyproof logic
	sLeft := start[0 : len(start)/2]
	// sRight := start[len(start)/2:]

	eLeft := end[0 : len(end)/2]
	// eRight := end[len(end)/2:]

	// Get numeric representations to work with
	sLeftI := ArrRuneToInt(sLeft)
	eLeftI := ArrRuneToInt(eLeft)

	// Since the definition of a doubled number is that sLeft and sRight are equal
	// we only need to increment sLeft. We'll increase it till it is larger than eLeft.
	possibleDoubles := [][]rune{}
	for sLeftI <= eLeftI {
		possibleDoubles = append(possibleDoubles, StrToArrRune(fmt.Sprintf("%d%d", sLeftI, sLeftI)))
		sLeftI++
	}

	doubles := [][]rune{}
	for _, r := range possibleDoubles {
		i := ArrRuneToInt(r)
		if i >= sNum && i <= eNum {
			doubles = append(doubles, r)
		}
	}

	return doubles
}

// CheckSpan2 will assess all numbers in a span, returning any that are combinations of the same digits repeated at least twice
func CheckSpan2(s *Span) [][]rune {
	possibleReps := [][]rune{}
	start := ArrRuneToInt(s.start)
	end := ArrRuneToInt(s.end)
	for start <= end {
		if IsInvalidId(IntToArrRune(start)) {
			possibleReps = append(possibleReps, IntToArrRune(start))
		}
		start++
	}

	sNum := ArrRuneToInt(s.start)
	eNum := ArrRuneToInt(s.end)
	// Limit to min/max & two digits
	limited := limitMinMax(possibleReps, sNum, eNum)

	// Dedupe
	deduped := dedupeRuneSlices(limited)

	// Sort the results for test consistency
	sortRuneSlicesByIntValue(deduped)

	return deduped
}

func IsInvalidId(r []rune) bool {
	n := len(r)
	// A single digit can't be invalid
	if n < 2 {
		return false
	}

	// Can't be all zeros
	if isAllZeros(r) {
		return false
	}

	// Try all possible chunk sizes that could repeat to form the whole ID.
	for size := 1; size <= n/2; size++ {
		if n%size != 0 {
			continue // chunk size must divide the total length
		}

		repeats := n / size
		if repeats < 2 {
			continue // must repeat at least once
		}

		// Check if this is a repeated pattern
		if isRepeatedPattern(r, size) {
			return true
		}
	}

	return false // no repeated pattern found => valid
}

func isAllZeros(r []rune) bool {
	for _, d := range r {
		if d != '0' {
			return false
		}
	}
	return len(r) > 0
}

// isRepeatedPattern returns true if r is made of r[0:size] repeated.
func isRepeatedPattern(r []rune, size int) bool {
	for i := size; i < len(r); i++ {
		if r[i] != r[i%size] {
			return false
		}
	}
	return true
}

func StrToArrRune(s string) []rune {
	return []rune(s)
}

func ArrRuneToInt(r []rune) int {
	i, err := strconv.Atoi(string(r))
	if err != nil {
		panic(fmt.Sprintf("failed to convert rune slice %v to integer: %v", r, err))
	}
	return i
}

func IntToArrRune(i int) []rune {
	return []rune(strconv.Itoa(i))
}

func StrToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to convert string %q to integer: %v", s, err))
	}
	return num
}

func ArrRuneToStr(r []rune) string {
	return string(r)
}

func limitMinMax(reps [][]rune, sNum, eNum int) [][]rune {
	result := [][]rune{}

	for _, r := range reps {
		id := ArrRuneToInt(r)
		if id < sNum || id > eNum || id < 10 {
			continue
		}
		result = append(result, r)
	}

	return result
}

func dedupeRuneSlices(reps [][]rune) [][]rune {
	seen := make(map[int]struct{})
	result := [][]rune{}

	for _, r := range reps {
		id := ArrRuneToInt(r)
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, r)
	}

	return result
}

func sortRuneSlicesByIntValue(slices [][]rune) {
	sort.Slice(slices, func(i, j int) bool {
		return ArrRuneToInt(slices[i]) < ArrRuneToInt(slices[j])
	})
}
