package cpu

var Instructions = map[byte]Instruction{
	0x00: {1, 4, "NOP", i_nop},

	0xF3: {1, 4, "DI", i_set_interrupt(false)},
	0xFB: {1, 4, "EI", i_set_interrupt(true)},

	0x2F: {1, 4, "CPL", i_cpl},

	0x3C: {1, 4, "INC A", i_inc_a},
	0x04: {1, 4, "INC B", i_inc_n('B')},
	0x0C: {1, 4, "INC C", i_inc_n('C')},
	0x14: {1, 4, "INC D", i_inc_n('D')},
	0x1C: {1, 4, "INC E", i_inc_n('E')},
	0x24: {1, 4, "INC H", i_inc_n('H')},
	0x2C: {1, 4, "INC L", i_inc_n('L')},

	0x3D: {1, 4, "DEC A", i_dec_a},
	0x05: {1, 4, "DEC B", i_dec_n('B')},
	0x0D: {1, 4, "DEC C", i_dec_n('C')},
	0x15: {1, 4, "DEC D", i_dec_n('D')},
	0x1D: {1, 4, "DEC E", i_dec_n('E')},
	0x25: {1, 4, "DEC H", i_dec_n('H')},
	0x2D: {1, 4, "DEC L", i_dec_n('L')},

	0x03: {1, 8, "INC BC", i_inc_nn("BC")},
	0x13: {1, 8, "INC DE", i_inc_nn("DE")},
	0x23: {1, 8, "INC HL", i_inc_nn("HL")},
	0x33: {1, 8, "INC SP", i_inc_nn("SP")},

	0x0B: {1, 8, "DEC BC", i_dec_nn("BC")},
	0x1B: {1, 8, "DEC DE", i_dec_nn("DE")},
	0x2B: {1, 8, "DEC HL", i_dec_nn("HL")},
	0x3B: {1, 8, "DEC SP", i_dec_nn("SP")},

	0x80: {1, 4, "ADD A,B", i_add_a_n('B')},
	0x81: {1, 4, "ADD A,C", i_add_a_n('C')},
	0x82: {1, 4, "ADD A,D", i_add_a_n('D')},
	0x83: {1, 4, "ADD A,E", i_add_a_n('E')},
	0x84: {1, 4, "ADD A,H", i_add_a_n('H')},
	0x85: {1, 4, "ADD A,L", i_add_a_n('L')},
	0x87: {1, 4, "ADD A,A", i_add_a_n('A')},

	0x86: {1, 8, "ADD A,(HL)", i_add_a_phl},

	0x09: {1, 8, "ADD HL,BC", i_add_hl_nn("BC")},
	0x19: {1, 8, "ADD HL,DE", i_add_hl_nn("DE")},
	0x29: {1, 8, "ADD HL,HL", i_add_hl_nn("HL")},
	0x39: {1, 8, "ADD HL,SP", i_add_hl_nn("SP")},

	0x97: {1, 4, "SUB A", i_sub_n('A')},
	0x90: {1, 4, "SUB B", i_sub_n('B')},
	0x91: {1, 4, "SUB C", i_sub_n('C')},
	0x92: {1, 4, "SUB D", i_sub_n('D')},
	0x93: {1, 4, "SUB E", i_sub_n('E')},
	0x94: {1, 4, "SUB H", i_sub_n('H')},
	0x95: {1, 4, "SUB L", i_sub_n('L')},

	0xA8: {1, 4, "XOR B", i_xor_n('B')},
	0xA9: {1, 4, "XOR C", i_xor_n('C')},
	0xAA: {1, 4, "XOR D", i_xor_n('D')},
	0xAB: {1, 4, "XOR E", i_xor_n('E')},
	0xAC: {1, 4, "XOR H", i_xor_n('H')},
	0xAD: {1, 4, "XOR L", i_xor_n('L')},
	0xAF: {1, 4, "XOR A", i_xor_n('A')},

	0xA0: {1, 4, "AND B", i_and_n('B')},
	0xA1: {1, 4, "AND C", i_and_n('C')},
	0xA2: {1, 4, "AND D", i_and_n('D')},
	0xA3: {1, 4, "AND E", i_and_n('E')},
	0xA4: {1, 4, "AND H", i_and_n('H')},
	0xA5: {1, 4, "AND L", i_and_n('L')},
	0xA7: {1, 4, "AND A", i_and_n('A')},

	0xB0: {1, 4, "OR B", i_or_n('B')},
	0xB1: {1, 4, "OR C", i_or_n('C')},
	0xB2: {1, 4, "OR D", i_or_n('D')},
	0xB3: {1, 4, "OR E", i_or_n('E')},
	0xB4: {1, 4, "OR H", i_or_n('H')},
	0xB5: {1, 4, "OR L", i_or_n('L')},
	0xB7: {1, 4, "OR A", i_or_n('A')},

	0xE6: {2, 8, "AND $%02X", i_and},
	0xEE: {2, 8, "XOR $%02X", i_xor},
	0xF6: {2, 8, "OR $%02X", i_or},

	0x17: {1, 4, "RLA", i_rla},

	0x3E: {2, 8, "LD A,$%02X", i_ld_a_nn},
	0xFA: {2, 16, "LD A,($%02X)", i_ld_a_pnn},
	0x06: {2, 8, "LD B,$%02X", i_ld_n('B')},
	0x0E: {2, 8, "LD C,$%02X", i_ld_n('C')},
	0x16: {2, 8, "LD D,$%02X", i_ld_n('D')},
	0x1E: {2, 8, "LD E,$%02X", i_ld_n('E')},
	0x26: {2, 8, "LD H,$%02X", i_ld_n('H')},
	0x2E: {2, 8, "LD L,$%02X", i_ld_n('L')},

	0x0A: {1, 8, "LD A,(BC)", i_ld_n_pnn('A', "BC")},
	0x1A: {1, 8, "LD A,(DE)", i_ld_n_pnn('A', "DE")},
	0x7E: {1, 8, "LD A,(HL)", i_ld_n_pnn('A', "HL")},
	0x46: {1, 8, "LD B,(HL)", i_ld_n_pnn('B', "HL")},
	0x4E: {1, 8, "LD C,(HL)", i_ld_n_pnn('C', "HL")},
	0x56: {1, 8, "LD D,(HL)", i_ld_n_pnn('D', "HL")},
	0x5E: {1, 8, "LD E,(HL)", i_ld_n_pnn('E', "HL")},
	0x66: {1, 8, "LD H,(HL)", i_ld_n_pnn('H', "HL")},
	0x6E: {1, 8, "LD L,(HL)", i_ld_n_pnn('L', "HL")},

	0x2A: {1, 8, "LDI A,(HL)", i_ldi_a_phl},

	0x40: {1, 4, "LD B,B", i_ld_n_n('B', 'B')},
	0x41: {1, 4, "LD B,C", i_ld_n_n('B', 'C')},
	0x42: {1, 4, "LD B,D", i_ld_n_n('B', 'D')},
	0x43: {1, 4, "LD B,E", i_ld_n_n('B', 'E')},
	0x44: {1, 4, "LD B,H", i_ld_n_n('B', 'H')},
	0x45: {1, 4, "LD B,L", i_ld_n_n('B', 'L')},
	0x47: {1, 4, "LD B,A", i_ld_n_n('B', 'A')},

	0x48: {1, 4, "LD C,B", i_ld_n_n('C', 'B')},
	0x49: {1, 4, "LD C,C", i_ld_n_n('C', 'C')},
	0x4A: {1, 4, "LD C,D", i_ld_n_n('C', 'D')},
	0x4B: {1, 4, "LD C,E", i_ld_n_n('C', 'E')},
	0x4C: {1, 4, "LD C,H", i_ld_n_n('C', 'H')},
	0x4D: {1, 4, "LD C,L", i_ld_n_n('C', 'L')},
	0x4F: {1, 4, "LD C,A", i_ld_n_n('C', 'A')},

	0x50: {1, 4, "LD D,B", i_ld_n_n('D', 'B')},
	0x51: {1, 4, "LD D,C", i_ld_n_n('D', 'C')},
	0x52: {1, 4, "LD D,D", i_ld_n_n('D', 'D')},
	0x53: {1, 4, "LD D,E", i_ld_n_n('D', 'E')},
	0x54: {1, 4, "LD D,H", i_ld_n_n('D', 'H')},
	0x55: {1, 4, "LD D,L", i_ld_n_n('D', 'L')},
	0x57: {1, 4, "LD D,A", i_ld_n_n('D', 'A')},

	0x58: {1, 4, "LD E,B", i_ld_n_n('E', 'B')},
	0x59: {1, 4, "LD E,C", i_ld_n_n('E', 'C')},
	0x5A: {1, 4, "LD E,D", i_ld_n_n('E', 'D')},
	0x5B: {1, 4, "LD E,E", i_ld_n_n('E', 'E')},
	0x5C: {1, 4, "LD E,H", i_ld_n_n('E', 'H')},
	0x5D: {1, 4, "LD E,L", i_ld_n_n('E', 'L')},
	0x5F: {1, 4, "LD E,A", i_ld_n_n('E', 'A')},

	0x60: {1, 4, "LD H,B", i_ld_n_n('H', 'B')},
	0x61: {1, 4, "LD H,C", i_ld_n_n('H', 'C')},
	0x62: {1, 4, "LD H,D", i_ld_n_n('H', 'D')},
	0x63: {1, 4, "LD H,E", i_ld_n_n('H', 'E')},
	0x64: {1, 4, "LD H,H", i_ld_n_n('H', 'H')},
	0x65: {1, 4, "LD H,L", i_ld_n_n('H', 'L')},
	0x67: {1, 4, "LD H,A", i_ld_n_n('H', 'A')},

	0x68: {1, 4, "LD L,B", i_ld_n_n('L', 'B')},
	0x69: {1, 4, "LD L,C", i_ld_n_n('L', 'C')},
	0x6A: {1, 4, "LD L,D", i_ld_n_n('L', 'D')},
	0x6B: {1, 4, "LD L,E", i_ld_n_n('L', 'E')},
	0x6C: {1, 4, "LD L,H", i_ld_n_n('L', 'H')},
	0x6D: {1, 4, "LD L,L", i_ld_n_n('L', 'L')},
	0x6F: {1, 4, "LD L,A", i_ld_n_n('L', 'A')},

	0x78: {1, 4, "LD A,B", i_ld_n_n('A', 'B')},
	0x79: {1, 4, "LD A,C", i_ld_n_n('A', 'C')},
	0x7A: {1, 4, "LD A,D", i_ld_n_n('A', 'D')},
	0x7B: {1, 4, "LD A,E", i_ld_n_n('A', 'E')},
	0x7C: {1, 4, "LD A,H", i_ld_n_n('A', 'H')},
	0x7D: {1, 4, "LD A,L", i_ld_n_n('A', 'L')},
	0x7F: {1, 4, "LD A,A", i_ld_n_n('A', 'A')},

	0x01: {3, 12, "LD BC,$%02X%02X", i_ld_nn("BC")},
	0x11: {3, 12, "LD DE,$%02X%02X", i_ld_nn("DE")},
	0x21: {3, 12, "LD HL,$%02X%02X", i_ld_nn("HL")},
	0x31: {3, 12, "LD SP,$%02X%02X", i_ld_nn("SP")},

	0xE2: {1, 8, "LD (C),A", i_ld_pc_a},
	0xEA: {3, 16, "LD ($%02X%02X),A", i_ld_pn_a},

	0x02: {1, 8, "LD (BC),A", i_ld_pnn_a("BC")},
	0x12: {1, 8, "LD (DE),A", i_ld_pnn_a("DE")},
	0x77: {1, 8, "LD (HL),A", i_ld_pnn_a("HL")},

	0x36: {2, 12, "LD (HL),%02X", i_ld_phl_n},

	0x22: {1, 8, "LDI (HL),A", i_ldi_phl_a},
	0x32: {1, 8, "LDD (HL),A", i_ldd_phl_a},

	0xE0: {2, 12, "LDH ($%02X),A", i_ldh_pn_a},
	0xF0: {2, 12, "LDH A,($%02X)", i_ldh_a_pn},

	0xBE: {1, 8, "CP (HL)", i_cp_phl},
	0xFE: {2, 8, "CP $%02X", i_cp_n},

	0x18: {2, 8, "JR,$%02X", i_jr},
	0x20: {2, 8, "JR NZ,$%02X", i_jr_nz},
	0x28: {2, 8, "JR Z,$%02X", i_jr_z},

	0xC5: {1, 16, "PUSH BC", i_push_nn("BC")},
	0xD5: {1, 16, "PUSH DE", i_push_nn("DE")},
	0xE5: {1, 16, "PUSH HL", i_push_nn("HL")},
	0xF5: {1, 16, "PUSH AF", i_push_af},

	0xC1: {1, 12, "POP BC", i_pop_nn("BC")},
	0xD1: {1, 12, "POP DE", i_pop_nn("DE")},
	0xE1: {1, 12, "POP HL", i_pop_nn("HL")},

	0xC3: {3, 12, "JP $%02X%02X", i_jp_nn},
	0xE9: {1, 4, "JP (HL)", i_jp_hl}, // weird mnemonic, we go to HL, not (HL)
	0xCB: {1, 4, "PREFIX CB", i_prefix_cb},
	0xCD: {3, 24, "CALL $%02X%02X", i_call},
	0xC9: {1, 8, "RET", i_ret},

	0xC7: {1, 16, "RST,$00", i_rst(0x00)},
	0xCF: {1, 16, "RST,$08", i_rst(0x08)},
	0xD7: {1, 16, "RST,$10", i_rst(0x10)},
	0xDF: {1, 16, "RST,$18", i_rst(0x18)},
	0xE7: {1, 16, "RST,$20", i_rst(0x20)},
	0xEF: {1, 16, "RST,$28", i_rst(0x28)},
	0xF7: {1, 16, "RST,$30", i_rst(0x30)},
	0xFF: {1, 16, "RST,$38", i_rst(0x38)},
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

// Increments register nn
func i_inc_nn(name string) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		r := c.GetRegisterAddress(name)
		*r++
	}
}

