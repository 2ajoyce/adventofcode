package main

import (
	"day7/internal"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func writeInputToFile(INPUT_FILE string, inputData []string, t *testing.T) {
	err := os.WriteFile(INPUT_FILE, []byte(strings.Join(inputData, "")), 0644)
	if err != nil {
		t.Errorf("Failed to write input data: %v", err)
	}
}

func validateOutput(t *testing.T, content string, expectedContent string) bool {
	if content != expectedContent {
		t.Errorf("Expected \n%s\n but got \n%s\n", expectedContent, content)
		return false
	}
	return true
}

func SetUpTestInput(t *testing.T, INPUT_FILE string) (total big.Int, inputData []string) {
	// Write the input data to input.txt
	inputData = append(inputData, "190: 10 19\n")
	inputData = append(inputData, "3267: 81 40 27\n")
	inputData = append(inputData, "83: 17 5\n")
	inputData = append(inputData, "156: 15 6\n")
	inputData = append(inputData, "7290: 6 8 6 15\n")
	inputData = append(inputData, "161011: 16 10 13\n")
	inputData = append(inputData, "192: 17 8 14\n")
	inputData = append(inputData, "21037: 9 7 18 13\n")
	inputData = append(inputData, "292: 11 6 16 20")
	//total = *big.NewInt(3749) // Part 1
	total = *big.NewInt(11387)
	return total, inputData
}

func SetUpFuzzyInput(t *testing.T, INPUT_FILE string, PROOF_FILE string) (total big.Int, inputData []string) {
	// Set the random seed and print it to the console
	randSeed := int64(rand.Intn(1000))
	rand.Seed(randSeed)
	fmt.Printf("Random seed: %d\n", randSeed)

	NUMBER_OF_LINES := 1000
	MAX_NUMBER_OF_INPUTS_PER_LINE := 10
	MAX_SIZE_OF_INPUT := 50
	fmt.Printf("Generating %d lines of random input...\n", NUMBER_OF_LINES)

	total = *big.NewInt(0)
	// This slice of strings will store the solved form of each line
	proof := make([]string, NUMBER_OF_LINES)

	// Generate 1000 random lines of input
	for i := 0; i < NUMBER_OF_LINES; i++ {
		// GENERATE NUMBERS
		numbersCount := rand.Intn(MAX_NUMBER_OF_INPUTS_PER_LINE-1) + 2 // The number of input numbersStr in the line (2, MAX_NUMBER_OF_INPUTS_PER_LINE]
		numbersStr := make([]string, numbersCount)                     // The string form of the input numbers
		numbersInt := make([]int, numbersCount)                        // Store the numerical values of the input numbers

		// Fill in the numbers slice with random numbers
		for j := 0; j < numbersCount; j++ {
			numbersInt[j] = rand.Intn(MAX_SIZE_OF_INPUT)
			numbersStr[j] = fmt.Sprintf("%d", numbersInt[j])
		}

		// GENERATE OPERATORS
		operatorsCount := numbersCount - 1
		operators := make([]internal.Operator, operatorsCount)
		for j := 0; j < operatorsCount; j++ {
			opType := rand.Intn(3)
			switch opType {
			case 0:
				operators[j] = internal.Add
			case 1:
				operators[j] = internal.Multiply
			case 2:
				operators[j] = internal.Or
			default:
				t.Errorf("unknown operator type: %d, should be 0, 1 or 2", opType)
			}
		}

		// GENERATE TOTAL
		lineTotal := *big.NewInt(int64(numbersInt[0]))
		for j := 1; j < len(numbersStr); j++ {
			switch operators[j-1] {
			case internal.Add:
				lineTotal.Add(&lineTotal, big.NewInt(int64(numbersInt[j])))
			case internal.Multiply:
				lineTotal.Mul(&lineTotal, big.NewInt(int64(numbersInt[j])))
			case internal.Or:
				concatenatedStr := fmt.Sprintf("%s%d", &lineTotal, numbersInt[j])
				_, success := lineTotal.SetString(concatenatedStr, 10)
				if !success {
					t.Errorf("failed to set string for big.Int")
				}
			default:
				t.Errorf("unknown operator: %b, should be %b or %b", operators[j-1], internal.Add, internal.Multiply)
			}
		}

		// GENERATE PROOF
		e := internal.NewEquation(lineTotal, numbersInt)
		valid, err := e.Validate(operators)
		if err != nil || !valid {
			t.Errorf("equation %d is not valid: %s", i, e.String())
			t.FailNow()
		}
		proof[i] = fmt.Sprintln(e)

		// GENERATE INPUT
		line := fmt.Sprintf("%s: %s\n", &lineTotal, strings.Join(numbersStr, " "))
		total.Add(&total, &lineTotal)
		inputData = append(inputData, line)
	}
	err := WriteOutput(PROOF_FILE, proof)
	if err != nil {
		t.Errorf("error writing proof file: %s", err)
	}
	return total, inputData
}

func TestMain(t *testing.T) {
	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"

	// Don't forget to clean up!
	os.Setenv("INPUT_FILE", INPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	defer os.Unsetenv("OUTPUT_FILE")
	os.Setenv("PARALLELISM", "1")
	defer os.Unsetenv("PARALLELISM")
	//os.Setenv("DEBUG", "true")
	//defer os.Unsetenv("DEBUG")

	// Set up the input data
	total, inputData := SetUpTestInput(t, INPUT_FILE)
	writeInputToFile(INPUT_FILE, inputData, t)

	// Run the main function
	main()

	// Read the content of output.txt
	data, err := os.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
	}
	expectedContent := fmt.Sprintf("Calibration Result: %s", &total)
	valid := validateOutput(t, string(data), expectedContent)
	if valid {
		os.Remove(INPUT_FILE)
		os.Remove(OUTPUT_FILE)
	}
}

func TestMainParallel(t *testing.T) {
	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"

	// Don't forget to clean up!
	os.Setenv("INPUT_FILE", INPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	defer os.Unsetenv("OUTPUT_FILE")
	os.Setenv("PARALLELISM", "1")
	defer os.Unsetenv("PARALLELISM")
	//os.Setenv("DEBUG", "true")
	//defer os.Unsetenv("DEBUG")

	// Set up the input data
	total, inputData := SetUpTestInput(t, INPUT_FILE)
	writeInputToFile(INPUT_FILE, inputData, t)

	// Run the main function
	main()

	// Read the content of output.txt
	data, err := os.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
	}
	expectedContent := fmt.Sprintf("Calibration Result: %s", &total)
	valid := validateOutput(t, string(data), expectedContent)
	if valid {
		os.Remove(INPUT_FILE)
		os.Remove(OUTPUT_FILE)
	}
}

func TestMainFuzzy(t *testing.T) {
	// Disable cache for fuzzy testing
	os.Setenv("CACHE", "false")
	defer os.Unsetenv("CACHE")

	INPUT_FILE := "test_input.txt"
	OUTPUT_FILE := "test_output.txt"
	PROOF_FILE := "test_proof.txt"

	// Don't forget to clean up!
	os.Setenv("INPUT_FILE", INPUT_FILE)
	defer os.Unsetenv("INPUT_FILE")
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	defer os.Unsetenv("OUTPUT_FILE")
	os.Setenv("PARALLELISM", "1")
	defer os.Unsetenv("PARALLELISM")
	// os.Setenv("DEBUG", "true")
	// defer os.Unsetenv("DEBUG")

	// Set up the input data
	total, inputData := SetUpFuzzyInput(t, INPUT_FILE, PROOF_FILE)
	fmt.Println("Writing input to file")
	writeInputToFile(INPUT_FILE, inputData, t)

	// Run the main function
	main()

	// Read the content of output.txt
	data, err := os.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
	}
	expectedContent := fmt.Sprintf("Calibration Result: %s", &total)
	valid := validateOutput(t, string(data), expectedContent)
	if valid {
		defer os.Remove(INPUT_FILE)
		defer os.Remove(OUTPUT_FILE)
		defer os.Remove(PROOF_FILE)
	}
}
