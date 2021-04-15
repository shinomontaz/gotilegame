package main

import (
	"fmt"
	"math"
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
	//	hero.setPos(win.Bounds().Center())

	world := NewWorld()

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

	for !win.Closed() {

		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Whitesmoke)

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

		world.Draw(win)
		hero.Draw(win, win.Bounds().Center().Sub(hero.getPos()))
		win.Update()

		frames++
		select {
		case <-frametime:
			hero.Tick()
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
			fmt.Println(cam, hero.getPos())
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
