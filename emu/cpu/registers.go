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
