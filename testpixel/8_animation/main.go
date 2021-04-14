package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Animation",
		Bounds: pixel.R(0, 0, 500, 500),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//	win.SetSmooth(true)

	hero := NewHero()

	var (
		frames    = 0
		second    = time.Tick(time.Second)
		frametime = time.Tick(120 * time.Millisecond)
	)

	rand.Seed(time.Now().UTC().UnixNano())

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)

		if win.Pressed(pixelgl.KeyLeft) {
			hero.Notify(LEFT)
		}
		if win.Pressed(pixelgl.KeyRight) {
			hero.Notify(RIGHT)
		}
		if win.Pressed(pixelgl.KeyLeftControl) {
			hero.Notify(FIRE)
		}
		if win.Pressed(pixelgl.KeyRightControl) {
			hero.Notify(FIRE)
		}
		if win.Pressed(pixelgl.KeyEnter) {
			hero.Notify(ENTER)
		}

		hero.Draw(win, win.Bounds().Center())
		win.Update()

		frames++
		select {
		case <-frametime:
			hero.Tick()
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}