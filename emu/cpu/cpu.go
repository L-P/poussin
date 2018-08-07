package cpu

import (
	"errors"
	"fmt"
	"image"

	"home.leo-peltier.fr/poussin/emu/ppu"
	"home.leo-peltier.fr/poussin/emu/rom"
)

type CPU struct {
	Mem  [0xFFFF]byte
	Boot [256]byte
	ROM  [1024 * 1024 * 8]byte // Per Wikipedia, a GB ROM is 8 MiB max

	// For debugging purposes
	PreviousInstruction Instruction
	PreviousLowArg      byte
	PreviousHighArg     byte

	// Adressable at 0xFFFF
	InterruptEnable byte

	// Interrupt Master Flag, not addressable
	InterruptMaster bool

	PPU *ppu.PPU

	// Switching to CB opcode (0xCB for code bank?) is an instruction in
	// itself, when this flag is set it means the next opcode we read will be
	// from the CB opcodes
	NextOpcodeIsCB bool
	Cycle          int

	// Number of OP decoded and executed (debug)
	OPCount int

	// Registers
	A             uint8  // Accumulator
	FlagZero      bool   // True if last result was 0
	FlagSubstract bool   // True if last operation was a substraction / decrement
	FlagHalfCarry bool   // True if the last operation had a half-carry
	FlagCarry     bool   // True if the last operation had a carry
	BC            uint16 // General-purpose register
	DE            uint16 // General-purpose register
	HL            uint16 // General-purpose register that doubles as a faster memory pointer
	SP            uint16 // Stack pointer
	PC            uint16 // Program counter
}

func New(nextFrame chan<- *image.RGBA) CPU {
	return CPU{
		PPU: ppu.New(nextFrame),
	}
}

func (c *CPU) Step() error {
	opcode := c.Fetch(c.PC)
	ins, err := c.Decode(opcode)
	if err != nil {
		fmt.Printf(
			"%-22s %s\n",
			c.PreviousInstruction.String(c.PreviousLowArg, c.PreviousHighArg),
			c.String(),
		)
		return err
	}

	for i := byte(0); i < ins.Cycles; i++ {
		c.PPU.Cycle()
	}

	var l, h byte
	if ins.Length > 1 {
		l = c.Fetch(c.PC + 1)
	}
	if ins.Length > 2 {
		h = c.Fetch(c.PC + 2)
	}

	c.PreviousInstruction = ins
	c.PreviousLowArg = l
	c.PreviousHighArg = h

	/*
		if opcode != 0xCB { // don't clutter with PREFIX CB
			defer func() { fmt.Printf("%-22s %s\n", ins.String(l, h), c.String()) }()
		}
		// */

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
	c.OPCount++
	ins.Func(c, l, h)

	return nil
}

func (c *CPU) LoadBootROM(data []byte) error {
	if count := copy(c.Boot[:], data); count != 256 {
		return fmt.Errorf("did not copy 256 bytes: %d", count)
	}

	return nil
}

func (c *CPU) LoadROM(data []byte) error {
	copy(c.ROM[:], data)

	h := rom.NewHeader(data)
	fmt.Printf("ROM loaded: %s\n", h.String())

	if h.CGBOnly {
		return errors.New("only DMG games are supported")
	}

	return nil
}
