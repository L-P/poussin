package cpu

import "math/bits"

var cbInstructionsMap = map[byte]Instruction{
	0x00: {1, 4, "RLC B", i_cb_rlc_n('B')},
	0x01: {1, 4, "RLC C", i_cb_rlc_n('C')},
	0x02: {1, 4, "RLC D", i_cb_rlc_n('D')},
	0x03: {1, 4, "RLC E", i_cb_rlc_n('E')},
	0x04: {1, 4, "RLC H", i_cb_rlc_n('H')},
	0x05: {1, 4, "RLC L", i_cb_rlc_n('L')},
	0x07: {1, 4, "RLC A", i_cb_rlc_n('A')},

	0x3F: {1, 4, "SRL A", i_cb_srl_n('A')},
	0x38: {1, 4, "SRL B", i_cb_srl_n('B')},
	0x39: {1, 4, "SRL C", i_cb_srl_n('C')},
	0x3A: {1, 4, "SRL D", i_cb_srl_n('D')},
	0x3B: {1, 4, "SRL E", i_cb_srl_n('E')},
	0x3C: {1, 4, "SRL H", i_cb_srl_n('H')},
	0x3D: {1, 4, "SRL L", i_cb_srl_n('L')},

	0x10: {1, 4, "RL B", i_cb_rl_n('B')},
	0x11: {1, 4, "RL C", i_cb_rl_n('C')},
	0x12: {1, 4, "RL D", i_cb_rl_n('D')},
	0x13: {1, 4, "RL E", i_cb_rl_n('E')},
	0x14: {1, 4, "RL H", i_cb_rl_n('H')},
	0x15: {1, 4, "RL L", i_cb_rl_n('L')},
	0x17: {1, 4, "RL A", i_cb_rl_n('A')},

	0x18: {1, 4, "RR B", i_cb_rr_n('B')},
	0x19: {1, 4, "RR C", i_cb_rr_n('C')},
	0x1A: {1, 4, "RR D", i_cb_rr_n('D')},
	0x1B: {1, 4, "RR E", i_cb_rr_n('E')},
	0x1C: {1, 4, "RR H", i_cb_rr_n('H')},
	0x1D: {1, 4, "RR L", i_cb_rr_n('L')},
	0x1F: {1, 4, "RR A", i_cb_rr_n('A')},

	0x27: {1, 4, "SLA A", i_cb_sla_n('A')},
	0x20: {1, 4, "SLA B", i_cb_sla_n('B')},
	0x21: {1, 4, "SLA C", i_cb_sla_n('C')},
	0x22: {1, 4, "SLA D", i_cb_sla_n('D')},
	0x23: {1, 4, "SLA E", i_cb_sla_n('E')},
	0x24: {1, 4, "SLA H", i_cb_sla_n('H')},
	0x25: {1, 4, "SLA L", i_cb_sla_n('L')},

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

	0x80: {1, 4, "RES 0,B", i_cb_res_x_n(0, 'B')},
	0x81: {1, 4, "RES 0,C", i_cb_res_x_n(0, 'C')},
	0x82: {1, 4, "RES 0,D", i_cb_res_x_n(0, 'D')},
	0x83: {1, 4, "RES 0,E", i_cb_res_x_n(0, 'E')},
	0x84: {1, 4, "RES 0,H", i_cb_res_x_n(0, 'H')},
	0x85: {1, 4, "RES 0,L", i_cb_res_x_n(0, 'L')},
	0x87: {1, 4, "RES 0,A", i_cb_res_x_n(0, 'A')},
	0x88: {1, 4, "RES 1,B", i_cb_res_x_n(1, 'B')},
	0x89: {1, 4, "RES 1,C", i_cb_res_x_n(1, 'C')},
	0x8A: {1, 4, "RES 1,D", i_cb_res_x_n(1, 'D')},
	0x8B: {1, 4, "RES 1,E", i_cb_res_x_n(1, 'E')},
	0x8C: {1, 4, "RES 1,H", i_cb_res_x_n(1, 'H')},
	0x8D: {1, 4, "RES 1,L", i_cb_res_x_n(1, 'L')},
	0x8F: {1, 4, "RES 1,A", i_cb_res_x_n(1, 'A')},
	0x90: {1, 4, "RES 2,B", i_cb_res_x_n(2, 'B')},
	0x91: {1, 4, "RES 2,C", i_cb_res_x_n(2, 'C')},
	0x92: {1, 4, "RES 2,D", i_cb_res_x_n(2, 'D')},
	0x93: {1, 4, "RES 2,E", i_cb_res_x_n(2, 'E')},
	0x94: {1, 4, "RES 2,H", i_cb_res_x_n(2, 'H')},
	0x95: {1, 4, "RES 2,L", i_cb_res_x_n(2, 'L')},
	0x97: {1, 4, "RES 2,A", i_cb_res_x_n(2, 'A')},
	0x98: {1, 4, "RES 3,B", i_cb_res_x_n(3, 'B')},
	0x99: {1, 4, "RES 3,C", i_cb_res_x_n(3, 'C')},
	0x9A: {1, 4, "RES 3,D", i_cb_res_x_n(3, 'D')},
	0x9B: {1, 4, "RES 3,E", i_cb_res_x_n(3, 'E')},
	0x9C: {1, 4, "RES 3,H", i_cb_res_x_n(3, 'H')},
	0x9D: {1, 4, "RES 3,L", i_cb_res_x_n(3, 'L')},
	0x9F: {1, 4, "RES 3,A", i_cb_res_x_n(3, 'A')},
	0xA0: {1, 4, "RES 4,B", i_cb_res_x_n(4, 'B')},
	0xA1: {1, 4, "RES 4,C", i_cb_res_x_n(4, 'C')},
	0xA2: {1, 4, "RES 4,D", i_cb_res_x_n(4, 'D')},
	0xA3: {1, 4, "RES 4,E", i_cb_res_x_n(4, 'E')},
	0xA4: {1, 4, "RES 4,H", i_cb_res_x_n(4, 'H')},
	0xA5: {1, 4, "RES 4,L", i_cb_res_x_n(4, 'L')},
	0xA7: {1, 4, "RES 4,A", i_cb_res_x_n(4, 'A')},
	0xA8: {1, 4, "RES 5,B", i_cb_res_x_n(5, 'B')},
	0xA9: {1, 4, "RES 5,C", i_cb_res_x_n(5, 'C')},
	0xAA: {1, 4, "RES 5,D", i_cb_res_x_n(5, 'D')},
	0xAB: {1, 4, "RES 5,E", i_cb_res_x_n(5, 'E')},
	0xAC: {1, 4, "RES 5,H", i_cb_res_x_n(5, 'H')},
	0xAD: {1, 4, "RES 5,L", i_cb_res_x_n(5, 'L')},
	0xAF: {1, 4, "RES 5,A", i_cb_res_x_n(5, 'A')},
	0xB0: {1, 4, "RES 6,B", i_cb_res_x_n(6, 'B')},
	0xB1: {1, 4, "RES 6,C", i_cb_res_x_n(6, 'C')},
	0xB2: {1, 4, "RES 6,D", i_cb_res_x_n(6, 'D')},
	0xB3: {1, 4, "RES 6,E", i_cb_res_x_n(6, 'E')},
	0xB4: {1, 4, "RES 6,H", i_cb_res_x_n(6, 'H')},
	0xB5: {1, 4, "RES 6,L", i_cb_res_x_n(6, 'L')},
	0xB7: {1, 4, "RES 6,A", i_cb_res_x_n(6, 'A')},
	0xB8: {1, 4, "RES 7,B", i_cb_res_x_n(7, 'B')},
	0xB9: {1, 4, "RES 7,C", i_cb_res_x_n(7, 'C')},
	0xBA: {1, 4, "RES 7,D", i_cb_res_x_n(7, 'D')},
	0xBB: {1, 4, "RES 7,E", i_cb_res_x_n(7, 'E')},
	0xBC: {1, 4, "RES 7,H", i_cb_res_x_n(7, 'H')},
	0xBD: {1, 4, "RES 7,L", i_cb_res_x_n(7, 'L')},
	0xBF: {1, 4, "RES 7,A", i_cb_res_x_n(7, 'A')},

	0xC0: {1, 4, "SET 0,B", i_cb_set_x_n(0, 'B')},
	0xC1: {1, 4, "SET 0,C", i_cb_set_x_n(0, 'C')},
	0xC2: {1, 4, "SET 0,D", i_cb_set_x_n(0, 'D')},
	0xC3: {1, 4, "SET 0,E", i_cb_set_x_n(0, 'E')},
	0xC4: {1, 4, "SET 0,H", i_cb_set_x_n(0, 'H')},
	0xC5: {1, 4, "SET 0,L", i_cb_set_x_n(0, 'L')},
	0xC7: {1, 4, "SET 0,A", i_cb_set_x_n(0, 'A')},
	0xC8: {1, 4, "SET 1,B", i_cb_set_x_n(1, 'B')},
	0xC9: {1, 4, "SET 1,C", i_cb_set_x_n(1, 'C')},
	0xCA: {1, 4, "SET 1,D", i_cb_set_x_n(1, 'D')},
	0xCB: {1, 4, "SET 1,E", i_cb_set_x_n(1, 'E')},
	0xCC: {1, 4, "SET 1,H", i_cb_set_x_n(1, 'H')},
	0xCD: {1, 4, "SET 1,L", i_cb_set_x_n(1, 'L')},
	0xCF: {1, 4, "SET 1,A", i_cb_set_x_n(1, 'A')},
	0xD0: {1, 4, "SET 2,B", i_cb_set_x_n(2, 'B')},
	0xD1: {1, 4, "SET 2,C", i_cb_set_x_n(2, 'C')},
	0xD2: {1, 4, "SET 2,D", i_cb_set_x_n(2, 'D')},
	0xD3: {1, 4, "SET 2,E", i_cb_set_x_n(2, 'E')},
	0xD4: {1, 4, "SET 2,H", i_cb_set_x_n(2, 'H')},
	0xD5: {1, 4, "SET 2,L", i_cb_set_x_n(2, 'L')},
	0xD7: {1, 4, "SET 2,A", i_cb_set_x_n(2, 'A')},
	0xD8: {1, 4, "SET 3,B", i_cb_set_x_n(3, 'B')},
	0xD9: {1, 4, "SET 3,C", i_cb_set_x_n(3, 'C')},
	0xDA: {1, 4, "SET 3,D", i_cb_set_x_n(3, 'D')},
	0xDB: {1, 4, "SET 3,E", i_cb_set_x_n(3, 'E')},
	0xDC: {1, 4, "SET 3,H", i_cb_set_x_n(3, 'H')},
	0xDD: {1, 4, "SET 3,L", i_cb_set_x_n(3, 'L')},
	0xDF: {1, 4, "SET 3,A", i_cb_set_x_n(3, 'A')},
	0xE0: {1, 4, "SET 4,B", i_cb_set_x_n(4, 'B')},
	0xE1: {1, 4, "SET 4,C", i_cb_set_x_n(4, 'C')},
	0xE2: {1, 4, "SET 4,D", i_cb_set_x_n(4, 'D')},
	0xE3: {1, 4, "SET 4,E", i_cb_set_x_n(4, 'E')},
	0xE4: {1, 4, "SET 4,H", i_cb_set_x_n(4, 'H')},
	0xE5: {1, 4, "SET 4,L", i_cb_set_x_n(4, 'L')},
	0xE7: {1, 4, "SET 4,A", i_cb_set_x_n(4, 'A')},
	0xE8: {1, 4, "SET 5,B", i_cb_set_x_n(5, 'B')},
	0xE9: {1, 4, "SET 5,C", i_cb_set_x_n(5, 'C')},
	0xEA: {1, 4, "SET 5,D", i_cb_set_x_n(5, 'D')},
	0xEB: {1, 4, "SET 5,E", i_cb_set_x_n(5, 'E')},
	0xEC: {1, 4, "SET 5,H", i_cb_set_x_n(5, 'H')},
	0xED: {1, 4, "SET 5,L", i_cb_set_x_n(5, 'L')},
	0xEF: {1, 4, "SET 5,A", i_cb_set_x_n(5, 'A')},
	0xF0: {1, 4, "SET 6,B", i_cb_set_x_n(6, 'B')},
	0xF1: {1, 4, "SET 6,C", i_cb_set_x_n(6, 'C')},
	0xF2: {1, 4, "SET 6,D", i_cb_set_x_n(6, 'D')},
	0xF3: {1, 4, "SET 6,E", i_cb_set_x_n(6, 'E')},
	0xF4: {1, 4, "SET 6,H", i_cb_set_x_n(6, 'H')},
	0xF5: {1, 4, "SET 6,L", i_cb_set_x_n(6, 'L')},
	0xF7: {1, 4, "SET 6,A", i_cb_set_x_n(6, 'A')},
	0xF8: {1, 4, "SET 7,B", i_cb_set_x_n(7, 'B')},
	0xF9: {1, 4, "SET 7,C", i_cb_set_x_n(7, 'C')},
	0xFA: {1, 4, "SET 7,D", i_cb_set_x_n(7, 'D')},
	0xFB: {1, 4, "SET 7,E", i_cb_set_x_n(7, 'E')},
	0xFC: {1, 4, "SET 7,H", i_cb_set_x_n(7, 'H')},
	0xFD: {1, 4, "SET 7,L", i_cb_set_x_n(7, 'L')},
	0xFF: {1, 4, "SET 7,A", i_cb_set_x_n(7, 'A')},

	0x30: {1, 4, "SWAP B", i_cb_swap_n('B')},
	0x31: {1, 4, "SWAP C", i_cb_swap_n('C')},
	0x32: {1, 4, "SWAP D", i_cb_swap_n('D')},
	0x33: {1, 4, "SWAP E", i_cb_swap_n('E')},
	0x34: {1, 4, "SWAP H", i_cb_swap_n('H')},
	0x35: {1, 4, "SWAP L", i_cb_swap_n('L')},
	0x37: {1, 4, "SWAP A", i_cb_swap_n('A')},

	0x86: {1, 16, "RES 0,(HL)", i_cb_res_x_phl(0)},
	0x8E: {1, 16, "RES 1,(HL)", i_cb_res_x_phl(1)},
	0x96: {1, 16, "RES 2,(HL)", i_cb_res_x_phl(2)},
	0x9E: {1, 16, "RES 3,(HL)", i_cb_res_x_phl(3)},
	0xA6: {1, 16, "RES 4,(HL)", i_cb_res_x_phl(4)},
	0xAE: {1, 16, "RES 5,(HL)", i_cb_res_x_phl(5)},
	0xB6: {1, 16, "RES 6,(HL)", i_cb_res_x_phl(6)},
	0xBE: {1, 16, "RES 7,(HL)", i_cb_res_x_phl(7)},

	0xC6: {1, 16, "SET 0,(HL)", i_cb_set_x_phl(0)},
	0xCE: {1, 16, "SET 1,(HL)", i_cb_set_x_phl(1)},
	0xD6: {1, 16, "SET 2,(HL)", i_cb_set_x_phl(2)},
	0xDE: {1, 16, "SET 3,(HL)", i_cb_set_x_phl(3)},
	0xE6: {1, 16, "SET 4,(HL)", i_cb_set_x_phl(4)},
	0xEE: {1, 16, "SET 5,(HL)", i_cb_set_x_phl(5)},
	0xF6: {1, 16, "SET 6,(HL)", i_cb_set_x_phl(6)},
	0xFE: {1, 16, "SET 7,(HL)", i_cb_set_x_phl(7)},
}

