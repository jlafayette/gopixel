package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Vehicle calculates it's own desired force to get to a target.
type Vehicle struct {
	pos          pixel.Vec
	acc          pixel.Vec
	vel          pixel.Vec
	tgt          *Target
	col          pixel.RGBA
	colShade     pixel.RGBA
	velCol       pixel.RGBA
	accCol       pixel.RGBA
	maxSpeed     float64
	maxForce     float64
	arriveRadius float64
}

// NewVehicle instantiates a new vehicle
func NewVehicle(tgt *Target) Vehicle {

	return Vehicle{
		pos:          pixel.V(200, screenHeight/2),
		acc:          pixel.V(0, 0),
		vel:          pixel.V(0, 3),
		tgt:          tgt,
		col:          pixel.RGB(0, .8, 0),
		colShade:     pixel.RGB(0, .2, 0),
		velCol:       pixel.RGB(1, 0, 0),
		accCol:       pixel.RGB(0, 0, 1),
		maxSpeed:     3,
		maxForce:     .1,
		arriveRadius: 50,
	}
}

func (v *Vehicle) update() {
	v.seek(v.tgt.pos)
	v.pos = v.pos.Add(v.vel)
	v.vel = v.vel.Add(v.acc)
}

func (v *Vehicle) seek(tgt pixel.Vec) {
	toTgt := v.pos.To(tgt)
	var mag float64
	if toTgt.Len() < v.arriveRadius { // arrive without overshooting
		mag = mapRange(toTgt.Len(), 0, v.arriveRadius, 0, v.maxSpeed)
	} else {
		mag = v.maxSpeed
	}
	desired := pixel.Unit(toTgt.Angle()).Scaled(mag)
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
	imd.Color = v.colShade
	imd.Push(v.pos)
	imd.Circle(5, 1)
	imd.Color = v.velCol
	imd.Push(v.pos)
	imd.Push(v.pos.Add(v.vel.Scaled(5)))
	imd.Line(1)
	imd.Color = v.accCol
	imd.Push(v.pos)
	imd.Push(v.pos.Add(v.acc.Scaled(75)))
	imd.Line(1)
}

func mapRange(in, inMin, inMax, outMin, outMax float64) float64 {
	// If in X falls between A and B, and you want output Y to
	// fall between C and D, apply the following linear transform:
	// Y = (X-A)/(B-A) * (D-C) + C
	return (in-inMin)/(inMax-inMin)*(outMax-outMin) + outMin
}
