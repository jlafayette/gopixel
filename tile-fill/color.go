package main

import (
	"image/color"
	"math/rand"
)

// randomColor generates a random color
func randomColor(lo, hi int) color.NRGBA {
	return color.NRGBA{
		uint8(rand.Intn(hi) - lo),
		uint8(rand.Intn(hi) - lo),
		uint8(rand.Intn(hi) - lo),
		255,
	}
}
