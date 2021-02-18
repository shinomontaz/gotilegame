package main

import (
	"menutest/screens/menu"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Test menu",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	face, err := loadTTF("capture.ttf", 20)
	if err != nil {
		panic(err)
	}
	basicAtlas := text.NewAtlas(face, text.ASCII)
	//	basicTxt := text.New(pixel.V(100, 500), basicAtlas)
	//	basicTxt.Color = colornames.Red

	fscreen := menu.NewMain(basicAtlas)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		fscreen.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
