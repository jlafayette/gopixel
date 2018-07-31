package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

func randomPos(bounds pixel.Rect) pixel.Vec {
	x := bounds.Min.X + rand.Float64()*(bounds.Max.X-bounds.Min.X)
	y := bounds.Min.Y + rand.Float64()*(bounds.Max.Y-bounds.Min.Y)
	return pixel.V(x, y)
}

func randomVel(max float64) pixel.Vec {
	x := -max + rand.Float64()*(max-(-max))
	y := -max + rand.Float64()*(max-(-max))
	return pixel.V(x, y)
}
