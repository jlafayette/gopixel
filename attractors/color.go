package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

// randomColor generates a random color
func randomColor() pixel.RGBA {
	return pixel.RGB(
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
	)
}

// similar colors
func similarRandomColor(i int) pixel.RGBA {
	switch i {
	case 0:
		return pixel.RGB(
			randFloat(.33, 1),
			.1,
			rand.Float64(),
		)
	case 1:
		return pixel.RGB(
			.1,
			randFloat(.33, 1),
			randFloat(0, .5),
		)
	case 2:
		return pixel.RGB(
			rand.Float64(),
			randFloat(0, .33),
			randFloat(.33, 1),
		)
	}
	return randomColor()
}
