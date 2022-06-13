package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Back struct {
	width    float64
	height   float64
	p        pixel.Picture
	viewport pixel.Rect
	vector1  pixel.Vec
	vector2  pixel.Vec

	c *pixelgl.Canvas

	part1 *pixel.Sprite
	part2 *pixel.Sprite

	pos   pixel.Vec
	steps int
}

func NewBack(start pixel.Vec, viewport pixel.Rect, path string) *Back {
	width := viewport.W()
	height := viewport.H()

	//	canvas := pixelgl.NewCanvas(pixel.R(0, 0, width, height))
	canvas := pixelgl.NewCanvas(viewport)

	//	viewport = viewport.Moved(pixel.Vec{0, 100})
	x, y := viewport.Min.X, viewport.Min.Y
	b := Back{
		width:    width,
		height:   height,
		viewport: viewport,
		c:        canvas,
		vector1:  pixel.V(x+width/2+1, y+height/2),
		vector2:  pixel.V(x+3*width/2, y+height/2),

		pos:   start,
		steps: 0,
	}

	bg, err := LoadPicture(path)
	if err != nil {
		panic(err)
	}

	b.p = bg

	b.part1 = pixel.NewSprite(b.p, pixel.R(0, 0, bg.Bounds().W()/2, bg.Bounds().H()))
	b.part2 = pixel.NewSprite(b.p, pixel.R(bg.Bounds().W()/2, 0, bg.Bounds().W(), bg.Bounds().H()))

	return &b
}

func (b *Back) Draw(t pixel.Target, pos pixel.Vec, cam pixel.Vec) {
	b.c.Clear(pixel.RGB(0, 0, 0))
	x, _ := pos.XY()
	cam.Y = 0
	steps := int(math.Abs(x / b.width))

	if steps != b.steps {
		b.vector1, b.vector2 = b.vector2, b.vector1

		b.steps = steps
	}

	x = math.Mod(x, b.width)

	b.pos = pixel.V(x, 0)

	mtx1 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		b.width/b.part1.Frame().W(),
		b.height/b.part1.Frame().H(),
	)).Moved(b.vector1.Sub(b.pos).Sub(cam))

	mtx2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		b.width/b.part2.Frame().W(),
		b.height/b.part2.Frame().H(),
	)).Moved(b.vector2.Sub(b.pos).Sub(cam))

	//	b.part1.Draw(b.c, mtx1)
	//	b.part2.Draw(b.c, mtx2)

	b.part1.Draw(t, mtx1)
	b.part2.Draw(t, mtx2)

	//	fmt.Printf("back draw vec1 %v, vec2 %v, center %v, pos %v, cam %v\n", b.vector1, b.vector2, b.c.Bounds().Center(), pos, cam)

	//	b.c.Draw(t, pixel.IM.Moved(b.c.Bounds().Center()))
}
