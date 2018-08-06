package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Path ...
type Path struct {
	start       pixel.Vec
	end         pixel.Vec
	lineColor   pixel.RGBA
	radiusColor pixel.RGBA
	radius      float64
}

// NewPath ...
func NewPath() Path {
	return Path{
		start:       pixel.V(0, randFloat(100, screenHeight-100)),
		end:         pixel.V(screenWidth, randFloat(100, screenHeight-100)),
		lineColor:   pixel.RGB(0, 0, 0),
		radiusColor: pixel.RGB(.5, .5, .5),
		radius:      50,
	}
}

func (p *Path) update() {
}

func (p *Path) draw(imd *imdraw.IMDraw) {
	imd.Color = p.radiusColor
	imd.Push(p.start)
	imd.Push(p.end)
	imd.Line(p.radius)
	imd.Color = p.lineColor
	imd.Push(p.start)
	imd.Push(p.end)
	imd.Line(2)
}
