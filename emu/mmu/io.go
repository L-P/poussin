package mmu

const (
	IOP1             = uint16(0xFF00) // P1 Joypad (R/W)
	IODIV            = 0xFF04         // Divider Register (R/W*)
	IOSCY            = 0xFF42         // BG Scroll Y (R/W)
	IOSCX            = 0xFF43         // BG Scroll X (R/W)
	IOLY             = 0xFF44         // LCDC Y-Coordinate
	IODisableBootROM = 0xFF50
)

func (m *MMU) ReadIO(addr uint16) byte {
	switch addr {
	case IOLY:
		return 144 // DEBUG
	}

	return 0
}

func (m *MMU) SetIO(addr uint16, value byte) {
	switch addr {
	case IODisableBootROM:
		m.Mem[addr] = 1 // Boot ROM can never be re-enabled
	case IODIV:
		fallthrough
	case IOLY:
		m.Mem[addr] = 0
	}
}
