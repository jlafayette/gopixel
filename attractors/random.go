package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

func randomPos(bounds pixel.Rect) pixel.Vec {
	x := randFloat(bounds.Min.X, bounds.Max.X)
	y := randFloat(bounds.Min.Y, bounds.Max.Y)
	return pixel.V(x, y)
}

func randomVel(max float64) pixel.Vec {
	x := randFloat(-max, max)
	y := randFloat(-max, max)
	return pixel.V(x, y)
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
