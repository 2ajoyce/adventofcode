package day17

import (
	"errors"
	"fmt"
	"math/big"
)

// An opcode is a three bit integer (0-7)
type Opcode int

// GetInstruction returns the function associated with the opcode
func (o Opcode) GetInstruction() (func(comp *Computer, operand Opcode) error, error) {
	switch o {
	case 0:
		return adv, nil
	case 1:
		return bxl, nil
	case 2:
		return bst, nil
	case 3:
		return jnz, nil
	case 4:
		return bxc, nil
	case 5:
		return out, nil
	case 6:
		return bdv, nil
	case 7:
		return cdv, nil
	}
	return nil, errors.New(fmt.Sprintf("invalid opcode: %d", o))
}

func adv(comp *Computer, operand Opcode) error {
	comboOperand, err := comp.GetComboOperand(operand)
	if err != nil {
		return fmt.Errorf("error getting combo operand: %v", err)
	}
	numerator := comp.GetRegisterA()
	denominator := big.NewInt(0).Exp(big.NewInt(2), comboOperand, nil)
	quotient := big.NewInt(0).Div(numerator, denominator)
	comp.SetRegisterA(quotient)
	comp.SetInstructionPointer(comp.GetInstructionPointer() + 2)
	return nil
}

func bxl(comp *Computer, operand Opcode) error {
	registerB := comp.GetRegisterB()
	result := big.NewInt(0).Xor(registerB, big.NewInt(int64(operand)))
	comp.SetRegisterB(result)
	comp.SetInstructionPointer(comp.GetInstructionPointer() + 2)
	return nil
}

func bst(comp *Computer, operand Opcode) error {
	comboOperand, err := comp.GetComboOperand(operand)
	if err != nil {
		return fmt.Errorf("error getting combo operand: %v", err)
	}
	value := big.NewInt(0).Mod(comboOperand, big.NewInt(8))
	comp.SetRegisterB(value)
	comp.SetInstructionPointer(comp.GetInstructionPointer() + 2)
	return nil
}

func jnz(comp *Computer, operand Opcode) error {
	if comp.GetRegisterA().Cmp(big.NewInt(0)) == 0 {
		comp.SetInstructionPointer(comp.GetInstructionPointer() + 2)
		return nil
	}
	comp.SetInstructionPointer(int(operand))
	return nil
}

func bxc(comp *Computer, operand Opcode) error {
	registerB := comp.GetRegisterB()
	registerC := comp.GetRegisterC()
	result := big.NewInt(0).Xor(registerB, registerC)
	comp.SetRegisterB(result)
	comp.SetInstructionPointer(comp.GetInstructionPointer() + 2)
	return nil
}

func out(comp *Computer, operand Opcode) error {
	comboOperand, err := comp.GetComboOperand(operand)
	if err != nil {
		return fmt.Errorf("error getting combo operand: %v", err)
	}
	value := big.NewInt(0).Mod(comboOperand, big.NewInt(8))
	comp.Output <- value
	comp.SetInstructionPointer(comp.GetInstructionPointer() + 2)
	return nil
}

func bdv(comp *Computer, operand Opcode) error {
	comboOperand, err := comp.GetComboOperand(operand)
	if err != nil {
		return fmt.Errorf("error getting combo operand: %v", err)
	}
	numerator := comp.GetRegisterA()
	denominator := big.NewInt(0).Exp(big.NewInt(2), comboOperand, nil)
	quotient := big.NewInt(0).Div(numerator, denominator)
	comp.SetRegisterB(quotient)
	comp.SetInstructionPointer(comp.GetInstructionPointer() + 2)
	return nil
}

func cdv(comp *Computer, operand Opcode) error {
	comboOperand, err := comp.GetComboOperand(operand)
	if err != nil {
		return fmt.Errorf("error getting combo operand: %v", err)
	}
	numerator := comp.GetRegisterA()
	denominator := big.NewInt(0).Exp(big.NewInt(2), comboOperand, nil)
	quotient := big.NewInt(0).Div(numerator, denominator)
	comp.SetRegisterC(quotient)
	comp.SetInstructionPointer(comp.GetInstructionPointer() + 2)
	return nil
}
