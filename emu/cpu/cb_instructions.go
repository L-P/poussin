package cpu

var CBInstructions = map[byte]Instruction{
	0x11: {1, 8, "RL C", i_cb_rl_c},
	0x7C: {1, 8, "BIT 7,H", i_cb_bit_x_h(7)},
}

// Rotates C left through Carry flag
func i_cb_rl_c(cpu *CPU, _, _ byte) {
	C := uint8(cpu.Registers.BC & 0x00FF)

	oldCarry := uint8(0)
	if cpu.GetFlag(FlagC) {
		oldCarry = uint8(1)
	}

	if (C & (1 << 7)) > 0 {
		cpu.SetFlag(FlagC)
	} else {
		cpu.ClearFlag(FlagC)
	}

	C = (C << 1) | oldCarry
	cpu.Registers.BC &= 0xFF00
	cpu.Registers.BC |= uint16(C)

	if C == 0 {
		cpu.SetFlag(FlagZ)
	} else {
		cpu.ClearFlag(FlagZ)
	}

	cpu.ClearFlag(FlagN)
	cpu.ClearFlag(FlagH)
}

// Sets flag Z if the nth bit of H is not set
func i_cb_bit_x_h(bit uint) InstructionImplementation {
	return func(cpu *CPU, _, _ byte) {
		if (cpu.Registers.HL & (1 << (bit + 8))) > 0 {
			cpu.ClearFlag(FlagZ)
		} else {
			cpu.SetFlag(FlagZ)
		}

		cpu.ClearFlag(FlagN)
		cpu.SetFlag(FlagH)
	}
}
