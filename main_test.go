package main

import (
	"io/ioutil"
	"sync"
	"testing"

	"github.com/L-P/poussin/emu/cpu"
	"github.com/L-P/poussin/emu/ppu"
)

func TestCPUInstrs(t *testing.T) {
	roms := []struct {
		until uint16
		path  string
	}{
		{0xC7D2, "./tests/cpu_instrs/individual/01-special.gb"},
		{0xC7F4, "./tests/cpu_instrs/individual/02-interrupts.gb"},
		{0xCB44, "./tests/cpu_instrs/individual/03-op sp,hl.gb"},
		{0xCB35, "./tests/cpu_instrs/individual/04-op r,imm.gb"},
		{0xCB31, "./tests/cpu_instrs/individual/05-op rp.gb"},
		{0xCC5F, "./tests/cpu_instrs/individual/06-ld r,r.gb"},
		{0xCBB0, "./tests/cpu_instrs/individual/07-jr,jp,call,ret,rst.gb"},
		{0xCB91, "./tests/cpu_instrs/individual/08-misc instrs.gb"},
		{0xCE67, "./tests/cpu_instrs/individual/09-op r,r.gb"},
		{0xCF58, "./tests/cpu_instrs/individual/10-bit ops.gb"},
		{0xCC62, "./tests/cpu_instrs/individual/11-op a,(hl).gb"},
	}

	var wg sync.WaitGroup
	wg.Add(len(roms))

	for _, v := range roms {
		go func(path string, until uint16) {
			runTest(t, path, until)
			wg.Done()
		}(v.path, v.until)
	}

	wg.Wait()
}

func runTest(t *testing.T, romPath string, until uint16) {
	cpu := getCPU(t, romPath)
	for cpu.PC != until {
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
		t.Errorf("test failed for rom %s: %s", romPath, buf)
	}
}

func getCPU(t *testing.T, romPath string) cpu.CPU {
	c := cpu.New(ppu.New(nil), nil, true)
	c.SimulateBoot()

	rom, err := ioutil.ReadFile(romPath)
	if err != nil {
		t.Fatalf("unable to read ROM %s: %s", romPath, err)
	}
	c.LoadROM(rom)

	return c
}
