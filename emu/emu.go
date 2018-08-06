package emu

import (
	"fmt"
	"time"

	"home.leo-peltier.fr/poussin/emu/cpu"
)

type Gameboy struct {
	CPU cpu.CPU
}

func NewGameboy(nextFrame chan<- int) *Gameboy {
	return &Gameboy{
		CPU: cpu.New(nextFrame),
	}
}

func (g *Gameboy) LoadBootROM(rom []byte) error {
	return g.CPU.LoadBootROM(rom)
}

func (g *Gameboy) LoadROM(rom []byte) error {
	return g.CPU.LoadROM(rom)
}

func (g *Gameboy) Run() {
	lastOPsTime := time.Now()
	lastOPs := 0

	for true {
		if err := g.CPU.Step(); err != nil {
			panic(err)
		}

		delta := time.Now().Sub(lastOPsTime)
		if delta >= time.Duration(1*time.Second) {
			fmt.Printf("OP/s : %d\n", g.CPU.OPCount-lastOPs)
			lastOPsTime = time.Now()
			lastOPs = g.CPU.OPCount
		}
	}
}
