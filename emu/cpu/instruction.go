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

func (ins Instruction) String(l, h byte) string {
	switch ins.Length {
	case 1:
		return ins.Name
	case 2:
		return fmt.Sprintf(ins.Name, l)
	case 3:
		return fmt.Sprintf(ins.Name, h, l)
	}

	panic("unreachable")
}
