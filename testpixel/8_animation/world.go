package main

import (
	"math"

	"github.com/faiface/pixel"
)

type World struct {
	width  float64
	height float64
	//	bg     *pixel.Sprite
	bg       pixel.Picture
	viewport pixel.Rect
	start    pixel.Vec
}

func NewWorld(width, height float64, start pixel.Vec) *World {
	w := World{
		width:    width,
		height:   height,
		viewport: pixel.R(0, 0, width, height),
		start:    start,
	}
	bg, err := loadPicture("back.png")
	if err != nil {
		panic(err)
	}

	w.bg = bg //pixel.NewSprite(bg, pixel.R(bg.Bounds().Min.X, bg.Bounds().Min.Y, bg.Bounds().Max.X, bg.Bounds().Max.Y))

	return &w
}

func (w World) Draw(t pixel.Target, pos pixel.Vec) {
	//	r := w.bg.Frame()

	xProj := pos.Project(pixel.V(1.0, 0))
	x, _ := xProj.XY()
	//	jumps := int(x / (w.width / 2))
	x1 := math.Mod(x, w.width/2)

	var part1, part2 *pixel.Sprite

	part1 = pixel.NewSprite(w.bg, w.viewport)
	part2 = pixel.NewSprite(w.bg, w.viewport)

	w.viewport = w.viewport.Moved(pos)
	part1.Draw(t, pixel.IM.Moved(pixel.V(x1, w.height/2)))
	if x1 > 0 {
		part2.Draw(t, pixel.IM.Moved(pixel.V(x1+w.width, w.height/2)))
	} else {
		part2.Draw(t, pixel.IM.Moved(pixel.V(x1-w.width, w.height/2)))
	}
}
