package main

import (
	"day17/internal/day17"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"testing"
)

const INPUT_FILE = "test_input.txt"
const OUTPUT_FILE = "test_output.txt"

func validateOutput(t *testing.T, expectedContent string) bool {
	data, err := os.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
		return false
	}
	content := string(data)
	if content != expectedContent {
		t.Errorf("Expected output to contain '%s', but got: %v", expectedContent, content)
		return false
	}
	// If the validation fails, the input and output are retained for troubleshooting
	os.Remove(INPUT_FILE)
	os.Remove(OUTPUT_FILE)
	return true
}

func TestMain(m *testing.M) {
	// Set up environment variables here
	os.Setenv("INPUT_FILE", INPUT_FILE)
	os.Setenv("OUTPUT_FILE", OUTPUT_FILE)
	os.Setenv("PARALLELISM", "1")
	os.Setenv("DEBUG", "true")

	// Run all tests
	code := m.Run()

	// Clean up any resources if necessary
	os.Unsetenv("INPUT_FILE")
	os.Unsetenv("OUTPUT_FILE")
	os.Unsetenv("PARALLELISM")
	os.Unsetenv("DEBUG")

	// Exit with the same status as `go test`
	os.Exit(code)
}

func TestSolveComputer1(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(729))
	comp.SetOpcodes([]day17.Opcode{0, 1, 5, 4, 3, 0})
	expectedOutput := "4,6,3,5,6,3,5,2,1,0"
	output := ""

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for out := range comp.Output {
			output = output + fmt.Sprintf("%s,", out)
		}
		if len(output) > 0 {
			output = strings.TrimSuffix(output, ",")
		}
	}()

	err := SolveComputer(0, comp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	wg.Wait()

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s' but got: '%s'", expectedOutput, output)
	}
}

func TestSolveComputer2(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(10))
	comp.SetOpcodes([]day17.Opcode{5, 0, 5, 1, 5, 4})
	expectedOutput := "0,1,2"
	output := ""

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for out := range comp.Output {
			output = output + fmt.Sprintf("%s,", out)
		}
		// Remove the trailing comma
		if len(output) > 0 {
			output = strings.TrimSuffix(output, ",")
		}
	}()

	err := SolveComputer(0, comp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	wg.Wait()

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s', but got: %s", expectedOutput, output)
	}
}

func TestSolveComputer3(t *testing.T) {
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(2024))
	comp.SetOpcodes([]day17.Opcode{0, 1, 5, 4, 3, 0})
	expectedOutput := "4,2,5,6,7,7,7,7,3,1,0"
	output := ""

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for out := range comp.Output {
			output = output + fmt.Sprintf("%s,", out)
		}
		// Remove the trailing comma
		if len(output) > 0 {
			output = strings.TrimSuffix(output, ",")
		}
	}()

	err := SolveComputer(0, comp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	wg.Wait()

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s', but got: %s", expectedOutput, output)
	}
}

func TestSolveCanSolvePart2(t *testing.T) {
	// This test is checking the known example from the description to prove that the program
	// works as expected
	comp := day17.NewComputer()
	comp.SetRegisterA(big.NewInt(117440))
	comp.SetOpcodes([]day17.Opcode{0, 3, 5, 4, 3, 0})
	expectedOutput := "0,3,5,4,3,0"
	output := ""

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for out := range comp.Output {
			output = output + fmt.Sprintf("%s,", out)
		}
		// Remove the trailing comma
		if len(output) > 0 {
			output = strings.TrimSuffix(output, ",")
		}
	}()

	err := SolveComputer(0, comp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	wg.Wait()

	if output != expectedOutput {
		t.Errorf("Expected output to be '%s', but got: %s", expectedOutput, output)
	}
}