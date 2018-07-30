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
	pos      pixel.Vec
	prevPos  pixel.Vec
	acc      pixel.Vec
	vel      pixel.Vec
	mass     float64
	radius   float64
	moveable bool
	color    color.RGBA
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
		mass:     mass,
		radius:   r,
		moveable: true,
		color:    c,
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

func (p *Particle) draw(imd *imdraw.IMDraw) {
	if p.visible {
		imd.Color = p.color
		imd.Push(p.prevPos)
		imd.Push(p.pos)
		imd.Line(p.radius)
	}
}

func radiusFromMass(mass float64) float64 {
	// A/PI = r2
	r := (math.Cbrt(mass/math.Pi) + math.Sqrt(mass/math.Pi)) / 2
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

// similar colors
func similarRandomColor(i int) color.RGBA {
	switch i {
	case 0:
		return color.RGBA{
			uint8(rand.Intn(156) + 100),
			uint8(10),
			uint8(rand.Intn(256)),
			255,
		}
	case 1:
		return color.RGBA{
			uint8(10),
			uint8(rand.Intn(156) + 100),
			uint8(rand.Intn(175)),
			255,
		}
	case 2:
		return color.RGBA{
			uint8(rand.Intn(256)),
			uint8(rand.Intn(100)),
			uint8(rand.Intn(156) + 100),
			255,
		}
	}
	return randomColor()
}

func randomPos(bounds pixel.Rect) pixel.Vec {
	x := bounds.Min.X + rand.Float64()*(bounds.Max.X-bounds.Min.X)
	y := bounds.Min.Y + rand.Float64()*(bounds.Max.Y-bounds.Min.Y)
	return pixel.V(x, y)
}

func randomVel(max float64) pixel.Vec {
	x := -max + rand.Float64()*(max-(-max))
	y := -max + rand.Float64()*(max-(-max))
	return pixel.V(x, y)
}
