package main

import (
	"fmt"
	"time"
)

type Stage2 struct {
	ID      int
	done    chan struct{}
	isReady bool
	inform  Inform
	j       Job
	nextId  int
}

func NewStage2(id int, f Inform) *Stage2 {
	return &Stage2{
		ID:     id,
		done:   make(chan struct{}),
		inform: f,
	}
}

func (s *Stage2) GetID() int {
	return s.ID
}

func (s *Stage2) Init() {
	time.Sleep(10000)
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
	fmt.Println("stage started", s.ID)
}

func (s *Stage2) SetJob(j Job) {
	s.j = j
}

func (s *Stage2) Notify(event int) {
	s.inform(event)
}

func (s *Stage2) Run(dt float64) {
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

func (s *Stage2) SetNext(event, id int) {
	s.nextId = id
}

func (s *Stage2) GetNext(event int) int {
	return s.nextId
}
