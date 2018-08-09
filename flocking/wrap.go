package main

import (
	"math"

	"github.com/faiface/pixel"
)

func wrapSpace(a, b pixel.Vec) (pixel.Vec, pixel.Vec) {
	screen := pixel.R(0, 0, screenWidth, screenHeight)
	halfX := screen.Size().X * .5
	halfY := screen.Size().Y * .5
	if math.Abs(b.X-a.X) > halfX {
		if a.X < b.X {
			a = a.Add(pixel.V(screenWidth, 0))
		} else {
			b = b.Add(pixel.V(screenWidth, 0))
		}
	}
	if math.Abs(b.Y-a.Y) > halfY {
		if a.Y < b.Y {
			a = a.Add(pixel.V(0, screenHeight))
		} else {
			b = b.Add(pixel.V(0, screenHeight))
		}
	}
	return a, b
}

// wraps around screen edges
func distance(a, b pixel.Vec) float64 {
	a, b = wrapSpace(a, b)
	return a.To(b).Len()
}

// wraps around screen edges
func wrapTo(a, b pixel.Vec) pixel.Vec {
	a, b = wrapSpace(a, b)
	return a.To(b)
}
