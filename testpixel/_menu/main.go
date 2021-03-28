package main

import (
	"fmt"
	"log"
	"menutest/screens/game"
	"menutest/screens/menu"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
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

func (c *Ctrl) Sound(name string) {
	c.Mc.Effects <- name
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

	mctrl := NewSoundController()
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

func NewSoundController() *MusicController {
	// load music
	// load sounds

	list, baseSR := loadMusic()
	volumes := make(map[string]*effects.Volume)

	volumes["music"] = &effects.Volume{
		Base:   1,
		Volume: 0,
		Silent: false,
	}

	volumes["effects"] = &effects.Volume{
		Base:   1,
		Volume: 0,
		Silent: false,
	}

	mc := &MusicController{
		musics:         make(map[string]string),
		effects:        list,
		volumes:        volumes,
		Effects:        make(chan string),
		baseSampleRate: baseSR,
	}

	mc.musics["monster"] = "music\\monster.mp3"
	mc.musics["fear"] = "music\\fear.mp3"
	mc.musics["anxiety"] = "music\\anxiety.mp3"

	return mc
}

type MusicController struct {
	musics         map[string]string
	effects        map[string]*beep.Buffer
	volumes        map[string]*effects.Volume
	baseSampleRate beep.SampleRate
	currAmbient    string
	currFile       *os.File
	currMusic      beep.StreamSeekCloser
	ambient        string
	Effects        chan string
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

func (mc *MusicController) SetAmbient(name string) {
	var err error
	//	var format beep.Format
	if mc.currAmbient == name {
		return
	}

	if path, ok := mc.musics[name]; ok {
		if mc.currMusic != nil {
			mc.currMusic.Close()
		}

		mc.currFile, err = os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		mc.currMusic, _, err = mp3.Decode(mc.currFile)
		if err != nil {
			log.Fatal(err)
		}

		speaker.Lock()
		fmt.Println("SetAmbient", name)

		mc.volumes["music"].Streamer = &beep.Ctrl{Streamer: beep.Loop(-1, mc.currMusic), Paused: false}
		speaker.Unlock()

		speaker.Play(beep.ResampleRatio(3, 1, mc.volumes["music"]))
		mc.currAmbient = name
	}
}

func main() {
	pixelgl.Run(run)
}
