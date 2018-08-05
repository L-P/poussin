package cpu

func (cpu *CPU) StackPush16b(data ...uint16) {
	for _, v := range data {
		cpu.StackPush8b(
			uint8((v & 0xFF00 >> 8)),
			uint8(v&0x00FF),
		)
	}
}

func (cpu *CPU) StackPush8b(data ...byte) {
	for _, v := range data {
		cpu.SP--
		cpu.MMU.Set8b(cpu.SP, v)
	}
}
