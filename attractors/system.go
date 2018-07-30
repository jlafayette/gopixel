package main

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

func basic() []Particle {
	var particles []Particle
	a1 := NewParticle(pixel.V(screenWidth/2, screenHeight/2), pixel.V(0, 0), 5000)
	a1.color = colornames.White
	particles = append(particles, a1)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	for i := 0; i < rand.Intn(9)+1; i++ {
		m := .01 + rand.Float64()*(1.5-.01)
		p := NewOrbiter(a1, m, pixel.R(200, 200, screenWidth-200, screenHeight-200), .5, 1.05)
		particles = append(particles, p)
	}
	return particles
}

func gasGiant() []Particle {
	var particles []Particle
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	a1 := NewParticle(pixel.V(screenWidth/2, screenHeight/2), pixel.V(0, 0), 15000)
	a1.moveable = false
	a1.color = colornames.White
	particles = append(particles, a1)

	g1 := NewOrbiter(a1, 500, pixel.R(300, 200, screenWidth-300, screenHeight-200), .8, 1.0)
	particles = append(particles, g1)

	// a1.vel = g1.vel

	for i := 0; i < rand.Intn(5)+1; i++ {
		m := .001 + rand.Float64()*(.05-.001)
		p := NewOrbiter(g1, m, pixel.R(g1.pos.X+g1.radius+2, g1.pos.Y+g1.radius+2, g1.pos.X+g1.radius+10, g1.pos.Y+g1.radius+10), .95, 1.05)
		particles = append(particles, p)
	}

	for i := 0; i < rand.Intn(3)+1; i++ {
		m := .1 + rand.Float64()*(10-.1)
		p := NewOrbiter(a1, m, pixel.R(200, 200, screenWidth-200, screenHeight-200), .7, 1.1)
		particles = append(particles, p)
	}

	return particles
}

func random() []Particle {
	var particles []Particle
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	colorSeed := rand.Intn(3)

	var anchors []Particle
	for i := 0; i < rand.Intn(15)+5; i++ {
		pos := randomPos(pixel.R(100, 100, screenWidth-100, screenHeight-100))
		m := 500 + rand.Float64()*(1000-500)
		a1 := NewParticle(pos, pixel.V(0, 0), m)
		a1.moveable = false
		a1.visible = false
		a1.color = colornames.Black
		particles = append(particles, a1)
		anchors = append(anchors, a1)
	}
	for i := 0; i < rand.Intn(100)+10; i++ {
		m := .01 + rand.Float64()*(100-.01)
		p := NewOrbiter(anchors[rand.Intn(len(anchors))], m, pixel.R(250, 200, screenWidth-250, screenHeight-200), -1.5, 1.5)
		p.color = similarRandomColor(colorSeed)
		// pos := randomPos(pixel.R(450, 400, screenWidth-450, screenHeight-400))
		// vel := randomVel(1)
		// p := NewParticle(pos, vel, m)
		particles = append(particles, p)
	}
	return particles
}
