package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Boid ...
type Boid struct {
	pos           pixel.Vec
	acc           pixel.Vec
	vel           pixel.Vec
	col           pixel.RGBA
	colShade      pixel.RGBA
	velCol        pixel.RGBA
	accCol        pixel.RGBA
	maxSpeed      float64
	maxForce      float64
	alignDistance float64
}

// NewBoid ...
func NewBoid(pos pixel.Vec) Boid {
	maxSpeed := randFloat(3, 5)
	maxForce := randFloat(.2, .3)

	return Boid{
		pos:           pos,
		acc:           pixel.ZV,
		vel:           pixel.ZV,
		col:           pixel.RGB(0, .8, 0),
		colShade:      pixel.RGB(0, .2, 0),
		velCol:        pixel.RGB(1, 0, 0),
		accCol:        pixel.RGB(0, 0, 1),
		maxSpeed:      maxSpeed,
		maxForce:      maxForce,
		alignDistance: 50,
	}
}

func (b *Boid) update(bounds pixel.Rect, allboids []Boid) {
	desired := b.align(allboids)
	steering := desired.Sub(b.vel)
	if steering.Len() > b.maxForce {
		steering = pixel.Unit(steering.Angle()).Scaled(b.maxForce)
	}
	b.acc = steering
	b.pos = b.pos.Add(b.vel)
	if !bounds.Contains(b.pos) {
		if b.pos.X < 0 {
			b.pos = b.pos.Add(pixel.V(screenWidth, 0))
		} else if b.pos.X > screenWidth {
			b.pos = b.pos.Sub(pixel.V(screenWidth, 0))
		} else if b.pos.Y < 0 {
			b.pos = b.pos.Add(pixel.V(0, screenHeight))
		} else if b.pos.Y > screenHeight {
			b.pos = b.pos.Sub(pixel.V(0, screenHeight))
		}
	}
	b.vel = b.vel.Add(b.acc)
}

func (b *Boid) align(neighbors []Boid) pixel.Vec {
	var avgDir float64
	var avgSpeed float64
	count := 0
	for i := 0; i < len(neighbors); i++ {
		distance := b.pos.To(neighbors[i].pos).Len()
		if distance > 0 && distance < b.alignDistance {

			// strength is 1 when close, 0 when far
			// lesserMaxSpeed := math.Min(b.maxSpeed, neighbors[i].maxSpeed)
			// avgSpeed = avgSpeed + mapRange(distance, 0, b.alignDistance, neighbors[i].vel.Len(), lesserMaxSpeed)

			avgDir = avgDir + neighbors[i].vel.Angle()
			avgSpeed = avgSpeed + neighbors[i].vel.Len()
			count++
		}
	}
	if count > 0 {
		avgDir = avgDir / float64(count)
		avgSpeed = avgSpeed / float64(count)
		return pixel.Unit(avgDir).Scaled(avgSpeed)
	}
	return b.vel
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

func mapRange(in, inMin, inMax, outMin, outMax float64) float64 {
	return (in-inMin)/(inMax-inMin)*(outMax-outMin) + outMin
}
