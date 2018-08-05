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
	scale   int
	width   int
	height  int
}

// NewRule ...
func NewRule(rule uint8, scale int) Rule {
	width := screenWidth / scale
	height := screenHeight / scale
	return Rule{
		on:      color.RGBA{200, 200, 200, 255},
		off:     color.RGBA{100, 100, 100, 255},
		rule:    translateRule(rule),
		row:     make([]bool, width, width),
		prevRow: make([]bool, width, width),
		scale:   scale,
		width:   width,
		height:  height,
	}
}

func (r *Rule) update() {
}

func (r *Rule) draw(t pixel.Target) {
	img := image.NewRGBA(image.Rect(0, 0, r.width, r.height))

	// first row has one on pixel in the middle
	for x := 0; x < r.width; x++ {
		if x == r.width/2 {
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

	for y := 1; y < r.height; y++ {
		for x := 0; x < r.width; x++ {

			// determine r.row[x]
			var neighbors byte // arrangement in previous row: left, middle, right
			// middle
			if r.prevRow[x] {
				neighbors = neighbors | byte(2)
			}
			// left
			if x == 0 {
				if r.prevRow[r.width-1] {
					neighbors = neighbors | byte(4)
				}
			} else if r.prevRow[x-1] {
				neighbors = neighbors | byte(4)
			}
			// right
			if x == r.width-1 {
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
	m := pixel.IM
	m = m.Moved(pic.Bounds().Center())
	m = m.ScaledXY(pixel.ZV, pixel.V(float64(r.scale), float64(r.scale)))

	sprite.Draw(t, m)
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
