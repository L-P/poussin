package cpu

import (
	"fmt"

	"home.leo-peltier.fr/poussin/emu/mmu"
)

type CPU struct {
	MMU       mmu.MMU
	Registers Registers

	// Switching to CB opcode (0xCB for code bank?) is an instruction in
	// itself, when this flag is set it means the next opcode we read will be
	// from the CB opcodes and will need to be reset
	NextOpcodeIsCB bool
	Cycle          int
}

// Flags (F 'register')
const (
	// Zero
	FlagZ = 1 << 7

	// Add/Sub
	FlagN = 1 << 6

	// Half Carry
	FlagH = 1 << 5

	// Carry
	FlagC = 1 << 4
)

type Registers struct {
	// Accumulator
	A uint8

	// Flags
	F uint8

	BC uint16
	DE uint16
	HL uint16

	// Stack pointer
	SP uint16

	// Program counter
	PC uint16
}

func (r Registers) String() string {
	flags := [5]byte{'-', '-', '-', '-', 0x00}
	if (r.F & FlagZ) > 0 {
		flags[0] = 'Z'
	}
	if (r.F & FlagN) > 0 {
		flags[1] = 'N'
	}
	if (r.F & FlagH) > 0 {
		flags[2] = 'H'
	}
	if (r.F & FlagC) > 0 {
		flags[3] = 'C'
	}

	return fmt.Sprintf(
		"A:%02X BC:%04X DE:%04X HL:%04X SP:%04X PC:%04X Flags:%s",
		r.A, r.BC, r.DE, r.HL, r.SP, r.PC, flags,
	)
}

func New() CPU {
	return CPU{
		MMU: mmu.New(),
	}
}

func (cpu *CPU) Step() error {
	opcode := cpu.MMU.Peek(cpu.Registers.PC)
	ins, err := cpu.Decode(opcode)
	if err != nil {
		return err
	}

	var l, h byte
	if ins.Length > 1 {
		l = cpu.MMU.Peek(cpu.Registers.PC + 1)
	}
	if ins.Length > 2 {
		h = cpu.MMU.Peek(cpu.Registers.PC + 2)
	}

	defer func() { fmt.Printf("%-22s %s\n", ins.String(l, h), cpu.Registers.String()) }()

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
	cpu.Registers.PC += uint16(ins.Length)
	ins.Func(cpu, l, h)

	return nil
}

func (cpu *CPU) SetFlag(flag uint8) {
	cpu.Registers.F |= flag
}

func (cpu *CPU) ClearFlag(flag uint8) {
	cpu.Registers.F &= ^flag
}

func (cpu *CPU) GetFlag(flag uint8) bool {
	return (cpu.Registers.F & flag) > 0
}
