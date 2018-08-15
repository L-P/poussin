package main

import (
	"io/ioutil"
	"testing"

	"home.leo-peltier.fr/poussin/emu/cpu"
	"home.leo-peltier.fr/poussin/emu/ppu"
)

func TestCPUInstrs01(t *testing.T) {
	cpu := getCPU(t, "./tests/cpu_instrs/individual/01-special.gb")
	for cpu.PC != 0xC7D2 {
		cycles, err := cpu.Step()
		if err != nil {
			t.Fatalf("%s", err)
		}
		cpu.MemIOBuffer.Reset()
		for i := 0; i < cycles; i++ {
			cpu.PPU.Cycle()
		}
	}

	buf := cpu.SBBuffer.String()
	if buf[len(buf)-7:len(buf)-1] != "Passed" {
		t.Error(buf)
	}
}

func getCPU(t *testing.T, romPath string) cpu.CPU {
	c := cpu.New(ppu.New(nil), true)
	c.SimulateBoot()

	rom, err := ioutil.ReadFile(romPath)
	if err != nil {
		t.Fatalf("unable to read ROM %s: %s", romPath, err)
	}
	c.LoadROM(rom)

	return c
}
