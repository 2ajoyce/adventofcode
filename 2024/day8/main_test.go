package main

import (
	"fmt"
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

func SetUpTestInput(t *testing.T, INPUT_FILE string) (total int, inputData []string) {
	inputData = append(inputData, "............\n")
	inputData = append(inputData, "........0...\n")
	inputData = append(inputData, ".....0......\n")
	inputData = append(inputData, ".......0....\n")
	inputData = append(inputData, "....0.......\n")
	inputData = append(inputData, "......A.....\n")
	inputData = append(inputData, "............\n")
	inputData = append(inputData, "............\n")
	inputData = append(inputData, "........A...\n")
	inputData = append(inputData, ".........A..\n")
	inputData = append(inputData, "............\n")
	inputData = append(inputData, "............")
	total = 34
	return total, inputData
}

func SetUpFuzzyInput(t *testing.T, INPUT_FILE string, PROOF_FILE string) (total int, inputData []string) {
	// Set the random seed and print it to the console
	randSeed := int64(rand.Intn(1000))
	rand.Seed(randSeed)
	fmt.Printf("Random seed: %d\n", randSeed)

	NUMBER_OF_LINES := 1000

	fmt.Printf("Generating %d lines of random input...\n", NUMBER_OF_LINES)

	proof := make([]string, NUMBER_OF_LINES)

	for i := 0; i < NUMBER_OF_LINES; i++ {
		line := ""
		total++
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
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("DEBUG")

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
	expectedContent := fmt.Sprintf("Unique Antinode Locations: %d", total)
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
	os.Setenv("PARALLELISM", "2")
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
	expectedContent := fmt.Sprintf("Unique Antinode Locations: %d", total)
	valid := validateOutput(t, string(data), expectedContent)
	if valid {
		os.Remove(INPUT_FILE)
		os.Remove(OUTPUT_FILE)
	}
}

func TestMainFuzzy(t *testing.T) {
	// Disable cache for fuzzy testing
	// For some reason this only appears to work with "run test" not with the sidebar test runner
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
	expectedContent := fmt.Sprintf("Calibration Result: %d", &total)
	valid := validateOutput(t, string(data), expectedContent)
	if valid {
		defer os.Remove(INPUT_FILE)
		defer os.Remove(OUTPUT_FILE)
		defer os.Remove(PROOF_FILE)
	}
}
