package start

import (
	"fmt"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
)

type Screen struct {
	w   *win.Win
	mux *gui.Mux
	env gui.Env
}

func New(w *win.Win, mux *gui.Mux, env gui.Env) *Screen {
	return &Screen{w: w, mux: mux, env: env}
}

func (s *Screen) Run() {
	for event := range s.env.Events() {
		switch event.(type) {
		case win.WiClose:
			close(s.env.Draw())
		case win.MoDown:
			e := event.(win.MoDown)
			if e.Button == win.ButtonLeft {
				fmt.Println("left click")
			} else {
				fmt.Println("another click")
			}
		}
	}
}
