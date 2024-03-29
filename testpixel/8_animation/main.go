package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const WIDTH = 640.0
const HEIGTH = 500.0

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Animation",
		Bounds: pixel.R(0, 0, WIDTH, HEIGTH),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	hero := NewHero()
	world := NewWorld(WIDTH, HEIGTH, hero.getPos())

	var (
		camPos = pixel.ZV
		// camSpeed     = 500.0
		// camZoom      = 1.0
		// camZoomSpeed = 1.2
		frames    = 0
		second    = time.Tick(time.Second)
		frametime = time.Tick(120 * time.Millisecond)
	)

	rand.Seed(time.Now().UTC().UnixNano())
	last := time.Now()

	rgba := color.RGBA{205, 231, 244, 1}

	for !win.Closed() {

		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(rgba)

		camPos = pixel.Lerp(camPos, hero.getPos(), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos)

		win.SetMatrix(cam)

		if win.Pressed(pixelgl.KeyRightControl) {
			hero.Notify(FIRE)
		} else if win.Pressed(pixelgl.KeyLeft) {
			hero.Notify(LEFT)
		} else if win.Pressed(pixelgl.KeyRight) {
			hero.Notify(RIGHT)
		} else if win.Pressed(pixelgl.KeyLeftControl) {
			hero.Notify(FIRE)
		} else if win.Pressed(pixelgl.KeyEnter) {
			hero.Notify(ENTER)
		}

		world.Draw(win, hero.getPos(), camPos)
		//		hero.Draw(win, win.Bounds().Center().Sub(hero.getPos()))
		hero.Draw(win, win.Bounds().Center().Sub(hero.getPos()).Sub(pixel.V(0, 140.0)))

		win.Update()

		frames++
		select {
		case <-frametime:
			hero.Tick()
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			fmt.Println("hero.getPos()", hero.getPos())

			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
