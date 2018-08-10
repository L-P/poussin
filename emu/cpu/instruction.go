package cpu

import (
	"fmt"
)

type InstructionImplementation func(*CPU, byte, byte)

type Instruction struct {
	Length byte
	Cycles byte
	Name   string

	// l is the low byte of a 16b argument
	// h is the high byte of a 16b argument
	Func InstructionImplementation
}

func (i Instruction) String(l, h byte) string {
	switch i.Length {
	case 1:
		return i.Name
	case 2:
		return fmt.Sprintf(i.Name, l)
	case 3:
		return fmt.Sprintf(i.Name, h, l)
	}

	panic("unreachable")
}

func (i Instruction) Valid() bool {
	return i.Cycles > 0
}

func Decode(opcode byte, cb bool) Instruction {
	if cb {
		return CBInstructions[opcode]
	}

	return Instructions[opcode]
}
