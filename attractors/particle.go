package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Particle is a dynamic point that is being simulated and responds to forces.
type Particle struct {
	pos     pixel.Vec
	prevPos pixel.Vec
	acc     pixel.Vec
	vel     pixel.Vec
	mass    float64
	radius  float64
	color   color.RGBA
}

// NewParticle instantiates a new particle
func NewParticle(pos, vel pixel.Vec, mass float64) Particle {
	r := radiusFromMass(mass)
	c := randomColor()
	return Particle{
		pos:     pos,
		prevPos: pos,
		acc:     pixel.V(0, 0),
		vel:     vel,
		mass:    mass,
		radius:  r,
		color:   c,
	}
}

// NewOrbiter ...
func NewOrbiter(a Particle) Particle {

	// randomized position
	edgeOffset := 200.0
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
	min := .5
	max := 1.05
	r := min + rand.Float64()*(max-min)
	magnitude = magnitude * r

	vel := pixel.Unit(angle).Scaled(magnitude)
	return NewParticle(pos, vel, .1)
}

func (p *Particle) update() {
	p.prevPos = p.pos
	p.pos = p.pos.Add(p.vel)
	p.vel = p.vel.Add(p.acc)
}

func (p *Particle) draw(imd *imdraw.IMDraw) {
	imd.Color = p.color
	imd.Push(p.prevPos)
	imd.Push(p.pos)
	imd.Line(p.radius)
}

func radiusFromMass(mass float64) float64 {
	// A/PI = r2
	r := math.Sqrt(mass / math.Pi)
	return math.Max(r, 1)
}

// randomColor generates a random color
func randomColor() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
}
