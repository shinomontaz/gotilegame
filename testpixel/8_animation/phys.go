package main

import (
	"github.com/faiface/pixel"
)

type Phys struct {
	gravity   float64
	runSpeed  float64
	walkSpeed float64

	jumpSpeed float64

	rect pixel.Rect
	vel  pixel.Vec

	pos    pixel.Vec
	ground bool
}

func NewPhys() *Phys {
	p := Phys{
		speed: 5,
		pos:   pixel.V(0.0, 0.0),
	}

	p.rect = pixel.R(0, 0, xstep, ystep)
	return &p
}

func (h *Hero) Draw(t pixel.Target, pos pixel.Vec) {
	if h.frame >= len(h.frames[h.state]) {
		h.frame = 0
	}
	rect := h.frames[h.state][h.frame]
	soldier := pixel.NewSprite(h.spritesheet, rect)
	soldier.Draw(t, pixel.IM.ScaledXY(pixel.ZV, pixel.V(h.dir, 1)).Scaled(pixel.ZV, 1.5).Moved(pos))
}

func (h *Hero) Update(platforms []Platform) {
	// check collisions
	h.rect = h.rect.Moved(h.pos)

	// check collisions against each platform
	h.ground = false
	if h.state == WALKING {
		for _, p := range platforms {
			if h.rect.Max.X <= p.Rect.Min.X || h.rect.Min.X >= p.Rect.Max.X {
				continue
			}
			if h.rect.Min.Y > p.Rect.Max.Y || h.rect.Min.Y < p.Rect.Max.Y+h.pos.Y {
				continue
			}
			h.rect = h.rect.Moved(pixel.V(0, p.Rect.Max.Y-h.rect.Min.Y))
			h.ground = true
		}
	}
}

func (h *Hero) Tick() {
	if h.state != DEAD && h.state != DYING {
		h.state = STANDING
	}

	h.frame++
	if h.state == DYING && h.frame >= len(h.frames[h.state]) {
		h.state = DEAD
	}
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
		h.dir = TO_LEFT
		h.pos = h.pos.Add(pixel.V(1.0, 0))
	case RIGHT:
		h.state = WALKING
		h.dir = TO_RIGHT
		h.pos = h.pos.Add(pixel.V(-1.0, 0))
	case ENTER:
		h.state = DYING
	}

	if h.prevstate != h.state {
		h.frame = 0
	}

	h.prevstate = h.state
}

func (h *Hero) setPos(pos pixel.Vec) {
	h.pos = pos
}

func (h *Hero) getPos() pixel.Vec {
	return h.pos
}
