package main

import (
	"image/color"
	"math/rand"
)

func dark() color.NRGBA {
	return color.NRGBA{10, 10, 10, 255}
}

func red() color.NRGBA {
	return color.NRGBA{255, 50, 50, 255}
}

// randomColor generates a random color
func randomColor(lo, hi int) color.NRGBA {
	return color.NRGBA{
		uint8(rand.Intn(hi) - lo),
		uint8(rand.Intn(hi) - lo),
		uint8(rand.Intn(hi) - lo),
		255,
	}
}
