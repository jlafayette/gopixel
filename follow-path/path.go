package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Path ...
type Path struct {
	points      [8]pixel.Vec
	lineColor   pixel.RGBA
	radiusColor pixel.RGBA
	radius      float64
}

// NewPath ...
func NewPath() Path {
	return Path{
		points:      randomPoints(),
		lineColor:   pixel.RGB(.5, .5, .5),
		radiusColor: pixel.RGB(.8, .8, .8),
		radius:      25,
	}
}

func randomPoints() [8]pixel.Vec {
	return [8]pixel.Vec{
		pixel.V(0, randFloat(100, screenHeight-100)),
		pixel.V(1*(screenWidth/7), randFloat(100, screenHeight-100)),
		pixel.V(2*(screenWidth/7), randFloat(100, screenHeight-100)),
		pixel.V(3*(screenWidth/7), randFloat(100, screenHeight-100)),
		pixel.V(4*(screenWidth/7), randFloat(100, screenHeight-100)),
		pixel.V(5*(screenWidth/7), randFloat(100, screenHeight-100)),
		pixel.V(6*(screenWidth/7), randFloat(100, screenHeight-100)),
		pixel.V(screenWidth, randFloat(100, screenHeight-100)),
	}
}

func (p *Path) drawClosest(v pixel.Vec, imd *imdraw.IMDraw) {
	// debug function, practice getting closest point to path.
	imd.Color = pixel.RGB(0, 0, 0)
	imd.Push(v)
	imd.Circle(5, 0)

	closestPt := pixel.V(-1000, -1000)
	start := p.points[0]
	// end := p.points[1]
	for i := 0; i < len(p.points)-1; i++ {
		pt := closestPoint(v, p.points[i], p.points[i+1])
		segLen := p.points[i].To(p.points[i+1]).Len()
		if p.points[i].To(pt).Len() < segLen && p.points[i+1].To(pt).Len() < segLen {
			if v.To(pt).Len() < v.To(closestPt).Len() {
				closestPt = pt
				start = p.points[i]
				// end = p.points[i+1]
			}
		}
	}
	for i := 0; i < len(p.points); i++ {
		if v.To(p.points[i]).Len() < v.To(closestPt).Len() {
			closestPt = p.points[i]
			start = p.points[i]
		}
	}
	imd.Push(start)
	imd.Push(v)
	imd.Line(1)

	imd.Color = pixel.RGB(1, 0, 0)
	imd.Push(closestPt)
	imd.Circle(5, 0)
	imd.Color = pixel.RGB(.5, 0, 0)
	imd.Push(closestPt)
	imd.Circle(5, 1)
}

func (p *Path) draw(imd *imdraw.IMDraw) {
	imd.Color = p.radiusColor
	for i := 0; i < len(p.points)-1; i++ {
		imd.Push(p.points[i])
		imd.Push(p.points[i+1])
		imd.Line(p.radius * 2)
	}

	imd.Color = p.lineColor
	for i := 0; i < len(p.points)-1; i++ {
		imd.Push(p.points[i])
		imd.Push(p.points[i+1])
		imd.Line(2)
	}
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
