package cpu

// Fetch reads a byte from mapped memory
func (c *CPU) Fetch(addr uint16) byte {
	var v byte
	switch AddrToMemType(addr) {
	case ROM0:
		v = c.FetchROM0(addr)
	case ROMX:
		v = c.FetchROMX(addr)
	case VRAM:
		v = c.PPU.Fetch(addr)
	case IO:
		v = c.FetchIO(addr)
	case IERegister:
		v = c.FetchIE()
	default:
		v = c.Mem[addr]
	}

	if c.EnableDebug && c.InCycle {
		c.MemIOBuffer.WriteByte(byte(c.PC & 0x00FF))
		c.MemIOBuffer.WriteByte(byte((c.PC & 0xFF00) >> 8))
		c.MemIOBuffer.WriteByte(0x01)
		c.MemIOBuffer.WriteByte(byte(addr & 0x00FF))
		c.MemIOBuffer.WriteByte(byte((addr & 0xFF00) >> 8))
		c.MemIOBuffer.WriteByte(v)
	}

	return v
}

// Writes a byte to mapped memory
func (c *CPU) Write(addr uint16, b byte) {
	if c.EnableDebug && c.InCycle {
		c.MemIOBuffer.WriteByte(byte(c.PC & 0x00FF))
		c.MemIOBuffer.WriteByte(byte((c.PC & 0xFF00) >> 8))
		c.MemIOBuffer.WriteByte(0x02)
		c.MemIOBuffer.WriteByte(byte(addr & 0x00FF))
		c.MemIOBuffer.WriteByte(byte((addr & 0xFF00) >> 8))
		c.MemIOBuffer.WriteByte(b)
	}

	switch AddrToMemType(addr) {
	case IO:
		c.WriteIO(addr, b)
	case VRAM:
		c.PPU.Write(addr, b)
	case IERegister:
		c.WriteIE(b)
	default:
		c.Mem[addr] = b // TODO remove default case
	}
}

func (c *CPU) FetchROM0(addr uint16) byte {
	// During bootstrap 0x0000-0x00FF is mapped to boot ROM
	if c.FetchIO(IODisableBootROM) == 0 && addr <= 0xFF {
		return c.Boot[addr]
	}

	return c.ROM[addr]
}

// TODO: bank switch
func (c *CPU) FetchROMX(addr uint16) byte {
	return c.ROM[addr]
}

type MemType int

const (
	ROM0       = MemType(iota) // Non-switchable ROM bank
	ROMX                       // Switchable ROM bank
	VRAM                       // VRAM Vido RAM, switchable (0-1) in GBC mode
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

const (
	IEVBlank  = 1 << 0
	IELCDSTAT = 1 << 1
	IETimer   = 1 << 2
	IESerial  = 1 << 3
	IEJoypad  = 1 << 4
)

// AddrToMemType returns the MemType the given address belongs to
func AddrToMemType(addr uint16) MemType {
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF:
		return ROM0
	case addr >= 0x4000 && addr <= 0x7FFF:
		return ROMX
	case addr >= 0x8000 && addr <= 0x9FFF:
		return VRAM
	case addr >= 0xA000 && addr <= 0xBFFF:
		return SRAM
	case addr >= 0xC000 && addr <= 0xCFFF:
		return WRAM0
	case addr >= 0xD000 && addr <= 0xDFFF:
		return WRAMX
	case addr >= 0xE000 && addr <= 0xFDFF:
		return Echo
	case addr >= 0xFE00 && addr <= 0xFE9F:
		return OAM
	case addr >= 0xFEA0 && addr <= 0xFEFF:
		return Unused
	case addr >= 0xFF00 && addr <= 0xFF7F:
		return IO
	case addr >= 0xFF80 && addr <= 0xFFFE:
		return HRAM
	case addr >= 0xFFFF && addr <= 0xFFFF:
		return IERegister
	}

	panic("unreachable")
}

// MemTypeName returns a MemType name as a string
func MemTypeName(t MemType) string {
	switch t {
	case ROM0:
		return "ROM0"
	case ROMX:
		return "ROMX"
	case VRAM:
		return "VRAM"
	case SRAM:
		return "SRAM"
	case WRAM0:
		return "WRAM0"
	case WRAMX:
		return "WRAMX"
	case Echo:
		return "Echo"
	case OAM:
		return "OAM"
	case Unused:
		return "Unused"
	case IO:
		return "IO"
	case HRAM:
		return "HRAM"
	case IERegister:
		return "IERegister"
	}

	panic("unreachable")
}

func (c *CPU) WriteIE(value byte) {
	c.InterruptEnable = value
}

func (c *CPU) FetchIE() byte {
	return c.InterruptEnable
}

func (c *CPU) IEEnabled(ie byte) bool {
	return (c.FetchIE() & ie) != 0x00
}
