package ppu

import (
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
		} else if p.Cycles >= 84 && p.Cycles < 448 {
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
	if lcdY > DotMatrixHeight || lcdX > DotMatrixWidth {
		return
	}

	x := (lcdX + p.SCX) % DotMatrixWidth
	y := (lcdY + p.SCY) % DotMatrixHeight

	dataAddr := p.GetTileDataAddress(x, y)
	tileData := p.GetTileData(dataAddr, x, y)

	c := Colorize(p.Palettize(tileData))

	p.BackBuffer().SetRGBA(int(x), int(y), color.RGBA{c, c, c, 255})
}

func Colorize(b byte) byte {
	switch b {
	case 0x00:
		return 0xFF
	case 0x01:
		return 0xC0
	case 0x02:
		return 0x40
	case 0x03:
		return 0x00
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

// Returns the adress of the tile data for pixel at x,y
func (p *PPU) GetTileDataAddress(x, y byte) uint16 {
	mapOffset, end := p.GetBGTileMapRange()
	tileX, tileY := uint16(x/8), uint16(y/8)
	tileIDAddr := mapOffset + tileX + (tileY * 8)
	if tileIDAddr > end {
		panic("fetching tile outside of tile map")
	}

	tileID := p.FetchVRAM(tileIDAddr)

	dataOffset, _ := p.GetBGWindowTileDataRange()
	addr := dataOffset + uint16(tileID)

	return addr
}

// Returns the tile data for pixel x,y of tile at address addr
func (p *PPU) GetTileData(addr uint16, x, y byte) byte {
	dX, dY := byte(x%8), byte(y%8)
	b := p.FetchVRAM(addr + uint16(dY))

	return (b >> dX) & 0x03
}
