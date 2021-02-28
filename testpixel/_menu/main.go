package main

import (
	"fmt"
	"menutest/screens/game"
	"menutest/screens/menu"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

const (
	STATE_GAME = iota
	STATE_NOGAME
)

type Ctrl struct {
	Exit      bool
	nav       chan string
	log       chan string
	Current   string
	GameState int
	Mc        *MusicController
}

func (c *Ctrl) Navigate(to string) {
	c.Current = to

	if to == "firstcreen" {
		c.Mc.SetAmbient("monster")
	}
	if to == "game" {
		c.Mc.SetAmbient("anxiety")
	}

}

func (c *Ctrl) Quit() {
	c.Exit = true
}

func (c *Ctrl) Log(msg string) {
	c.log <- msg
}

func NewController(mc *MusicController) *Ctrl {
	return &Ctrl{
		GameState: STATE_NOGAME,
		nav:       make(chan string),
		log:       make(chan string),
		Mc:        mc,
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

	mctrl := NewMusicController()
	ctrl := NewController(mctrl)
	go mctrl.run()

	registry := make(map[string]Screen)
	registry["firstcreen"] = menu.NewMain(basicAtlas, ctrl)
	registry["secondcreen"] = menu.NewSec(basicAtlas, ctrl)
	registry["game"] = game.New(ctrl, basicAtlas)

	ctrl.Navigate("firstcreen")

	for !win.Closed() && !ctrl.Exit {
		win.Clear(colornames.Whitesmoke)
		registry[ctrl.Current].Draw(win)
		win.Update()
	}
}

func NewMusicController() *MusicController {
	list := loadMusic()
	volumes := make(map[string]*effects.Volume)

	volumes["music"] = &effects.Volume{
		Base:   2,
		Volume: 0,
		Silent: false,
	}

	volumes["effects"] = &effects.Volume{
		Base:   2,
		Volume: 0,
		Silent: false,
	}

	mc := &MusicController{
		effects: list,
		volumes: volumes,
		Effects: make(chan string),
	}

	return mc
}

func (mc *MusicController) run() {
	for cmd := range mc.Effects {
		if _, ok := mc.effects[cmd]; ok {
			sound := mc.effects[cmd].Streamer(0, mc.effects[cmd].Len())
			mc.volumes["effects"].Streamer = sound
			speaker.Play(mc.volumes["effects"])
		}
	}
}

type MusicController struct {
	//	musics map[string]*beep.Buffer
	effects map[string]*beep.Buffer
	volumes map[string]*effects.Volume
	ambient string
	Effects chan string
}

func (mc *MusicController) SetAmbient(name string) {
	if _, ok := mc.effects[name]; ok {
		speaker.Lock()
		fmt.Println("SetAmbient", name)
		sound := mc.effects[name].Streamer(0, mc.effects[name].Len())
		loopedSound := beep.Loop(-1, sound)
		mc.volumes["music"].Streamer = loopedSound
		speaker.Unlock()

		speaker.Play(mc.volumes["music"])
	}
}

func main() {
	pixelgl.Run(run)
}
