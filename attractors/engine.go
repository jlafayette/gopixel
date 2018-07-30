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
	particles []Particle
}

// NewEngine initializes a new physics engine
func NewEngine(particles []Particle) Engine {
	return Engine{
		particles: particles,
	}
}

func (e *Engine) update() {
	for i := 0; i < len(e.particles); i++ {
		totalAcceleration := pixel.V(0, 0)
		for j := 0; j < len(e.particles); j++ {
			if i != j {
				// Calculate force direction vector
				dir := e.particles[i].pos.Add(e.particles[i].vel).To(e.particles[j].pos)

				// Get distance squared (minimum distance is radius of attractor)
				distanceSq := math.Pow(math.Max(dir.Len(), e.particles[j].radius*1.5), 2)

				// Calcuate magnitude:
				//   F = G * (m1 * m2)/d2
				//   F = M * A  ->  A = F / M
				magnitude := G * ((e.particles[i].mass * e.particles[j].mass) / distanceSq)

				// Alternat formula
				// F = - G*M*m*r^(-2)
				// magnitude := G * e.particles[i].mass * e.particles[j].mass * math.Pow(dir.Len(), -2)

				acceleration := magnitude / e.particles[i].mass

				// e.particles[i].acc = pixel.Unit(dir.Angle()).Scaled(acceleration)
				totalAcceleration = totalAcceleration.Add(pixel.Unit(dir.Angle()).Scaled(acceleration))
			}
		}
		e.particles[i].acc = totalAcceleration
	}
	for i := 0; i < len(e.particles); i++ {
		e.particles[i].update()
	}
}

func (e *Engine) draw(imd *imdraw.IMDraw) {
	for i := 0; i < len(e.particles); i++ {
		e.particles[i].draw(imd)
	}
}
