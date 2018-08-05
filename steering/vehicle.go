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
	maxSpeed float64
	maxForce float64
}

// NewVehicle instantiates a new vehicle
func NewVehicle() Vehicle {

	return Vehicle{
		pos:      pixel.V(200, screenHeight/2),
		acc:      pixel.V(0, 0),
		vel:      pixel.V(1, 2),
		maxSpeed: 3,
		maxForce: 1,
	}
}

func (v *Vehicle) update() {
	v.pos = v.pos.Add(v.vel)
	v.vel = v.vel.Add(v.acc)
}

func (v *Vehicle) draw(imd *imdraw.IMDraw) {
	// TODO: Use a sprite to draw this
	imd.Push(v.pos)
	imd.Circle(5, 0)
}
