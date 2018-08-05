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

func (cpu *CPU) Step() error {
	opcode := cpu.MMU.Peek(cpu.PC)
	ins, err := cpu.Decode(opcode)
	if err != nil {
		return err
	}

	var l, h byte
	if ins.Length > 1 {
		l = cpu.MMU.Peek(cpu.PC + 1)
	}
	if ins.Length > 2 {
		h = cpu.MMU.Peek(cpu.PC + 2)
	}

	defer func() { fmt.Printf("%-22s %s\n", ins.String(l, h), cpu.String()) }()

	return cpu.Execute(ins, l, h)
}

func (cpu *CPU) Decode(opcode byte) (Instruction, error) {
	bank := Instructions
	if cpu.NextOpcodeIsCB {
		bank = CBInstructions
		defer func() { cpu.NextOpcodeIsCB = false }()
	}

	ins, ok := bank[opcode]
	if !ok {
		return Instruction{}, fmt.Errorf(
			"opcode not found: 0x%02X (CB: %t)",
			opcode,
			cpu.NextOpcodeIsCB,
		)
	}

	return ins, nil
}

func (cpu *CPU) Execute(ins Instruction, l, h byte) error {
	if ins.Func == nil {
		return fmt.Errorf("no function defined for %s", ins.Name)
	}

	cpu.Cycle += int(ins.Cycles)
	cpu.PC += uint16(ins.Length)
	ins.Func(cpu, l, h)

	return nil
}
