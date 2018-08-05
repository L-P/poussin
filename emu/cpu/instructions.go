package cpu

var Instructions = map[byte]Instruction{
	0x00: {1, 4, "NOP", i_nop},
	0x0C: {1, 4, "INC C", i_inc_c},
	0x0E: {2, 8, "LD C,$%02X", i_ld_c},
	0x20: {2, 8, "JR NZ,$%02X", i_jr_nz},

	0x11: {3, 12, "LD DE,$%02X%02X", i_ld_de},
	0x21: {3, 12, "LD HL,$%02X%02X", i_ld_hl},
	0x31: {3, 12, "LD SP,$%02X%02X", i_ld_sp},

	0x1A: {1, 8, "LD A,(DE)", i_ld_a_pde},

	0x32: {1, 8, "LDD (HL),A", i_ldd_phl_a},
	0x3E: {2, 8, "LD A,$%02X", i_ld_a},
	0x77: {1, 8, "LD (HL),A", i_ld_phl_a},
	0xAF: {1, 4, "XOR A", i_xor_a},
	0xCB: {1, 4, "PREFIX CB", i_cb},
	0xE0: {2, 12, "LDH ($%02X),A", i_ldh_pn_a},
	0xE2: {1, 8, "LD (C),A", i_ld_c_a},
}

func i_nop(*CPU, byte, byte) {}

// Increment register C
func i_inc_c(cpu *CPU, _, _ byte) {
	C := uint8(cpu.Registers.BC & 0x00FF)

	if (((C & 0xF) + 1) & 0x10) > 0 {
		cpu.SetFlag(FlagH)
	} else {
		cpu.ClearFlag(FlagH)
	}

	C++
	cpu.Registers.BC = (cpu.Registers.BC & 0xFF00) | uint16(C)

	if C == 0 {
		cpu.SetFlag(FlagZ)
	} else {
		cpu.ClearFlag(FlagZ)
	}

	cpu.ClearFlag(FlagN)
}

// Load A into 0xFF00 + C
func i_ld_c_a(cpu *CPU, _, _ byte) {
	addr := 0xFF00 | (cpu.Registers.BC & 0x00FF)
	cpu.MMU.Set8b(addr, cpu.Registers.A)
}

// Load 8b value into C
func i_ld_c(cpu *CPU, l, _ byte) {
	cpu.Registers.BC &= 0xFF00
	cpu.Registers.BC |= uint16(l)
}

// Load 8b value into A
func i_ld_a(cpu *CPU, l, _ byte) {
	cpu.Registers.A = l
}

// Load 16b value into stack pointer
func i_ld_sp(cpu *CPU, l, h byte) {
	cpu.Registers.SP = (uint16(h) << 8) | uint16(l)
}

// Load 16b value into HL register
func i_ld_hl(cpu *CPU, l, h byte) {
	cpu.Registers.HL = (uint16(h) << 8) | uint16(l)
}

// Load 16b value into DE register
func i_ld_de(cpu *CPU, l, h byte) {
	cpu.Registers.DE = (uint16(h) << 8) | uint16(l)
}

// Put A into address pointed by HL and decrement HL
func i_ldd_phl_a(cpu *CPU, l, _ byte) {
	cpu.MMU.Set8b(cpu.Registers.HL, cpu.Registers.A)
	cpu.Registers.HL--
}

// Put A into address pointed by HL
func i_ld_phl_a(cpu *CPU, l, _ byte) {
	cpu.MMU.Set8b(cpu.Registers.HL, cpu.Registers.A)
}

// Put A into address 0xFF00+l
func i_ldh_pn_a(cpu *CPU, l, _ byte) {
	cpu.MMU.Set8b(0xFF00+uint16(l), cpu.Registers.A)
}

// XOR A against itself, effectively clearing it and all flags
func i_xor_a(cpu *CPU, _, _ byte) {
	cpu.Registers.F = 0
}

// Tells our virtual CPU the next instruction is from the CB block
func i_cb(cpu *CPU, _, _ byte) {
	cpu.NextOpcodeIsCB = true
}

// Load the value at address pointed by DE in A
func i_ld_a_pde(cpu *CPU, _, _ byte) {
	cpu.Registers.A = cpu.MMU.Get8b(cpu.Registers.DE)
}

// Jump to signed addr offset if Z flag is not set
func i_jr_nz(cpu *CPU, l, _ byte) {
	if !cpu.GetFlag(FlagZ) {
		addr := int16(cpu.Registers.PC) + int16(int8(l))
		cpu.Registers.PC = uint16(addr)
	}
}
