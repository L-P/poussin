package cpu

var CBInstructions = map[byte]Instruction{
	0x11: {1, 8, "RL C", i_cb_rl_c},
	0x7C: {1, 8, "BIT 7,H", i_cb_bit_x_h(7)},
}

// Rotates C left through Carry flag
func i_cb_rl_c(cpu *CPU, _, _ byte) {
	C := uint8(cpu.BC & 0x00FF)

	oldCarry := uint8(0)
	if cpu.FlagCarry {
		oldCarry = uint8(1)
	}

	cpu.FlagCarry = (C & (1 << 7)) > 0

	C = (C << 1) | oldCarry

	cpu.SetC(C)
	cpu.FlagZero = C == 0
	cpu.FlagSubstract = false
	cpu.FlagHalfCarry = false
}

// Sets flag Z if the nth bit of H is not set
func i_cb_bit_x_h(bit uint) InstructionImplementation {
	return func(cpu *CPU, _, _ byte) {
		cpu.FlagZero = (cpu.HL & (1 << (bit + 8))) == 0
		cpu.FlagSubstract = false
		cpu.FlagHalfCarry = true
	}
}
