package cpu

var Instructions = map[byte]Instruction{
	0x00: {1, 4, "NOP", i_nop},

	0x3C: {1, 4, "INC A", i_inc_a},
	0x04: {1, 4, "INC B", i_inc_n('B')},
	0x0C: {1, 4, "INC C", i_inc_n('C')},
	0x14: {1, 4, "INC D", i_inc_n('D')},
	0x1C: {1, 4, "INC E", i_inc_n('E')},
	0x24: {1, 4, "INC H", i_inc_n('H')},
	0x2C: {1, 4, "INC L", i_inc_n('L')},

	0x13: {1, 8, "INC DE", i_inc_de},
	0x23: {1, 8, "INC HL", i_inc_hl},

	0x3D: {1, 4, "DEC A", i_dec_a},
	0x05: {1, 4, "DEC B", i_dec_n('B')},
	0x0D: {1, 4, "DEC C", i_dec_n('C')},
	0x15: {1, 4, "DEC D", i_dec_n('D')},
	0x1D: {1, 4, "DEC E", i_dec_n('E')},
	0x25: {1, 4, "DEC H", i_dec_n('H')},
	0x2D: {1, 4, "DEC L", i_dec_n('L')},

	0x86: {1, 8, "ADD A,(HL)", i_add_a_phl},

	0x97: {1, 4, "SUB A", i_sub_n('A')},
	0x90: {1, 4, "SUB B", i_sub_n('B')},
	0x91: {1, 4, "SUB C", i_sub_n('C')},
	0x92: {1, 4, "SUB D", i_sub_n('D')},
	0x93: {1, 4, "SUB E", i_sub_n('E')},
	0x94: {1, 4, "SUB H", i_sub_n('H')},
	0x95: {1, 4, "SUB L", i_sub_n('L')},

	0xAF: {1, 4, "XOR A", i_xor_a},
	0x17: {1, 4, "RLA", i_rla},

	0x3E: {2, 8, "LD A,$%02X", i_ld_a},
	0x06: {2, 8, "LD B,$%02X", i_ld_n('B')},
	0x0E: {2, 8, "LD C,$%02X", i_ld_n('C')},
	0x16: {2, 8, "LD D,$%02X", i_ld_n('D')},
	0x1E: {2, 8, "LD E,$%02X", i_ld_n('E')},
	0x26: {2, 8, "LD H,$%02X", i_ld_n('H')},
	0x2E: {2, 8, "LD L,$%02X", i_ld_n('L')},

	0x1A: {1, 8, "LD A,(DE)", i_ld_a_pde},

	0x7F: {1, 4, "LD A,A", i_ld_n_n('A', 'A')},
	0x78: {1, 4, "LD A,B", i_ld_n_n('A', 'B')},
	0x79: {1, 4, "LD A,C", i_ld_n_n('A', 'C')},
	0x7A: {1, 4, "LD A,D", i_ld_n_n('A', 'D')},
	0x7B: {1, 4, "LD A,E", i_ld_n_n('A', 'E')},
	0x7C: {1, 4, "LD A,H", i_ld_n_n('A', 'H')},
	0x7D: {1, 4, "LD A,L", i_ld_n_n('A', 'L')},

	0x4F: {1, 4, "LD C,A", i_ld_n_n('C', 'A')},
	0x57: {1, 4, "LD D,A", i_ld_n_n('D', 'A')},
	0x67: {1, 4, "LD H,A", i_ld_n_n('H', 'A')},

	0x11: {3, 12, "LD DE,$%02X%02X", i_ld_de},
	0x21: {3, 12, "LD HL,$%02X%02X", i_ld_hl},
	0x31: {3, 12, "LD SP,$%02X%02X", i_ld_sp},

	0xE2: {1, 8, "LD (C),A", i_ld_pc_a},
	0xEA: {3, 16, "LD ($%02X%02X),A", i_ld_pn_a},
	0x77: {1, 8, "LD (HL),A", i_ld_phl_a},

	0x22: {1, 8, "LDI (HL),A", i_ldi_phl_a},
	0x32: {1, 8, "LDD (HL),A", i_ldd_phl_a},

	0xE0: {2, 12, "LDH ($%02X),A", i_ldh_pn_a},
	0xF0: {2, 12, "LDH A,($%02X)", i_ldh_a_pn},

	0xBE: {1, 8, "CP (HL)", i_cp_phl},
	0xFE: {2, 8, "CP $%02X", i_cp_n},

	0x18: {2, 8, "JR,$%02X", i_jr},
	0x20: {2, 8, "JR NZ,$%02X", i_jr_nz},
	0x28: {2, 8, "JR Z,$%02X", i_jr_z},

	0xC1: {1, 12, "POP BC", i_pop_bc},
	0xC5: {1, 16, "PUSH BC", i_push_bc},
	0xCB: {1, 4, "PREFIX CB", i_prefix_cb},
	0xCD: {3, 12, "CALL $%02X%02X", i_call},
	0xC9: {1, 8, "RET", i_ret},
}

