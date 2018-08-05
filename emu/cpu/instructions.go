package cpu

var Instructions = map[byte]Instruction{
	0x00: {1, 4, "NOP", i_nop},

	0x0C: {1, 4, "INC C", i_inc_c},
	0xAF: {1, 4, "XOR A", i_xor_a},
	0x17: {1, 4, "RLA", i_rla},

	0x06: {2, 8, "LD B,$%02X", i_ld_b},
	0x0E: {2, 8, "LD C,$%02X", i_ld_c},
	0x3E: {2, 8, "LD A,$%02X", i_ld_a},
	0x1A: {1, 8, "LD A,(DE)", i_ld_a_pde},

	0x4F: {1, 4, "LD C,A", i_ld_c_a},

	0x11: {3, 12, "LD DE,$%02X%02X", i_ld_de},
	0x21: {3, 12, "LD HL,$%02X%02X", i_ld_hl},
	0x31: {3, 12, "LD SP,$%02X%02X", i_ld_sp},

	0xE2: {1, 8, "LD (C),A", i_ld_pc_a},
	0x77: {1, 8, "LD (HL),A", i_ld_phl_a},
	0xE0: {2, 12, "LDH ($%02X),A", i_ldh_pn_a},
	0x32: {1, 8, "LDD (HL),A", i_ldd_phl_a},

	0x20: {2, 8, "JR NZ,$%02X", i_jr_nz},

	0xC5: {1, 16, "PUSH BC", i_push_bc},
	0xCB: {1, 4, "PREFIX CB", i_prefix_cb},
	0xCD: {3, 12, "CALL $%02X%02X", i_call},
}

func i_nop(*CPU, byte, byte) {}

// Increment register C
func i_inc_c(cpu *CPU, _, _ byte) {
	C := cpu.GetC()

	cpu.FlagHalfCarry = (((C & 0xF) + 1) & 0x10) > 0

	C++
	cpu.SetC(C)
	cpu.FlagZero = C == 0
	cpu.FlagSubstract = false
}

// Load A into 0xFF00 + C
func i_ld_pc_a(cpu *CPU, _, _ byte) {
	cpu.MMU.Set8b(0xFF00|uint16(cpu.GetC()), cpu.A)
}

// Load 8b value into C
func i_ld_c(cpu *CPU, l, _ byte) {
	cpu.SetC(l)
}

// Load 8b value into B
func i_ld_b(cpu *CPU, l, _ byte) {
	cpu.SetB(l)
}

// Load 8b value into A
func i_ld_a(cpu *CPU, l, _ byte) {
	cpu.A = l
}

// Load 16b value into stack pointer
func i_ld_sp(cpu *CPU, l, h byte) {
	cpu.SP = (uint16(h) << 8) | uint16(l)
}

// Load 16b value into HL register
func i_ld_hl(cpu *CPU, l, h byte) {
	cpu.HL = (uint16(h) << 8) | uint16(l)
}

// Load 16b value into DE register
func i_ld_de(cpu *CPU, l, h byte) {
	cpu.DE = (uint16(h) << 8) | uint16(l)
}

// Put A into address pointed by HL and decrement HL
func i_ldd_phl_a(cpu *CPU, l, _ byte) {
	cpu.MMU.Set8b(cpu.HL, cpu.A)
	cpu.HL--
}

// Put A into address pointed by HL
func i_ld_phl_a(cpu *CPU, l, _ byte) {
	cpu.MMU.Set8b(cpu.HL, cpu.A)
}

// Put A into address 0xFF00+l
func i_ldh_pn_a(cpu *CPU, l, _ byte) {
	cpu.MMU.Set8b(0xFF00+uint16(l), cpu.A)
}

// XOR A against itself, effectively clearing it and all flags
func i_xor_a(cpu *CPU, _, _ byte) {
	cpu.ClearFlags()
}

// Tells our virtual CPU the next instruction is from the CB block
func i_prefix_cb(cpu *CPU, _, _ byte) {
	cpu.NextOpcodeIsCB = true
}

// Pushes the address of the next instruction onto the stack and jump
func i_call(cpu *CPU, l, h byte) {
	cpu.StackPush16b(cpu.PC)
	cpu.PC = (uint16(h) << 8) | uint16(l)
}

// Pushes BC to the stack
func i_push_bc(cpu *CPU, l, h byte) {
	cpu.StackPush16b(cpu.BC)
}

// Load the value at address pointed by DE in A
func i_ld_a_pde(cpu *CPU, _, _ byte) {
	cpu.A = cpu.MMU.Get8b(cpu.DE)
}

// Load the value of C into A
func i_ld_c_a(cpu *CPU, _, _ byte) {
	cpu.A = cpu.GetC()
}

// Jump to signed addr offset if Z flag is not set
func i_jr_nz(cpu *CPU, l, _ byte) {
	if !cpu.FlagZero {
		addr := int16(cpu.PC) + int16(int8(l))
		cpu.PC = uint16(addr)
	}
}

func i_rla(cpu *CPU, _, _ byte) {
	oldCarry := uint8(0)
	if cpu.FlagCarry {
		oldCarry = uint8(1)
	}

	cpu.FlagCarry = (cpu.A & (1 << 7)) > 0

	cpu.A = (cpu.A << 1) | oldCarry

	// GBCPUman says the flag depends on the final value, other sources says
	// RLA always clear the flags, no sure who to trust.
	cpu.FlagZero = false
	cpu.FlagSubstract = false
	cpu.FlagHalfCarry = false
}
