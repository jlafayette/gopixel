package main

import (
	"math"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/imdraw"
)

// G is the gravitational constant
const G float64 = 6.674

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
		totalAcceleration := pixel.V(0, 0)
		for j := 0; j < len(e.attractors); j++ {
			// Calculate force direction vector
			dir := e.particles[i].pos.Add(e.particles[i].vel).To(e.attractors[j].pos)

			// Get distance squared (minimum distance is radius of attractor)
			distanceSq := math.Pow(math.Max(dir.Len(), e.attractors[j].radius*2), 2)

			// Calcuate magnitude:
			//   F = G * (m1 * m2)/d2
			//   F = M * A  ->  A = F / M
			magnitude := G * ((e.particles[i].mass * e.attractors[j].mass) / distanceSq)

			// Alternat formula
			// F = - G*M*m*r^(-2)
			// magnitude := G * e.particles[i].mass * e.attractors[j].mass * math.Pow(dir.Len(), -2)

			acceleration := magnitude / e.particles[i].mass

			// e.particles[i].acc = pixel.Unit(dir.Angle()).Scaled(acceleration)
			totalAcceleration = totalAcceleration.Add(pixel.Unit(dir.Angle()).Scaled(acceleration))
		}
		for k := 0; k < len(e.particles); k++ {
			if i != k {
				dir := e.particles[i].pos.Add(e.particles[i].vel).To(e.particles[k].pos)
				distance := math.Max(dir.Len(), e.particles[i].radius*2+e.particles[k].radius*2)
				distanceSq := math.Pow(distance, 2)
				magnitude := G * ((e.particles[i].mass * e.particles[k].mass) / distanceSq)
				acceleration := magnitude / e.particles[i].mass
				totalAcceleration = totalAcceleration.Add(pixel.Unit(dir.Angle()).Scaled(acceleration))
			}
		}
		e.particles[i].acc = totalAcceleration
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
