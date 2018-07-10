package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/faiface/pixel"
)

func BenchmarkCells_generateVoronoi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Seed(time.Now().Unix())
		c := NewCells(25, pixel.R(0, 0, 400, 400))
		c.generateVoronoi()
	}
}
