package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Vehicle ...
type Vehicle struct {
	pos           pixel.Vec
	acc           pixel.Vec
	vel           pixel.Vec
	futurePos     pixel.Vec
	cp            pixel.Vec
	desiredV      pixel.Vec
	path          *Path
	col           pixel.RGBA
	colShade      pixel.RGBA
	velCol        pixel.RGBA
	accCol        pixel.RGBA
	maxSpeed      float64
	maxForce      float64
	futurePosMult float64
	pathLookAhead float64
}

// NewVehicle ...
func NewVehicle(pos pixel.Vec, path *Path) Vehicle {
	maxSpeed := randFloat(3, 5)
	// maxForce := randFloat(.15, .3)
	maxForce := randFloat(.2, .3)

	futurePosMult := mapRange(maxForce, .2, .3, 35, 15) // should be inverse of maxForce

	return Vehicle{
		pos:           pos,
		acc:           pixel.ZV,
		vel:           pixel.ZV,
		path:          path,
		col:           pixel.RGB(0, .8, 0),
		colShade:      pixel.RGB(0, .2, 0),
		velCol:        pixel.RGB(1, 0, 0),
		accCol:        pixel.RGB(0, 0, 1),
		maxSpeed:      maxSpeed,
		maxForce:      maxForce,
		futurePosMult: futurePosMult,
		pathLookAhead: 50,
	}
}

func (v *Vehicle) update() {
	v.desiredV = v.followPath()
	v.seek(v.desiredV)
	v.pos = v.pos.Add(v.vel)
	v.vel = v.vel.Add(v.acc)
}

func (v *Vehicle) followPath() pixel.Vec {
	v.futurePos = v.pos.Add(v.vel.Scaled(v.futurePosMult))

	// is future pos on the path?
	cp := pixel.V(-1000, -1000)
	start := v.path.points[0]
	seg := v.path.points[0].To(v.path.points[1])
	// end := p.points[1]
	for i := 0; i < len(v.path.points)-1; i++ {
		pt := closestPoint(v.futurePos, v.path.points[i], v.path.points[i+1])
		segLen := v.path.points[i].To(v.path.points[i+1]).Len()
		if v.path.points[i].To(pt).Len() < segLen && v.path.points[i+1].To(pt).Len() < segLen {
			if v.futurePos.To(pt).Len() < v.futurePos.To(cp).Len() {
				cp = pt
				start = v.path.points[i]
				seg = start.To(v.path.points[i+1])
				// end = p.points[i+1]
			}
		}
	}
	for i := 0; i < len(v.path.points); i++ {
		if v.futurePos.To(v.path.points[i]).Len() < v.futurePos.To(cp).Len() {
			cp = v.path.points[i]
			start = v.path.points[i]
			if i != len(v.path.points)-1 {
				seg = start.To(v.path.points[i+1])
			} else {
				seg = v.path.points[i-1].To(start)
			}
		}
	}
	// v.cp = closestPoint(v.futurePos, v.path.start, v.path.end)
	v.cp = cp

	if v.cp.To(v.futurePos).Len() < v.path.radius {
		// if yes, do nothing
		return v.futurePos
	}
	// if no, move along the path a bit
	alongPath := pixel.Unit(seg.Angle()).Scaled(v.pathLookAhead)
	return v.cp.Add(alongPath)
}

func (v *Vehicle) seek(tgt pixel.Vec) {
	toTgt := v.pos.To(tgt)
	desired := pixel.Unit(toTgt.Angle()).Scaled(v.maxSpeed)
	steering := desired.Sub(v.vel)
	if steering.Len() > v.maxForce {
		steering = pixel.Unit(steering.Angle()).Scaled(v.maxForce)
	}
	v.acc = steering
}

func (v *Vehicle) draw(imd *imdraw.IMDraw) {
	imd.Color = v.col
	imd.Push(v.pos)
	imd.Circle(5, 0)
	// imd.Push(v.futurePos)
	// imd.Circle(3, 0)
	imd.Color = v.colShade
	imd.Push(v.pos)
	imd.Circle(5, 1)
	imd.Color = v.velCol
	imd.Push(v.pos)
	imd.Push(v.pos.Add(v.vel.Scaled(5)))
	imd.Line(1)
	// imd.Push(v.cp)
	// imd.Circle(3, 0)
	imd.Color = v.accCol
	imd.Push(v.pos)
	imd.Push(v.pos.Add(v.acc.Scaled(35)))
	imd.Line(1)
	// imd.Push(v.desiredV)
	// imd.Circle(3, 0)
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func mapRange(in, inMin, inMax, outMin, outMax float64) float64 {
	return (in-inMin)/(inMax-inMin)*(outMax-outMin) + outMin
}