// Decrements register nn
func i_dec_nn(name string) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		r := c.GetRegisterAddress(name)
		*r--
	}
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
	c.Write(0xFF00|uint16(c.GetC()), c.A)
}

// Loads 8b value into n
func i_ld_n(name byte) InstructionImplementation {
	return func(c *CPU, l, _ byte) {
		_, set := c.GetRegisterCallbacks(name)
		set(l)
	}
}

// Loads 8b value into A
func i_ld_a_nn(c *CPU, l, _ byte) {
	c.A = l
}

// Loads value at given address into A
func i_ld_a_pnn(c *CPU, l, h byte) {
	c.A = c.Fetch((uint16(h) << 8) | uint16(l))
}

// Loads 16b value into register
func i_ld_nn(name string) InstructionImplementation {
	return func(c *CPU, l, h byte) {
		r := c.GetRegisterAddress(name)
		*r = (uint16(h) << 8) | uint16(l)
	}
}

// Puts A into address pointed by HL and decrement HL
func i_ldd_phl_a(c *CPU, _, _ byte) {
	c.Write(c.HL, c.A)
	c.HL--
}

// Puts A into address pointed by HL and increment HL
func i_ldi_phl_a(c *CPU, _, _ byte) {
	c.Write(c.HL, c.A)
	c.HL++
}

