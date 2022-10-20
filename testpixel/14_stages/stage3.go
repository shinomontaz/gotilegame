package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Stage3 struct {
	id       int
	done     chan struct{}
	isReady  bool
	inform   Inform
	j        Job
	eventMap map[int]int
	sprite   *pixel.Sprite
}

func NewStage3(f Inform) *Stage3 {
	return &Stage3{
		id:       STAGE_GAME,
		done:     make(chan struct{}),
		inform:   f,
		eventMap: map[int]int{EVENT_QUIT: STAGE_MENU},
	}
}

func (s *Stage3) GetID() int {
	return s.id
}

func (s *Stage3) Init() {
	pic, err := loadPicture("test_sprite.png")
	if err != nil {
		panic(err)
	}

	s.sprite = pixel.NewSprite(pic, pic.Bounds())

	time.Sleep(2 * time.Second)

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
	fmt.Println("stage started", s.id)
}

func (s *Stage3) SetJob(j Job) {
	s.j = j
}

func (s *Stage3) Notify(event int) {
	s.inform(event)
}

func (s *Stage3) Run(win *pixelgl.Window, dt float64) {
	if !s.isReady {
		s.Notify(EVENT_NOTREADY)
		return
	}
	select {
	case <-s.done:
		s.Notify(EVENT_DONE)
	default:
		s.draw(win, dt)
	}
}

func (s *Stage3) draw(win *pixelgl.Window, dt float64) {
	win.Clear(colornames.Skyblue)
	s.sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	if win.Pressed(pixelgl.KeyEnter) {
		s.Notify(EVENT_ENTER)
	}
	if win.Pressed(pixelgl.KeyEscape) {
		s.Notify(EVENT_QUIT)
	}
}

func (s *Stage3) SetNext(event, id int) {
	s.eventMap[event] = id
}

func (s *Stage3) GetNext(event int) (int, bool) {
	next, ok := s.eventMap[event]
	return next, ok
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
