package internal

import (
	"fmt"
	"math"
)

type Stone = int

// Custom Error type representing an odd-length number
type OddLengthError struct {
	Length int
}

func (e OddLengthError) Error() string {
	return fmt.Sprintf("odd length number: %d", e.Length)
}

func CharCount(s Stone) int {
	return int(math.Log10(float64(s))) + 1
}

func IsEven(s Stone) bool {
	return CharCount(s)%2 == 0
}

// Divide the stone into two equal parts.
// Throw an error if the stone has an odd number of digits
// The left stone will have the left half of digits
// The right stone will have the right half of digits
// Example 1:
//
//	Input: Stone{value: 1234567890}
//	Output:
//	  left:  Stone{value: 12345}
//	  right: Stone{value: 67890}
//
// Example 2:
//
//	Input: Stone{value: 123004}
//	Output:
//	  left:  Stone{value: 123}
//	  right: Stone{value: 4}
func Split(s int) (left, right Stone, err error) {
	if !IsEven(s) {
		return 0, 0, OddLengthError{CharCount(s)}
	}
	k := CharCount(s) / 2
	power := int(math.Pow(10, float64(k)))
	left = s / power
	right = s % power
	return
}

func PrintStones(stones []Stone) {
	fmt.Printf("Stones(%d): [", len(stones))
	for i, stone := range stones {
		fmt.Printf("%d", stone)
		if i < len(stones)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Println("]")
}
