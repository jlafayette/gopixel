package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Particle is a dynamic point that is being simulated and responds to forces.
type Particle struct {
	pos      pixel.Vec
	prevPos  pixel.Vec
	acc      pixel.Vec
	vel      pixel.Vec
	color    pixel.RGBA
	mass     float64
	radius   float64
	moveable bool
	visible  bool
}

// NewParticle instantiates a new particle
func NewParticle(pos, vel pixel.Vec, mass float64) Particle {
	r := radiusFromMass(mass)
	c := randomColor()
	return Particle{
		pos:      pos,
		prevPos:  pos,
		acc:      pixel.V(0, 0),
		vel:      vel,
		color:    c,
		mass:     mass,
		radius:   r,
		moveable: true,
		visible:  true,
	}
}

// NewOrbiter ...
func NewOrbiter(a Particle, mass float64, bounds pixel.Rect, minOffset, maxOffset float64) Particle {
	pos := randomPos(bounds)

	toAttractor := pos.To(a.pos)
	angle := toAttractor.Normal().Angle()

	// equation for circular orbit.
	magnitude := math.Sqrt((G * (1 + a.mass)) / toAttractor.Len())

	// random velocity offset
	r := minOffset + rand.Float64()*(maxOffset-minOffset)
	magnitude = magnitude * r

	vel := a.vel.Add(pixel.Unit(angle).Scaled(magnitude))
	return NewParticle(pos, vel, mass)
}

func (p *Particle) update() {
	p.prevPos = p.pos
	if p.moveable {
		p.pos = p.pos.Add(p.vel)
		p.vel = p.vel.Add(p.acc)
	}
}

func (p *Particle) drawTrail(imd *imdraw.IMDraw) {
	if p.visible {
		imd.Color = p.color
		imd.Push(p.prevPos)
		imd.Push(p.pos)
		imd.Line(p.radius)
	}
}

func (p *Particle) drawPos(imd *imdraw.IMDraw) {
	if p.visible {
		imd.Color = p.color
		imd.Push(p.pos)
		imd.Circle(p.radius, 0)
	}
}

func radiusFromMass(mass float64) float64 {
	// A/PI = r2
	r := (math.Cbrt(mass/math.Pi) + math.Sqrt(mass/math.Pi)) / 2
	return math.Max(r, 1)
}
