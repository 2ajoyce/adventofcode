package main

import (
	"day21/internal/aocUtils"
	"day21/internal/day21"
	"fmt"
	"os"
	"strconv"
)

func main() {

	////////////////////////////////////////////////////////////////////
	// ENVIRONMENT SETUP
	////////////////////////////////////////////////////////////////////

	//os.Setenv("DEBUG", "true")
	INPUT_FILE := os.Getenv("INPUT_FILE")
	OUTPUT_FILE := os.Getenv("OUTPUT_FILE")
	PARALLELISM, err := strconv.Atoi(os.Getenv("PARALLELISM"))
	if PARALLELISM < 1 || err != nil {
		PARALLELISM = 1
	}
	fmt.Printf("PARALLELISM: %d\n\n", PARALLELISM)

	if INPUT_FILE == "" || OUTPUT_FILE == "" {
		fmt.Println("INPUT_FILE and OUTPUT_FILE environment variables not set")
		fmt.Println("Defaulting to input.txt and output.txt")
		INPUT_FILE = "input.txt"
		OUTPUT_FILE = "output.txt"
	}

	////////////////////////////////////////////////////////////////////
	// READ INPUT FILE
	////////////////////////////////////////////////////////////////////

	lines, err := aocUtils.ReadFile(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error reading from %s: %v", INPUT_FILE, err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// SOLUTION LOGIC
	////////////////////////////////////////////////////////////////////

	input, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve(input)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// WRITE OUTPUT FILE
	////////////////////////////////////////////////////////////////////

	err = aocUtils.WriteToFile(OUTPUT_FILE, results)
	if err != nil {
		fmt.Printf("Error writing to %s: %v", OUTPUT_FILE, err)
		return
	}

	fmt.Printf("Successfully processed %s and created %s\n", INPUT_FILE, OUTPUT_FILE)
}

func parseLines(lines []string) ([]string, error) {
	// DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	codes := make([]string, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		codes = append(codes, line)
	}

	return codes, nil
}

func solve(codes []string) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning solve...")

	if DEBUG {
		fmt.Println("Codes:")
		for _, code := range codes {
			fmt.Printf("    %s\n", code)
		}
	}

	// Create a map of the shortest possible outputs for each integer 0-9
	optimizedValues := make(map[day21.Coord]map[rune]string)
	for x := 0; x < 3; x++ {
		for y := 0; y < 4; y++ {
			c := day21.Coord{X: x, Y: y}
			optimizedValues[c] = make(map[rune]string)
			for i := 0; i < 10; i++ {
				optimizedValues[c][rune(i+48)] = complexMachine1(c, rune(i+48))
			}
			// Find the optimal values for A also
			optimizedValues[c]['A'] = complexMachine1(c, 'A')
		}
	}
	for k, v := range optimizedValues[day21.Coord{X: 2, Y: 3}] {
		fmt.Printf("Key: %c, Value: %s\n", k, v)
	}

	totalCost := 0
	for code := range codes {
		input := ""
		nk := day21.NewNumericKeypad()
		currentLocation := nk.GetCurrentPosition()
		for _, c := range codes[code] {
			input += optimizedValues[currentLocation][c]
			currentLocation = nk.GetPosition(c)
		}
		fmt.Printf("Code: %s, Input: %s\n", codes[code], input)
		cost, err := calculateCost(codes[code], len(input))
		if err != nil {
			return nil, fmt.Errorf("error calculating cost: %v", err)
		}
		totalCost += cost
	}

	totalCostStr := strconv.Itoa(totalCost)
	return []string{totalCostStr}, nil
}

type type1 map[day21.Coord]map[rune]string

func complexMachine1(c day21.Coord, input rune) string {
	nk := day21.NewNumericKeypad()
	nk.SetCurrentPosition(c.X, c.Y)
	nkmArray := nk.CalculateMovements(input)
	if len(nkmArray) == 0 {
		return ""
	}

	optimalValues := make(map[day21.Coord]map[rune]string)
	possibleRunes := []rune{'<', '>', '^', 'v', 'A'}
	for x := 0; x < 3; x++ {
		for y := 0; y < 2; y++ {
			c := day21.Coord{X: x, Y: y}
			optimalValues[c] = make(map[rune]string)
			for _, r := range possibleRunes {
				optimalValues[c][r] = complexMachine3(day21.Coord{X: x, Y: y}, r)
			}
		}
	}

	smallestOutput := ""
	for _, nkm := range nkmArray {
		output := ""
		dk1 := day21.NewDirectionalKeypad()
		for _, nkmChar := range nkm {
			dk2 := day21.NewDirectionalKeypad()
			dk1Movement := optimalValues[dk1.GetCurrentPosition()][nkmChar]
			for _, dk1Char := range dk1Movement {
				dk2Movement := optimalValues[dk2.GetCurrentPosition()][dk1Char]
				output += dk2Movement
				dk2.Move(dk2Movement)
			}
			dk1.Move(dk1Movement)
		}
		if len(output) < len(smallestOutput) || smallestOutput == "" {
			smallestOutput = output
		}
	}
	return smallestOutput
}

func complexMachine2(c day21.Coord, input rune) string {
	dk := day21.NewDirectionalKeypad()
	dk.SetCurrentPosition(c.X, c.Y)
	dkArray := dk.CalculateMovements(input)
	if len(dkArray) == 0 {
		return "A"
	}
	smallestOutput := []rune{}
	for _, dkChar := range dkArray {
		output := []rune{}
		for _, dkChar := range dkChar {
			output = append(output, dkChar)
		}
		if len(output) < len(smallestOutput) || len(smallestOutput) == 0 {
			smallestOutput = output
		}
	}
	return string(smallestOutput)
}

func complexMachine3(c day21.Coord, input rune) string {
	dk := day21.NewDirectionalKeypad()
	dk.SetCurrentPosition(c.X, c.Y)
	dkArray := dk.CalculateMovements(input)
	if len(dkArray) == 0 {
		return "A"
	}
	smallestOutput := []rune{}
	smallest2xDeepOutput := []rune{}

	for _, dkm := range dkArray {
		output := []rune{}
		dk1 := day21.NewDirectionalKeypad()
		for _, dkChar := range dkm {
			path := complexMachine2(dk1.GetCurrentPosition(), dkChar)
			for _, char := range path {
				output = append(output, char)
			}
			dk1.Move(path)
		}
		if len(output) < len(smallest2xDeepOutput) || len(smallest2xDeepOutput) == 0 {
			smallest2xDeepOutput = output
			smallestOutput = []rune(dkm)
		}
	}
	return string(smallestOutput)
}

func calculateCost(code string, inputLen int) (int, error) {
	DEBUG := os.Getenv("DEBUG") == "true"

	// Remove the last character from the code
	c := code[:len(code)-1]
	codeInt, err := strconv.Atoi(c)
	if err != nil {
		return 0, fmt.Errorf("error converting code to int: %v", err)
	}

	if DEBUG {
		fmt.Printf("Input Length: %d, Code Int: %d\n", inputLen, codeInt)
	}

	cost := inputLen * codeInt
	if DEBUG {
		fmt.Printf("Cost: %d\n", cost)
	}

	return cost, nil
}
