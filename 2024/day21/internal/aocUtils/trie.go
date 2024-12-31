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

// Find a substitution in the Trie
func (t *Trie) Substitute(input string, fallback func(string) string) string {
	re := regexp.MustCompile(`([\^v<>]*A)`)
	matches := re.FindAllStringIndex(input, -1)
	if matches == nil {
		return input // No matches, return the original input
	}

	var result []byte
	lastIndex := 0

	for _, match := range matches {
		start, end := match[0], match[1]

		// Append unmatched portion
		if start > lastIndex {
			result = append(result, input[lastIndex:start]...)
		}

		// Process the matched substring
		substring := input[start:end]
		node := t.root
		for _, char := range substring {
			charStr := string(char)
			if child, exists := node.children[charStr]; exists {
				node = child
			} else {
				break
			}
		}

		if node.value != "" {
			// Found a replacement
			result = append(result, node.value...)
		} else {
			// Fallback for unmatched substring
			calculated := fallback(substring)
			result = append(result, calculated...)
		}

		lastIndex = end
	}

	// Append any remaining portion of the input string
	if lastIndex < len(input) {
		result = append(result, input[lastIndex:]...)
	}

	return string(result)
}
