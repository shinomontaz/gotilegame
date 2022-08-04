package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	b := pixel.R(100, 100, 600, 450)
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, canvas!",
		Bounds: b,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture("gamebackground.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	canvas := pixelgl.NewCanvas(b)
	canvas2 := pixelgl.NewCanvas(b)
	//	last := time.Now()
	for !win.Closed() {
		//		dt := time.Since(last).Seconds()

		if win.Pressed(pixelgl.KeyLeft) {
			b = b.Moved(pixel.Vec{-1, 0})
		}
		if win.Pressed(pixelgl.KeyRight) {
			b = b.Moved(pixel.Vec{1, 0})
		}

		fmt.Println(b)

		canvas.Clear(colornames.Black)
		canvas2.Clear(colornames.Black)

		canvas.SetBounds(b)
		canvas2.SetBounds(b)

		sprite.Draw(canvas, pixel.IM.Moved(canvas.Bounds().Center()))
		canvas.Draw(canvas2, pixel.IM.Moved(canvas2.Bounds().Center()))

		win.Clear(colornames.Skyblue)
		canvas2.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		win.Update()
		//		last = time.Now()
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
