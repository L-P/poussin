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

	if len(flag.Args()) != 2 {
		fmt.Println("Usage: poussin [-cpuprofile FILE] [-memprofile FILE] BOOTROM ROM")
		os.Exit(1)
	}

	nextFrame := make(chan *image.RGBA, 1)
	gb := emu.NewGameboy(nextFrame)
	defer gb.Close()

	bootRom, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	if err := gb.LoadBootROM(bootRom); err != nil {
		panic(err)
	}

	rom, err := ioutil.ReadFile(flag.Args()[1])
	if err != nil {
		panic(err)
	}
	if err := gb.LoadROM(rom); err != nil {
		panic(err)
	}

	quit := make(chan bool)
	go gb.Run(quit)
	go gl.Run(nextFrame, quit)

	select {
	case <-quit:
	}
	close(quit)

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
