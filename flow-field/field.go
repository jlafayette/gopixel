package main

import (
	"math/rand"

	"github.com/faiface/pixel"
)

const (
	numX = 40
	numY = 40
)

// Field ...
type Field struct {
	cells [numX][numY]pixel.Vec
	color pixel.RGBA
}

// NewField ...
func NewField() Field {
	f := Field{
		cells: [numX][numY]pixel.Vec{},
		color: pixel.RGB(.5, .5, .5),
	}
	f.randomizeFlow()
	return f
}

func (f *Field) randomizeFlow() {
	// visited := [numX][numY]bool{}
}

func (f *Field) update() {
}

func (f *Field) draw() {

}

func (f *Field) lookup(pos pixel.Vec) pixel.Vec {
	return pixel.Unit(0)
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
