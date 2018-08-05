package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Target is a spot that the vehicle is trying to get to.
type Target struct {
	pos pixel.Vec
}

// NewTarget instantiates a new target
func NewTarget(pos pixel.Vec) Target {
	return Target{
		pos: pos,
	}
}

func (t *Target) update() {
}

func (t *Target) draw(imd *imdraw.IMDraw) {
	imd.Push(t.pos)
	imd.Circle(6, 2)
	imd.Push(t.pos.Sub(pixel.V(6, 0)))
	imd.Push(t.pos.Add(pixel.V(6, 0)))
	imd.Line(1)
	imd.Push(t.pos.Sub(pixel.V(0, 6)))
	imd.Push(t.pos.Add(pixel.V(0, 6)))
	imd.Line(1)
}
