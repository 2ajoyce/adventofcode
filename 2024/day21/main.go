package main

import (
	"day21/internal/aocUtils"
	"day21/internal/day21"
	"fmt"
	"os"
	"slices"
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

// Map of starting coordinate, rune input, depth, and optimal output
type optimalValueMap map[day21.Coord]map[rune]map[int]string

func solve(codes []string) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	DEPTH := os.Getenv("DEPTH")
	depth, err := strconv.Atoi(DEPTH)
	if err != nil {
		depth = 1
	}
	fmt.Println("Beginning solve...")

	if DEBUG {
		fmt.Println("Codes:")
		for _, code := range codes {
			fmt.Printf("    %s\n", code)
		}
	}

	var totalCost int = 0
	for d := range depth {
		totalCost = 0
		// fmt.Println("Generating optimal directional values...")
		optimalDirectionalValues := generateOptimalDirectionalValues(d)
		// fmt.Println("Generated optimal directional values")
		optimizedValues := generateOptimalNumericValues(optimalDirectionalValues, d)
		for code := range codes {
			input := ""
			nk := day21.NewNumericKeypad()
			currentLocation := nk.GetCurrentPosition()
			for _, c := range codes[code] {
				input += optimizedValues[currentLocation][c][0]
				currentLocation = nk.GetPosition(c)
			}
			cost, err := calculateCost(codes[code], len(input))
			if err != nil {
				return nil, fmt.Errorf("error calculating cost: %v", err)
			}
			fmt.Printf("Depth: %d, Code: %s, InputLen: %d, Cost: %d\n", d, codes[code], len(input), cost)
			totalCost += cost
		}
	}

	totalCostStr := strconv.Itoa(totalCost)
	fmt.Printf("Total Cost: %s\n", totalCostStr)
	return []string{totalCostStr}, nil
}

func generateOptimalNumericValues(optimalDirectionalValues optimalValueMap, depth int) optimalValueMap {
	// Create a map of the shortest possible outputs for each integer 0-9
	optimizedValues := make(map[day21.Coord]map[rune]map[int]string)
	// Start a timer
	// start := time.Now()
	for x := 0; x < 3; x++ {
		for y := 0; y < 4; y++ {
			// fmt.Printf("Generating optimal numeric values for %d, %d : %.2fs\n", x, y, time.Since(start).Seconds())
			c := day21.Coord{X: x, Y: y}
			optimizedValues[c] = make(map[rune]map[int]string)
			for i := 0; i < 10; i++ {
				r := rune(i + 48)
				optimizedValues[c][r] = make(map[int]string)
				optimizedValues[c][r][0] = generateOptimalNumericValuesForCoord(optimalDirectionalValues, c, r, depth)
			}
			// Find the optimal values for A also
			optimizedValues[c]['A'] = make(map[int]string)
			optimizedValues[c]['A'][0] = generateOptimalNumericValuesForCoord(optimalDirectionalValues, c, 'A', depth)
		}
	}
	return optimizedValues
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
			output += optimalValues[dk1.GetCurrentPosition()][nkmChar][maxDepth]
			dk1.Move(optimalValues[dk1.GetCurrentPosition()][nkmChar][0])
		}
		if len(output) < len(smallestOutput) || smallestOutput == "" {
			smallestOutput = output
		}
	}
	return smallestOutput
}

