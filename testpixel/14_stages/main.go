// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func inform(e int) {
	switch e {
	case EVENT_DONE:
		fmt.Println("event done")
		next, ok := currStage.GetNext(EVENT_DONE)
		if ok {
			setStage(next)
		}
	case EVENT_ENTER:
		fmt.Println("event enter")
		next, ok := currStage.GetNext(EVENT_ENTER)
		if ok {
			setStage(next)
		}
	case EVENT_QUIT:
		fmt.Println("event quit")
		next, ok := currStage.GetNext(EVENT_QUIT)
		if ok {
			setStage(next)
		} else {
			isquit = true
		}
	case EVENT_NOTREADY:
		fmt.Println("event not ready")
		loadingStage.SetUp(WithJob(currStage.Init), WithNext(EVENT_DONE, currStage.GetID()))
		setStage(loadingStage.GetID())
	}
}

func setStage(id int) {
	currStage = stages[id]
	currStage.Start()
}

var currStage IStage
var isquit bool
var stages map[int]IStage
var loadingStage IStage

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Hello, World!",
		Bounds: pixel.R(0, 0, 500, 350),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	stages = make(map[int]IStage, 0)
	loadingStage = NewStage1(inform)
	stages[STAGE_LOADING] = loadingStage
	stages[STAGE_MENU] = NewStage2(inform)
	//	stages[2].SetUp(WithNext(EVENT_ENTER, 3))
	stages[STAGE_GAME] = NewStage3(inform)

	currStage = stages[STAGE_LOADING]

	currStage.SetUp(WithJob(stages[STAGE_MENU].Init), WithNext(EVENT_DONE, STAGE_MENU))
	currStage.Init()
	currStage.Start()

	last := time.Now()

	for !win.Closed() && !isquit {
		dt := time.Since(last).Seconds()
		currStage.Run(win, dt)
		win.Update()
	}

}
