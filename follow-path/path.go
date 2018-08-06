package main

import (
	"github.com/faiface/pixel"
)

// Path ...
type Path struct {
	start pixel.Vec
	end   pixel.Vec
	r     float64
}

// NewPath ...
func NewPath() Path {
	return Path{}
}
