package main

import (
	"image/color"
	"math/rand"
)

// randomColor generates a random color
func randomColor() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
}

// similar colors
func similarRandomColor(i int) color.RGBA {
	switch i {
	case 0:
		return color.RGBA{
			uint8(rand.Intn(156) + 100),
			uint8(10),
			uint8(rand.Intn(256)),
			255,
		}
	case 1:
		return color.RGBA{
			uint8(10),
			uint8(rand.Intn(156) + 100),
			uint8(rand.Intn(175)),
			255,
		}
	case 2:
		return color.RGBA{
			uint8(rand.Intn(256)),
			uint8(rand.Intn(100)),
			uint8(rand.Intn(156) + 100),
			255,
		}
	}
	return randomColor()
}
