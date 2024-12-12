package internal

import (
	"errors"
	"fmt"
	"math/big"
)

type Stone struct {
	Value big.Int
}

func NewStone(value *big.Int) *Stone {
	return &Stone{Value: *value}
}

func (s *Stone) ChangeValue(newValue big.Int) *Stone {
	s.Value = newValue
	return s
}

// Custom Error type representing an odd-length number
type OddLengthError struct {
	Length int
}

func (e OddLengthError) Error() string {
	return fmt.Sprintf("odd length number: %d", e.Length)
}

func (s *Stone) IsEven() bool {
	valueStr := s.Value.String()
	length := len(valueStr)
	return length%2 == 0
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
func (s *Stone) Split() (left, right *Stone, err error) {
	// Convert the big.Int to a string for easier manipulation.
	valueStr := s.Value.String()
	length := len(valueStr)
	if !s.IsEven() {
		return nil, nil, OddLengthError{length}
	}
	halfLength := length / 2
	leftValue, success := big.NewInt(0).SetString(valueStr[:halfLength], 10)
	if !success {
		return nil, nil, errors.New("failed to convert left half of the number")
	}
	rightValue, success := big.NewInt(0).SetString(valueStr[halfLength:], 10)
	if !success {
		return nil, nil, errors.New("failed to convert right half of the number")
	}
	left = NewStone(leftValue)
	right = NewStone(rightValue)
	return
}

func PrintStones(stones []Stone) {
	fmt.Printf("Stones(%d): [", len(stones))
	for i, stone := range stones {
		fmt.Printf("%s", stone.Value.String())
		if i < len(stones)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Println("]")
}