func i_nop(*CPU, byte, byte) {}

func i_sub_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, _ := c.GetRegisterCallbacks(name)

		b := c.A - get()

		c.FlagHalfCarry = (c.A & 0xF) < (get() & 0xF)
		c.FlagCarry = b > c.A

		c.A = b
		c.FlagZero = b == 0
		c.FlagSubstract = true
	}
}

// Increments register B
func i_inc_b(c *CPU, _, _ byte) {
	var B byte

	B, c.FlagHalfCarry = increment(c.GetB())

	c.SetB(B)
	c.FlagZero = B == 0
	c.FlagSubstract = false
}

// Increments register C
func i_inc_c(c *CPU, _, _ byte) {
	var C byte

	C, c.FlagHalfCarry = increment(c.GetC())

	c.SetC(C)
	c.FlagZero = C == 0
	c.FlagSubstract = false
}

// Increments register DE
func i_inc_de(c *CPU, _, _ byte) {
	c.DE++
}

// Increments register HL
func i_inc_hl(c *CPU, _, _ byte) {
	c.HL++
}

// Decrements register n
func i_dec_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)
		var b byte

		b, c.FlagHalfCarry = decrement(get())

		set(b)
		c.FlagZero = b == 0
		c.FlagSubstract = true
	}
}

// Increments register n
func i_inc_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)
		var b byte

		b, c.FlagHalfCarry = increment(get())

		set(b)
		c.FlagZero = b == 0
		c.FlagSubstract = false
	}
}

// Increments register A
func i_inc_a(c *CPU, _, _ byte) {
	c.A, c.FlagHalfCarry = increment(c.A)
	c.FlagZero = c.A == 0
	c.FlagSubstract = false
}

// Decrements register A
func i_dec_a(c *CPU, _, _ byte) {
	c.A, c.FlagHalfCarry = decrement(c.A)
	c.FlagZero = c.A == 0
	c.FlagSubstract = true
}

// Loads A into 0xFF00 + C
func i_ld_pc_a(c *CPU, _, _ byte) {
	c.MMU.Set8b(0xFF00|uint16(c.GetC()), c.A)
}

// Loads 8b value into n
func i_ld_n(name byte) InstructionImplementation {
	return func(c *CPU, l, _ byte) {
		_, set := c.GetRegisterCallbacks(name)
		set(l)
	}
}

// Loads 8b value into A
func i_ld_a(c *CPU, l, _ byte) {
	c.A = l
}

// Loads 16b value into stack pointer
func i_ld_sp(c *CPU, l, h byte) {
	c.SP = (uint16(h) << 8) | uint16(l)
}

// Loads 16b value into HL register
func i_ld_hl(c *CPU, l, h byte) {
	c.HL = (uint16(h) << 8) | uint16(l)
}

// Loads 16b value into DE register
func i_ld_de(c *CPU, l, h byte) {
	c.DE = (uint16(h) << 8) | uint16(l)
}

