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

type Term string
type Sentence string

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
}

func (node *TrieNode) Insert(term string) {
	current := node
	for _, char := range term {
		if _, exists := current.children[char]; !exists {
			current.children[char] = NewTrieNode()
		}
		current = current.children[char]
	}
	current.isEnd = true
}

func (node *TrieNode) FindPrefixes(sentence string, start int) []int {
	prefixes := []int{}
	current := node
	for i := start; i < len(sentence); i++ {
		char := rune(sentence[i])
		if child, exists := current.children[char]; exists {
			current = child
			if current.isEnd {
				// i+1 because slicing is non-inclusive at the end index
				prefixes = append(prefixes, i+1)
			}
		} else {
			break
		}
	}
	return prefixes
}

func solve(terms []Term, sentences []Sentence) ([]string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Beginning solve...")
	bar := progressbar.Default(int64(len(sentences)))

	if DEBUG {
		fmt.Printf("Terms: %v\n", terms)
		fmt.Printf("Sentences: %v\n", sentences)
	}

	// Build the Trie with all terms
	trie := NewTrieNode()
	for _, term := range terms {
		trie.Insert(string(term))
		if DEBUG {
			fmt.Printf("Inserted term into Trie: %s\n", term)
		}
	}

	var validSentences int = 0
	var totalCombinations int = 0
	for _, sentence := range sentences {
		if DEBUG {
			fmt.Printf("Processing sentence: %s\n", sentence)
		}
		if canDecompose(string(sentence), trie) {
			if DEBUG {
				fmt.Printf("Sentence '%s' is valid.\n", sentence)
			}
			validSentences++
			count := countDecompositions(string(sentence), trie)
			totalCombinations += count
			if DEBUG {
				fmt.Printf("Sentence '%s' can be decomposed in %d ways.\n", sentence, count)
			}
		} else {
			if DEBUG {
				fmt.Printf("Sentence '%s' is invalid.\n", sentence)
			}
		}
		bar.Add(1)
	}

	result := []string{strconv.Itoa(validSentences)}
	fmt.Printf("Found %d valid sentences with %d combinations\n", validSentences, totalCombinations)
	return result, nil
}

func canDecompose(sentence string, trie *TrieNode) bool {
	DEBUG := os.Getenv("DEBUG") == "true"
	n := len(sentence)
	if n == 0 {
		if DEBUG {
			fmt.Println("Encountered empty sentence during decomposition.")
		}
		return true
	}

	// dp[i] is true if sentence[0:i] can be decomposed into valid terms
	dp := make([]bool, n+1)
	dp[0] = true // Empty string

	for i := 0; i < n; i++ {
		if !dp[i] {
			continue
		}
		prefixEndIndices := trie.FindPrefixes(sentence, i)
		if DEBUG && len(prefixEndIndices) > 0 {
			fmt.Printf("At index %d, found prefixes ending at indices: %v\n", i, prefixEndIndices)
		}
		for _, end := range prefixEndIndices {
			if DEBUG {
				fmt.Printf("Marking dp[%d] as true because sentence[%d:%d] is a valid term.\n", end, i, end)
			}
			dp[end] = true
		}
	}

	if DEBUG {
		fmt.Printf("DP Array for sentence '%s': %v\n", sentence, dp)
	}

	return dp[n]
}

// countDecompositions returns the number of unique ways to decompose the sentence into terms
func countDecompositions(sentence string, trie *TrieNode) int {
	DEBUG := os.Getenv("DEBUG") == "true"
	n := len(sentence)
	if n == 0 {
		if DEBUG {
			fmt.Println("Encountered empty sentence during decomposition.")
		}
		return 1 // One way to decompose an empty string
	}

	// dp[i] represents the number of ways to decompose sentence[0:i]
	dp := make([]int, n+1)
	dp[0] = 1 // Base case: one way to decompose an empty string

	for i := 0; i < n; i++ {
		if dp[i] == 0 {
			continue // No valid decompositions ending at i
		}
		prefixEndIndices := trie.FindPrefixes(sentence, i)
		if DEBUG && len(prefixEndIndices) > 0 {
			fmt.Printf("At index %d, found prefixes ending at indices: %v\n", i, prefixEndIndices)
		}
		for _, end := range prefixEndIndices {
			dp[end] += dp[i]
			if DEBUG {
				fmt.Printf("Adding %d to dp[%d], total now: %d\n", dp[i], end, dp[end])
			}
		}
	}

	if DEBUG {
		fmt.Printf("DP Array for sentence '%s': %v\n", sentence, dp)
	}

	return dp[n]
}

// decomposeSentence decomposes the sentence into all possible sequences of terms
// This function does not functionally work at large scale. It is too slow.
func decomposeSentence(sentence string, trie *TrieNode, memo map[string][][]Term) [][]Term {
	DEBUG := os.Getenv("DEBUG") == "true"
	if DEBUG {
		fmt.Printf("Decomposing sentence: %s\n", sentence)
	}

	// If the sentence is empty, return an empty slice indicating a valid decomposition
	if len(sentence) == 0 {
		if DEBUG {
			fmt.Println("Sentence is empty")
		}
		return [][]Term{}
	}

	// Check if the result is already memoized
	if result, found := memo[sentence]; found {
		if DEBUG {
			fmt.Printf("Returning memoized result for sentence '%s': %v\n", sentence, result)
		}
		return result
	}

	decompositions := [][]Term{}

	// Iterate over the sentence to find all prefix matches using the Trie
	prefixEndIndices := trie.FindPrefixes(sentence, 0)
	if DEBUG && len(prefixEndIndices) > 0 {
		fmt.Printf("Found prefixes in sentence '%s' ending at indices: %v\n", sentence, prefixEndIndices)
	}

	for _, end := range prefixEndIndices {
		currentTerm := Term(sentence[:end])
		if DEBUG {
			fmt.Printf("Found term '%s' in sentence '%s'\n", currentTerm, sentence)
		}

		remainingSentence := sentence[end:]
		if len(remainingSentence) == 0 {
			if DEBUG {
				fmt.Printf("The sentence '%s' is fully matched by term '%s'\n", sentence, currentTerm)
			}
			decompositions = append(decompositions, []Term{currentTerm})
			continue
		}
		if DEBUG {
			fmt.Printf("Remaining sentence after term '%s': %s\n", currentTerm, remainingSentence)
		}

		subCompositions := decomposeSentence(remainingSentence, trie, memo)
		if subCompositions == nil {
			if DEBUG {
				fmt.Printf("No decompositions found for remaining sentence '%s'\n", remainingSentence)
			}
			continue
		}

		for _, subComposition := range subCompositions {
			combined := append([]Term{currentTerm}, subComposition...)
			decompositions = append(decompositions, combined)
			if DEBUG {
				fmt.Printf("Combined decomposition: %v\n", combined)
			}
		}
	}

	// Memoize the result
	memo[sentence] = decompositions
	if DEBUG {
		fmt.Printf("Sentence '%s' can be decomposed in %d ways: %v\n", sentence, len(decompositions), decompositions)
	}

	return decompositions
}
