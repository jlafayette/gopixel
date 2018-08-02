package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
)

func basic() []Particle {
	var particles []Particle
	a1 := NewParticle(pixel.V(screenWidth/2, screenHeight/2), pixel.V(0, 0), randFloat(3000, 7000))
	a1.color = pixel.RGB(1, 1, 1)
	a1.moveable = false
	particles = append(particles, a1)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	for i := 0; i < rand.Intn(9)+1; i++ {
		m := randFloat(.01, 1.5)
		p := NewOrbiter(a1, m, screenSafeDistance(150), .5, 1.05)
		p.color.A = .25
		particles = append(particles, p)
	}
	return particles
}

func basic2() []Particle {
	var particles []Particle
	a1 := NewParticle(pixel.V(screenWidth/2, screenHeight/2), pixel.V(0, 0), randFloat(2000, 8000))
	a1.color = pixel.RGB(1, 1, 1)
	a1.moveable = false
	particles = append(particles, a1)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	for i := 0; i < rand.Intn(9)+2; i++ {
		m := randFloat(.05, 100)
		p := NewOrbiter(a1, m, screenSafeDistance(100), .5, 1.05)
		p.color.A = .25
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
	a1.color = pixel.RGB(1, 1, 1)
	particles = append(particles, a1)

	g1 := NewOrbiter(a1, 500, screenSafeDistance(50), .8, 1.0)
	g1.color.A = .001
	particles = append(particles, g1)

	for i := 0; i < rand.Intn(5)+1; i++ {
		m := .001 + rand.Float64()*(.05-.001)
		p := NewOrbiter(g1, m, 10, .95, 1.05)
		p.color.A = .35
		particles = append(particles, p)
	}

	for i := 0; i < rand.Intn(3)+1; i++ {
		m := .1 + rand.Float64()*(10-.1)
		p := NewOrbiter(a1, m, screenSafeDistance(150), .7, 1.1)
		p.color.A = .25
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
		pos := randomPos(pixel.R(50, 50, screenWidth-50, screenHeight-50))
		m := 500 + rand.Float64()*(1000-500)
		a1 := NewParticle(pos, pixel.V(0, 0), m)
		a1.moveable = false
		a1.visible = false
		particles = append(particles, a1)
		anchors = append(anchors, a1)
	}
	for i := 0; i < rand.Intn(100)+10; i++ {
		m := .01 + rand.Float64()*(100-.01)
		p := NewOrbiter(anchors[rand.Intn(len(anchors))], m, randFloat(0, 400), -1.5, 1.5)
		p.color = similarRandomColor(colorSeed)
		p.color.A = .1
		particles = append(particles, p)
	}
	return particles
}

func gravityPaths() []Particle {
	var anchors []Particle
	var particles []Particle
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	colorSeed := rand.Intn(3)

	for i := 0; i < rand.Intn(15)+5; i++ {
		pos := randomPos(pixel.R(50, 50, screenWidth-50, screenHeight-50))
		m := 500 + rand.Float64()*(1000-500)
		a1 := NewParticle(pos, pixel.V(0, 0), m)
		a1.moveable = false
		a1.visible = false
		particles = append(particles, a1)
		anchors = append(anchors, a1)
	}

	for i := 0; i < rand.Intn(1)+500; i++ {
		// position
		x := randFloat(-2000, 0)
		y := randFloat(0, screenHeight)
		pos := pixel.V(x, y)

		// velocity
		vel := pixel.V(randFloat(10, 30), 0)

		// mass
		m := randFloat(.01, 10)

		p := NewParticle(pos, vel, m)
		p.color = similarRandomColor(colorSeed)
		p.color.A = .1
		particles = append(particles, p)
	}
	return particles
}

func screenSafeDistance(buffer float64) float64 {
	// assumes that particle to orbit is at the center of the screen.
	return math.Min(screenWidth, screenHeight)/2 - buffer
}
