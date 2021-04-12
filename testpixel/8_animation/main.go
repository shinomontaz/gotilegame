package main

import (
	"fmt"
	"image/png"
	"math/rand"
	"os"
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

	spritesheet, err := loadPicture("hero_spritesheet.png")
	if err != nil {
		panic(err)
	}

	var (
		walkFrames  []pixel.Rect
		shootFrame  pixel.Rect
		standFrame  pixel.Rect
		deathFrames []pixel.Rect
		deadFrame   pixel.Rect
		frames      = 0
		second      = time.Tick(time.Second)
		frametime   = time.Tick(120 * time.Millisecond)
	)

	xstep := 80.0
	ystep := 94.0
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Min.X+6*xstep; x += xstep {
		walkFrames = append(walkFrames, pixel.R(x, spritesheet.Bounds().Max.Y-ystep, x+xstep, spritesheet.Bounds().Max.Y-2*ystep))
	}
	shootFrame = pixel.R(spritesheet.Bounds().Min.X, spritesheet.Bounds().Max.Y-2*ystep, spritesheet.Bounds().Min.X+xstep, spritesheet.Bounds().Max.Y-3*ystep)
	standFrame = pixel.R(spritesheet.Bounds().Min.X, spritesheet.Bounds().Max.Y, spritesheet.Bounds().Min.X+xstep, spritesheet.Bounds().Max.Y-ystep)
	for x := spritesheet.Bounds().Min.X + 3*xstep; x < spritesheet.Bounds().Min.X+6*xstep; x += xstep {
		deathFrames = append(deathFrames, pixel.R(x, spritesheet.Bounds().Max.Y-3*ystep, x+xstep, spritesheet.Bounds().Max.Y-4*ystep))
	}
	deadFrame = pixel.R(spritesheet.Bounds().Min.X+6*xstep, spritesheet.Bounds().Max.Y-3*ystep, spritesheet.Bounds().Min.X+6*xstep+ystep, spritesheet.Bounds().Max.Y-4*ystep)

	rand.Seed(time.Now().UTC().UnixNano())
	i := 0
	j := 0
	lastState := 0
	dir := 1.0 // direction to the right
	var frame pixel.Rect

	deadState := 5
	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		state := 0

		if lastState == deadState {
			state = 5
		}
		if lastState == 4 {
			state = 4
		}

		if win.Pressed(pixelgl.KeyLeft) && state != deadState && state != 4 {
			state = 1 // walk left
			dir = -1  // set direactino as left
		}
		if win.Pressed(pixelgl.KeyRight) && state != deadState && state != 4 {
			state = 2 // walk right
			dir = 1
		}

		if win.Pressed(pixelgl.KeyLeftControl) && state != deadState && state != 4 {
			state = 3 // shoot
		}
		if win.Pressed(pixelgl.KeyRightControl) && state != deadState && state != 4 {
			state = 3 // shoot
		}
		if win.Pressed(pixelgl.KeyEnter) && state != deadState {
			state = 4 // death
		}

		if lastState != state {
			i = 0
			j = 0
		}

		lastState = state

		if lastState == 4 && j == len(deathFrames)-1 {
			lastState = deadState
		}

		switch state {
		case 0:
			frame = standFrame
		case 1:
			frame = walkFrames[i]
		case 2:
			frame = walkFrames[i]
		case 3:
			frame = shootFrame
		case 4:
			frame = deathFrames[j]
		case 5:
			frame = deadFrame
		}

		soldier := pixel.NewSprite(spritesheet, frame)
		soldier.Draw(win, pixel.IM.ScaledXY(pixel.ZV, pixel.V(dir, 1)).Scaled(pixel.ZV, 1.5).Moved(win.Bounds().Center()))

		win.Update()

		frames++
		select {
		case <-frametime:
			i++
			i = i % len(walkFrames)
			j++
			j = j % len(deathFrames)

		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

func main() {
	pixelgl.Run(run)
}
