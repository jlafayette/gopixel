package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Walker will wander around the screen.
type Walker struct {
	pos   pixel.Vec
	steps []pixel.Vec
}

// NewWalker creates a new Walker at the middle of the screen.
func NewWalker(stepsPerFrame int) Walker {
	return Walker{
		pos:   pixel.V(screenWidth/2, screenHeight/2),
		steps: make([]pixel.Vec, stepsPerFrame, stepsPerFrame),
	}
}

func (w *Walker) update() {
	for i := 0; i < len(w.steps); i++ {
		move := pixel.Unit(randFloat(0, 2*math.Pi)).Scaled(1)
		w.pos = w.pos.Add(move)
		w.steps[i] = w.pos
	}
}

func (w *Walker) draw(imd *imdraw.IMDraw) {
	for i := 0; i < len(w.steps); i++ {
		imd.Push(w.steps[i])
	}
	imd.Circle(1, 0)
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
