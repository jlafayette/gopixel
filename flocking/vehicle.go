package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Boid ...
type Boid struct {
	pos              pixel.Vec
	acc              pixel.Vec
	vel              pixel.Vec
	col              pixel.RGBA
	colShade         pixel.RGBA
	velCol           pixel.RGBA
	accCol           pixel.RGBA
	maxSpeed         float64
	maxForce         float64
	neighborDistance float64
}

// NewBoid ...
func NewBoid(pos pixel.Vec) Boid {
	maxSpeed := randFloat(3, 4)
	maxForce := randFloat(.15, .3)

	return Boid{
		pos:              pos,
		acc:              pixel.ZV,
		vel:              pixel.ZV,
		col:              pixel.RGB(0, .8, 0),
		colShade:         pixel.RGB(0, .2, 0),
		velCol:           pixel.RGB(1, 0, 0),
		accCol:           pixel.RGB(0, 0, 1),
		maxSpeed:         maxSpeed,
		maxForce:         maxForce,
		neighborDistance: 75,
	}
}

func (b *Boid) update(bounds pixel.Rect, allboids []Boid) {
	var neighbors []Boid
	for i := 0; i < len(allboids); i++ {
		distance := distance(b.pos, allboids[i].pos)
		if distance > 0 && distance < b.neighborDistance {
			neighbors = append(neighbors, allboids[i])
		}
	}

	// weight
	align := b.align(neighbors).Scaled(1)
	separate := b.separate(neighbors).Scaled(1)
	cohesion := b.cohere(neighbors).Scaled(.05)
	bias := pixel.V(-.005, 0)
	b.acc = align.Add(separate).Add(cohesion).Add(bias)

	// limit acc
	if b.acc.Len() > b.maxForce {
		b.acc = pixel.Unit(b.acc.Angle()).Scaled(b.maxForce)
	}

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

	// limit vel
	if b.vel.Len() > b.maxSpeed {
		b.vel = pixel.Unit(b.vel.Angle()).Scaled(b.maxSpeed)
	}
}

func (b *Boid) align(neighbors []Boid) pixel.Vec {
	var avgDir float64
	var avgSpeed float64
	count := 0
	for i := 0; i < len(neighbors); i++ {

		// strength is 1 when close, 0 when far
		// lesserMaxSpeed := math.Min(b.maxSpeed, neighbors[i].maxSpeed)
		// avgSpeed = avgSpeed + mapRange(distance, 0, b.alignDistance, neighbors[i].vel.Len(), lesserMaxSpeed)

		avgDir = avgDir + neighbors[i].vel.Angle()
		avgSpeed = avgSpeed + neighbors[i].vel.Len()
		count++

	}
	if count > 0 {
		avgDir = avgDir / float64(count)
		avgSpeed = avgSpeed / float64(count)
		desired := pixel.Unit(avgDir).Scaled(avgSpeed)
		return desired.Sub(b.vel)
	}
	return b.vel
}

func (b *Boid) separate(neighbors []Boid) pixel.Vec {
	desiredSeparation := 75.0

	// This blend determines how bouncy and chaotic the simulation
	// feels. Putting this all way to zero makes everything stop
	// after separating. 0.9 seems to work well if only separate is
	// used, but combined with other behaviors, any value less than
	// 1 will suck all the movement out of the group.

	// desired := pixel.Lerp(pixel.ZV, b.vel, .85)
	desired := b.vel

	count := 0
	for i := 0; i < len(neighbors); i++ {
		sepImpulse := wrapTo(neighbors[i].pos, b.pos)
		if sepImpulse.Len() < desiredSeparation {
			a := sepImpulse.Angle()
			m := sepImpulse.Len()
			strength := mapRange(desiredSeparation-m, 8, desiredSeparation, 0, b.maxSpeed)
			desired = desired.Add(pixel.Unit(a).Scaled(strength))
			count++
		}
	}
	if count > 0 {
		return desired.Sub(b.vel)
	}
	return b.vel
}

func (b *Boid) cohere(neighbors []Boid) pixel.Vec {

	if len(neighbors) == 0 {
		return b.pos
	}
	screen := pixel.R(0, 0, screenWidth, screenHeight)
	halfX := screen.Size().X * .5
	halfY := screen.Size().Y * .5
	avgX := b.pos.X
	avgY := b.pos.Y
	for i := 0; i < len(neighbors); i++ {
		if math.Abs(b.pos.X-neighbors[i].pos.X) > halfX {
			if b.pos.X > neighbors[i].pos.X {
				avgX = avgX + (neighbors[i].pos.X + screenWidth)
			} else {
				avgX = avgX + (neighbors[i].pos.X - screenWidth)
			}
		} else {
			avgX = avgX + neighbors[i].pos.X
		}
		if math.Abs(b.pos.Y-neighbors[i].pos.Y) > halfY {
			if b.pos.Y > neighbors[i].pos.Y {
				avgY = avgY + (neighbors[i].pos.Y + screenHeight)
			} else {
				avgY = avgY + (neighbors[i].pos.Y - screenHeight)
			}
		} else {
			avgY = avgY + neighbors[i].pos.Y
		}
	}
	avgX = avgX / float64(1+len(neighbors))
	avgY = avgY / float64(1+len(neighbors))
	desired := wrapTo(b.pos, pixel.V(avgX, avgY))

	// Usually the force should equal desired minus velocity, but in
	// this case it would suck all the life out of the simulation and
	// tend to become static.

	// return desired.Sub(b.vel)
	return desired
}

func (b *Boid) draw(imd *imdraw.IMDraw) {
	imd.Color = b.col
	imd.Push(b.pos)
	imd.Circle(5, 0)
	// imd.Color = b.colShade
	// imd.Push(b.pos)
	// imd.Circle(5, 1)
	// imd.Color = b.velCol
	// imd.Push(b.pos)
	// imd.Push(b.pos.Add(b.vel.Scaled(5)))
	// imd.Line(1)
	// imd.Color = b.accCol
	// imd.Push(b.pos)
	// imd.Push(b.pos.Add(b.acc.Scaled(35)))
	// imd.Line(1)
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func mapRange(in, inMin, inMax, outMin, outMax float64) float64 {
	return (in-inMin)/(inMax-inMin)*(outMax-outMin) + outMin
}
