package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"home.leo-peltier.fr/poussin/emu"
)

func main() {
	/*
		nextFrame := make(chan int, 1)
		quit := make(chan int, 1)
		fb := image.NewRGBA(image.Rect(0, 0, 160, 144))

		for x := 0; x < fb.Rect.Size().X; x++ {
			for y := 0; y < fb.Rect.Size().Y; y++ {
				fb.SetRGBA(x, y, color.RGBA{255, 0, 255, 255})
			}
		}

		go gl.Run(fb, nextFrame, quit)

		select {
		case <-quit:
		}
	// */

	if len(os.Args) != 2 {
		fmt.Println("Usage: poussin BOOTROM")
		os.Exit(1)
	}

	gb := emu.NewGameboy()
	bootRom, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := gb.Boot(bootRom); err != nil {
		panic(err)
	}

	if err := gb.Run(); err != nil {
		panic(err)
	}
}
