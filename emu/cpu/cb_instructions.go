package cpu

var CBInstructions = map[byte]Instruction{
	0x7C: {1, 8, "BIT 7,H", i_cb_bitxh(7)},
}

func i_cb_bitxh(bit uint) InstructionImplementation {
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
