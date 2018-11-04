package cpu

import (
	"github.com/L-P/poussin/emu/ppu"
)

const (
	// IOP1 Joypad (R/W)
	IOP1 = 0xFF00

	// IOSB Serial transfer data
	IOSB = 0xFF01

	// IODisableBootROM (R/W once)
	IODisableBootROM = 0xFF50

	// IOIF Interrupt flag
	IOIF = 0xFF0F

	// {{{ Timer registers

	// IODIV Divider register (R/W*), see InternalDIV
	IODIV = 0xFF04

	// IOTIMA Timer counter (R/W)
	IOTIMA = 0xFF05

	// IOTMA Timer modulo (R/W)
	IOTMA = 0xFF06

	// IOTAC Timer control (R/W)
	IOTAC = 0xFF07

	// }}} Timer registers

	// {{{ Sound registers

	// IONR10 TODO
	IONR10 = 0xFF10

	// IONR11 TODO
	IONR11 = 0xFF11

	// IONR12 TODO
	IONR12 = 0xFF12

	// IONR13 TODO
	IONR13 = 0xFF13

	// IONR14 TODO
	IONR14 = 0xFF14

	// IONR21 TODO
	IONR21 = 0xFF16

	// IONR22 TODO
	IONR22 = 0xFF17

	// IONR23 TODO
	IONR23 = 0xFF18

	// IONR24 TODO
	IONR24 = 0xFF19

	// IONR30 TODO
	IONR30 = 0xFF1A

	// IONR31 TODO
	IONR31 = 0xFF1B

	// IONR32 TODO
	IONR32 = 0xFF1C

	// IONR33 TODO
	IONR33 = 0xFF1D

	// IONR34 TODO
	IONR34 = 0xFF1E

	// IONR41 TODO
	IONR41 = 0xFF20

	// IONR42 TODO
	IONR42 = 0xFF21

	// IONR43 TODO
	IONR43 = 0xFF22

	// IONR44 TODO
	IONR44 = 0xFF23

	// IONR50 TODO
	IONR50 = 0xFF24

	// IONR51 TODO
	IONR51 = 0xFF25

	// IONR52 TODO
	IONR52 = 0xFF26

	// IOWaveStart TODO
	IOWaveStart = 0xFF30

	// IOWaveEnd   TODO
	IOWaveEnd = 0xFF3F

	// }}} Sound registers
)

// FetchIO returns a byte value from a hardware register.
func (c *CPU) FetchIO(addr uint16) byte {
	if ppu.IsPPUIO(addr) {
		return c.PPU.Fetch(addr)
	}

	switch addr {
	case IODisableBootROM:
		return c.Mem[IODisableBootROM]
	case IOP1:
		return c.FetchIOP1()
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

// WriteIOP1 writes to the Joypad register 2 out bits.
func (c *CPU) WriteIOP1(b byte) {
	const mask byte = 0x30
	c.Mem[IOP1] = (c.Mem[IOP1] &^ mask) | (b & mask) | 0xC0
}

// FetchIOP1 returns the joypad register state.
func (c *CPU) FetchIOP1() byte {
	var (
		b0 byte = 1
		b1 byte = 1
		b2 byte = 1
		b3 byte = 1
	)

	// Buttons
	if c.Mem[IOP1]&(1<<4) == 0 {
		if c.Joypad.A {
			b0 = 0
		}
		if c.Joypad.B {
			b1 = 0
		}
		if c.Joypad.Select {
			b2 = 0
		}
		if c.Joypad.Start {
			b3 = 0
		}
	}

	// Arrows
	if c.Mem[IOP1]&(1<<5) == 0 {
		if c.Joypad.Right {
			b0 = 0
		}
		if c.Joypad.Left {
			b1 = 0
		}
		if c.Joypad.Up {
			b2 = 0
		}
		if c.Joypad.Down {
			b3 = 0
		}
	}

	return 0xC0 | (c.Mem[IOP1] & 0xF0) | (b0 << 0) | (b1 << 1) | (b2 << 2) | (b3 << 3)
}

// WriteIO writes a byte to a hardware register.
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
	case IOP1:
		c.WriteIOP1(value)
	default:
		c.Mem[addr] = value
		// fmt.Printf("unhandled I/O write at %02X\n", addr)
	}
}

// WriteIF writes the interrupt flag value.
func (c *CPU) WriteIF(value byte) {
	c.Mem[IOIF] = 0xE0 | value
}

// FetchIF returns the interrupt flag value.
func (c *CPU) FetchIF() byte {
	// First 3 bits are always set
	return c.Mem[IOIF] | 0xE0
}

// IFIsSet returns true if the given interrupt mask/flag is set.
func (c *CPU) IFIsSet(mask byte) bool {
	return c.Mem[IOIF]&mask == mask
}

// SetIF sets a single interrupt flag/mask.
func (c *CPU) SetIF(mask byte) {
	c.Mem[IOIF] |= mask
}

// UnSetIF unsets a single interrupt flag/mask.
func (c *CPU) UnSetIF(mask byte) {
	c.Mem[IOIF] &^= mask
}

const (
	// TACEnable is TAC enable mask.
	TACEnable = (1 << 2)

	// TACSpeed4 00: 4.096 KHz, every 1024 cycle
	TACSpeed4 = 0x00

	// TACSpeed262 01: 262.144 Khz, every 16 cycle
	TACSpeed262 = 0x01

	// TACSpeed65 10: 65.536 KHz, every 64 cycle
	TACSpeed65 = 0x02

	// TACSpeed16 11: 16.384 KHz, every 256 cycle
	TACSpeed16 = 0x03

	// TACSpeedMask wher the speed information is in the TAC register.
	TACSpeedMask = 0x03
)

// WriteTAC writes to the timer interrupt control register.
func (c *CPU) WriteTAC(value byte) {
	c.Mem[IOTAC] = value | 0xF8
}

// IsTACEnabled returns true if the timer interrupt is enabled.
func (c *CPU) IsTACEnabled() bool {
	return c.Mem[IOTAC]&TACEnable == TACEnable
}

// SetTACEnabled enables the timer interrupt.
func (c *CPU) SetTACEnabled() {
	c.Mem[IOTAC] |= TACEnable
}

// UnSetTACEnabled disables the timer interrupt.
func (c *CPU) UnSetTACEnabled() {
	c.Mem[IOTAC] &^= TACEnable
}

// FetchTAC returns the value of the timer interrupt control register.
func (c *CPU) FetchTAC() byte {
	// First 3 bits are always set
	return c.Mem[IOTAC] | 0xF8
}
