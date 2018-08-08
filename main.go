package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"

	"home.leo-peltier.fr/poussin/emu"
	"home.leo-peltier.fr/poussin/renderer/gl"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: poussin BOOTROM ROM")
		os.Exit(1)
	}

	nextFrame := make(chan *image.RGBA, 1)
	gb := emu.NewGameboy(nextFrame)
	defer gb.Close()

	bootRom, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	if err := gb.LoadBootROM(bootRom); err != nil {
		panic(err)
	}

	rom, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}
	if err := gb.LoadROM(rom); err != nil {
		panic(err)
	}

	quit := make(chan int)
	go gb.Run(quit)
	go gl.Run(nextFrame, quit)

	select {
	case <-quit:
	}
}
