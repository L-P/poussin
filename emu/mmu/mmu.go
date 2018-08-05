package mmu

import "fmt"

type MMU struct {
	Mem [0xFFFF]byte
}

func New() MMU {
	return MMU{}
}

func (m *MMU) LoadBootROM(rom []byte) error {
	if len(rom) != 256 {
		return fmt.Errorf("boot ROM should be 256 bytes, got %d", len(rom))
	}

	if count := copy(m.Mem[:], rom); count != 256 {
		return fmt.Errorf("copied less than 256 bytes: %d", count)
	}

	return nil
}

// Writes a single byte of data
func (m *MMU) Set8b(addr uint16, value byte) {
	m.Mem[addr] = value
}

// Peek returns a single byte from raw memory, should not be used by instructions.
func (m *MMU) Peek(addr uint16) byte {
	return m.Mem[addr]
}
