package day17

import (
	"math/big"
	"testing"
)

func TestNewComputer(t *testing.T) {
	comp := NewComputer()
	if comp == nil {
		t.Fatal("expected non-nil Computer instance")
	}
	if len(comp.opcodes) != 0 {
		t.Fatalf("expected empty opcodes, got %d", len(comp.opcodes))
	}
	if comp.a.Cmp(big.NewInt(0)) != 0 || comp.b.Cmp(big.NewInt(0)) != 0 || comp.c.Cmp(big.NewInt(0)) != 0 {
		t.Fatal("expected all registers to be initialized to 0")
	}
}

func TestGetSetInstructionPointer(t *testing.T) {
	comp := NewComputer()
	expected := 5
	comp.SetInstructionPointer(expected)
	if comp.GetInstructionPointer() != expected {
		t.Fatalf("expected instruction pointer to be %d, got %d", expected, comp.GetInstructionPointer())
	}
	comp.SetInstructionPointer(10)
	expected = 10
	if comp.GetInstructionPointer() != expected {
		t.Fatalf("expected instruction pointer to be %d, got %d", expected, comp.GetInstructionPointer())
	}
}

func TestSetOpcodes(t *testing.T) {
	comp := NewComputer()
	opcodes := []Opcode{0, 1, 2, 3, 4, 5, 6, 7}
	err := comp.SetOpcodes(opcodes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comp.opcodes) != len(opcodes) {
		t.Fatalf("expected %d opcodes, got %d", len(opcodes), len(comp.opcodes))
	}
}

func TestSetOpcodesInvalid(t *testing.T) {
	comp := NewComputer()
	opcodes := []Opcode{8}
	err := comp.SetOpcodes(opcodes)
	if err == nil {
		t.Fatal("expected error for invalid opcode, got nil")
	}
}

func TestGetOpcodes(t *testing.T) {
	comp := NewComputer()
	opcodes := []Opcode{0, 1, 2, 3}
	comp.SetOpcodes(opcodes)
	retrieved := comp.GetOpcodes()
	if len(retrieved) != len(opcodes) {
		t.Fatalf("expected %d opcodes, got %d", len(opcodes), len(retrieved))
	}
	for i, opcode := range retrieved {
		if opcode != opcodes[i] {
			t.Fatalf("expected opcode %d at index %d, got %d", opcodes[i], i, opcode)
		}
	}
}

func TestGetSetRegisters(t *testing.T) {
	comp := NewComputer()
	a := big.NewInt(10)
	b := big.NewInt(20)
	c := big.NewInt(30)

	comp.SetRegisterA(a)
	comp.SetRegisterB(b)
	comp.SetRegisterC(c)

	if comp.GetRegisterA().Cmp(a) != 0 {
		t.Fatalf("expected register A to be %s, got %s", a.String(), comp.GetRegisterA().String())
	}
	if comp.GetRegisterB().Cmp(b) != 0 {
		t.Fatalf("expected register B to be %s, got %s", b.String(), comp.GetRegisterB().String())
	}
	if comp.GetRegisterC().Cmp(c) != 0 {
		t.Fatalf("expected register C to be %s, got %s", c.String(), comp.GetRegisterC().String())
	}
}

func TestString(t *testing.T) {
	comp := NewComputer()
	expected := "Computer: A=0 B=0 C=0"
	if comp.String() != expected {
		t.Fatalf("expected %s, got %s", expected, comp.String())
	}
}

func TestOutputChannel(t *testing.T) {
	comp := NewComputer()
	expected := big.NewInt(42)
	go func() {
		comp.Output <- expected
	}()
	result := <-comp.Output
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected output to be %s, got %s", expected.String(), result.String())
	}
}