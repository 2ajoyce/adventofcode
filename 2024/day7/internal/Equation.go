package internal

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type Operator rune

const (
	Add      Operator = '0'
	Multiply Operator = '1'
	Or       Operator = '2' // Concatenation operator for integers
)

func convertOperatorsToSymbols(operators []Operator) ([]string, error) {
	symbols := make([]string, len(operators))
	for i, op := range operators {
		switch op {
		case Add:
			symbols[i] = "+"
		case Multiply:
			symbols[i] = "*"
		case Or:
			symbols[i] = "||"
		default:
			return nil, fmt.Errorf("unknown operator: %c", op)
		}
	}
	return symbols, nil
}

type Equation interface {
	Validate(operators []Operator) (bool, error)
	Total() big.Int
	Numbers() []int
	Solve() (bool, error)
	IsSolved() bool
	IsValid() bool
	String() string
}

type equation struct {
	total     big.Int
	numbers   []int
	operators []Operator
	valid     bool
	solved    bool
}

func (e *equation) Total() big.Int {
	return e.total
}

func (e *equation) Numbers() []int {
	return e.numbers
}

type EquationValidationError struct {
	Message string
}

func (v *EquationValidationError) Error() string {
	return v.Message
}

func (e *equation) Validate(operators []Operator) (bool, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	// First, make sure that the number of operators is one less than the number of numbers
	// Second, evaluate the equation from left to right ignoring typical precedence rules
	// Third, compare the result with the total and return the result
	symbols, err := convertOperatorsToSymbols(operators)
	if err != nil {
		return false, err
	}
	if DEBUG {
		fmt.Printf("Validating equation with operators: %s\n", strings.Join(symbols, " "))
		numberStrings := make([]string, len(e.numbers))
		for i, number := range e.numbers {
			numberStrings[i] = strconv.Itoa(number)
		}
		fmt.Printf("Validating equation with numbers: %s\n", strings.Join(numberStrings, "  "))
		fmt.Printf("Validating equation with total: %s\n", &e.total)
	}

	if len(operators) != len(e.numbers)-1 {
		return false, &EquationValidationError{"invalid number of operators"}
	}
	result := *big.NewInt(int64(e.numbers[0]))
	for i, operator := range operators {
		switch operator {
		case Add:
			result.Add(&result, big.NewInt(int64(e.numbers[i+1])))
		case Multiply:
			result.Mul(&result, big.NewInt(int64(e.numbers[i+1])))
		case Or:
			concatenatedStr := fmt.Sprintf("%s%d", &result, e.numbers[i+1])
			_, success := result.SetString(concatenatedStr, 10)
			if !success {
				return false, fmt.Errorf("failed to set string for big.Int")
			}
		default:
			return false, &EquationValidationError{"invalid operator"}

		}
	}

	if result.Cmp(&e.total) == 0 {
		e.valid = true
		e.operators = operators
		if DEBUG {
			fmt.Printf("Valid Equation: %v\n\n", e)
		}
		return true, nil
	}
	if DEBUG {
		fmt.Printf("Invalid Equation: %v\n\n", e)
	}
	return false, nil
}

// IsSolved returns whether the function has been solved. It may or may not be a valid equation.
func (e *equation) IsSolved() bool {
	return e.solved
}

// IsValid returns whether the equation is valid. It may or may not have been solved.
func (e *equation) IsValid() bool {
	return e.valid
}

func toBase3(n, length int) string {
	if n == 0 {
		return fmt.Sprintf("%0*s", length, "0")
	}
	digits := ""
	for n > 0 {
		digits = fmt.Sprintf("%d", n%3) + digits
		n /= 3
	}
	// Pad with leading zeros to ensure the string has the required length
	if len(digits) < length {
		digits = fmt.Sprintf("%0*s", length, digits)
	}
	return digits
}

func (e *equation) Solve() (bool, error) {
	lenOperatorSets := len(e.numbers) - 1
	if lenOperatorSets == 0 {
		return false, fmt.Errorf("no operators needed for a single number")
	}

	// Calculate the total number of operator sets as 3^lenOperatorSets
	numOperatorSets := int(math.Pow(3, float64(lenOperatorSets)))
	operatorSets := make([]string, numOperatorSets)

	// Generate all possible trinary strings of length lenOperatorSets
	for i := 0; i < numOperatorSets; i++ {
		operatorSets[i] = toBase3(i, lenOperatorSets)
	}

	// Iterate through each operator set
	for _, operatorSet := range operatorSets {
		operators := make([]Operator, lenOperatorSets)

		// Convert each character to an operator
		for i, operatorChar := range operatorSet {
			switch operatorChar {
			case '0':
				operators[i] = Add // Define Add appropriately
			case '1':
				operators[i] = Multiply // Define Multiply appropriately
			case '2':
				operators[i] = Or // Define Or appropriately
			default:
				return false, fmt.Errorf("unknown operator type '%c', should be '0', '1', or '2'", operatorChar)
			}
		}

		// Validate the operator combination
		valid, err := e.Validate(operators)
		if err != nil {
			return false, err
		}

		if valid {
			e.solved = true
			return true, nil
		}
	}

	// If no valid combination is found
	e.solved = false
	return false, nil
}

// Outputs the equation in the form 10 = 5 + 3 * 2
func (e *equation) String() string {
	symbols, err := convertOperatorsToSymbols(e.operators)
	if err != nil {
		fmt.Printf("Failed to convert operators to symbols: %v\n", err)
		return ""
	}

	result := fmt.Sprintf("%s = ", &e.total)
	for i, number := range e.numbers {
		result += fmt.Sprintf("%s", &number)
		if i < len(symbols) {
			result += fmt.Sprintf(" %s ", symbols[i])
		}
		if len(symbols) == 0 && i < len(e.numbers)-1 {
			result += " ? "
		}
	}
	return result
}

func NewEquation(total big.Int, numbers []int) *equation {
	return &equation{total: total, numbers: numbers, valid: false, solved: false, operators: nil}
}
