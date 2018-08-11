package ppu

import "fmt"

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
