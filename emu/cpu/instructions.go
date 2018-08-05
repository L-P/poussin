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
func i_inc_c(c *CPU, _, _ byte) {
	C := c.GetC()

	c.FlagHalfCarry = (((C & 0xF) + 1) & 0x10) > 0

	C++
	c.SetC(C)
	c.FlagZero = C == 0
	c.FlagSubstract = false
}

// Load A into 0xFF00 + C
func i_ld_pc_a(c *CPU, _, _ byte) {
	c.MMU.Set8b(0xFF00|uint16(c.GetC()), c.A)
}

// Load 8b value into C
func i_ld_c(c *CPU, l, _ byte) {
	c.SetC(l)
}

// Load 8b value into B
func i_ld_b(c *CPU, l, _ byte) {
	c.SetB(l)
}

// Load 8b value into A
func i_ld_a(c *CPU, l, _ byte) {
	c.A = l
}

// Load 16b value into stack pointer
func i_ld_sp(c *CPU, l, h byte) {
	c.SP = (uint16(h) << 8) | uint16(l)
}

// Load 16b value into HL register
func i_ld_hl(c *CPU, l, h byte) {
	c.HL = (uint16(h) << 8) | uint16(l)
}

// Load 16b value into DE register
func i_ld_de(c *CPU, l, h byte) {
	c.DE = (uint16(h) << 8) | uint16(l)
}

// Put A into address pointed by HL and decrement HL
func i_ldd_phl_a(c *CPU, l, _ byte) {
	c.MMU.Set8b(c.HL, c.A)
	c.HL--
}

// Put A into address pointed by HL
func i_ld_phl_a(c *CPU, l, _ byte) {
	c.MMU.Set8b(c.HL, c.A)
}

// Put A into address 0xFF00+l
func i_ldh_pn_a(c *CPU, l, _ byte) {
	c.MMU.Set8b(0xFF00+uint16(l), c.A)
}

// XOR A against itself, effectively clearing it and all flags
func i_xor_a(c *CPU, _, _ byte) {
	c.ClearFlags()
}

// Tells our virtual CPU the next instruction is from the CB block
func i_prefix_cb(c *CPU, _, _ byte) {
	c.NextOpcodeIsCB = true
}

// Pushes the address of the next instruction onto the stack and jump
func i_call(c *CPU, l, h byte) {
	c.StackPush16b(c.PC)
	c.PC = (uint16(h) << 8) | uint16(l)
}

// Pushes BC to the stack
func i_push_bc(c *CPU, l, h byte) {
	c.StackPush16b(c.BC)
}

// Load the value at address pointed by DE in A
func i_ld_a_pde(c *CPU, _, _ byte) {
	c.A = c.MMU.Get8b(c.DE)
}

// Load the value of C into A
func i_ld_c_a(c *CPU, _, _ byte) {
	c.A = c.GetC()
}

// Jump to signed addr offset if Z flag is not set
func i_jr_nz(c *CPU, l, _ byte) {
	if !c.FlagZero {
		addr := int16(c.PC) + int16(int8(l))
		c.PC = uint16(addr)
	}
}

func i_rla(c *CPU, _, _ byte) {
	c.A, c.FlagCarry = rotateLeftWithCarry(c.A, c.FlagCarry)

	// GBCPUman says the Z flag depends on the final value, other sources says
	// RLA always clear the flags, no sure who to trust.
	c.FlagZero = false
	c.FlagSubstract = false
	c.FlagHalfCarry = false
}
