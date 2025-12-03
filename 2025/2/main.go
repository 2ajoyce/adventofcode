package main

import (
	"bufio"
	"fmt"
	"os"
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
		invalidIds := CheckSpan(span)
		for _, id := range invalidIds {
			result += ArrRuneToInt(id)
		}
	}
	return fmt.Sprintf("%d", result), nil
}

// CheckSpan will assess all numbers in a span, returning any that are "doubles"
// This version takes the easy approach iterating over every number in the span.
func CheckSpan(s *Span) [][]rune {
	fmt.Printf("CheckSpan: %c | %c\n", s.start, s.end)
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
	sRight := start[len(start)/2:]
	fmt.Printf("Start: %c | %c\n", sLeft, sRight)

	eLeft := end[0 : len(end)/2]
	eRight := end[len(end)/2:]
	fmt.Printf("End: %c | %c\n", eLeft, eRight)

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

// CheckSpan will assess all numbers in a span, returning any that are "doubles"
// This version attempts to be smarter than brute force, but is not yet functional
func CheckSpan1(s *Span) []string {
	fmt.Printf("CheckSpan: %c | %c\n", s.start, s.end)
	// Doubled numbers can't be odd so we can rule out any span where the start and end have N digits and N is odd
	if len(s.start) == len(s.end) && len(s.start)%2 != 0 {
		return []string{}
	}
	// If we've reached this point there are two possible cases
	// - The start and end have length N and N is even
	// - The start and end have different lengths

	// Case1: The start and end have length N and N is even
	// Whether doubles will exist in this range can be determined by walking the digits from left to right
	//   ie, 10, 19: 1 -> 1 gives us no flexibility in the first place
	//   ie, 10, 99: 1 -> 9 gives us 8 digits of flexibility in the first place

	// Once we reach a digit that has flexibility we can start to construct the double
	//   ie:  10, 19: The first half is 1 with no flexibility. That means the only possible double is 11
	//   We can check if 11 is less than the end. 11 < 19, so we can stop searching and return that.

	//   ie: 10, 99: The first half is 1, 2, 3, 4, 5, 6, 7, 8, 9. The doubles would be 11, 22, 33, 44, 55, 66, 77, 88, 99.
	//   We can compare each double to the end(inclusive) and stop searching. All doubles can be returned.

	//   ie: 1000, 1012: The first half is 10:10, 1->1, 0->0 have no flexibility. The only possible couple is 1010
	//   We can compare 1010 to the end(inclusive) and stop searching. 1010 can be returned.

	// Split the number in half to dummyproof logic
	sLeft := s.start[0 : len(s.start)/2]
	sRight := s.start[len(s.start)/2:]
	fmt.Printf("Start: %c | %c\n", sLeft, sRight)

	for i, d := range sLeft { // For every digit in the left half of the start
		if d == rune(sRight[i]) {
			fmt.Printf("%c & %c are the same\n", d, rune(sRight[i]))
			continue
		} else {

		}
	}

	// Case 2: The start and end have different lengths
	// WIP

	return []string{}
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
