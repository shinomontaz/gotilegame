package main

import "github.com/faiface/pixel"

type World struct {
	width  int
	height int
	bg     *pixel.Sprite
}

func NewWorld() *World {
	w := World{}
	bg, err := loadPicture("back.png")
	if err != nil {
		panic(err)
	}

	w.bg = pixel.NewSprite(bg, pixel.R(bg.Bounds().Min.X, bg.Bounds().Min.Y, bg.Bounds().Max.X, bg.Bounds().Max.Y))

	return &w
}

func (w World) Draw(t pixel.Target) {
	r := w.bg.Frame()
	w.bg.Draw(t, pixel.IM.Moved(pixel.V(r.Max.X/2, r.Max.Y/2)))
}
