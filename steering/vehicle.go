package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Vehicle calculates it's own desired force to get to a target.
type Vehicle struct {
	pos      pixel.Vec
	acc      pixel.Vec
	vel      pixel.Vec
	tgt      *Target
	maxSpeed float64
	maxForce float64
}

// NewVehicle instantiates a new vehicle
func NewVehicle(tgt *Target) Vehicle {

	return Vehicle{
		pos:      pixel.V(200, screenHeight/2),
		acc:      pixel.V(0, 0),
		vel:      pixel.V(0, 3),
		tgt:      tgt,
		maxSpeed: 3,
		maxForce: .1,
	}
}

func (v *Vehicle) update() {
	v.seek(v.tgt.pos)
	v.pos = v.pos.Add(v.vel)
	v.vel = v.vel.Add(v.acc)
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
	// TODO: Use a sprite to draw this
	imd.Push(v.pos)
	imd.Circle(5, 0)
}