// Puts A into address pointed by nn
func i_ld_pnn_a(name string) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		r := c.GetRegisterAddress(name)
		c.Write(*r, c.A)
	}
}

// Puts n into address pointed by HL
func i_ld_phl_n(c *CPU, l, _ byte) {
	c.Write(c.HL, l)
}

// Puts A into given address pointed by HL
func i_ld_pn_a(c *CPU, l, h byte) {
	c.Write(uint16(l)|(uint16(h)<<8), c.A)
}

// Puts A into address 0xFF00+l
func i_ldh_pn_a(c *CPU, l, _ byte) {
	c.Write(0xFF00+uint16(l), c.A)
}

// Puts value at 0xFF00+l into A
func i_ldh_a_pn(c *CPU, l, _ byte) {
	c.A = c.Fetch(0xFF00 + uint16(l))
}

// Performs a logical AND against A and l
func i_and(c *CPU, l, _ byte) {
	c.A &= l
	c.ClearFlags()
	c.FlagZero = c.A == 0
	c.FlagHalfCarry = true
}

// Performs a logical XOR against A and l
func i_xor(c *CPU, l, _ byte) {
	c.A ^= l
	c.ClearFlags()
	c.FlagZero = c.A == 0
}

// Performs a logical OR against A and l
func i_or(c *CPU, l, _ byte) {
	c.A |= l
	c.ClearFlags()
	c.FlagZero = c.A == 0
}

