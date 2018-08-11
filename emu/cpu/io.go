package cpu

import (
	"home.leo-peltier.fr/poussin/emu/ppu"
)

const (
	IOP1 = 0xFF00 // P1 Joypad (R/W)
	IOSB = 0xFF01 // Serial transfer data

	IODIV  = 0xFF04 // Divider register (R/W*)
	IOTIMA = 0xFF05 // Timer counter (R/W)
	IOTMA  = 0xFF06 // Timer modulo (R/W)
	IOTAC  = 0xFF07 // TImer control (R/W)

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
	case IOIF:
		return c.FetchIF()
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
		c.Mem[IODIV] = 0
		return
	case IODisableBootROM:
		c.Mem[IODisableBootROM] = 1 // Boot ROM can never be re-enabled
		return
	case IOIF:
		c.WriteIF(value)
		return
	default:
		c.Mem[addr] = value
		// fmt.Printf("unhandled I/O write at %02X\n", addr)
	}
}

func (c *CPU) WriteIF(value byte) {
	c.Mem[IOIF] = value
}

func (c *CPU) FetchIF() byte {
	// First 3 bits are always set
	return 0xE0 | c.Mem[IOIF]
}
