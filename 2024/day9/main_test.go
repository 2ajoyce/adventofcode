package main

import (
	"fmt"
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

func SetUpFullTestInput(t *testing.T, INPUT_FILE string) (total int, inputData []string) {
	inputData = append(inputData, "12345")
	total = 60
	// inputData = append(inputData, "2333133121414131402")
	// total = 1928
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
	total, inputData := SetUpFullTestInput(t, INPUT_FILE)
	writeInputToFile(INPUT_FILE, inputData, t)

	// Run the main function
	main()

	// Read the content of output.txt
	data, err := os.ReadFile(OUTPUT_FILE)
	if err != nil {
		t.Errorf("Failed to read %s: %v", OUTPUT_FILE, err)
	}
	expectedContent := fmt.Sprintf("Checksum: %d", total)
	valid := validateOutput(t, string(data), expectedContent)
	if valid {
		os.Remove(INPUT_FILE)
		os.Remove(OUTPUT_FILE)
	}
}

// func TestConvertInput1(t *testing.T) {
// 	inputData := []string{"12345"}
// 	expectedOutput := "0..111....22222"
// 	output := convertInput(inputData[0])
// 	if output != expectedOutput {
// 		t.Errorf("Expected %s but got %s", expectedOutput, output)
// 	}
// }

// func TestConvertInput2(t *testing.T) {
// 	inputData := []string{"2333133121414131402"}
// 	expectedOutput := "00...111...2...333.44.5555.6666.777.888899"
// 	output := convertInput(inputData[0])
// 	if output != expectedOutput {
// 		t.Errorf("Expected %s but got %s", expectedOutput, output)
// 	}
// }

// func TestSort1(t *testing.T) {
// 	inputData := []string{"0..111....22222"}
// 	expectedOutput := "022111222......"
// 	output := sortInput(inputData[0])
// 	if output != expectedOutput {
// 		t.Errorf("Expected %s but got %s", expectedOutput, output)
// 	}
// }

// func TestSort2(t *testing.T) {
// 	inputData := []string{"00...111...2...333.44.5555.6666.777.888899"}
// 	expectedOutput := "0099811188827773336446555566.............."
// 	output := sortInput(inputData[0])
// 	if output != expectedOutput {
// 		t.Errorf("Expected %s but got %s", expectedOutput, output)
// 	}
// }

// func TestChecksum1(t *testing.T) {
// 	inputData := []string{"022111222"}
// 	expectedOutput := "60" //"0*0+1*2+2*2+3*1+4*1+5*1+6*2+7*2+8*2"
// 	output := checksumInput(inputData[0])
// 	if output != expectedOutput {
// 		t.Errorf("Expected %s but got %s", expectedOutput, output)
// 	}
// }

// func TestChecksum2(t *testing.T) {
// 	inputData := []string{"0099811188827773336446555566"}
// 	expectedOutput := "1928"
// 	output := checksumInput(inputData[0])
// 	if output != expectedOutput {
// 		t.Errorf("Expected %s but got %s", expectedOutput, output)
// 	}
// }
