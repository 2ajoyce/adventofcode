package aocUtils

import (
    "fmt"
    "testing"
)

func TestTrieSubstitution(t *testing.T) {
    // Initialize the trie
    trie := NewTrie() // Assuming NewTrie() is the constructor for the trie

    // Add rules to the trie
    fmt.Println("Running test")
    trie.Insert("v<<A", "<vA<AA>>^A")
    trie.Insert("vA", "<vA>^A")
    trie.Insert("<A", "v<<A>>^A")
    trie.Insert("<vA", "v<<A>A>^A")
    trie.Insert("A", "A")
    fmt.Println("Rules added to the trie")

    // Define a fallback function
    fallback := func(substring string) string {
        result := ""
        for _, char := range substring {
            if char == 'A' {
                result += "B"
            } else {
                result += "*"
            }
        }
        return result
    }

    // Input and substitution
    input := "v<<A>A>^A"
    output := trie.Substitute(input, fallback)
    fmt.Printf("Input: %s, Output: %s\n", input, output)

    expectedOutput := "<vA<AA>>^A*B**B"
    if output == expectedOutput {
        fmt.Println("Test passed: Output matches expected output")
    } else {
        t.Errorf("Test failed: Expected '%s', got '%s'", expectedOutput, output)
    }
}

func TestTrieSubstitutionDouble(t *testing.T) {
    // Initialize the trie
    trie := NewTrie() // Assuming NewTrie() is the constructor for the trie

    // Add rules to the trie
    fmt.Println("Running test")
    trie.Insert("v<<A", "<vA<AA>>^A")
    trie.Insert("vA", "<vA>^A")
    trie.Insert("<A", "v<<A>>^A")
    trie.Insert("<vA", "v<<A>A>^A")
    trie.Insert("A", "A")
    fmt.Println("Rules added to the trie")

    // Define a fallback function
    fallback := func(substring string) string {
        result := ""
        for _, char := range substring {
            if char == 'A' {
                result += "B"
            } else {
                result += "*"
            }
        }
        return result
    }

    // Input and substitution
    input := "v<<Av<<A"
    output := trie.Substitute(input, fallback)
    fmt.Printf("Input: %s, Output: %s\n", input, output)

    expectedOutput := "<vA<AA>>^A<vA<AA>>^A"
    if output == expectedOutput {
        fmt.Println("Test passed: Output matches expected output")
    } else {
        t.Errorf("Test failed: Expected '%s', got '%s'", expectedOutput, output)
    }
}
