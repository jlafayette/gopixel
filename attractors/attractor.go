package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Attractor is a point that exerts attraction force to other particles.
type Attractor struct {
	pos  pixel.Vec
	mass float64
}

// NewAttractor instantiates a new attractor
func NewAttractor(pos pixel.Vec) Attractor {
	return Attractor{
		pos:  pos,
		mass: 10,
	}
}

func (p *Attractor) update() {

}

func (p *Attractor) draw(imd *imdraw.IMDraw) {
	imd.Push(p.pos)
	imd.Circle(10, 0)
}