// Rotates n left through Carry flag
func i_cb_rl_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)
		var b byte
		b, c.FlagCarry = rotateLeftWithCarry(get(), c.FlagCarry)

		set(b)
		c.FlagZero = b == 0
		c.FlagSubstract = false
		c.FlagHalfCarry = false
	}
}

// Rotates n right through Carry flag
func i_cb_rr_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)
		var b byte
		b, c.FlagCarry = rotateRightWithCarry(get(), c.FlagCarry)

		set(b)
		c.FlagZero = b == 0
		c.FlagSubstract = false
		c.FlagHalfCarry = false
	}
}

// Shifts n right into Carry. MSB set to 0.
func i_cb_srl_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)
		b := get()

		c.FlagCarry = b&0x01 == 0x01
		b = b >> 1

		set(b)
		c.FlagZero = b == 0
		c.FlagSubstract = false
		c.FlagHalfCarry = false
	}
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

// Resets bit x of n
func i_cb_res_x_n(bit uint, name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)
		set(get() &^ (1 << bit))
	}
}

// Resets bit x of the value pointed by HL
func i_cb_res_x_phl(bit uint) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		c.Write(c.HL, c.Fetch(c.HL)&^(1<<bit))
	}
}

// Sets bit x of the value pointed by HL
func i_cb_set_x_phl(bit uint) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		c.Write(c.HL, c.Fetch(c.HL)|(1<<bit))
	}
}

// Sets bit x of n
func i_cb_set_x_n(bit uint, name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)
		set(get() | (1 << bit))
	}
}

// Swaps high and low nibble of a register
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

// Shifts n into carry
func i_cb_sla_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)

		w := uint16(get()) << 1
		b := byte(w)
		set(b)

		c.ClearFlags()
		c.FlagZero = b == 0
		c.FlagCarry = (w & (1 << 8)) > 0
	}
}

// Rotates n left, old 7 bit to carry
func i_cb_rlc_n(name byte) InstructionImplementation {
	return func(c *CPU, _, _ byte) {
		get, set := c.GetRegisterCallbacks(name)

		b := get()
		c.FlagCarry = (b & (1 << 7)) > 0

		b = bits.RotateLeft8(b, 1)
		set(b)

		c.FlagZero = b == 0
		c.FlagHalfCarry = false
		c.FlagSubstract = false
	}
}
