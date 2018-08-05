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
	field    *Field
	col      pixel.RGBA
	colShade pixel.RGBA
	velCol   pixel.RGBA
	accCol   pixel.RGBA
	maxSpeed float64
	maxForce float64
}

// NewVehicle instantiates a new vehicle
func NewVehicle(pos pixel.Vec, field *Field) Vehicle {

	return Vehicle{
		pos:      pos,
		acc:      pixel.V(0, 0),
		vel:      pixel.V(0, 3),
		field:    field,
		col:      pixel.RGB(0, .8, 0),
		colShade: pixel.RGB(0, .2, 0),
		velCol:   pixel.RGB(1, 0, 0),
		accCol:   pixel.RGB(0, 0, 1),
		maxSpeed: 3,
		maxForce: .1,
	}
}

func (v *Vehicle) update() {
	tgt := v.pos.Add(v.field.lookup(v.pos))
	v.seek(tgt)
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
