package start

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"gotilegame/components/console"
	"gotilegame/types"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
)

type Screen struct {
	w      *win.Win
	mux    *gui.Mux
	env    gui.Env
	cv     image.Rectangle
	cfg    Config
	cmpnts []types.IComponent
}

type Config struct {
	X1    int
	Y1    int
	X2    int
	Y2    int
	Color color.Color
}

var chConsole chan string

func New(w *win.Win, mux *gui.Mux, env gui.Env) *Screen {
	chConsole = make(chan string, 100)
	s := &Screen{w: w, mux: mux, env: env, cmpnts: make([]types.IComponent, 0)}
	s.init()
	return s
}

func (s *Screen) Run() {
	for component := range s.cmpnts {
		go component.Run()
	}

	for event := range s.env.Events() {
		switch event.(type) {
		case win.WiClose:
			close(s.env.Draw())
		case win.MoDown:
			e := event.(win.MoDown)
			if e.Button == win.ButtonLeft {
				chConsole <- "left click"
			} else {
				chConsole <- "another click"
			}
		}
	}

	close(s.env.Draw())
}

func (s *Screen) init() { // create subcomponents
	// read json of geometry and styles
	file, _ := os.Open("screens/start/template.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&s.cfg)
	if err != nil {
		log.Fatal(err)
	}

	s.cv = image.Rect(s.cfg.X1, s.cfg.Y1, s.cfg.X2, s.cfg.Y2)
	// read subcomponents

	c := console.New(chConsole, image.Rect(0, 100, 350, 250), s.mux.MakeEnv())
	s.cmpnts = append(s.cmpnts, c)
}

//func(r image.Rectangle, images []image.Image) func(draw.Image) image.Rectangle
//func (s *Screen) Draw(drw draw.Image) image.Rectangle {
func (s *Screen) Draw(r image.Rectangle) func(draw.Image) image.Rectangle {
	return func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, r, &image.Uniform{s.cfg.Color}, image.ZP, draw.Src)
		newImage := image.NewRGBA(image.Rect(s.cfg.X1, s.cfg.Y1, s.cfg.X2, s.cfg.Y2))

		for component := range s.cmpnts {
			component.Draw(r)
		}

		bounds := newImage.Bounds()
		draw.Draw(drw, bounds.Intersect(r), newImage, bounds.Min, draw.Src)
		return r
	}
}

func (c *Config) UnmarshalJSON(b []byte) error {
	var tempJson struct {
		cl color.Color
	}
	// Unmarshal to our temp struct
	err := json.Unmarshal(b, &tempJson)
	if err != nil {
		return err
	}
	// convert our new friends O(n) to the interface type
	c.Color = color.Color(tempJson.cl)
	return nil
}
