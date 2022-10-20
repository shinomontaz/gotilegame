package main

import (
	"fmt"
	"time"
)

type Stage3 struct {
	ID      int
	done    chan struct{}
	isReady bool
	inform  Inform
	j       Job
	nextId  int
}

func NewStage3(id int, f Inform) *Stage3 {
	return &Stage3{
		ID:     id,
		done:   make(chan struct{}),
		inform: f,
	}
}

func (s *Stage3) GetID() int {
	return s.ID
}

func (s *Stage3) Init() {
	time.Sleep(10000)
	s.isReady = true
}

func (s *Stage3) SetUp(opts ...StageOpt) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Stage3) Start() {
	if !s.isReady {
		return
	}
	if s.j != nil {
		go func() {
			s.j()
			s.done <- struct{}{}
		}()

	}
	fmt.Println("stage started", s.ID)
}

func (s *Stage3) SetJob(j Job) {
	s.j = j
}

func (s *Stage3) Notify(event int) {
	s.inform(event)
}

func (s *Stage3) Run(dt float64) {
	if !s.isReady {
		s.Notify(EVENT_NOTREADY)
		return
	}
	select {
	case <-s.done:
		s.Notify(EVENT_DONE)
	default:
		//		fmt.Println(s.ID)
	}
}

func (s *Stage3) SetNext(event, id int) {
	s.nextId = id
}

func (s *Stage3) GetNext(event int) int {
	return s.nextId
}
