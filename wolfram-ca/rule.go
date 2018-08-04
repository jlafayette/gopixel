package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/faiface/pixel"
)

// Rule ....
type Rule struct {
}

// NewRule ...
func NewRule() Rule {
	return Rule{}
}

func (r *Rule) update() {
}

func (r *Rule) draw(t pixel.Target) {
	img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			img.Set(x, y, color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
		}
	}
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(t, pixel.IM.Moved(pic.Bounds().Center()))
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