// Performs a logical AND against A and n
func i_and_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, _ := c.GetRegisterCallbacks(name)

		c.A &= get()
		c.ClearFlags()
		c.FlagZero = c.A == 0
		c.FlagHalfCarry = true
	}
}

// Performs a logical XOR against A xor n
func i_xor_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, _ := c.GetRegisterCallbacks(name)

		c.A ^= get()
		c.ClearFlags()
		c.FlagZero = c.A == 0
	}
}

// Performs a logical OR against A and n
func i_or_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, _ := c.GetRegisterCallbacks(name)

		c.A |= get()
		c.ClearFlags()
		c.FlagZero = c.A == 0
	}
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

// Pushes a two-byte register to the stack
func i_push_nn(name string) InstructionImplementation {
	return func(c *CPU, l, h byte) {
		r := c.GetRegisterAddress(name)
		c.StackPush16b(*r)
	}
}

func i_push_af(c *CPU, _, _ byte) {
	c.StackPush16b((uint16(c.A) << 8) | uint16(c.GetFlags()))
}

// Pops two bytes from the stack to nn
func i_pop_nn(name string) InstructionImplementation {
	return func(c *CPU, l, h byte) {
		r := c.GetRegisterAddress(name)
		*r = c.StackPop16b()
	}
}

// Loads the value at address pointed by nn in n
func i_ld_n_pnn(dst byte, src string) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		_, set := c.GetRegisterCallbacks(dst)
		r := c.GetRegisterAddress(src)

		set(c.Fetch(*r))
	}
}

