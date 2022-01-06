package main

import (
	"fmt"
	_ "image/png"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

type screenLogger struct {
	bt       *text.Text
	ba       *text.Atlas
	canvas   *imdraw.IMDraw
	onBt     bool
	onCanvas bool
}

func (sc *screenLogger) draw(win *pixelgl.Window, cam pixel.Matrix) {
	if sc.onBt {
		sc.drawText(win, cam)
	}

	if sc.onCanvas {
		sc.drawCanvas(win)
	}
}

func (sc *screenLogger) drawText(win *pixelgl.Window, cam pixel.Matrix) {
	var position pixel.Vec = pixel.V(10, 900)

	sc.bt = text.New(position, sc.ba)
	sc.bt.Color = colornames.Whitesmoke
	fmt.Fprintf(sc.bt, "text postion: %v \r\n", position)
	sc.bt.Color = colornames.Whitesmoke

	fmt.Fprintf(sc.bt, "windows W:%v, H:%v \r\n", win.Bounds().W(), win.Bounds().H())

	sc.bt.Draw(win, pixel.IM.Scaled(sc.bt.Orig, 1))

	for _, platform := range platforms {
		sc.bt = text.New(platform.Min, sc.ba)
		sc.bt.Color = colornames.Pink
		fmt.Fprintf(sc.bt, platform.String())
		sc.bt.Draw(win, pixel.IM.Scaled(sc.bt.Orig, 1))
	}
}

func (sc *screenLogger) drawCanvas(win *pixelgl.Window) {
	sc.canvas.Draw(win)
}

func (sc *screenLogger) initCanvas() {
	for _, p := range platforms {
		sc.canvas.Color = colornames.Orange
		sc.canvas.EndShape = imdraw.RoundEndShape
		//imd.Push(pixel.V(0, 64), pixel.V(320, 64))
		sc.canvas.Push(p.Min, p.Max)
		sc.canvas.Rectangle(1)
		sc.canvas.Line(2)
	}
}
