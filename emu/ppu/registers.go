package ppu

import "fmt"

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