// Loads the value at address pointed by HL in A and increment HL
func i_ldi_a_phl(c *CPU, _, _ byte) {
	c.A = c.Fetch(c.HL)
	c.HL++
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

// Adds the value n to A
func i_add_a_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, _ := c.GetRegisterCallbacks(name)

		old := c.A
		add := get()

		c.A += add
		c.FlagZero = c.A == 0
		c.FlagSubstract = false
		c.FlagHalfCarry = (((old & 0xF) + (add & 0xF)) & 0x10) > 0
		c.FlagCarry = c.A < old
	}
}

// Adds the value at *HL to A
func i_add_a_phl(c *CPU, _, _ byte) {
	old := c.A
	add := c.Fetch(c.HL)

	c.A += add
	c.FlagZero = c.A == 0
	c.FlagSubstract = false
	c.FlagHalfCarry = (((old & 0xF) + (add & 0xF)) & 0x10) > 0
	c.FlagCarry = c.A < old
}

// Compare A with the value at *HL
func i_cp_phl(c *CPU, _, _ byte) {
	v := c.Fetch(c.HL)
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

// Jumps to given address
func i_jp_nn(c *CPU, l, h byte) {
	c.PC = uint16(l) | (uint16(h) << 8)
}

// Jumps to the address HL
func i_jp_hl(c *CPU, _, _ byte) {
	c.PC = c.HL
}

func i_set_interrupt(v bool) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		c.InterruptMaster = v
	}
}

// Flips all bits of A ("Complement")
func i_cpl(c *CPU, _, _ byte) {
	c.A = c.A ^ 0xFF
	c.FlagHalfCarry = true
	c.FlagSubstract = true
}

// Pushes PC onto the stack and jump to 0x0000 + l
func i_rst(l byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		c.StackPush16b(c.PC)
		c.PC = uint16(l)
	}
}

// Adds the value of register nn to HL
func i_add_hl_nn(name string) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		r := c.GetRegisterAddress(name)
		old := c.HL

		c.HL += *r

		c.FlagSubstract = false
		c.FlagCarry = c.HL < old
		// TODO c.HalfCarry
	}
}
