package main

import "github.com/faiface/pixel/imdraw"

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
