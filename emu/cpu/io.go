package cpu

import (
	"home.leo-peltier.fr/poussin/emu/ppu"
)

const (
	IOP1 = 0xFF00 // P1 Joypad (R/W)
	IOSB = 0xFF01 // Serial transfer data

	IODIV  = 0xFF04 // Divider register (R/W*), see InternalDIV
	IOTIMA = 0xFF05 // Timer counter (R/W)
	IOTMA  = 0xFF06 // Timer modulo (R/W)
	IOTAC  = 0xFF07 // Timer control (R/W)

	IODisableBootROM = 0xFF50
	IOIF             = 0xFF0F // Interrupt flag

	// {{{ Sound registers
	IONR10 = 0xFF10
	IONR11 = 0xFF11
	IONR12 = 0xFF12
	IONR13 = 0xFF13
	IONR14 = 0xFF14

	IONR21 = 0xFF16
	IONR22 = 0xFF17
	IONR23 = 0xFF18
	IONR24 = 0xFF19
	IONR30 = 0xFF1A

	IONR31 = 0xFF1B
	IONR32 = 0xFF1C
	IONR33 = 0xFF1D
	IONR34 = 0xFF1E

	IONR41 = 0xFF20
	IONR42 = 0xFF21
	IONR43 = 0xFF22
	IONR44 = 0xFF23

	IONR50 = 0xFF24
	IONR51 = 0xFF25
	IONR52 = 0xFF26

	IOWaveStart = 0xFF30
	IOWaveEnd   = 0xFF3F
	// }}} Sound registers
)

func (c *CPU) FetchIO(addr uint16) byte {
	if ppu.IsPPUIO(addr) {
		return c.PPU.Fetch(addr)
	}

	switch addr {
	case IODisableBootROM:
		return c.Mem[IODisableBootROM]
	case IOP1:
		return 0x0F // TODO handle input
	case IOIF:
		return c.FetchIF()
	case IODIV:
		return byte((c.InternalDIV & 0xFF00) >> 8)
	case IOTAC:
		return c.FetchTAC()
	default:
		// fmt.Printf("unhandled I/O fetch at %02X\n", addr)
		return c.Mem[addr]
	}
}

func (c *CPU) WriteIO(addr uint16, value byte) {
	if ppu.IsPPUIO(addr) {
		c.PPU.Write(addr, value)
		return
	}

	switch addr {
	case IODIV:
		c.InternalDIV = 0
	case IOSB:
		if c.EnableDebug {
			c.SBBuffer.WriteByte(value)
		}
		c.Mem[IOSB] = value
	case IODisableBootROM:
		c.Mem[IODisableBootROM] = 1 // Boot ROM can never be re-enabled
	case IOIF:
		c.WriteIF(value)
	case IOTAC:
		c.WriteTAC(value)
	default:
		c.Mem[addr] = value
		// fmt.Printf("unhandled I/O write at %02X\n", addr)
	}
}

func (c *CPU) WriteIF(value byte) {
	c.Mem[IOIF] = 0xE0 | value
}

func (c *CPU) FetchIF() byte {
	// First 3 bits are always set
	return c.Mem[IOIF] | 0xE0
}

func (c *CPU) IFIsSet(mask byte) bool {
	return c.Mem[IOIF]&mask == mask
}

func (c *CPU) SetIF(mask byte) {
	c.Mem[IOIF] |= mask
}

func (c *CPU) UnSetIF(mask byte) {
	c.Mem[IOIF] &^= mask
}

const (
	TACEnable    = (1 << 2)
	TACSpeed4    = 0x00 // 00: 4.096 KHz, every 1024 cycle
	TACSpeed262  = 0x01 // 01: 262.144 Khz, every 16 cycle
	TACSpeed65   = 0x02 // 10: 65.536 KHz, every 64 cycle
	TACSpeed16   = 0x03 // 11: 16.384 KHz, every 256 cycle
	TACSpeedMask = 0x03
)

func (c *CPU) WriteTAC(value byte) {
	c.Mem[IOTAC] = value | 0xF8
}

func (c *CPU) IsTACEnabled() bool {
	return c.Mem[IOTAC]&TACEnable == TACEnable
}

func (c *CPU) SetTACEnabled() {
	c.Mem[IOTAC] |= TACEnable
}

func (c *CPU) UnSetTACEnabled() {
	c.Mem[IOTAC] &^= TACEnable
}

func (c *CPU) FetchTAC() byte {
	// First 3 bits are always set
	return c.Mem[IOTAC] | 0xF8
}
