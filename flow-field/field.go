package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
)

const (
	numX = 20
	numY = 20
)

// Field ...
type Field struct {
	cells [numX][numY]cell
	color pixel.RGBA
}

type cell struct {
	force     pixel.Vec
	drawStart pixel.Vec
	drawEnd   pixel.Vec
}

func newCell(x, y int) cell {
	pixelPerCellX := screenWidth / numX
	halfX := pixelPerCellX / 2
	pixelPerCellY := screenHeight / numY
	halfY := pixelPerCellY / 2
	shortestHalf := math.Min(float64(halfX), float64(halfY)) * .75

	force := pixel.Unit(randFloat(0, 2*math.Pi))
	center := pixel.V(float64(pixelPerCellX*x+halfX), float64(pixelPerCellY*y+halfY))
	start := center.Sub(force.Scaled(shortestHalf))
	end := center.Add(force.Scaled(shortestHalf))
	return cell{
		force:     force,
		drawStart: start,
		drawEnd:   end,
	}
}

// NewField ...
func NewField() Field {
	f := Field{
		cells: [numX][numY]cell{},
		color: pixel.RGB(.5, .5, .5),
	}
	f.randomizeFlow()
	return f
}

func (f *Field) randomizeFlow() {
	for x := 0; x < numX; x++ {
		for y := 0; y < numY; y++ {
			f.cells[x][y] = newCell(x, y)
		}
	}
}

func (f *Field) update() {
}

func (f *Field) draw(imd *imdraw.IMDraw) {
	imd.Color = f.color
	for x := 0; x < numX; x++ {
		for y := 0; y < numY; y++ {
			imd.Push(f.cells[x][y].drawStart)
			imd.Push(f.cells[x][y].drawEnd)
			imd.Line(1)
			imd.Push(f.cells[x][y].drawStart)
			imd.Circle(3, 0)
		}
	}
}

func (f *Field) lookup(pos pixel.Vec) pixel.Vec {
	pixelPerCellX := screenWidth / numX
	pixelPerCellY := screenHeight / numY
	x := pos.X / float64(pixelPerCellX)
	y := pos.Y / float64(pixelPerCellY)
	x = pixel.Clamp(x, 0, numX-1)
	y = pixel.Clamp(y, 0, numY-1)
	return f.cells[int(x)][int(y)].force
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
