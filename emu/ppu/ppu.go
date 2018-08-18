package ppu

import (
	"fmt"
	"image"
	"image/color"
)

const (
	DotMatrixWidth  = 160
	DotMatrixHeight = 144
)

type PPU struct {
	VRAM [8192]byte

	// CPU Cycles since last full display
	Cycles int

	// True right _after_ the cycle that got to VBlank ran
	InterruptVBlank bool

	Buffers         [2]*image.RGBA
	BackBufferIndex int

	// We'll send to this when we're ready to display a frame
	NextFrame    chan<- *image.RGBA
	PushedFrames int

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

func New(nextFrame chan<- *image.RGBA) *PPU {
	p := PPU{
		NextFrame: nextFrame,
		Buffers: [2]*image.RGBA{
			image.NewRGBA(image.Rect(0, 0, DotMatrixWidth, DotMatrixHeight)),
			image.NewRGBA(image.Rect(0, 0, DotMatrixWidth, DotMatrixHeight)),
		},
		STAT: 1 << 7, // Bit is always set
	}

	return &p
}

// Runs the PPU for one cycle
func (p *PPU) Cycle() {
	p.Draw()
	p.Cycles = (p.Cycles + 1) % 456
	p.InterruptVBlank = false

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
		} else if p.Cycles > 84 && p.Cycles < 448 {
			p.SetSTATMode(ModeHBlank)
		}
	} else if p.LY == 144 {
		p.SetSTATMode(ModeHBlank)

		if p.Cycles >= 4 {
			if p.Cycles == 4 {
				p.InterruptVBlank = true
				p.SendFrame()
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

// SendFrame sends a frame to the renderer
func (p *PPU) SendFrame() {
	// Send the current back buffer and promote it to front
	// HACK: This call is blocking, thus ensuring we don't run the emulation
	// crazy fast, this only works if the receiver runs at 60hz of course
	// ie. my monitor and graphics card run at 60hz with vsync enabled so
	// this line is the actual simulation speed regulator
	if p.NextFrame != nil { // nil in tests because we don't care about pictures
		p.NextFrame <- p.Buffers[p.BackBufferIndex]
	}

	p.PushedFrames++
	p.BackBufferIndex = (p.BackBufferIndex + 1) % 2
}

func (p *PPU) BackBuffer() *image.RGBA {
	return p.Buffers[p.BackBufferIndex]
}

func (p *PPU) Draw() {
	lcdX := byte(p.Cycles / 2)
	lcdY := p.LY
	if lcdY >= DotMatrixHeight || lcdX >= DotMatrixWidth {
		return
	}

	x := (lcdX + p.SCX) % DotMatrixWidth
	y := (lcdY + p.SCY) % DotMatrixHeight
	p.BackBuffer().SetRGBA(int(x), int(y), Colorize(0))

	if p.LCDCHas(LCDCDisplayBGAndWindow) {
		p.DrawBackground(x, y)
	}
}

func (p *PPU) DrawBackground(x, y byte) {
	dataAddr := p.GetTileDataAddress(x, y)
	tileData := p.GetTileData(dataAddr, x, y)

	c := Colorize(p.Palettize(tileData))

	p.BackBuffer().SetRGBA(int(x), int(y), c)
}

func Colorize(b byte) color.RGBA {
	switch b {
	case 0x00:
		return color.RGBA{224, 248, 208, 255}
	case 0x01:
		return color.RGBA{136, 192, 112, 255}
	case 0x02:
		return color.RGBA{52, 104, 86, 255}
	case 0x03:
		return color.RGBA{8, 24, 32, 255}
	}

	panic("trying to color byte > 3")
}

func (p *PPU) Palettize(b byte) byte {
	switch b {
	case 0x00:
		return p.BGP & 0x03
	case 0x01:
		return (p.BGP & 0x0C) >> 2
	case 0x02:
		return (p.BGP & 0x30) >> 4
	case 0x03:
		return (p.BGP & 0xC0) >> 6
	}

	panic("trying to palette byte > 3")
}

// GetTileMapID returns the tile ID of the tile that should be displayed at pixel x, y
func (p *PPU) GetTileMapID(x, y byte) byte {
	mapOffset, end := p.GetBGTileMapRange()
	tileX, tileY := uint16(x/8), uint16(y/8)

	tileIDAddr := mapOffset + tileX + (tileY * 32)
	if tileIDAddr > end {
		panic("fetching tile outside of tile map")
	}

	return p.FetchVRAM(tileIDAddr)
}

// Returns the adress of the tile data for pixel at x,y
func (p *PPU) GetTileDataAddress(x, y byte) uint16 {
	tileID := p.GetTileMapID(x, y)

	dataOffset, end := p.GetBGWindowTileDataRange()
	addr := dataOffset
	if p.LCDCHas(LCDCBGWindowTileDataSelect) {
		addr += uint16(tileID) * 16
	} else {
		addr += uint16(int16(tileID) * 16)
	}

	if addr > end {
		panic(fmt.Errorf("%04X > %04X for tile #%02X", addr, end, tileID))
	}

	return addr
}

// Returns the tile data for pixel x,y of tile at address addr
func (p *PPU) GetTileData(addr uint16, x, y byte) byte {
	dX := byte(7 - (x % 8))   // one bit per pixel half-color, MSB first
	dY := uint16(2 * (y % 8)) // two bytes per row

	a := p.FetchVRAM(addr + dY)
	b := p.FetchVRAM(addr + dY + 1)

	bit1 := (a & (1 << dX)) >> dX
	bit2 := (b & (1 << dX)) >> dX

	return bit1 | (bit2 << 1)
}
