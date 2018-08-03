package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

// Walker will wander around the screen.
type Walker struct {
	pos    pixel.Vec
	steps  []pixel.Vec
	sprite *pixel.Sprite
	mask   pixel.RGBA
}

// NewWalker creates a new Walker at the middle of the screen.
func NewWalker(stepsPerFrame int, sprite *pixel.Sprite) Walker {
	return Walker{
		pos:    pixel.V(screenWidth/2, screenHeight/2),
		steps:  make([]pixel.Vec, stepsPerFrame, stepsPerFrame),
		sprite: sprite,
		mask:   pixel.Alpha(.25),
	}
}

func (w *Walker) update() {
	for i := 0; i < len(w.steps); i++ {
		move := pixel.Unit(randFloat(0, 2*math.Pi)).Scaled(randFloat(1, 7))
		pos := w.pos.Add(move)
		if pos.X > 0 && pos.Y > 0 && pos.X < screenWidth && pos.Y < screenHeight {
			w.pos = pos
			w.steps[i] = pos
		}
	}
}

func (w *Walker) draw(batch *pixel.Batch) {
	for i := 0; i < len(w.steps); i++ {
		mat := pixel.IM
		mat = mat.Moved(w.steps[i])
		w.sprite.DrawColorMask(batch, mat, w.mask)
	}
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
