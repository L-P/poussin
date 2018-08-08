package cpu

import "fmt"

// Low/High byte setters for 16b registers
func (c *CPU) SetB(b byte) {
	c.BC = (c.BC & 0x00FF) | (uint16(b) << 8)
}
func (c *CPU) SetC(b byte) {
	c.BC = (c.BC & 0xFF00) | uint16(b)
}
func (c *CPU) SetD(b byte) {
	c.DE = (c.DE & 0x00FF) | (uint16(b) << 8)
}
func (c *CPU) SetE(b byte) {
	c.DE = (c.DE & 0xFF00) | uint16(b)
}
func (c *CPU) SetH(b byte) {
	c.HL = (c.HL & 0x00FF) | (uint16(b) << 8)
}
func (c *CPU) SetL(b byte) {
	c.HL = (c.HL & 0xFF00) | uint16(b)
}

func (c *CPU) GetA() byte {
	return c.A
}
func (c *CPU) SetA(b byte) {
	c.A = b
}

// Low/High byte getters for 16b registers
func (c *CPU) GetB() byte {
	return byte((c.BC & 0xFF00) >> 8)
}
func (c *CPU) GetC() byte {
	return byte(c.BC & 0x00FF)
}
func (c *CPU) GetD() byte {
	return byte((c.DE & 0xFF00) >> 8)
}
func (c *CPU) GetE() byte {
	return byte(c.DE & 0x00FF)
}
func (c *CPU) GetH() byte {
	return byte((c.HL & 0xFF00) >> 8)
}
func (c *CPU) GetL() byte {
	return byte(c.HL & 0x00FF)
}

func (c *CPU) ClearFlags() {
	c.FlagZero = false
	c.FlagSubstract = false
	c.FlagHalfCarry = false
	c.FlagCarry = false
}

// Returns the flags as a byte (ie. the F register)
func (c *CPU) GetFlags() byte {
	f := byte(0x00)

	if c.FlagZero {
		f |= 1 << 7
	}
	if c.FlagSubstract {
		f |= 1 << 6
	}
	if c.FlagHalfCarry {
		f |= 1 << 5
	}
	if c.FlagCarry {
		f |= 1 << 4
	}

	return f
}

// Sets the flags from a byte (ie. write to F register)
func (c *CPU) SetFlags(b byte) {
	c.FlagZero = (b & (1 << 7)) > 0
	c.FlagSubstract = (b & (1 << 6)) > 0
	c.FlagHalfCarry = (b & (1 << 5)) > 0
	c.FlagCarry = (b & (1 << 4)) > 0
}

// Returns the get/set function for a register given by name (eg. 'H')
func (c *CPU) GetRegisterCallbacks(name byte) (get func() byte, set func(byte)) {
	switch name {
	case 'A':
		return c.GetA, c.SetA
	case 'B':
		return c.GetB, c.SetB
	case 'C':
		return c.GetC, c.SetC
	case 'D':
		return c.GetD, c.SetD
	case 'E':
		return c.GetE, c.SetE
	case 'H':
		return c.GetH, c.SetH
	case 'L':
		return c.GetL, c.SetL
	}

	panic("unreachable")
}

func (c *CPU) GetRegisterAddress(name string) *uint16 {
	switch name {
	case "BC":
		return &c.BC
	case "DE":
		return &c.DE
	case "HL":
		return &c.HL
	case "SP":
		return &c.SP
	}

	panic("unreachable")
}

func (c *CPU) String() string {
	flags := [4]byte{'-', '-', '-', '-'}
	if c.FlagZero {
		flags[0] = 'Z'
	}
	if c.FlagSubstract {
		flags[1] = 'N'
	}
	if c.FlagHalfCarry {
		flags[2] = 'H'
	}
	if c.FlagCarry {
		flags[3] = 'C'
	}

	return fmt.Sprintf(
		"A:%02X BC:%04X DE:%04X HL:%04X SP:%04X PC:%04X Flags:%s",
		c.A, c.BC, c.DE, c.HL, c.SP, c.PC, flags,
	)
}
