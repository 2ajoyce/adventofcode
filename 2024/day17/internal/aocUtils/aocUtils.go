package aocUtils

import (
	"bufio"
	"fmt"
	"os"
)

func ReadInput(INPUT_FILE string) ([]string, error) {
	inputFile, err := os.Open(INPUT_FILE)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %v", INPUT_FILE, err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func WriteToFile(FILE string, results []string) error {
	outputFile, err := os.Create(FILE)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", FILE, err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	// Write the results to output.txt, one line per result
	for i, res := range results {
		_, err := writer.WriteString(res)
		if err != nil {
			return fmt.Errorf("error writing value to %s: %v", FILE, err)
		}
		if i != len(results)-1 {
			_, err = writer.WriteString("\n")
			if err != nil {
				return fmt.Errorf("error writing newline to %s: %v", FILE, err)
			}
		}
	}

	// Flush the writer to ensure all data is written to output.txt
	writer.Flush()
	return nil
}

func AppendToFile(FILE string, results []string) error {
	outputFile, err := os.OpenFile(FILE, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", FILE, err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)

	// Append the results to the file, one line per result
	for _, res := range results {
		_, err := writer.WriteString(res + "\n")
		if err != nil {
			return fmt.Errorf("error writing value to %s: %v", FILE, err)
		}
	}

	// Flush the writer to ensure all data is written to the file
	writer.Flush()
	return nil
}
