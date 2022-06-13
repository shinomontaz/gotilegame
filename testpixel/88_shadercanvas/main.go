package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/go-gl/mathgl/mgl32"
)

var back *Back

var fragSource string
var uTime float32

var lastPos pixel.Vec

var uObjects []mgl32.Vec4

var uNumObjects int32
var uLight mgl32.Vec2
var b pixel.Rect
var canvas *pixelgl.Canvas

func gameloop(win *pixelgl.Window) {
	//	canvas = pixelgl.NewCanvas(b)
	canvas = pixelgl.NewCanvas(pixel.R(0, 0, b.W(), b.H()))

	canvas.SetUniform("uTime", &uTime)
	canvas.SetUniform("uLight", &uLight)
	//	canvas.SetUniform("uObjects", &uObjects)
	//	canvas.SetUniform("uNumObjects", &uNumObjects)

	rect := pixel.R(50, 50, 100, 80)
	rect2 := pixel.R(-100, -80, -50, -20)

	uObjects = []mgl32.Vec4{mgl32.Vec4{float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y)}}
	uObjects = append(uObjects, mgl32.Vec4{float32(rect2.Min.X), float32(rect2.Min.Y), float32(rect2.Max.X), float32(rect2.Max.Y)})

	uNumObjects = int32(len(uObjects))
	uLight = [2]float32{-50.0, 0}

	canvas.SetFragmentShader(fragSource)

	start := time.Now()
	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		mainStep(win)
		uTime = float32(time.Since(start).Seconds())
		win.Update()
	}
}

func mainStep(t *pixelgl.Window) { //t pixel.Target) {
	canvas.Clear(pixel.RGB(0, 0, 0))

	//	cam := pixel.ZV

	cam := lastPos.Sub(lastPos).Sub(pixel.V(0, 150))

	back.Draw(canvas, lastPos, cam)

	cent := pixel.Vec{canvas.Bounds().W() / 2, canvas.Bounds().H() / 2}
	fmt.Println("canvas.Bounds().Center()): %v", canvas.Bounds().Center())
	fmt.Println("cent: %v", cent)

	canvas.Draw(t, pixel.IM.Moved(t.Bounds().Center()))
	//	canvas.Draw(t, pixel.IM.Moved(cent))
}

func run() {
	b = pixel.R(100, 100, 500, 500)
	//	b = pixel.R(500, 500, 1000, 1000)

	lastPos = b.Center()
	cfg := pixelgl.WindowConfig{
		Title:  "Canvas, shaders and friends",
		Bounds: b,
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	fragSource, err = LoadFileToString("light.glsl")

	back = NewBack(lastPos, b, "gamebackground.png")

	if err != nil {
		panic(err)
	}

	gameloop(win)
}

func main() {
	pixelgl.Run(run)
}
