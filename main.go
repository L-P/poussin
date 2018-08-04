package main

import (
	"image"
	"image/color"

	"home.leo-peltier.fr/poussin/renderer/gl"
)

func main() {
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
}
