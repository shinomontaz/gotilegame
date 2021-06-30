package main

import (
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	windowWidth    = 500
	windowHeight   = 500
	linesPerSecond = 60
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, World!",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture("back.png")
	if err != nil {
		panic(err)
	}

	background1 := pixel.NewSprite(pic, pixel.R(0, 0, windowWidth, windowHeight))
	background2 := pixel.NewSprite(pic, pixel.R(windowWidth, 0, windowWidth*2, windowHeight))

	// In the beginning, vector1 will put background1 filling the whole window, while vector2 will
	// put background2 just at the right side of the window, out of view
	vector1 := pixel.V(windowWidth/2, windowHeight/2)
	vector2 := pixel.V(windowWidth+(windowWidth/2), windowHeight/2)

	i := float64(0)
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		// When one of the backgrounds has completely scrolled, we swap displacement vectors,
		// so the backgrounds will swap positions too regarding the previous iteration,
		// thus making the background endless.
		if i <= -windowWidth {
			i = 0
			vector1, vector2 = vector2, vector1
		}

		if win.Pressed(pixelgl.KeyEnter) {
			// This delta vector will move the backgrounds to the left
			d := pixel.V(-i, 0)
			background1.Draw(win, pixel.IM.Moved(vector1.Sub(d)))
			background2.Draw(win, pixel.IM.Moved(vector2.Sub(d)))
			i -= linesPerSecond * dt
		}

		win.Update()
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func main() {
	pixelgl.Run(run)
}
