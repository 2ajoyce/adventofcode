package day17

import (
	"errors"
	"fmt"
	"math/big"
)

////////////////////////////////////////
// Computer
////////////////////////////////////////

type Computer struct {
	opcodes            []Opcode
	instructionPointer int // index of current opcode
	a                  *big.Int
	b                  *big.Int
	c                  *big.Int
	Output             chan *big.Int
}

func NewComputer() *Computer {
	return &Computer{
		opcodes:            make([]Opcode, 0),
		instructionPointer: 0,
		a:                  big.NewInt(0),
		b:                  big.NewInt(0),
		c:                  big.NewInt(0),
		Output:             make(chan *big.Int),
	}
}

// // Opcodes // //
func (c *Computer) GetOpcodes() []Opcode {
	return append([]Opcode{}, c.opcodes...)
}

func (c *Computer) SetOpcodes(opcodes []Opcode) error {
	for _, opcode := range opcodes {
		if opcode < 0 || opcode > 7 {
			return errors.New(fmt.Sprintf("invalid opcode: %d", opcode))
		}
	}
	c.opcodes = append(c.opcodes, opcodes...)
	return nil
}

func (c *Computer) GetComboOperand(o Opcode) (*big.Int, error) {
	switch o {
	case 0, 1, 2, 3:
		return big.NewInt(int64(o)), nil
	case 4:
		return big.NewInt(0).Set(c.GetRegisterA()), nil
	case 5:
		return big.NewInt(0).Set(c.GetRegisterB()), nil
	case 6:
		return big.NewInt(0).Set(c.GetRegisterC()), nil
	case 7:
		return nil, errors.New("combo operand 7 is reserved and will never appear in valid programs")
	}
	return nil, errors.New(fmt.Sprintf("invalid opcode: %d", o))
}

// // Instruction Pointer // //
func (c *Computer) GetInstructionPointer() int {
	return c.instructionPointer
}

func (c *Computer) SetInstructionPointer(value int) {
	c.instructionPointer = value
}

// // Get Registers // //
func (c *Computer) GetRegisterA() *big.Int {
	return big.NewInt(0).Set(c.a)
}
func (c *Computer) GetRegisterB() *big.Int {
	return big.NewInt(0).Set(c.b)
}
func (c *Computer) GetRegisterC() *big.Int {
	return big.NewInt(0).Set(c.c)
}

// // Set Registers // //
func (c *Computer) SetRegisterA(value *big.Int) {
	c.a = big.NewInt(0).Set(value)
}
func (c *Computer) SetRegisterB(value *big.Int) {
	c.b = big.NewInt(0).Set(value)
}
func (c *Computer) SetRegisterC(value *big.Int) {
	c.c = big.NewInt(0).Set(value)
}

// // Stringer // //
func (c *Computer) String() string {
	return fmt.Sprintf("Computer: A=%s B=%s C=%s", c.a.String(), c.b.String(), c.c.String())
}
