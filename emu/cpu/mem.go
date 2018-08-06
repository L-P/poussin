package cpu

import "fmt"

// Reads a byte from mapped memory
func (c *CPU) Fetch(addr uint16) byte {
	switch AddrToMemType(addr) {
	case ROM0:
		return c.FetchROM0(addr)
	case IO:
		return c.FetchIO(addr)
	default:
		return c.Mem[addr]
	}
}

// Writes a byte to mapped memory
func (c *CPU) Write(addr uint16, b byte) {
	if AddrToMemType(addr) == IERegister {
		panic(fmt.Errorf("Set IERegister %02X\n", b))
	}

	c.Mem[addr] = b
}

func (c *CPU) FetchROM0(addr uint16) byte {
	// During bootstrap 0x0000-0x00FF is mapped to boot ROM
	if c.FetchIO(IODisableBootROM) == 0 && addr <= 0xFF {
		return c.Boot[addr]
	}

	return c.ROM[addr]
}

type MemType int

const (
	ROM0       = MemType(iota) // Non-switchable ROM bank
	ROMX                       // Switchable ROM bank
	VRAM                       // Vido RAM, switchable (0-1) in GBC mode
	SRAM                       // External RAM in cartridge
	WRAM0                      // Work ram
	WRAMX                      // Work ram, switchable (1-7) in GBC mode
	Echo                       // Mirrors other parts of RAM depending on mode
	OAM                        // Object attribute table (sprite info)
	Unused                     // Behavior depends on mode
	IO                         // I/O registers
	HRAM                       // Internal CPU RAM
	IERegister                 // Interrupt enable flags
)

// MemTypeName returns a MemType name as a string
func MemTypeName(t MemType) string {
	types := map[MemType]string{
		ROM0:       "ROM0",
		ROMX:       "ROMX",
		VRAM:       "VRAM",
		SRAM:       "SRAM",
		WRAM0:      "WRAM0",
		WRAMX:      "WRAMX",
		Echo:       "Echo",
		OAM:        "OAM",
		Unused:     "Unused",
		IO:         "IO",
		HRAM:       "HRAM",
		IERegister: "IERegister",
	}

	return types[t]
}

// AddrToMemType returns the MemType the given address belongs to
func AddrToMemType(addr uint16) MemType {
	ranges := map[MemType][2]uint16{
		ROM0:       {0x0000, 0x3FFF},
		ROMX:       {0x4000, 0x7FFF},
		VRAM:       {0x8000, 0x9FFF},
		SRAM:       {0xA000, 0xBFFF},
		WRAM0:      {0xC000, 0xCFFF},
		WRAMX:      {0xD000, 0xDFFF},
		Echo:       {0xE000, 0xFDFF},
		OAM:        {0xFE00, 0xFE9F},
		Unused:     {0xFEA0, 0xFEFF},
		IO:         {0xFF00, 0xFF7F},
		HRAM:       {0xFF80, 0xFFFE},
		IERegister: {0xFFFF, 0xFFFF},
	}

	for k, v := range ranges {
		if addr >= v[0] && addr <= v[1] {
			return k
		}
	}

	panic("unreachable")
}

const (
	IOP1             = 0xFF00 // P1 Joypad (R/W)
	IODIV            = 0xFF04 // Divider Register (R/W*)
	IOSCY            = 0xFF42 // BG Scroll Y (R/W)
	IOSCX            = 0xFF43 // BG Scroll X (R/W)
	IOLY             = 0xFF44 // LCDC Y-Coordinate
	IODisableBootROM = 0xFF50
)

func (c *CPU) FetchIO(addr uint16) byte {
	switch addr {
	case IODisableBootROM:
		return c.Mem[addr]
	}

	panic(fmt.Errorf("unhandled I/O read at %02X", addr))
}

func (c *CPU) SetIO(addr uint16, value byte) {
	switch addr {
	case IODisableBootROM:
		c.Mem[addr] = 1 // Boot ROM can never be re-enabled
		return
	}

	panic(fmt.Errorf("unhandled I/O read at %02X", addr))
}