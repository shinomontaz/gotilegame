package main

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/pixelgl"
)

type Stage2 struct {
	id       int
	done     chan struct{}
	isReady  bool
	inform   Inform
	j        Job
	eventMap map[int]int
}

func NewStage2(f Inform) *Stage2 {
	return &Stage2{
		id:       STAGE_MENU,
		done:     make(chan struct{}),
		inform:   f,
		eventMap: map[int]int{EVENT_ENTER: STAGE_GAME},
	}
}

func (s *Stage2) GetID() int {
	return s.id
}

func (s *Stage2) Init() {
	time.Sleep(2 * time.Second)
	s.isReady = true
}

func (s *Stage2) SetUp(opts ...StageOpt) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Stage2) Start() {
	if !s.isReady {
		return
	}

	if s.j != nil {
		go func() {
			s.j()
			s.done <- struct{}{}
		}()

	}
	fmt.Println("stage started", s.id)
}

func (s *Stage2) SetJob(j Job) {
	s.j = j
}

func (s *Stage2) Notify(event int) {
	s.inform(event)
}

func (s *Stage2) Run(win *pixelgl.Window, dt float64) {
	if !s.isReady {
		s.Notify(EVENT_NOTREADY)
		return
	}
	select {
	case <-s.done:
		s.Notify(EVENT_DONE)
	default:
		win.Clear(colornames.Aqua)
		if win.Pressed(pixelgl.KeyEnter) {
			s.Notify(EVENT_ENTER)
		}
		if win.Pressed(pixelgl.KeyEscape) {
			s.Notify(EVENT_QUIT)
		}
	}
}

func (s *Stage2) SetNext(event, id int) {
	s.eventMap[event] = id
}

func (s *Stage2) GetNext(event int) (int, bool) {
	next, ok := s.eventMap[event]
	return next, ok
}
