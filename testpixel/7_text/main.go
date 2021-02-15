package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, World!",
		Bounds: pixel.R(0, 0, 1024, 768),
		//		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	face, err := loadTTF("intuitive.ttf", 20)
	if err != nil {
		panic(err)
	}
	basicAtlas := text.NewAtlas(face, text.ASCII)
	basicTxt := text.New(pixel.V(100, 500), basicAtlas)
	basicTxt.Color = colornames.Red

	// lines := []string{
	// 	"This is a very very very long line",
	// 	"Short line",
	// 	"--=!@#$^&*()_+=--",
	// }

	// basicTxt.Color = colornames.Red

	// for _, line := range lines {
	// 	basicTxt.Dot.X -= basicTxt.BoundsOf(line).W() / 2
	// 	fmt.Fprintln(basicTxt, line)
	// }

	timer := time.Tick(time.Second / 120)
	second := time.Tick(time.Second * 2)
	fps := 0
	for !win.Closed() {
		basicTxt.WriteString(win.Typed())
		if win.JustPressed(pixelgl.KeyEnter) || win.Repeated(pixelgl.KeyEnter) {
			basicTxt.WriteRune('\n')
		}
		win.Clear(colornames.Black)

		basicTxt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(basicTxt.Bounds().Center())))
		win.Update()
		<-timer

		fps++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("FPS: %d", fps/2))
			fps = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
