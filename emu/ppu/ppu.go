package ppu

import "fmt"

type PPU struct {
	VRAM [8192]byte

	Cycles int

	// Registers mapped to FF40-FF4B
	LCDC byte
	STAT byte
	SCY  byte
	SCX  byte
	LY   byte
	LYC  byte
	DMA  byte
	BGP  byte
	OBP0 byte
	OBP1 byte
	WY   byte
	WX   byte
}

func New() *PPU {
	return &PPU{
		STAT: 1 << 7, // Bit is always set
	}
}

// Runs the PPU for one cycle
func (p *PPU) Cycle() {
	p.Cycles = (p.Cycles + 1) % 456

	if p.Cycles == 0 {
		p.LY = (p.LY + 1) % 154
	}
	p.SetSTATLYC(false)

	if p.LY >= 0 && p.LY <= 143 {
		p.SetSTATMode(ModeHBlank)

		if p.Cycles == 4 {
			p.SetSTATLYC(p.LY == p.LYC)
		} else if p.Cycles >= 80 && p.Cycles <= 84 {
			p.SetSTATMode(ModeOAM)
		} else if p.Cycles >= 84 && p.Cycles < 448 {
			p.SetSTATMode(ModeHBlank)
		}
	} else if p.LY == 144 {
		p.SetSTATMode(ModeHBlank)

		if p.Cycles >= 4 {
			if p.Cycles == 4 {
				p.SetSTATLYC(p.LY == p.LYC)
			}

			p.SetSTATMode(ModeVBlank)
		}
	} else if p.LY >= 144 && p.LY < 153 {
		p.SetSTATMode(ModeVBlank)
		if p.Cycles == 4 {
			p.SetSTATLYC(p.LY == p.LYC)
		}
	} else if p.LY == 153 {
		p.SetSTATMode(ModeVBlank)
		if p.Cycles == 4 || p.Cycles == 12 {
			p.SetSTATLYC(p.LY == p.LYC)
		}
	}
}

// Bits 0-1 of the STAT register are the mode
const (
	ModeHBlank   = byte(0x00)
	ModeVBlank   = 0x01
	ModeOAM      = 0x02
	ModeTransfer = 0x03
)

func (p *PPU) GetMode() byte {
	return p.STAT & 0x03
}

func (p *PPU) SetSTATMode(mode byte) {
	p.STAT = (p.STAT & 0xFC) | mode
}

func (p *PPU) SetSTATLYC(on bool) {
	if on {
		p.STAT |= 1 << 6
	} else {
		p.STAT &= ^uint8(1 << 6)
	}
}

// Returns true if the address is in the range of the PPU registers
func IsPPUIO(addr uint16) bool {
	return addr >= 0xFF40 && addr <= 0xFF4B
}

// Returns true if the address is in the range of VRAM
func IsVRAM(addr uint16) bool {
	return addr >= 0x8000 && addr <= 0x9FFF
}

func (p *PPU) FetchVRAM(addr uint16) byte {
	if !IsVRAM(addr) {
		panic(fmt.Errorf("PPU VRAM fetch outside of 0x800-0x9FFF: %02X", addr))
	}

	return p.VRAM[addr-0x8000]
}

func (p *PPU) WriteVRAM(addr uint16, b byte) {
	if !IsVRAM(addr) {
		panic(fmt.Errorf("PPU VRAM write outside of 0x800-0x9FFF: %02X", addr))
	}

	p.VRAM[addr-0x8000] = b
}

func (p *PPU) Fetch(addr uint16) byte {
	if IsPPUIO(addr) {
		return p.FetchRegister(addr)
	}

	return p.FetchVRAM(addr)
}

func (p *PPU) Write(addr uint16, b byte) {
	if IsPPUIO(addr) {
		p.WriteRegister(addr, b)
		return
	}

	p.WriteVRAM(addr, b)
}

// Reads a byte from the registers
func (p *PPU) FetchRegister(addr uint16) byte {
	switch addr {
	case 0xFF40:
		return p.LCDC
	case 0xFF41:
		return p.STAT
	case 0xFF42:
		return p.SCY
	case 0xFF43:
		return p.SCX
	case 0xFF44:
		return p.LY
	case 0xFF45:
		return p.LYC
	case 0xFF46:
		return p.DMA
	case 0xFF47:
		return p.BGP
	case 0xFF48:
		return p.OBP0
	case 0xFF49:
		return p.OBP1
	case 0xFF4A:
		return p.WY
	case 0xFF4B:
		return p.WX
	}

	// This is not an emulation problem, somone dun' goofed.
	panic(fmt.Errorf("PPU register fetch outside of 0xFF40-0xFF4B: %02X", addr))
}

func (p *PPU) WriteRegister(addr uint16, b byte) {
	if !IsPPUIO(addr) {
		// This is not an emulation problem, somone dun' goofed.
		panic(fmt.Errorf("PPU register write outside of 0xFF40-0xFF4B: %02X", addr))
	}

	switch addr {
	case 0xFF40:
		p.LCDC = b
	case 0xFF41:
		p.STAT = (b & 0x78) | p.STAT // mask 0b01111000
	case 0xFF42:
		p.SCY = b
	case 0xFF43:
		p.SCX = b
	case 0xFF44:
		p.LY = 0
	case 0xFF45:
		p.LYC = b
	case 0xFF46:
		p.DMA = b
	case 0xFF47:
		p.BGP = b
	case 0xFF48:
		p.OBP0 = b
	case 0xFF49:
		p.OBP1 = b
	case 0xFF4A:
		p.WY = b
	case 0xFF4B:
		p.WX = b
	}
}