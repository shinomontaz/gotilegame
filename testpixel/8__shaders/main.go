package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/go-gl/mathgl/mgl32"
)

var gopherimg *pixel.Sprite

var fragSource string
var uTime float32

//var uObjects []float32 // = 250 rectangles
var uObjects []mgl32.Vec4 // = 250 rectangles

var uNumObjects int32
var uLight mgl32.Vec2
var b pixel.Rect
var canvas *pixelgl.Canvas

func gameloop(win *pixelgl.Window) {
	uLight = [2]float32{-50.0, 0.0}

	canvas = pixelgl.NewCanvas(b)
	canvas.SetUniform("uTime", &uTime)
	canvas.SetUniform("uLight", &uLight)
	canvas.SetUniform("uObjects", &uObjects)
	canvas.SetUniform("uNumObjects", &uNumObjects)

	rect := pixel.R(50, 50, 100, 80)
	rect2 := pixel.R(-100, -80, -50, -20)
	// uObjects = []float32{float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y)}
	// uObjects = append(uObjects, float32(rect2.Min.X), float32(rect2.Min.Y), float32(rect2.Max.X), float32(rect2.Max.Y))

	uObjects = []mgl32.Vec4{mgl32.Vec4{float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y)}}
	uObjects = append(uObjects, mgl32.Vec4{float32(rect2.Min.X), float32(rect2.Min.Y), float32(rect2.Max.X), float32(rect2.Max.Y)})

	uNumObjects = int32(len(uObjects))

	canvas.SetFragmentShader(fragSource)

	start := time.Now()
	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			fmt.Println(win.MousePosition())
			pos := win.MousePosition()
			uLight = [2]float32{float32(pos.X - b.Center().X), float32(pos.Y - b.Center().Y)}
		}

		mainStep(win)
		uTime = float32(time.Since(start).Seconds())
		win.Update()
	}
}

func mainStep(t pixel.Target) {
	canvas.Clear(pixel.RGB(0, 0, 0))

	gopherimg.Draw(canvas, pixel.IM.Moved(b.Center()))

	canvas.Draw(t, pixel.IM.Moved(canvas.Bounds().Center()))
}

func run() {
	b = pixel.R(0, 0, 500, 500)
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, shaders!",
		Bounds: b,
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	f, err := os.Open("./gamebackground.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	pd := pixel.PictureDataFromImage(img)
	gopherimg = pixel.NewSprite(pd, pd.Bounds())

	//	fragSource, err = LoadFileToString("4lfyDM.glsl")
	fragSource, err = LoadFileToString("light.glsl")

	//	fragSource, err = LoadFileToString("slyGDR_2.glsl")

	if err != nil {
		panic(err)
	}

	gameloop(win)
}

func main() {
	pixelgl.Run(run)
}

func LoadFileToString(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
