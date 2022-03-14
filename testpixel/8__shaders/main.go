package main

import (
	"image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var gopherimg *pixel.Sprite

var fsGrayscale = `
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

	// Sum our 3 color channels
	float sum  = texture(uTexture, t).r;
	      sum += texture(uTexture, t).g;
	      sum += texture(uTexture, t).b;

	// Divide by 3, and set the output to the result
	vec4 color = vec4( sum/3, sum/3, sum/3, 1.0);
	fragColor = color;
}
`

var fsWater = `
#version 330 core

in vec2 vTexCoords;
out vec4 fragColor;

uniform sampler2D uTexture;
uniform vec4 uTexBounds;

// custom uniforms
uniform float uSpeed;
uniform float uTime;

void main() {
    vec2 t = vTexCoords / uTexBounds.zw;
	vec3 influence = texture(uTexture, t).rgb;

    if (influence.r + influence.g + influence.b > 0.3) {
		t.y += cos(t.x * 40.0 + (uTime * uSpeed))*0.005;
		t.x += cos(t.y * 40.0 + (uTime * uSpeed))*0.01;
	}

    vec3 col = texture(uTexture, t).rgb;
	fragColor = vec4(col * vec3(0.6, 0.6, 1.2),1.0);
}
`

var uTime, uSpeed float32
var b pixel.Rect
var canvas, canvas1 *pixelgl.Canvas

func gameloop(win *pixelgl.Window) {
	canvas1 = pixelgl.NewCanvas(b)
	canvas = pixelgl.NewCanvas(b)
	canvas.SetUniform("uTime", &uTime)
	canvas.SetUniform("uSpeed", &uSpeed)
	uSpeed = 5.0

	canvas.SetFragmentShader(fsWater)

	start := time.Now()
	for !win.Closed() {
		win.Clear(pixel.RGB(0, 0, 0))

		mainStep(win)
		uTime = float32(time.Since(start).Seconds())
		if win.Pressed(pixelgl.KeyRight) {
			uSpeed += 0.1
		}
		if win.Pressed(pixelgl.KeyLeft) {
			uSpeed -= 0.1
		}

		win.Update()
	}
}

func mainStep(t pixel.Target) {
	canvas1.Clear(pixel.RGB(0, 0, 0))
	canvas.Clear(pixel.RGB(0, 0, 0))

	gopherimg.Draw(canvas1, pixel.IM.Moved(pixel.Vec{100, 100}))
	gopherimg.Draw(canvas1, pixel.IM.Scaled(pixel.ZV, 2).Moved(canvas1.Bounds().Center())) // pixel.IM.Scaled(pixel.ZV, 2).Moved(canvas.Bounds().Center())
	canvas1.Draw(canvas, pixel.IM.Moved(canvas.Bounds().Center()))

	canvas.Draw(t, pixel.IM.Moved(canvas.Bounds().Center()))
}

func run() {
	b = pixel.R(0, 0, 325, 348)
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, shaders!",
		Bounds: b,
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	f, err := os.Open("./test_sprite.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	pd := pixel.PictureDataFromImage(img)
	gopherimg = pixel.NewSprite(pd, pd.Bounds())
	gameloop(win)
}

func main() {
	pixelgl.Run(run)
}
