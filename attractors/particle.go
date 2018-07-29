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

// NewOrbiter ...
func NewOrbiter(pos pixel.Vec, a Attractor) Particle {
	// pos := pixel.V(300, 400)
	toAttractor := pos.To(a.pos)
	angle := toAttractor.Normal().Angle()

	// distanceSq := math.Pow(math.Max(toAttractor.Len(), a.radius), 2)
	// magnitude := G * ((1 * a.mass) / distanceSq) * 120

	vel := pixel.Unit(angle).Scaled(.7)
	return NewParticle(pos, vel)
}

func (p *Particle) update() {
	p.pos = p.pos.Add(p.vel)
	p.vel = p.vel.Add(p.acc)
}

func (p *Particle) draw(imd *imdraw.IMDraw) {
	imd.Push(p.pos)
	imd.Circle(1, 0)
}
