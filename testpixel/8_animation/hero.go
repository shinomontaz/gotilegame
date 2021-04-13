package main

import (
	"github.com/faiface/pixel"
)

const (
	FIRE = iota
	LEFT
	RIGHT
	ENTER
)

const (
	WALKING = iota
	STANDING
	FIRING
	DYING
	DEAD
)

const (
	TO_LEFT  = -1.0
	TO_RIGHT = 1.0
)

type Hero struct {
	walkFrames  []pixel.Rect
	shootFrame  pixel.Rect
	standFrame  pixel.Rect
	deathFrames []pixel.Rect
	deadFrame   pixel.Rect
	prevstate   int
	state       int
	dir         float64
	frame       int
	frameLimit  int
	spritesheet pixel.Picture
}

func NewHero() *Hero {
	h := Hero{
		state: STANDING,
		dir:   TO_RIGHT,
	}

	spritesheet, err := loadPicture("hero_spritesheet.png")
	if err != nil {
		panic(err)
	}

	h.spritesheet = spritesheet

	xstep := 80.0
	ystep := 94.0
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Min.X+6*xstep; x += xstep {
		h.walkFrames = append(h.walkFrames, pixel.R(x, spritesheet.Bounds().Max.Y-ystep, x+xstep, spritesheet.Bounds().Max.Y-2*ystep))
	}
	h.shootFrame = pixel.R(spritesheet.Bounds().Min.X, spritesheet.Bounds().Max.Y-2*ystep, spritesheet.Bounds().Min.X+xstep, spritesheet.Bounds().Max.Y-3*ystep)
	h.standFrame = pixel.R(spritesheet.Bounds().Min.X, spritesheet.Bounds().Max.Y, spritesheet.Bounds().Min.X+xstep, spritesheet.Bounds().Max.Y-ystep)
	for x := spritesheet.Bounds().Min.X + 3*xstep; x < spritesheet.Bounds().Min.X+6*xstep; x += xstep {
		h.deathFrames = append(h.deathFrames, pixel.R(x, spritesheet.Bounds().Max.Y-3*ystep, x+xstep, spritesheet.Bounds().Max.Y-4*ystep))
	}
	h.deadFrame = pixel.R(spritesheet.Bounds().Min.X+6*xstep, spritesheet.Bounds().Max.Y-3*ystep, spritesheet.Bounds().Min.X+6*xstep+ystep, spritesheet.Bounds().Max.Y-4*ystep)
	return &h
}

func (h *Hero) Draw(t pixel.Target, matrix pixel.Matrix) {
	// get frame
	var rect pixel.Rect
	switch h.state {
	case STANDING:
		rect = h.standFrame
	case WALKING:
		rect = h.walkFrames[h.frame]
	case FIRING:
		rect = h.shootFrame
	case DYING:
		rect = h.deathFrames[h.frame]
	case DEAD:
		rect = h.deadFrame
	}

	soldier := pixel.NewSprite(h.spritesheet, rect)
	soldier.Draw(t, matrix)
}

func (h *Hero) Tick() {
	if h.state != DEAD && h.state != DYING {
		h.state = STANDING
	}
	if h.frameLimit == 0 {
		return
	}

	h.frame++
	if h.state == DYING && h.frame >= h.frameLimit {
		h.state = DEAD
	}
	h.frame = h.frame % h.frameLimit
}

func (h *Hero) Notify(action int) {
	if h.state == DYING || h.state == DEAD {
		return
	}

	switch action {
	case FIRE:
		h.state = FIRING
	case LEFT:
		h.state = WALKING
		h.frameLimit = len(h.walkFrames)
		h.dir = TO_LEFT
	case RIGHT:
		h.state = WALKING
		h.frameLimit = len(h.walkFrames)
		h.dir = TO_RIGHT
	case ENTER:
		h.state = DYING
		h.frameLimit = len(h.deathFrames)
	}

	if h.prevstate != h.state {
		h.frame = 0
	}

	h.prevstate = h.state
}
