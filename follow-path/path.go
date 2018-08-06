package main

import (
	"math"

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
		radius:      25,
	}
}

func (p *Path) drawClosest(v pixel.Vec, imd *imdraw.IMDraw) {
	// debug function, practice getting closest point to path.
	imd.Color = pixel.RGB(0, 0, 0)
	imd.Push(v)
	imd.Circle(5, 0)
	imd.Push(p.start)
	imd.Push(v)
	imd.Line(1)

	pt := closestPoint(v, p.start, p.end)
	imd.Color = pixel.RGB(1, 0, 0)
	imd.Push(pt)
	imd.Circle(5, 0)
	imd.Color = pixel.RGB(.5, 0, 0)
	imd.Push(pt)
	imd.Circle(5, 1)
}

func (p *Path) draw(imd *imdraw.IMDraw) {
	imd.Color = p.radiusColor
	imd.Push(p.start)
	imd.Push(p.end)
	imd.Line(p.radius * 2)
	imd.Color = p.lineColor
	imd.Push(p.start)
	imd.Push(p.end)
	imd.Line(2)
}

func angleBetween(a, b pixel.Vec) float64 {
	d := a.Dot(b)
	return math.Acos(d / (a.Len() * b.Len()))
}

func closestPoint(p, a, b pixel.Vec) pixel.Vec {
	ap := a.To(p)
	ab := a.To(b)
	ab = pixel.Unit(ab.Angle())
	ab = ab.Scaled(ap.Dot(ab))
	normalPoint := a.Add(ab)
	return normalPoint
}
