package console

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/faiface/gui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Console struct {
	cv     image.Rectangle
	mess   <-chan string
	images []image.Image
	env    gui.Env
}

func New(mess <-chan string, cv image.Rectangle, env gui.Env) *Console {
	return &Console{cv: cv, mess: mess, images: make([]image.Image, 0, 3), env: env}
}

func (c *Console) Draw(r image.Rectangle) func(draw.Image) image.Rectangle {
	return func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, r, &image.Uniform{color.Black}, image.ZP, draw.Src)
		consoleImage := image.NewRGBA(image.Rect(0, 0, width, lineHeight*(len(images)+2)))
		for i := range images {
			r := image.Rect(
				0, lineHeight*i,
				width, lineHeight*(i+1),
			)

			DrawLeftCentered(consoleImage, r, images[i], draw.Over)
		}

		f, err := os.Create("consoleImage.png")
		if err != nil {
			// Handle error
		}
		defer f.Close()
		err = png.Encode(f, consoleImage)

		bounds := consoleImage.Bounds()
		leftCenter := image.Pt(bounds.Min.X, (bounds.Min.Y+bounds.Max.Y)/2)
		target := image.Pt(r.Min.X, (r.Min.Y+r.Max.Y)/2)
		delta := target.Sub(leftCenter)
		draw.Draw(drw, bounds.Add(delta).Intersect(r), consoleImage, bounds.Min, draw.Src)
		return r
	}
}

func (c *Console) Run() {
	lineHeight := 20
	width := 100

	// first we draw a white rectangle
	//	c.env.Draw() <- redraw(c.cv, c.images)
	c.env.Draw() <- c.Draw()

	for mess := range c.mess {
		fmt.Println(mess)
		c.images = append(c.images, MakeTextImage(mess))
		if len(c.images) > 10 {
			c.images = c.images[1:]
		}
		//		c.env.Draw() <- redraw(c.cv, c.images)
		c.env.Draw() <- c.Draw()
	}

	close(c.env.Draw())
}

func MakeTextImage(text string) image.Image {
	drawer := &font.Drawer{
		Src:  &image.Uniform{color.RGBA{G: 0xFF, A: 0xFF}},
		Face: basicfont.Face7x13,
		Dot:  fixed.P(0, 0),
	}
	b26_6, _ := drawer.BoundString(text)
	bounds := image.Rect(
		b26_6.Min.X.Floor(),
		b26_6.Min.Y.Floor(),
		b26_6.Max.X.Ceil(),
		b26_6.Max.Y.Ceil(),
	)
	drawer.Dst = image.NewRGBA(bounds)
	drawer.DrawString(text)

	// f, err := os.Create(text + "outimage.png")
	// if err != nil {
	// 	// Handle error
	// }
	// defer f.Close()
	// err = png.Encode(f, drawer.Dst)

	return drawer.Dst
}

func DrawLeftCentered(dst draw.Image, r image.Rectangle, src image.Image, op draw.Op) {
	if src == nil {
		return
	}
	bounds := src.Bounds()
	leftCenter := image.Pt(bounds.Min.X, (bounds.Min.Y+bounds.Max.Y)/2)
	target := image.Pt(r.Min.X, (r.Min.Y+r.Max.Y)/2)
	delta := target.Sub(leftCenter)
	draw.Draw(dst, bounds.Add(delta).Intersect(r), src, bounds.Min, op)

	f, err := os.Create("outimage.png")
	if err != nil {
		// Handle error
	}
	defer f.Close()
	png.Encode(f, dst)
}
