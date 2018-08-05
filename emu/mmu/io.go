package mmu

import "fmt"

const (
	IOP1             = 0xFF00 // P1 Joypad (R/W)
	IODIV            = 0xFF04 // Divider Register (R/W*)
	IOSCY            = 0xFF42 // BG Scroll Y (R/W)
	IOSCX            = 0xFF43 // BG Scroll X (R/W)
	IOLY             = 0xFF44 // LCDC Y-Coordinate
	IODisableBootROM = 0xFF50
)

func (m *MMU) ReadIO(addr uint16) byte {
	switch addr {
	case IODisableBootROM:
		return m.Mem[addr]
	}

	panic(fmt.Errorf("unhandled I/O read at %02X", addr))
}

func (m *MMU) SetIO(addr uint16, value byte) {
	switch addr {
	case IODisableBootROM:
		m.Mem[addr] = 1 // Boot ROM can never be re-enabled
		return
	}

	panic(fmt.Errorf("unhandled I/O read at %02X", addr))
}
