package main

import "fmt"

type Stage1 struct {
	ID      int
	done    chan struct{}
	isReady bool
	inform  Inform
	j       Job
	nextId  int
}

func NewStage1(id int, f Inform) *Stage1 {
	return &Stage1{
		ID:     id,
		done:   make(chan struct{}),
		inform: f,
	}
}

func (s *Stage1) GetID() int {
	return s.ID
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
	fmt.Println("stage started", s.ID)
}

func (s *Stage1) SetJob(j Job) {
	s.j = j
}

func (s *Stage1) Notify(event int) {
	s.inform(event)
}

func (s *Stage1) Run(dt float64) {
	select {
	case <-s.done:
		s.Notify(EVENT_DONE)
	default:
		fmt.Println("Stage1 run", s.ID)
	}
}

func (s *Stage1) SetNext(event, id int) {
	s.nextId = id
}

func (s *Stage1) GetNext(event int) int {
	return s.nextId
}
