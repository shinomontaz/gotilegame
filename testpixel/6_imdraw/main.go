package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, World!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	imd := imdraw.New(nil)

	imd.Color = colornames.Darkred //pixel.RGB(1, 0, 0)
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(pixel.V(200, 100), pixel.V(0, 200))
	imd.Color = colornames.Darkgreen //pixel.RGB(1, 0, 0)
	imd.EndShape = imdraw.SharpEndShape
	imd.Push(pixel.V(800, 100), pixel.V(500, 700))
	imd.Color = colornames.Steelblue //pixel.RGB(1, 0, 0)
	//	imd.Push(pixel.V(500, 700))
	imd.Line(30)

	imd.Color = colornames.Limegreen
	imd.Push(pixel.V(500, 500))
	imd.Ellipse(pixel.V(120, 80), 0)

	imd.Color = colornames.Red
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(pixel.V(500, 350))
	imd.CircleArc(150, -math.Pi/3, 0, 20)

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
