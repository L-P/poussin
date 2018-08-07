package cpu

var CBInstructions = map[byte]Instruction{
	0x11: {1, 4, "RL C", i_cb_rl_c},

	0x40: {1, 4, "BIT 0,B", i_cb_bit_x_n(0, 'B')},
	0x41: {1, 4, "BIT 0,C", i_cb_bit_x_n(0, 'C')},
	0x42: {1, 4, "BIT 0,D", i_cb_bit_x_n(0, 'D')},
	0x43: {1, 4, "BIT 0,E", i_cb_bit_x_n(0, 'E')},
	0x44: {1, 4, "BIT 0,H", i_cb_bit_x_n(0, 'H')},
	0x45: {1, 4, "BIT 0,L", i_cb_bit_x_n(0, 'L')},
	0x47: {1, 4, "BIT 0,A", i_cb_bit_x_n(0, 'A')},
	0x48: {1, 4, "BIT 1,B", i_cb_bit_x_n(1, 'B')},
	0x49: {1, 4, "BIT 1,C", i_cb_bit_x_n(1, 'C')},
	0x4A: {1, 4, "BIT 1,D", i_cb_bit_x_n(1, 'D')},
	0x4B: {1, 4, "BIT 1,E", i_cb_bit_x_n(1, 'E')},
	0x4C: {1, 4, "BIT 1,H", i_cb_bit_x_n(1, 'H')},
	0x4D: {1, 4, "BIT 1,L", i_cb_bit_x_n(1, 'L')},
	0x4F: {1, 4, "BIT 1,A", i_cb_bit_x_n(1, 'A')},

	0x50: {1, 4, "BIT 2,B", i_cb_bit_x_n(2, 'B')},
	0x51: {1, 4, "BIT 2,C", i_cb_bit_x_n(2, 'C')},
	0x52: {1, 4, "BIT 2,D", i_cb_bit_x_n(2, 'D')},
	0x53: {1, 4, "BIT 2,E", i_cb_bit_x_n(2, 'E')},
	0x54: {1, 4, "BIT 2,H", i_cb_bit_x_n(2, 'H')},
	0x55: {1, 4, "BIT 2,L", i_cb_bit_x_n(2, 'L')},
	0x57: {1, 4, "BIT 2,A", i_cb_bit_x_n(2, 'A')},
	0x58: {1, 4, "BIT 3,B", i_cb_bit_x_n(3, 'B')},
	0x59: {1, 4, "BIT 3,C", i_cb_bit_x_n(3, 'C')},
	0x5A: {1, 4, "BIT 3,D", i_cb_bit_x_n(3, 'D')},
	0x5B: {1, 4, "BIT 3,E", i_cb_bit_x_n(3, 'E')},
	0x5C: {1, 4, "BIT 3,H", i_cb_bit_x_n(3, 'H')},
	0x5D: {1, 4, "BIT 3,L", i_cb_bit_x_n(3, 'L')},
	0x5F: {1, 4, "BIT 3,A", i_cb_bit_x_n(3, 'A')},

	0x60: {1, 4, "BIT 4,B", i_cb_bit_x_n(4, 'B')},
	0x61: {1, 4, "BIT 4,C", i_cb_bit_x_n(4, 'C')},
	0x62: {1, 4, "BIT 4,D", i_cb_bit_x_n(4, 'D')},
	0x63: {1, 4, "BIT 4,E", i_cb_bit_x_n(4, 'E')},
	0x64: {1, 4, "BIT 4,H", i_cb_bit_x_n(4, 'H')},
	0x65: {1, 4, "BIT 4,L", i_cb_bit_x_n(4, 'L')},
	0x67: {1, 4, "BIT 4,A", i_cb_bit_x_n(4, 'A')},
	0x68: {1, 4, "BIT 5,B", i_cb_bit_x_n(5, 'B')},
	0x69: {1, 4, "BIT 5,C", i_cb_bit_x_n(5, 'C')},
	0x6A: {1, 4, "BIT 5,D", i_cb_bit_x_n(5, 'D')},
	0x6B: {1, 4, "BIT 5,E", i_cb_bit_x_n(5, 'E')},
	0x6C: {1, 4, "BIT 5,H", i_cb_bit_x_n(5, 'H')},
	0x6D: {1, 4, "BIT 5,L", i_cb_bit_x_n(5, 'L')},
	0x6F: {1, 4, "BIT 5,A", i_cb_bit_x_n(5, 'A')},

	0x70: {1, 4, "BIT 6,B", i_cb_bit_x_n(6, 'B')},
	0x71: {1, 4, "BIT 6,C", i_cb_bit_x_n(6, 'C')},
	0x72: {1, 4, "BIT 6,D", i_cb_bit_x_n(6, 'D')},
	0x73: {1, 4, "BIT 6,E", i_cb_bit_x_n(6, 'E')},
	0x74: {1, 4, "BIT 6,H", i_cb_bit_x_n(6, 'H')},
	0x75: {1, 4, "BIT 6,L", i_cb_bit_x_n(6, 'L')},
	0x77: {1, 4, "BIT 6,A", i_cb_bit_x_n(6, 'A')},
	0x78: {1, 4, "BIT 7,B", i_cb_bit_x_n(7, 'B')},
	0x79: {1, 4, "BIT 7,C", i_cb_bit_x_n(7, 'C')},
	0x7A: {1, 4, "BIT 7,D", i_cb_bit_x_n(7, 'D')},
	0x7B: {1, 4, "BIT 7,E", i_cb_bit_x_n(7, 'E')},
	0x7C: {1, 4, "BIT 7,H", i_cb_bit_x_n(7, 'H')},
	0x7D: {1, 4, "BIT 7,L", i_cb_bit_x_n(7, 'L')},
	0x7F: {1, 4, "BIT 7,A", i_cb_bit_x_n(7, 'A')},

	0x30: {1, 4, "SWAP B", i_cb_swap_n('B')},
	0x31: {1, 4, "SWAP C", i_cb_swap_n('C')},
	0x32: {1, 4, "SWAP D", i_cb_swap_n('D')},
	0x33: {1, 4, "SWAP E", i_cb_swap_n('E')},
	0x34: {1, 4, "SWAP H", i_cb_swap_n('H')},
	0x35: {1, 4, "SWAP L", i_cb_swap_n('L')},
	0x37: {1, 4, "SWAP A", i_cb_swap_n('A')},
}

// Rotates C left through Carry flag
func i_cb_rl_c(c *CPU, _, _ byte) {
	var C byte
	C, c.FlagCarry = rotateLeftWithCarry(c.GetC(), c.FlagCarry)

	c.SetC(C)
	c.FlagZero = C == 0
	c.FlagSubstract = false
	c.FlagHalfCarry = false
}

// Sets flag Z if the nth bit of n is not set
func i_cb_bit_x_n(bit uint, name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, _ := c.GetRegisterCallbacks(name)
		c.FlagZero = (get() & (1 << bit)) == 0
		c.FlagSubstract = false
		c.FlagHalfCarry = true
	}
}

func i_cb_swap_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)

		b := get()
		b = ((b & 0x0F) << 4) | ((b & 0xF0) >> 4)
		set(b)

		c.ClearFlags()
		c.FlagZero = b == 0
	}
}
