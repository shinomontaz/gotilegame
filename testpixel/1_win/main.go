package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, World!",
		Bounds: pixel.R(0, 0, 500, 350),
		VSync:  true,
	}

	// fmt.Println("dfdf: ", int(pixelgl.KeySpace))
	// fmt.Println("dfdf2: ", pixelgl.Button(341).String())

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	for !win.Closed() {
		if win.Pressed(pixelgl.KeySpace) {
			fmt.Println("dfdf: ", int(pixelgl.KeySpace))
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
