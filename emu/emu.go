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

	debugger *debugger.Debugger
}

func NewGameboy(nextFrame chan<- *image.RGBA) (*Gameboy, error) {
	gb := Gameboy{
		ppu: ppu.New(nextFrame),
	}

	gb.cpu = cpu.New(gb.ppu)

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

	for !g.debugger.Closed() {
		cycles, err := g.cpu.Step()
		for i := 0; i < cycles; i++ {
			g.ppu.Cycle()
		}
		g.debugger.Update()

		if err != nil {
			g.debugger.Panic(err)
		}
	}

	close(closed)
}

func (g *Gameboy) Close() {
	g.debugger.Close()
}
