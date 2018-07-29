package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Particle is a dynamic point that is being simulated and responds to forces.
type Particle struct {
	pos  pixel.Vec
	acc  pixel.Vec
	vel  pixel.Vec
	mass float64
}

// NewParticle instantiates a new particle
func NewParticle(pos, vel pixel.Vec) Particle {
	return Particle{
		pos:  pos,
		acc:  pixel.V(0, 0),
		vel:  vel,
		mass: 1,
	}
}

func (p *Particle) update() {
	p.pos = p.pos.Add(p.vel)
	p.vel = p.vel.Add(p.acc)
}

func (p *Particle) draw(imd *imdraw.IMDraw) {
	imd.Push(p.pos)
	imd.Circle(3, 0)
}
