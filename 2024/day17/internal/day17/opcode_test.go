package day17

import (
	"math/big"
	"testing"
)

func TestAdv(t *testing.T) {
	tests := []struct {
		initialA  int64
		operand   Opcode
		expectedA int64
	}{
		// Keep operand in range 0-3 to avoid combo operand logic
		{10, 2, 2}, // 10 / (2^2) = 2.5 (integer division)
		{16, 2, 4}, // 16 / (2^2) = 4
		{72, 3, 9}, // 72 / (2^3) = 9
	}

	for _, test := range tests {
		comp := NewComputer()
		comp.SetRegisterA(big.NewInt(test.initialA))

		err := adv(comp, test.operand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if comp.GetRegisterA().Cmp(big.NewInt(test.expectedA)) != 0 {
			t.Errorf("expected RegA to be %d, got %d", test.expectedA, comp.GetRegisterA().Int64())
		}
		if comp.GetInstructionPointer() != 2 {
			t.Errorf("expected instruction pointer to be increased to 2, got %d", comp.GetInstructionPointer())
		}
	}
}

func TestBxl(t *testing.T) {
	tests := []struct {
		initialB  int64
		operand   Opcode
		expectedB int64
	}{
		{10, 2, 8},  // 10 ^ 2 = 8
		{15, 1, 14}, // 15 ^ 1 = 14
		{7, 3, 4},   // 7 ^ 3 = 4
	}

	for _, test := range tests {
		comp := NewComputer()
		comp.SetRegisterB(big.NewInt(test.initialB))

		err := bxl(comp, test.operand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if comp.GetRegisterB().Cmp(big.NewInt(test.expectedB)) != 0 {
			t.Errorf("expected RegB to be %d, got %d", test.expectedB, comp.GetRegisterB().Int64())
		}
		if comp.GetInstructionPointer() != 2 {
			t.Errorf("expected instruction pointer to be increased to 2, got %d", comp.GetInstructionPointer())
		}
	}
}

// Keep operand in range 0-3 to avoid combo operand logic
func TestBst(t *testing.T) {
	tests := []struct {
		operand   Opcode
		expectedB int64
	}{
		{0, 0}, // 0 % 8 = 0
		{1, 1}, // 1 % 8 = 1
		{2, 2}, // 2 % 8 = 2
	}

	for _, test := range tests {
		comp := NewComputer()

		err := bst(comp, test.operand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if comp.GetRegisterB().Cmp(big.NewInt(test.expectedB)) != 0 {
			t.Errorf("expected RegB to be %d, got %d", test.expectedB, comp.GetRegisterB().Int64())
		}
		if comp.GetInstructionPointer() != 2 {
			t.Errorf("expected instruction pointer to be increased to 2, got %d", comp.GetInstructionPointer())
		}
	}
}
func TestJnz(t *testing.T) {
	tests := []struct {
		initialA   int64
		operand    Opcode
		expectedIP int
	}{
		{0, 5, 2},  // Register A is 0, IP should increment normally (by 2)
		{10, 3, 3}, // Register A is non-zero, IP should be set to operand
		{7, 7, 7},  // Register A is non-zero, IP should be set to operand
	}

	for _, test := range tests {
		comp := NewComputer()
		comp.SetRegisterA(big.NewInt(test.initialA))
		if comp.GetInstructionPointer() != 0 {
			t.Errorf("expected instruction pointer to be 0, got %d", comp.GetInstructionPointer())
		}

		err := jnz(comp, test.operand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if comp.GetInstructionPointer() != test.expectedIP {
			t.Errorf("expected instruction pointer to be %d, got %d", test.expectedIP, comp.GetInstructionPointer())
		}
	}
}

func TestBxc(t *testing.T) {
	tests := []struct {
		initialB  int64
		initialC  int64
		operand   Opcode
		expectedB int64
	}{
		{10, 2, 0, 8},  // 10 ^ 2 = 8
		{15, 1, 0, 14}, // 15 ^ 1 = 14
		{7, 3, 0, 4},   // 7 ^ 3 = 4
	}

	for _, test := range tests {
		comp := NewComputer()
		comp.SetRegisterB(big.NewInt(test.initialB))
		comp.SetRegisterC(big.NewInt(test.initialC))

		err := bxc(comp, test.operand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if comp.GetRegisterB().Cmp(big.NewInt(test.expectedB)) != 0 {
			t.Errorf("expected RegB to be %d, got %d", test.expectedB, comp.GetRegisterB().Int64())
		}
		if comp.GetInstructionPointer() != 2 {
			t.Errorf("expected instruction pointer to be increased to 2, got %d", comp.GetInstructionPointer())
		}
	}
}
func TestOut(t *testing.T) {
	tests := []struct {
		operand  Opcode
		expected int64
	}{
		// Keep operand in range 0-3 to avoid combo operand logic
		{0, 0}, // 0 % 8 = 0
		{2, 2}, // 2 % 8 = 2
		{3, 3}, // 3 % 8 = 3
	}

	for _, test := range tests {
		comp := NewComputer()
		comp.Output = make(chan *big.Int, 1)

		err := out(comp, test.operand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		output := <-comp.Output
		if output.Cmp(big.NewInt(test.expected)) != 0 {
			t.Errorf("expected output to be %d, got %d", test.expected, output.Int64())
		}
		if comp.GetInstructionPointer() != 2 {
			t.Errorf("expected instruction pointer to be increased to 2, got %d", comp.GetInstructionPointer())
		}
	}
}
func TestBdv(t *testing.T) {
	tests := []struct {
		initialA  int64
		operand   Opcode
		expectedB int64
	}{
		// Keep operand in range 0-3 to avoid combo operand logic
		{10, 2, 2}, // 10 / (2^2) = 2.5 (integer division)
		{16, 2, 4}, // 16 / (2^2) = 4
		{72, 3, 9}, // 72 / (2^3) = 9
	}

	for _, test := range tests {
		comp := NewComputer()
		comp.SetRegisterA(big.NewInt(test.initialA))

		err := bdv(comp, test.operand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if comp.GetRegisterB().Cmp(big.NewInt(test.expectedB)) != 0 {
			t.Errorf("expected RegB to be %d, got %d", test.expectedB, comp.GetRegisterB().Int64())
		}
		if comp.GetInstructionPointer() != 2 {
			t.Errorf("expected instruction pointer to be increased to 2, got %d", comp.GetInstructionPointer())
		}
	}
}
func TestCdv(t *testing.T) {
	tests := []struct {
		initialA  int64
		operand   Opcode
		expectedC int64
	}{
		// Keep operand in range 0-3 to avoid combo operand logic
		{10, 2, 2}, // 10 / (2^2) = 2.5 (integer division)
		{16, 2, 4}, // 16 / (2^2) = 4
		{72, 3, 9}, // 72 / (2^3) = 9
	}

	for _, test := range tests {
		comp := NewComputer()
		comp.SetRegisterA(big.NewInt(test.initialA))

		err := cdv(comp, test.operand)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if comp.GetRegisterC().Cmp(big.NewInt(test.expectedC)) != 0 {
			t.Errorf("expected RegC to be %d, got %d", test.expectedC, comp.GetRegisterC().Int64())
		}
		if comp.GetInstructionPointer() != 2 {
			t.Errorf("expected instruction pointer to be increased to 2, got %d", comp.GetInstructionPointer())
		}
	}
}

func TestsFromProblem(t *testing.T) {
	tests := []struct {
		initialA       int64
		initialB       int64
		initialC       int64
		opcodes        []Opcode
		expectedA      int64
		expectedB      int64
		expectedC      int64
		expectedOutput []int64
	}{
		{ // If register C contains 9, the program 2,6 would set register B to 1.
			initialA:       0,
			initialB:       0,
			initialC:       9,
			opcodes:        []Opcode{2, 6},
			expectedA:      0,
			expectedB:      1,
			expectedC:      9,
			expectedOutput: []int64{},
		},
		{ // If register B contains 29, the program 1,7 would set register B to 26.
			initialA:       0,
			initialB:       29,
			initialC:       0,
			opcodes:        []Opcode{1, 7},
			expectedA:      0,
			expectedB:      26,
			expectedC:      0,
			expectedOutput: []int64{},
		},
		{ // If register B contains 2024 and register C contains 43690, the program 4,0 would set register B to 44354.
			initialA:       0,
			initialB:       2024,
			initialC:       43690,
			opcodes:        []Opcode{4, 0},
			expectedA:      0,
			expectedB:      44354,
			expectedC:      43690,
			expectedOutput: []int64{},
		},
		{ // If register A contains 2024, the program 0,1,5,4,3,0 would output 4,2,5,6,7,7,7,7,3,1,0 and leave 0 in register A.
			initialA:       2024,
			initialB:       0,
			initialC:       0,
			opcodes:        []Opcode{0, 1, 5, 4, 3, 0},
			expectedA:      0,
			expectedB:      0,
			expectedC:      0,
			expectedOutput: []int64{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0},
		},
	}

	for _, test := range tests {
		comp := NewComputer()
		comp.SetRegisterA(big.NewInt(test.initialA))
		comp.SetRegisterB(big.NewInt(test.initialB))
		comp.SetRegisterC(big.NewInt(test.initialC))
		comp.SetOpcodes(test.opcodes)
		comp.Output = make(chan *big.Int, len(test.expectedOutput))

		for comp.GetInstructionPointer() < len(test.opcodes) {
			opcode := test.opcodes[comp.GetInstructionPointer()]
			fn, err := opcode.GetInstruction()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				break
			}
			err = fn(comp, test.opcodes[comp.GetInstructionPointer()+1])
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				break
			}
		}

		if comp.GetRegisterA().Cmp(big.NewInt(test.expectedA)) != 0 {
			t.Errorf("expected RegA to be %d, got %d", test.expectedA, comp.GetRegisterA().Int64())
		}
		if comp.GetRegisterB().Cmp(big.NewInt(test.expectedB)) != 0 {
			t.Errorf("expected RegB to be %d, got %d", test.expectedB, comp.GetRegisterB().Int64())
		}
		if comp.GetRegisterC().Cmp(big.NewInt(test.expectedC)) != 0 {
			t.Errorf("expected RegC to be %d, got %d", test.expectedC, comp.GetRegisterC().Int64())
		}

		close(comp.Output)
		var output []int64
		for out := range comp.Output {
			output = append(output, out.Int64())
		}

		if len(output) != len(test.expectedOutput) {
			t.Errorf("expected output length to be %d, got %d", len(test.expectedOutput), len(output))
		}
		for i, v := range output {
			if v != test.expectedOutput[i] {
				t.Errorf("expected output[%d] to be %d, got %d", i, test.expectedOutput[i], v)
			}
		}
	}
}
