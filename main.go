package main

import (
	"fmt"
	"math/rand"
	"time"

	"gotilegame/config"
	"gotilegame/screens/start"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"github.com/faiface/mainthread"
)

var ev *config.Env
var chErrors chan error

func init() {
	rand.Seed(time.Now().UnixNano())
	ev = config.NewEnv("./config")
	ev.InitLog()

	chErrors = make(chan error, 1000)

	go func() {
		for err := range chErrors {
			fmt.Println("Error", err)
		}
	}()
}

func run() {
	w, err := win.New(win.Title(ev.Cfg.Title), win.Size(800, 600))
	if err != nil {
		panic(err)
	}

	mux, env := gui.NewMux(w)
	ss := start.New(w, mux, env)
	ss.Run()
}

func main() {
	mainthread.Run(run)
}
