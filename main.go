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

	"home.leo-peltier.fr/poussin/emu"
	"home.leo-peltier.fr/poussin/renderer/gl"
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
	if len(flag.Args()) != 2 {
		fmt.Println("Usage: poussin [-cpuprofile FILE] [-memprofile FILE] BOOTROM ROM")
		os.Exit(1)
	}

	nextFrame := make(chan *image.RGBA, 1)
	gb, err := emu.NewGameboy(nextFrame)
	if err != nil {
		return err
	}
	defer gb.Close()

	bootRom, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		return err
	}
	if err := gb.LoadBootROM(bootRom); err != nil {
		return err
	}

	rom, err := ioutil.ReadFile(flag.Args()[1])
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

	quit := make(chan bool)

	go gb.Run(quit)
	go r.Run(nextFrame, quit)

	select {
	case <-quit:
	}
	close(quit)

	return nil
}
