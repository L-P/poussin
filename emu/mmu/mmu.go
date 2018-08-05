package mmu

import "fmt"

type MMU struct {
	Mem [0xFFFF]byte
}

type Type int

const (
	ROM0       = Type(iota) // Non-switchable ROM bank
	ROMX                    // Switchable ROM bank
	VRAM                    // Vido RAM, switchable (0-1) in GBC mode
	SRAM                    // External RAM in cartridge
	WRAM0                   // Work ram
	WRAMX                   // Work ram, switchable (1-7) in GBC mode
	Echo                    // Mirrors other parts of RAM depending on mode
	OAM                     // Object attribute table (sprite info)
	Unused                  // Behavior depends on mode
	IO                      // I/O registers
	HRAM                    // Internal CPU RAM
	IERegister              // Interrupt enable flags
)

func MemoryTypeName(t Type) string {
	types := map[Type]string{
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

func AddressToMemoryType(addr uint16) Type {
	ranges := map[Type][2]uint16{
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

// Returns a single byte of data
func (m *MMU) Get8b(addr uint16) byte {
	fmt.Printf("Read at %04X from %s\n", addr, MemoryTypeName(AddressToMemoryType(addr)))
	switch AddressToMemoryType(addr) {
	case IO:
		return m.ReadIO(addr)
	default:
		return m.Mem[addr]
	}
}

// Peek returns a single byte from raw memory, should not be used by instructions.
func (m *MMU) Peek(addr uint16) byte {
	return m.Mem[addr]
}
