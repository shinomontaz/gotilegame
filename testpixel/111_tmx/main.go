package main

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var clearColor = colornames.Skyblue

var platforms []pixel.Rect
var wm *Worldmap

//var initialMove pixel.Vec

func gameloop(win *pixelgl.Window) {
	mainRect := wm.Data()
	var (
		camPos   = mainRect.Center()
		camSpeed = 1000.0
	)

	fmt.Println(camPos, win.Bounds(), mainRect)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Moved(win.Bounds().Center().Sub(camPos))

		win.SetMatrix(cam)

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}

		win.Clear(clearColor)

		// Draw tiles
		wm.Draw(win)

		win.Update()
	}
}

func run() {
	wm = NewWorldmap("assets", "map.tmx")
	platforms = wm.Geometry()
	mainRect := wm.Data()

	// Create the window with OpenGL
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel tmx with lafriks",
		Bounds: pixel.R(0.0, 0.0, mainRect.W(), mainRect.H()),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	panicIfErr(err)
	gameloop(win)
}

func main() {
	pixelgl.Run(run)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
