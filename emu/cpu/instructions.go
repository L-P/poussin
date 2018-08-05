package cpu

var Instructions = map[byte]Instruction{
	0x00: {1, 4, "NOP", i_nop},
	0x20: {2, 8, "JR NZ,$%02X", i_jrnz},
	0x21: {3, 12, "LD HL,$%02X%02X", i_ldhl},
	0x31: {3, 12, "LD SP,$%02X%02X", i_ldsp},
	0x32: {1, 8, "LD HL,A", i_ldhla},
	0xAF: {1, 4, "XOR A", i_xora},
	0xCB: {1, 4, "PREFIX CB", i_cb},
}

func i_nop(*CPU, byte, byte) {}

// Load 16b value into stack pointer
func i_ldsp(cpu *CPU, l byte, h byte) {
	cpu.Registers.SP = (uint16(h) << 8) | uint16(l)
}

// Load 16b value into HL register
func i_ldhl(cpu *CPU, l, h byte) {
	cpu.Registers.HL = (uint16(h) << 8) | uint16(l)
}

// Put A into address pointed by HL and decrement HL
func i_ldhla(cpu *CPU, l, _ byte) {
	cpu.MMU.Set8b(cpu.Registers.HL, cpu.Registers.A)
	cpu.Registers.HL--
}

// XOR A against itself, effectively clearing it and all flags
func i_xora(cpu *CPU, _, _ byte) {
	cpu.Registers.F = 0
}

// Tells our virtual CPU the next instruction is from the CB block
func i_cb(cpu *CPU, _, _ byte) {
	cpu.NextOpcodeIsCB = true
}

// Jump to signed addr offset if Z flag is not set
func i_jrnz(cpu *CPU, l, _ byte) {
	if !cpu.GetFlag(FlagZ) {
		addr := int16(cpu.Registers.PC) + int16(int8(l))
		cpu.Registers.PC = uint16(addr)
	}
}