// Puts A into address pointed by HL and decrement HL
func i_ldd_phl_a(c *CPU, l, _ byte) {
	c.MMU.Set8b(c.HL, c.A)
	c.HL--
}

// Puts A into address pointed by HL and increment HL
func i_ldi_phl_a(c *CPU, l, _ byte) {
	c.MMU.Set8b(c.HL, c.A)
	c.HL++
}

// Puts A into address pointed by HL
func i_ld_phl_a(c *CPU, l, _ byte) {
	c.MMU.Set8b(c.HL, c.A)
}

// Puts A into given address pointed by HL
func i_ld_pn_a(c *CPU, l, h byte) {
	c.MMU.Set8b(uint16(l)|(uint16(h)<<8), c.A)
}

// Puts A into address 0xFF00+l
func i_ldh_pn_a(c *CPU, l, _ byte) {
	c.MMU.Set8b(0xFF00+uint16(l), c.A)
}

// Puts value at 0xFF00+l into A
func i_ldh_a_pn(c *CPU, l, _ byte) {
	c.A = c.MMU.Get8b(0xFF00 + uint16(l))
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

// Pops a two bytes address stack and jump to it
func i_ret(c *CPU, _, _ byte) {
	c.PC = c.StackPop16b()
}

// Pushes BC to the stack
func i_push_bc(c *CPU, l, h byte) {
	c.StackPush16b(c.BC)
}

// Pops two bytes from the stack to BC
func i_pop_bc(c *CPU, l, h byte) {
	c.BC = c.StackPop16b()
}

// Loads the value at address pointed by DE in A
func i_ld_a_pde(c *CPU, _, _ byte) {
	c.A = c.MMU.Get8b(c.DE)
}

// Loads the value of register n into n
func i_ld_n_n(dst, src byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		_, set := c.GetRegisterCallbacks(dst)
		get, _ := c.GetRegisterCallbacks(src)
		set(get())
	}
}

// Jumps to the given address offset
func i_jr(c *CPU, l, _ byte) {
	c.PC = signedOffset(c.PC, l)
}

// Jumps to signed addr offset if Z flag is not set
func i_jr_nz(c *CPU, l, _ byte) {
	if !c.FlagZero {
		c.PC = signedOffset(c.PC, l)
	}
}

// Jumps to signed addr offset if Z is set
func i_jr_z(c *CPU, l, _ byte) {
	if c.FlagZero {
		c.PC = signedOffset(c.PC, l)
	}
}

// Compare A with the given value
func i_cp_n(c *CPU, l, _ byte) {
	c.FlagZero = c.A-l == 0
	c.FlagSubstract = true
	c.FlagHalfCarry = (c.A & 0xF) < (l & 0xF)
	c.FlagCarry = c.A < l
}

// Adds the value at *HL to A
func i_add_a_phl(c *CPU, _, _ byte) {
	old := c.A
	add := c.MMU.Get8b(c.HL)

	c.A += add
	c.FlagZero = c.A == 0
	c.FlagSubstract = false
	c.FlagHalfCarry = (((old & 0xF) + (add & 0xF)) & 0x10) > 0
	c.FlagCarry = c.A < old
}

// Compare A with the value at *HL
func i_cp_phl(c *CPU, _, _ byte) {
	v := c.MMU.Get8b(c.HL)
	c.FlagZero = c.A-v == 0
	c.FlagSubstract = true
	c.FlagHalfCarry = (c.A & 0xF) < (v & 0xF)
	c.FlagCarry = c.A < v
}

// Rotates register A left through carry
func i_rla(c *CPU, _, _ byte) {
	c.A, c.FlagCarry = rotateLeftWithCarry(c.A, c.FlagCarry)

	// GBCPUman says the Z flag depends on the final value, other sources says
	// RLA always clear the flags, no sure who to trust.
	c.FlagZero = false
	c.FlagSubstract = false
	c.FlagHalfCarry = false
}
