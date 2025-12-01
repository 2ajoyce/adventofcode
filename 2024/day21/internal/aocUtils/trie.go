package aocUtils

import (
	"regexp"
)

type TrieNode struct {
	value    string
	children map[string]*TrieNode
}

type Trie struct {
	root *TrieNode
}

// Create a new Trie
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[string]*TrieNode),
		},
	}
}

// Insert a rule into the Trie
func (t *Trie) Insert(key, value string) {
	node := t.root
	for _, char := range key {
		charStr := string(char)
		if _, exists := node.children[charStr]; !exists {
			node.children[charStr] = &TrieNode{
				children: make(map[string]*TrieNode),
			}
		}
		node = node.children[charStr]
	}
	node.value = value
}

func (t *Trie) Substitute(input string, fallback func(string) string) string {
	re := regexp.MustCompile(`([\^v<>]*A)`)
	matches := re.FindAllStringIndex(input, -1)
	if matches == nil {
		return input // No matches, return the original input
	}

	var result []byte
	lastIndex := 0

	for len(matches) > 0 {
		found := false
		// Try combinations of matches in decreasing size
		for size := len(matches); size > 0 && !found; size-- {
			for i := 0; i <= len(matches)-size; i++ {
				// Create a subset of matches
				subset := matches[i : i+size]

				// Extract the corresponding substring
				substring := input[subset[0][0]:subset[len(subset)-1][1]]

				// Attempt to find a substitution for the subset
				node := t.root
				for _, char := range substring {
					charStr := string(char)
					if child, exists := node.children[charStr]; exists {
						node = child
					} else {
						node = nil // This is not a valid path
						break
					}
				}

				// Ensure we've reached a terminal node with a value
				if node != nil && node.value != "" {
					// Found a replacement for this subset
					// Append unmatched portion before the subset
					if subset[0][0] > lastIndex {
						result = append(result, input[lastIndex:subset[0][0]]...)
					}
					// Append the substitution value
					result = append(result, node.value...)
					// Update lastIndex and remove processed matches
					lastIndex = subset[len(subset)-1][1]
					matches = matches[i+size:]
					found = true
					break
				}
			}
		}

		if !found {
			// No subset matched; process the first match using fallback
			firstMatch := input[matches[0][0]:matches[0][1]]
			if matches[0][0] > lastIndex {
				result = append(result, input[lastIndex:matches[0][0]]...)
			}
			calculated := fallback(firstMatch)
			result = append(result, calculated...)
			lastIndex = matches[0][1]
			matches = matches[1:]
		}
	}

	// Append any remaining portion of the input string
	if lastIndex < len(input) {
		result = append(result, input[lastIndex:]...)
	}

	return string(result)
}
