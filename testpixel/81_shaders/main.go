package main

import (
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
var uObject mgl32.Mat2
var uLight mgl32.Vec2
var b pixel.Rect
var canvas *pixelgl.Canvas

func gameloop(win *pixelgl.Window) {
	canvas = pixelgl.NewCanvas(b)
	canvas.SetUniform("uTime", &uTime)
	canvas.SetUniform("uLight", &uLight)
	canvas.SetUniform("uObject", &uObject)

	rect := pixel.R(50, 50, 100, 80)
	uObject = [4]float32{float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y)}
	uLight = [2]float32{-50.0, 0}

	//	fmt.Println(uObject)

	canvas.SetFragmentShader(fragSource)

	start := time.Now()
	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

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

	fragSource, err = LoadFileToString("test.glsl")

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
