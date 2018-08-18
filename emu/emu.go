package emu

import (
	"image"

	"github.com/L-P/poussin/emu/cpu"
	"github.com/L-P/poussin/emu/debugger"
	"github.com/L-P/poussin/emu/ppu"
)

type Gameboy struct {
	cpu cpu.CPU
	ppu *ppu.PPU

	debugger *debugger.Debugger
}

func NewGameboy(nextFrame chan<- *image.RGBA) (*Gameboy, error) {
	gb := Gameboy{
		ppu: ppu.New(nextFrame),
	}

	gb.cpu = cpu.New(gb.ppu, true)

	var err error
	gb.debugger, err = debugger.New(&gb.cpu, gb.ppu)
	if err != nil {
		return nil, err
	}

	return &gb, nil
}

func (g *Gameboy) LoadBootROM(rom []byte) error {
	return g.cpu.LoadBootROM(rom)
}

func (g *Gameboy) LoadROM(rom []byte) error {
	return g.cpu.LoadROM(rom)
}

func (g *Gameboy) Run(shouldClose <-chan bool, closed chan<- bool) {
	go g.debugger.RunGUI(shouldClose)
	defer close(closed)

	for !g.debugger.Closed() {
		cycles, err := g.cpu.Step()
		for i := 0; i < cycles; i++ {
			g.ppu.Cycle()
		}

		if g.cpu.EnableDebug {
			g.debugger.Update()

			if err != nil {
				g.debugger.Panic(err)
			}
		} else {
			if err != nil {
				panic(err)
			}
		}
	}
}

func (g *Gameboy) Close() {
	g.debugger.Close()
}

func (g *Gameboy) SimulateBoot() {
	g.cpu.SimulateBoot()
}
