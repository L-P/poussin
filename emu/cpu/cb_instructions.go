package cpu

var CBInstructions = map[byte]Instruction{
	0x11: {1, 8, "RL C", i_cb_rl_c},
	0x7C: {1, 8, "BIT 7,H", i_cb_bit_x_h(7)},
}

// Rotates C left through Carry flag
func i_cb_rl_c(cpu *CPU, _, _ byte) {
	var C byte
	C, cpu.FlagCarry = rotateLeftWithCarry(cpu.GetC(), cpu.FlagCarry)

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
