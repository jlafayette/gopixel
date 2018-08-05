package main

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
)

// Rule ....
type Rule struct {
	on      color.RGBA
	off     color.RGBA
	rule    [8]bool
	row     []bool
	prevRow []bool
}

// NewRule ...
func NewRule(rule uint8) Rule {
	return Rule{
		on:      color.RGBA{255, 255, 255, 255},
		off:     color.RGBA{0, 0, 0, 255},
		rule:    translateRule(rule),
		row:     make([]bool, screenWidth, screenWidth),
		prevRow: make([]bool, screenWidth, screenWidth),
	}
}

func (r *Rule) update() {
}

func (r *Rule) draw(t pixel.Target) {
	img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	// first row has one on pixel in the middle
	for x := 0; x < screenWidth; x++ {
		if x == screenWidth/2 {
			img.Set(x, 0, r.on)
			r.row[x] = true
		} else {
			img.Set(x, 0, r.off)
			r.row[x] = false
		}
	}

	// swap row with previous
	for i := 0; i < len(r.row); i++ {
		r.prevRow[i] = r.row[i]
	}

	for y := 1; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {

			// determine r.row[x]
			var neighbors byte // arrangement in previous row: left, middle, right
			// middle
			if r.prevRow[x] {
				neighbors = neighbors | byte(2)
			}
			// left
			if x == 0 {
				if r.prevRow[screenWidth-1] {
					neighbors = neighbors | byte(4)
				}
			} else if r.prevRow[x-1] {
				neighbors = neighbors | byte(4)
			}
			// right
			if x == screenWidth-1 {
				if r.prevRow[0] {
					neighbors = neighbors | byte(1)
				}
			} else if r.prevRow[x+1] {
				neighbors = neighbors | byte(1)
			}
			r.row[x] = applyRule(neighbors, r.rule)

			// set the current pixel to on or off
			if r.row[x] {
				img.Set(x, y, r.on)
			} else {
				img.Set(x, y, r.off)
			}
		}

		// swap row with previous
		for i := 0; i < len(r.row); i++ {
			r.prevRow[i] = r.row[i]
		}

	}
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(t, pixel.IM.Moved(pic.Bounds().Center()))
}

func applyRule(prev byte, rule [8]bool) bool {
	return rule[uint8(prev)]
}

func translateRule(rule uint8) [8]bool {
	var r [8]bool
	for i := uint8(0); i < 8; i++ {
		n := 1 << i
		r[i] = uint8(byte(rule)&byte(n)) > 0
	}
	return r
}
