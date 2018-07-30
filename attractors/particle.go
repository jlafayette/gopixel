package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Particle is a dynamic point that is being simulated and responds to forces.
type Particle struct {
	pos    pixel.Vec
	acc    pixel.Vec
	vel    pixel.Vec
	mass   float64
	radius float64
}

// NewParticle instantiates a new particle
func NewParticle(pos, vel pixel.Vec) Particle {
	return Particle{
		pos:    pos,
		acc:    pixel.V(0, 0),
		vel:    vel,
		mass:   .1,
		radius: 1,
	}
}

// NewOrbiter ...
func NewOrbiter(a Attractor) Particle {

	// randomized position
	edgeOffset := 100.0
	xMin := edgeOffset
	yMin := edgeOffset
	xMax := screenWidth - edgeOffset
	yMax := screenHeight - edgeOffset
	x := xMin + rand.Float64()*(xMax-xMin)
	y := yMin + rand.Float64()*(yMax-yMin)
	pos := pixel.V(x, y)

	toAttractor := pos.To(a.pos)
	angle := toAttractor.Normal().Angle()

	// equation for circular orbit.
	magnitude := math.Sqrt((G * (1 + a.mass)) / toAttractor.Len())

	// random velocity offset
	// min := .5
	// max := 1.07
	min := .99
	max := 1.01
	r := min + rand.Float64()*(max-min)
	magnitude = magnitude * r

	vel := pixel.Unit(angle).Scaled(magnitude)
	return NewParticle(pos, vel)
}

func (p *Particle) update() {
	p.pos = p.pos.Add(p.vel)
	p.vel = p.vel.Add(p.acc)
}

func (p *Particle) draw(imd *imdraw.IMDraw) {
	imd.Push(p.pos)
	imd.Circle(p.radius, 0)
}
