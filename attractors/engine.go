package main

import (
	"math"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/imdraw"
)

// G is the gravitational constant
const G float64 = 6.67408

// Engine will keep track of all physics objects and calcuate forces
type Engine struct {
	attractors []Attractor
	particles  []Particle
}

// NewEngine initializes a new physics engine
func NewEngine(attractors []Attractor, particles []Particle) Engine {
	return Engine{
		attractors: attractors,
		particles:  particles,
	}
}

func (e *Engine) update() {
	for i := 0; i < len(e.particles); i++ {
		for j := 0; j < len(e.attractors); j++ {
			// Calculate force direction vector
			dir := e.particles[i].pos.To(e.attractors[j].pos)

			// Get distance squared
			distanceSq := math.Pow(dir.Len(), 2)

			// Calcuate magnitude:
			//   F = G * (m1 * m2)/d2
			//   F = M * A  ->  A = F / M
			magnitude := G * ((e.particles[i].mass * e.attractors[j].mass) / distanceSq)
			acceleration := magnitude / e.particles[i].mass

			e.particles[i].acc = pixel.Unit(dir.Angle()).Scaled(acceleration)
		}
	}

	for i := 0; i < len(e.attractors); i++ {
		e.attractors[i].update()
	}
	for i := 0; i < len(e.particles); i++ {
		e.particles[i].update()
	}
}

func (e *Engine) draw(imd *imdraw.IMDraw) {
	for i := 0; i < len(e.attractors); i++ {
		e.attractors[i].draw(imd)
	}
	for i := 0; i < len(e.particles); i++ {
		e.particles[i].draw(imd)
	}
}
