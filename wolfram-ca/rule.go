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
func NewRule() Rule {
	return Rule{
		on:      color.RGBA{255, 255, 255, 255},
		off:     color.RGBA{0, 0, 0, 255},
		rule:    translateRule(30),
		row:     make([]bool, screenWidth, screenWidth),
		prevRow: make([]bool, screenWidth, screenWidth),
	}
}

func (r *Rule) update() {
}

func (r *Rule) draw(t pixel.Target) {
	img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	// First row has one on pixel in the middle
	for x := 0; x < screenWidth; x++ {
		if x == screenWidth/2 {
			img.Set(x, 0, r.on)
			r.row[x] = true
		} else {
			img.Set(x, 0, r.off)
			r.row[x] = false
		}
	}
	for i := 0; i < len(r.row); i++ {
		r.prevRow[i] = r.row[i]
	}
	// r.prevRow = r.row

	for y := 1; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {

			// determine r.row[x]
			me := r.prevRow[x]
			var left, right bool
			if x == 0 {
				left = r.prevRow[screenWidth-1]
			} else {
				left = r.prevRow[x-1]
			}
			if x == screenWidth-1 {
				right = r.prevRow[0]
			} else {
				right = r.prevRow[x+1]
			}
			// fmt.Printf("%d %d%d%d\n", x, bToI(left), bToI(me), bToI(right))
			// r.row[x] = r.whatAmI(left, me, right)
			// fmt.Printf("   %d \n", bToI(r.row[x]))
			n := 0
			if left {
				n = n + 4
			}
			if me {
				n = n + 2
			}
			if right {
				n = n + 1
			}
			r.row[x] = applyRule(byte(n), r.rule)

			if r.row[x] {
				img.Set(x, y, r.on)
			} else {
				img.Set(x, y, r.off)
			}
		}

		for i := 0; i < len(r.row); i++ {
			r.prevRow[i] = r.row[i]
		}

	}
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(t, pixel.IM.Moved(pic.Bounds().Center()))
}

func bToI(b bool) int {
	if b {
		return 1
	}
	return 0
}

func applyRule(prev byte, rule [8]bool) bool {
	return rule[uint8(prev)]
	// for i := 0; i < 8; i++ {
	// 	fmt.Printf("%03b %v\n", byte(i), r[i])
	// }
}

func translateRule(rule uint8) [8]bool {
	//fmt.Println(rule)
	var r [8]bool
	for i := uint8(0); i < 8; i++ {
		n := 1 << i
		//fmt.Printf("%d -> %d", i, 1<<i)
		//fmt.Printf("r: %08b\n", byte(rule)&byte(n))
		r[i] = uint8(byte(rule)&byte(n)) > 0
	}
	return r
}
