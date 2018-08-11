package ppu

const (
	LCDCDisplayBGAndWindow     = 1 << 0
	LCDCDisplaySprite          = 1 << 1
	LCDCSpriteSize             = 1 << 2
	LCDCBGTileMapSelect        = 1 << 3
	LCDCBGWindowTileDataSelect = 1 << 4
	LCDCBGWindowDisplay        = 1 << 5
	LCDCBGWindowTileMapSelect  = 1 << 6
	LCDCControl                = 1 << 7
)

func (p *PPU) GetWindowTileMapRange() (uint16, uint16) {
	if p.LCDCHas(LCDCBGWindowTileMapSelect) {
		return 0x9C00, 0x9FFF
	}

	return 0x9800, 0x9BFF
}

func (p *PPU) GetBGWindowTileDataRange() (uint16, uint16) {
	if p.LCDCHas(LCDCBGWindowTileDataSelect) {
		return 0x8000, 0x8FFF
	}

	return 0x8800, 0x97FF
}

func (p *PPU) GetBGTileMapRange() (uint16, uint16) {
	if p.LCDCHas(LCDCBGTileMapSelect) {
		return 0x9C00, 0x9FFF
	}

	return 0x9800, 0x9BFF
}

func (p *PPU) GetSpriteSize() (byte, byte) {
	if p.LCDCHas(LCDCSpriteSize) {
		return 8, 16
	}

	return 8, 8
}

func (p *PPU) LCDCHas(mask byte) bool {
	return (p.LCDC & mask) > 0
}
