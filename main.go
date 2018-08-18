package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/L-P/poussin/emu"
	"github.com/L-P/poussin/renderer/gl"
)

func main() {
	log.SetOutput(os.Stderr)

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}

func run() error {
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: poussin [-cpuprofile FILE] [-memprofile FILE] [BOOTROM] ROM")
		os.Exit(1)
	}

	nextFrame := make(chan *image.RGBA, 1)
	gb, err := emu.NewGameboy(nextFrame)
	if err != nil {
		return err
	}
	defer gb.Close()

	var bootRomPath string
	var romPath string
	if len(flag.Args()) == 2 {
		bootRomPath = flag.Args()[0]
		romPath = flag.Args()[1]
	} else {
		romPath = flag.Args()[0]
	}

	if bootRomPath != "" {
		bootRom, err := ioutil.ReadFile(bootRomPath)
		if err != nil {
			return err
		}
		if err := gb.LoadBootROM(bootRom); err != nil {
			return err
		}
	} else {
		gb.SimulateBoot()
	}

	rom, err := ioutil.ReadFile(romPath)
	if err != nil {
		return err
	}
	if err := gb.LoadROM(rom); err != nil {
		return err
	}

	r, err := gl.New()
	if err != nil {
		return err
	}
	defer r.Close()

	emuClosed := make(chan bool)
	rendererClosed := make(chan bool)
	closeEmu := make(chan bool)
	closeRenderer := make(chan bool)

	go gb.Run(closeEmu, emuClosed)
	go r.Run(nextFrame, closeRenderer, rendererClosed)

	select {
	case <-emuClosed:
		closeRenderer <- true
		<-rendererClosed
	case <-rendererClosed:
		closeEmu <- true
		<-emuClosed
	}

	return nil
}
