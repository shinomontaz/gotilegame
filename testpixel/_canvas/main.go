package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var uLight mgl32.Vec2

func run() {
	b := pixel.R(0, 0, 500, 500)
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
	pic2, err := loadPicture("test_sprite.png")
	if err != nil {
		panic(err)
	}

	fragSource, err := LoadFileToString("light.glsl")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite2 := pixel.NewSprite(pic2, pic2.Bounds())

	var (
		camPos   = pixel.ZV
		camSpeed = 500.0
	)
	uLight = [2]float32{float32(-b.W()/2 + 200.0), float32(b.H()/2 - 250.0)}

	canvas := pixelgl.NewCanvas(b)
	canvas2 := pixelgl.NewCanvas(b)

	canvas.SetUniform("uLight", &uLight)
	canvas.SetFragmentShader(fragSource)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		cam := pixel.IM.Moved(pixel.ZV.Sub(camPos))
		//		win.SetMatrix(cam)

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}

		canvas.Clear(colornames.Black)
		canvas.SetBounds(b)

		canvas2.Clear(colornames.Black)
		canvas2.SetBounds(b)

		canvas2.SetMatrix(cam)

		sprite.Draw(canvas2, pixel.IM.Moved(canvas2.Bounds().Center()))
		sprite2.Draw(canvas2, pixel.IM.Moved(canvas2.Bounds().Center()))

		canvas2.Draw(canvas, pixel.IM.Moved(canvas.Bounds().Center()))

		win.Clear(colornames.Skyblue)
		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

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

func LoadFileToString(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
