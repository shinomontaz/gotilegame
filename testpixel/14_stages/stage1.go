package main

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Stage1 struct {
	id       int
	done     chan struct{}
	isReady  bool
	inform   Inform
	j        Job
	eventMap map[int]int
}

func NewStage1(f Inform) *Stage1 {
	return &Stage1{
		id:       STAGE_LOADING,
		done:     make(chan struct{}),
		inform:   f,
		eventMap: make(map[int]int),
	}
}

func (s *Stage1) GetID() int {
	return s.id
}

func (s *Stage1) Init() {
	s.isReady = true
}

func (s *Stage1) SetUp(opts ...StageOpt) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Stage1) Start() {
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

func (s *Stage1) SetJob(j Job) {
	s.j = j
}

func (s *Stage1) Notify(event int) {
	s.inform(event)
}

func (s *Stage1) Run(win *pixelgl.Window, dt float64) {
	select {
	case <-s.done:
		s.Notify(EVENT_DONE)
	default:
		win.Clear(colornames.Black)
	}
}

func (s *Stage1) SetNext(event, id int) {
	s.eventMap[event] = id
}

func (s *Stage1) GetNext(event int) (int, bool) {
	next, ok := s.eventMap[event]
	return next, ok
}