func generateOptimalDirectionalValues(depth int) optimalValueMap {
	optimalDirectionalValues := make(optimalValueMap)
	possibleRunes := []rune{'<', '>', '^', 'v', 'A'}
	for x := 0; x < 3; x++ {
		for y := 0; y < 2; y++ {
			c := day21.Coord{X: x, Y: y}
			optimalDirectionalValues[c] = make(map[rune]map[int]string)
			for _, r := range possibleRunes {
				optimalDirectionalValues[c][r] = make(map[int]string)
				result := generateOptimalDirectionalValuesForCoord(day21.Coord{X: x, Y: y}, r, 1)
				optimalDirectionalValues[c][r][1] = result[1]
			}
		}
	}
	if depth == 1 {
		return optimalDirectionalValues
	}

	a1 := make(map[day21.Coord]map[rune]string)
	for x := 0; x < 3; x++ {
		for y := 0; y < 2; y++ {
			c := day21.Coord{X: x, Y: y}
			a1[c] = make(map[rune]string)
			for _, r := range possibleRunes {
				a1[c][r] = optimalDirectionalValues[c][r][1]
			}
		}
	}

	trie := aocUtils.NewTrie()
	// t := time.Now()
	for d := 2; d <= depth; d++ {
		// fmt.Printf("Generating optimal directional values for depth %d: %.2fs\n", d, time.Since(t).Seconds())
		for x := 0; x < 3; x++ {
			for y := 0; y < 2; y++ {
				if x == 0 && y == 0 {
					continue
				}
				c := day21.Coord{X: x, Y: y}
				// t2 := time.Now()
				for _, r := range possibleRunes {
					movement := a1[c][r]
					sub := trie.Substitute(movement, func(sub string) string {
						keyboard := day21.NewDirectionalKeypad()
						currentLocation := keyboard.GetCurrentPosition()
						replacement := ""
						for i, move := range sub {
							if i > 0 {
								currentLocation = keyboard.GetPosition(rune(movement[i-1]))
							}
							subMove := optimalDirectionalValues[currentLocation][move][1]
							replacement += subMove
						}
						return replacement
					})

					// fmt.Printf("\rGenerating optimal directional values for %d, %d, %s (%d): %.2fs", x, y, string(r), len(movement), time.Since(t2).Seconds())
					trie.Insert(movement, sub)
					a1[c][r] = sub
				}
			}
		}
		// fmt.Printf("\r")
		// Copy the values from a1 to optimalDirectionalValues
		for k, v := range a1 {
			for k2, v2 := range v {
				optimalDirectionalValues[k][k2][d] = v2
			}
		}
	}

	return optimalDirectionalValues
}

func generateOptimalDirectionalValuesForCoord(c day21.Coord, input rune, depth int) map[int]string {
	results := generateDirectionalValuesForCoordAtDepth(c, input, depth)
	if depth == 1 {
		slices.Sort(results[1])
		return map[int]string{1: results[1][0]}
	}
	optimalPath := make(map[int]string)
	for i := 1; i < depth; i++ {
		slices.Sort(results[i])
		optimalPath[i] = results[i][0]
	}
	return optimalPath
}

func generateDirectionalValuesForCoordAtDepth(c day21.Coord, input rune, depth int) map[int][]string {
	r1 := generateDirectionalValuesForCoord(c, input)
	r2 := make([]string, 0)
	resultsByDepth := make(map[int][]string)
	if depth == 1 {
		resultsByDepth[1] = r1
	}
	for i := 1; i < depth; i++ {
		r2 = make([]string, 0) // Clear r2
		for _, s := range r1 { // For every string in r1
			// Assign the results of promoting the directional value to r2
			r2 = append(r1, promoteDirectionalValue(s)...)
		}
		// Prune r2 for duplicates or strings longer than the shortest option
		r2ByLen := make(map[int][]string)
		shortestLen := len(r2[0])
		for i := 0; i < len(r2); i++ {
			if len(r2[i]) < shortestLen {
				shortestLen = len(r2[i])
			}
			r2ByLen[len(r2[i])] = append(r2ByLen[len(r2[i])], r2[i])
		}
		uniqueR2 := make(map[string]bool)
		for _, s := range r2ByLen[shortestLen] {
			uniqueR2[s] = true
		}
		r1 = make([]string, 0)
		for k := range uniqueR2 {
			r1 = append(r1, k)
		}
		resultsByDepth[i] = r1
	}
	return resultsByDepth
}

func promoteDirectionalValue(input string) []string {
	results := make([]string, 0)
	dk := day21.NewDirectionalKeypad()
	for _, c := range input {
		sub := generateDirectionalValuesForCoord(dk.GetCurrentPosition(), c)
		for _, s := range sub {
			results = append(results, input+s)
		}
	}
	return results
}

func generateDirectionalValuesForCoord(c day21.Coord, input rune) []string {
	dk := day21.NewDirectionalKeypad()
	dk.SetCurrentPosition(c.X, c.Y)
	dkArray := dk.CalculateMovements(input)
	if len(dkArray) == 0 {
		return []string{"A"}
	}

	return dkArray
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
