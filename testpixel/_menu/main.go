package main

import (
	"menutest/screens/game"
	"menutest/screens/menu"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Ctrl struct {
	Exit    bool
	nav     chan string
	log     chan string
	Current string
}

func (c *Ctrl) Navigate(to string) {
	c.Current = to
}

func (c *Ctrl) Quit() {
	c.Exit = true
}

func (c *Ctrl) Log(msg string) {
	c.log <- msg
}

func NewController() *Ctrl {
	return &Ctrl{
		nav: make(chan string),
		log: make(chan string),
	}
}

type Screen interface {
	Draw(win *pixelgl.Window)
}

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Test menu",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	face, err := loadTTF("capture.ttf", 20)
	if err != nil {
		panic(err)
	}
	basicAtlas := text.NewAtlas(face, text.ASCII)

	ctrl := NewController()

	//	go ctrl.Run()

	registry := make(map[string]Screen)
	registry["firstcreen"] = menu.NewMain(basicAtlas, ctrl)
	registry["secondcreen"] = menu.NewSec(basicAtlas, ctrl)
	registry["game"] = game.New(ctrl)

	ctrl.Navigate("firstcreen")

	for !win.Closed() && !ctrl.Exit {
		win.Clear(colornames.Whitesmoke)
		registry[ctrl.Current].Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
