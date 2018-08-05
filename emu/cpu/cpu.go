package cpu

import (
	"fmt"

	"home.leo-peltier.fr/poussin/emu/mmu"
)

type CPU struct {
	MMU mmu.MMU

	// Switching to CB opcode (0xCB for code bank?) is an instruction in
	// itself, when this flag is set it means the next opcode we read will be
	// from the CB opcodes and will need to be reset
	NextOpcodeIsCB bool
	Cycle          int

	// Registers
	A             uint8 // Accumulator
	FlagZero      bool
	FlagSubstract bool
	FlagHalfCarry bool
	FlagCarry     bool
	BC            uint16
	DE            uint16
	HL            uint16
	SP            uint16 // Stack pointer
	PC            uint16 // Program counter
}

func New() CPU {
	return CPU{
		MMU: mmu.New(),
	}
}

func (c *CPU) Step() error {
	opcode := c.MMU.Peek(c.PC)
	ins, err := c.Decode(opcode)
	if err != nil {
		return err
	}

	var l, h byte
	if ins.Length > 1 {
		l = c.MMU.Peek(c.PC + 1)
	}
	if ins.Length > 2 {
		h = c.MMU.Peek(c.PC + 2)
	}

	defer func() { fmt.Printf("%-22s %s\n", ins.String(l, h), c.String()) }()

	return c.Execute(ins, l, h)
}

func (c *CPU) Decode(opcode byte) (Instruction, error) {
	bank := Instructions
	if c.NextOpcodeIsCB {
		bank = CBInstructions
		defer func() { c.NextOpcodeIsCB = false }()
	}

	ins, ok := bank[opcode]
	if !ok {
		return Instruction{}, fmt.Errorf(
			"opcode not found: 0x%02X (CB: %t)",
			opcode,
			c.NextOpcodeIsCB,
		)
	}

	return ins, nil
}

func (c *CPU) Execute(ins Instruction, l, h byte) error {
	if ins.Func == nil {
		return fmt.Errorf("no function defined for %s", ins.Name)
	}

	c.Cycle += int(ins.Cycles)
	c.PC += uint16(ins.Length)
	ins.Func(c, l, h)

	return nil
}
