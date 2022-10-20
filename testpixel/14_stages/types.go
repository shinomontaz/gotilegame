package main

import (
	"github.com/faiface/pixel/pixelgl"
)

const (
	EVENT_DONE = iota
	EVENT_ENTER
	EVENT_QUIT
	EVENT_MENU
	EVENT_NOTREADY
)

const (
	STAGE_LOADING = iota
	STAGE_MENU
	STAGE_GAME
	STAGE_RECORDS
)

type Job func()
type Inform func(e int)

type StageOpt func(s IStage)

type IStage interface {
	GetID() int
	Run(win *pixelgl.Window, dt float64)
	Start()
	SetJob(j Job)
	Init()
	GetNext(event int) (int, bool)
	SetNext(event, id int)
	SetUp(opts ...StageOpt)
}
