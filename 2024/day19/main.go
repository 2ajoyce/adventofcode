package main

import (
	"day19/internal/aocUtils"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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

	terms, sentences, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve(terms, sentences)
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

type Term string
type Sentence string

func parseLines(lines []string) ([]Term, []Sentence, error) {
	// DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, nil, fmt.Errorf("input is empty")
	}

	// The first line is the list of comma separated terms
	termStrings := strings.Split(lines[0], ", ")

	terms := make([]Term, 0)
	for _, termString := range termStrings {
		terms = append(terms, Term(termString))
	}

	fmt.Printf("Ingested %d terms\n", len(terms))

	// The rest of the lines are sentences
	lines = lines[1:]

	sentences := make([]Sentence, 0)
	for _, line := range lines {
		if line == "" { // Skip blank lines
			continue
		}

		sentences = append(sentences, Sentence(line))
	}

	fmt.Printf("Ingested %d sentences\n", len(sentences))

	return terms, sentences, nil
}

func solve(terms []Term, sentences []Sentence) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning single-threaded solve")

	if DEBUG {
		fmt.Printf("Terms: %v\n", terms)
		fmt.Printf("Sentences: %v\n", sentences)
	}

	var validSentences int = 0
	for _, sentence := range sentences {
		decomposedSentence, err := decomposeSentence(terms, sentence)
		if err != nil {
			return nil, fmt.Errorf("Error decomposing sentence: %v", err)
		}
		if decomposedSentence == nil {
			fmt.Printf("Sentence %s returned nil terms\n", sentence)
			continue
		}
		if len(decomposedSentence) == 0 {
			fmt.Printf("Sentence %s returned no decomposed terms\n", sentence)
			continue
		}
		if DEBUG {
			fmt.Printf("Decomposed sentence: %v\n", decomposedSentence)
		}
		validSentences++
	}
	result := []string{strconv.Itoa(validSentences)}
	fmt.Printf("Found %d valid sentences\n", validSentences)
	return result, nil
}

func decomposeSentence(terms []Term, sentence Sentence) ([]Term, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	if DEBUG {
		fmt.Printf("Decomposing sentence: %s\n", sentence)
	}

	// If the sentence is empty, return an empty slice
	if len(sentence) == 0 {
		return []Term{}, nil
	}

	// If the sentence is not empty, but the terms are empty, return nil
	if len(terms) == 0 {
		return nil, nil
	}

	var decomposedTerms []Term = make([]Term, 0)

	// Sort the terms largest to smallest so that we can greedily match the largest terms first
	termSort := func(i, j Term) int {
		return len(j) - len(i)
	}
	slices.SortFunc(terms, termSort)

	// Iterate over the sentence and try to match the terms
	for i := 0; i < len(sentence); i++ {
		// Iterate over the terms
		for j := 0; j < len(terms); j++ {
			// If the term matches the sentence, add it to the decomposed terms
			if strings.HasPrefix(string(sentence[i:]), string(terms[j])) {
				decomposedTerms = append(decomposedTerms, terms[j])
				i += len(terms[j]) - 1
				break
			}
		}
	}

	// Verify that the decomposed terms can be recombined to form the original sentence
	var recombinedSentence Sentence
	for _, term := range decomposedTerms {
		recombinedSentence += Sentence(term)
	}
	if recombinedSentence != sentence {
		return nil, nil
	}

	return decomposedTerms, nil
}
