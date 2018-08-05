package cpu

func (c *CPU) StackPush16b(data ...uint16) {
	for _, v := range data {
		c.StackPush8b(
			uint8((v & 0xFF00 >> 8)),
			uint8(v&0x00FF),
		)
	}
}

func (c *CPU) StackPush8b(data ...byte) {
	for _, v := range data {
		c.SP--
		c.MMU.Set8b(c.SP, v)
	}
}

func (c *CPU) StackPop16b() uint16 {
	return uint16(c.StackPop8b()) | (uint16(c.StackPop8b()) << 8)
}

func (c *CPU) StackPop8b() byte {
	c.SP++
	return c.MMU.Get8b(c.SP - 1)
}
