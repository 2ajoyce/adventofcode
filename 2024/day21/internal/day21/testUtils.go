package day21

import "testing"

func compareMovementStrings(t *testing.T, output string, expectedOutput string) bool {
	if len(output) != len(expectedOutput) {
		t.Errorf("Expected output to be '%d' characters[%s], but got '%d' characters[%s]", len(expectedOutput), expectedOutput, len(output), output)
		return false
	}

	// Validate that every character in the expected output is in the actual output
	// The only exception is the 'A' character, which must be in the specified position
	countUp := 0
	countDown := 0
	countLeft := 0
	countRight := 0
	for i, c := range expectedOutput {
		if c == 'A' {
			if output[i] != 'A' {
				t.Errorf("Expected 'A' at position %d, but got '%c'", i, output[i])
			}
		} else {
			switch c {
			case '^':
				countUp++
			case 'v':
				countDown++
			case '<':
				countLeft++
			case '>':
				countRight++
			}
		}
	}
	for _, c := range output {
		switch c {
		case '^':
			countUp--
		case 'v':
			countDown--
		case '<':
			countLeft--
		case '>':
			countRight--
		}
	}
	if countUp != 0 {
		t.Errorf("Vertical movement error, Expected: %s Received: %s", expectedOutput, output)
		return false
	}
	if countDown != 0 {
		t.Errorf("Vertical movement error, Expected: %s Received: %s", expectedOutput, output)
		return false
	}
	if countLeft != 0 {
		t.Errorf("Horizontal movement error, Expected: %s Received: %s", expectedOutput, output)
		return false
	}
	if countRight != 0 {
		t.Errorf("Horizontal movement error, Expected: %s Received: %s", expectedOutput, output)
		return false
	}
	return true
}

// v<<A>>^A<A>AvA<^AA>A<vAAA>^A
// <v<AA<AAvA<A<AA<vA<vA<vAA