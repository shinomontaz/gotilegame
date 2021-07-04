package main

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
)

type World struct {
	width    float64
	height   float64
	bg       pixel.Picture
	viewport pixel.Rect
	vector1  pixel.Vec
	vector2  pixel.Vec
	part1    *pixel.Sprite
	part2    *pixel.Sprite
	pos      pixel.Vec
	steps    int
}

func NewWorld(width, height float64, start pixel.Vec) *World {
	w := World{
		width:    width,
		height:   height,
		viewport: pixel.R(0, 0, width, height),
		vector1:  pixel.V(width/2, height/2+140.0),
		vector2:  pixel.V(3*width/2, height/2+140.0),
		pos:      start,
		steps:    0,
	}
	bg, err := loadPicture("assets/DarkVaniaAssets/background/cloudySky_640x360px.png")
	if err != nil {
		panic(err)
	}

	w.bg = bg

	w.part1 = pixel.NewSprite(w.bg, w.viewport)
	//	w.part2 = pixel.NewSprite(w.bg, pixel.R(width, 0, 2*width, height))
	w.part2 = pixel.NewSprite(w.bg, w.viewport)

	return &w
}

func (w *World) Draw(t pixel.Target, pos pixel.Vec, cam pixel.Vec) {
	x, _ := pos.XY()
	steps := int(math.Abs(x / w.width))

	if steps != w.steps {
		fmt.Println("w.steps, steps", w.steps, steps)
		w.vector1, w.vector2 = w.vector2, w.vector1
		w.steps = steps
	}

	x = math.Mod(x, w.width)

	w.pos = pixel.V(-x, 0)

	w.part1.Draw(t, pixel.IM.Moved(w.vector1.Sub(w.pos).Sub(cam)))
	w.part2.Draw(t, pixel.IM.Moved(w.vector2.Sub(w.pos).Sub(cam)))
}
