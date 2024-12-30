package main

import (
	"day21/internal/aocUtils"
	"day21/internal/day21"
	"fmt"
	"os"
	"strconv"

	"github.com/schollz/progressbar/v3"
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

// Map of starting coordinate, rune input, depth, and optimal output
type optimalValueMap map[day21.Coord]map[rune]map[int]string

func solve(codes []string) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning solve...")

	if DEBUG {
		fmt.Println("Codes:")
		for _, code := range codes {
			fmt.Printf("    %s\n", code)
		}
	}

	depth := 1
	optimizedValues := generateOptimalNumericValues(depth)

	for k, v := range optimizedValues[day21.Coord{X: 2, Y: 3}] {
		fmt.Printf("Key: %c, Value: %s\n", k, v[depth])
	}

	totalCost := 0
	for code := range codes {
		input := ""
		nk := day21.NewNumericKeypad()
		currentLocation := nk.GetCurrentPosition()
		for _, c := range codes[code] {
			input += optimizedValues[currentLocation][c][depth]
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
	fmt.Printf("Total Cost: %s\n", totalCostStr)
	return []string{totalCostStr}, nil
}

func generateOptimalNumericValues(depth int) optimalValueMap {
	optimalDirectionalValues := generateOptimalDirectionalValues()

	bar := progressbar.Default(int64(3 * 4 * 11))
	// Create a map of the shortest possible outputs for each integer 0-9
	optimizedValues := make(map[day21.Coord]map[rune]map[int]string)
	for x := 0; x < 3; x++ {
		for y := 0; y < 4; y++ {
			c := day21.Coord{X: x, Y: y}
			optimizedValues[c] = make(map[rune]map[int]string)
			for i := 0; i < 10; i++ {
				optimizedValues[c][rune(i+48)] = make(map[int]string)
				optimizedValues[c][rune(i+48)][1] = generateOptimalNumericValuesForCoord(optimalDirectionalValues, c, rune(i+48), depth)
			}
			// Find the optimal values for A also
			optimizedValues[c]['A'] = make(map[int]string)
			optimizedValues[c]['A'][1] = generateOptimalNumericValuesForCoord(optimalDirectionalValues, c, 'A', depth)
			bar.Add(1)
		}
	}
	return optimizedValues
}

func recursiveFunc(optimalValues optimalValueMap, dk1 *day21.DirectionalKeypad, priorChar rune, depth int, maxDepth int) string {
	output := ""
	currentCoord := dk1.GetCurrentPosition()
	dk1Movement := optimalValues[currentCoord][priorChar][1]
	if depth > maxDepth {
		dk1.Move(dk1Movement)
		return dk1Movement
	}
	dk2 := day21.NewDirectionalKeypad()
	for _, dk1Char := range dk1Movement {
		newCoord := dk2.GetCurrentPosition()
		if _, ok := optimalValues[newCoord][dk1Char][depth+1]; !ok {
			optimalValues[newCoord][dk1Char][depth+1] = recursiveFunc(optimalValues, dk2, dk1Char, depth+1, maxDepth)
		} else {
			dk2.Move(optimalValues[newCoord][dk1Char][depth+1])
		}
		output += optimalValues[newCoord][dk1Char][depth+1]
	}
	dk1.Move(dk1Movement)
	return output
}

func generateOptimalNumericValuesForCoord(optimalValues optimalValueMap, c day21.Coord, input rune, maxDepth int) string {
	nk := day21.NewNumericKeypad()
	nk.SetCurrentPosition(c.X, c.Y)
	nkmArray := nk.CalculateMovements(input)
	if len(nkmArray) == 0 {
		return ""
	}

	smallestOutput := ""
	for _, nkm := range nkmArray {
		output := ""
		dk1 := day21.NewDirectionalKeypad()
		for _, nkmChar := range nkm {
			depth := 1
			output += recursiveFunc(optimalValues, dk1, nkmChar, depth, maxDepth)
		}
		if len(output) < len(smallestOutput) || smallestOutput == "" {
			smallestOutput = output
		}
	}

	return smallestOutput
}

func generateOptimalDirectionalValues() optimalValueMap {
	optimalDirectionalValues := make(optimalValueMap)
	possibleRunes := []rune{'<', '>', '^', 'v', 'A'}
	for x := 0; x < 3; x++ {
		for y := 0; y < 2; y++ {
			c := day21.Coord{X: x, Y: y}
			optimalDirectionalValues[c] = make(map[rune]map[int]string)
			for _, r := range possibleRunes {
				optimalDirectionalValues[c][r] = make(map[int]string)
				optimalDirectionalValues[c][r][1] = generateOptimalDirectionalValuesForCoord(day21.Coord{X: x, Y: y}, r)
			}
		}
	}
	return optimalDirectionalValues
}

func generateOptimalDirectionalValuesForCoord(c day21.Coord, input rune) string {
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
			path := generateDirectionalValuesForCoord(dk1.GetCurrentPosition(), dkChar)
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

func generateDirectionalValuesForCoord(c day21.Coord, input rune) string {
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
