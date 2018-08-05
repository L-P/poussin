package emu

import (
	"fmt"

	"home.leo-peltier.fr/poussin/emu/cpu"
)

type Gameboy struct {
	ROM []byte
	CPU cpu.CPU
}

func NewGameboy() *Gameboy {
	return &Gameboy{}
}

func (g *Gameboy) LoadRom(rom []byte) {
	g.ROM = make([]byte, len(rom))
	copy(g.ROM, rom)
	fmt.Printf("ROM Loaded: %d bytes\n", len(g.ROM))
}

func (g *Gameboy) LoadBootROM(rom []byte) error {
	return g.CPU.MMU.LoadBootROM(rom)
}

func (g *Gameboy) LoadROM(rom []byte) error {
	return g.CPU.MMU.LoadROM(rom)
}

func (g *Gameboy) Run() error {
	for true {
		if err := g.CPU.Step(); err != nil {
			return err
		}
	}

	return nil
}
