package main

import (
	"day19/internal/aocUtils"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	bar := progressbar.Default(int64(len(sentences)))

	// Memoization map to store results for sentences
	memo := make(map[Sentence][][]Term)

	var validSentences int = 0
	for _, sentence := range sentences {
		decomposedSentence := decomposeSentence(terms, sentence, memo)
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
		bar.Add(1)
	}
	result := []string{strconv.Itoa(validSentences)}
	fmt.Printf("Found %d valid sentences\n", validSentences)
	return result, nil
}

func decomposeSentence(terms []Term, sentence Sentence, memo map[Sentence][][]Term) [][]Term {
	DEBUG := os.Getenv("DEBUG") == "true"
	if DEBUG {
		fmt.Printf("Decomposing sentence: %s\n", sentence)
	}

	// If the sentence is empty, return an empty slice
	if len(sentence) == 0 {
		if DEBUG {
			fmt.Println("Sentence is empty")
		}
		return [][]Term{}
	}

	// If the sentence is not empty, but the terms are empty, return nil
	if len(terms) == 0 {
		if DEBUG {
			fmt.Printf("Terms are empty for sentence %s\n", sentence)
		}
		memo[sentence] = nil
		return nil
	}

	// Check if the result is already memoized
	if result, found := memo[sentence]; found {
		if DEBUG {
			fmt.Printf("Memoized result for sentence %s: %v\n", sentence, result)
		}
		return result
	}

	// Iterate over the terms to match the sentence
	decompositions := [][]Term{}
	for _, term := range terms {
		if !strings.HasPrefix(string(sentence), string(term)) {
			continue
		}
		if DEBUG {
			fmt.Printf("Found term %s in sentence %s\n", term, sentence)
		}

		// Recursively decompose the remaining sentence
		remainingSentence := Sentence(sentence[len(term):])
		if len(remainingSentence) == 0 {
			if DEBUG {
				fmt.Printf("The sentence %s is fully matched by term %s\n", sentence, term)
			}
			decompositions = append(decompositions, []Term{term})
			continue
		}
		if DEBUG {
			fmt.Printf("Remaining sentence: %s\n", remainingSentence)
		}
		subCompositions := decomposeSentence(terms, remainingSentence, memo)
		if DEBUG {
			fmt.Printf("Remaining sentence %s can be decomposed %d ways: %v\n", remainingSentence, len(subCompositions), subCompositions)
		}
		for _, subComposition := range subCompositions {
			decompositions = append(decompositions, append([]Term{term}, subComposition...))
		}
	}
	// Store the result in the memoization map
	memo[sentence] = decompositions
	if DEBUG {
		fmt.Printf("Sentence %s can be decomposed %d ways: %v\n", sentence, len(memo[sentence]), memo[sentence])
	}
	return memo[sentence]
}
