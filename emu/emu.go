package emu

import (
	"image"

	"github.com/L-P/poussin/emu/cpu"
	"github.com/L-P/poussin/emu/debugger"
	"github.com/L-P/poussin/emu/ppu"
)

// Gameboy is a DMG emulator.
type Gameboy struct {
	cpu cpu.CPU
	ppu *ppu.PPU

	debugger *debugger.Debugger
}

// NewGameboy creates a new Gameboy.
func NewGameboy(
	nextFrame chan<- *image.RGBA,
	input <-chan cpu.JoypadState,
) (*Gameboy, error) {
	gb := Gameboy{
		ppu: ppu.New(nextFrame),
	}

	gb.cpu = cpu.New(gb.ppu, input, true)

	var err error
	gb.debugger, err = debugger.New(&gb.cpu, gb.ppu)
	if err != nil {
		return nil, err
	}

	return &gb, nil
}

// LoadBootROM puts a boot rom in the 256 first bytes or RAM.
func (g *Gameboy) LoadBootROM(rom []byte) error {
	return g.cpu.LoadBootROM(rom)
}

// LoadROM loads a ROM in RAM.
func (g *Gameboy) LoadROM(rom []byte) error {
	return g.cpu.LoadROM(rom)
}

// Run the emulation and the debugger.
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

// Close frees up all resources used by the emulator.
func (g *Gameboy) Close() {
	g.debugger.Close()
}

// SimulateBoot puts the CPU in the same state it would be after running the Nintendo boot ROM.
func (g *Gameboy) SimulateBoot() {
	g.cpu.SimulateBoot()
}
