package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Boid ...
type Boid struct {
	pos      pixel.Vec
	acc      pixel.Vec
	vel      pixel.Vec
	col      pixel.RGBA
	colShade pixel.RGBA
	velCol   pixel.RGBA
	accCol   pixel.RGBA
	maxSpeed float64
	maxForce float64
}

// NewBoid ...
func NewBoid(pos pixel.Vec) Boid {
	maxSpeed := randFloat(3, 5)
	maxForce := randFloat(.2, .3)

	return Boid{
		pos:      pos,
		acc:      pixel.ZV,
		vel:      pixel.ZV,
		col:      pixel.RGB(0, .8, 0),
		colShade: pixel.RGB(0, .2, 0),
		velCol:   pixel.RGB(1, 0, 0),
		accCol:   pixel.RGB(0, 0, 1),
		maxSpeed: maxSpeed,
		maxForce: maxForce,
	}
}

func (b *Boid) update() {
	b.pos = b.pos.Add(b.vel)
	b.vel = b.vel.Add(b.acc)
}

func (b *Boid) draw(imd *imdraw.IMDraw) {
	imd.Color = b.col
	imd.Push(b.pos)
	imd.Circle(5, 0)
	imd.Color = b.colShade
	imd.Push(b.pos)
	imd.Circle(5, 1)
	imd.Color = b.velCol
	imd.Push(b.pos)
	imd.Push(b.pos.Add(b.vel.Scaled(5)))
	imd.Line(1)
	imd.Color = b.accCol
	imd.Push(b.pos)
	imd.Push(b.pos.Add(b.acc.Scaled(35)))
	imd.Line(1)
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
