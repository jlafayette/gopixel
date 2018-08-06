package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Vehicle ...
type Vehicle struct {
	pos      pixel.Vec
	acc      pixel.Vec
	vel      pixel.Vec
	path     *Path
	col      pixel.RGBA
	colShade pixel.RGBA
	velCol   pixel.RGBA
	accCol   pixel.RGBA
	maxSpeed float64
	maxForce float64
}

// NewVehicle ...
func NewVehicle(pos pixel.Vec, path *Path) Vehicle {
	return Vehicle{
		pos:      pos,
		acc:      pixel.ZV,
		vel:      pixel.ZV,
		path:     path,
		col:      pixel.RGB(0, .8, 0),
		colShade: pixel.RGB(0, .2, 0),
		velCol:   pixel.RGB(1, 0, 0),
		accCol:   pixel.RGB(0, 0, 1),
		maxSpeed: randFloat(3, 5),
		maxForce: randFloat(.15, .3),
	}
}

func (v *Vehicle) update() {
	tgt := v.followPath()
	v.seek(tgt)
	v.pos = v.pos.Add(v.vel)
	v.vel = v.vel.Add(v.acc)
}

func (v *Vehicle) followPath() pixel.Vec {
	return pixel.ZV
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
	imd.Push(v.pos.Add(v.acc.Scaled(35)))
	imd.Line(1)
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
