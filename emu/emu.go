package emu

import (
	"home.leo-peltier.fr/poussin/emu/cpu"
)

type Gameboy struct {
	CPU cpu.CPU
}

func NewGameboy() *Gameboy {
	return &Gameboy{
		CPU: cpu.New(),
	}
}

func (g *Gameboy) LoadBootROM(rom []byte) error {
	return g.CPU.LoadBootROM(rom)
}

func (g *Gameboy) LoadROM(rom []byte) error {
	return g.CPU.LoadROM(rom)
}

func (g *Gameboy) Run() error {
	for true {
		if err := g.CPU.Step(); err != nil {
			return err
		}
	}

	return nil
}
