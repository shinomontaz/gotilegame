package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	b := pixel.R(0, 0, 500, 350)
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, World!",
		Bounds: b,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture("test_sprite.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	canvas := pixelgl.NewCanvas(b)

	canvas.Clear(colornames.Black)
	sprite.Draw(canvas, pixel.IM.Moved(win.Bounds().Center()))

	win.Clear(colornames.Skyblue)
	canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	//	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	for !win.Closed() {
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
