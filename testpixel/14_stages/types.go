package main

const (
	EVENT_DONE = iota
	EVENT_ENTER
	EVENT_NOTREADY
)

type Job func()
type Inform func(e int)

type StageOpt func(s IStage)

type IStage interface {
	GetID() int
	Run(dt float64)
	Start()
	SetJob(j Job)
	Init()
	GetNext(event int) int
	SetNext(event, id int)
	SetUp(opts ...StageOpt)
}
