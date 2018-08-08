package emu

import (
	"image"

	"home.leo-peltier.fr/poussin/emu/cpu"
	"home.leo-peltier.fr/poussin/emu/debugger"
	"home.leo-peltier.fr/poussin/emu/ppu"
)

type Gameboy struct {
	cpu cpu.CPU
	ppu *ppu.PPU

	debugger debugger.Debugger
}

func NewGameboy(nextFrame chan<- *image.RGBA) *Gameboy {
	gb := Gameboy{
		ppu: ppu.New(nextFrame),
	}
	gb.cpu = cpu.New(gb.ppu)
	gb.debugger = debugger.New(&gb.cpu, gb.ppu)

	return &gb
}

func (g *Gameboy) LoadBootROM(rom []byte) error {
	return g.cpu.LoadBootROM(rom)
}

func (g *Gameboy) LoadROM(rom []byte) error {
	return g.cpu.LoadROM(rom)
}

func (g *Gameboy) Run(quit chan int) {
	go g.debugger.Run(quit)

	for true {
		cycles, err := g.cpu.Step()
		for i := 0; i < cycles; i++ {
			g.ppu.Cycle()
		}
		g.debugger.Update()

		if err != nil {
			g.debugger.Update()
			g.debugger.Panic(err)

			// Wait for someone to tell us to quit (ie. the renderer)
			select {
			case <-quit:
				return
			}
			return
		}

		select {
		case <-quit:
			return
		default:
		}
	}
}

func (g *Gameboy) Close() {
	g.debugger.Close()
}
