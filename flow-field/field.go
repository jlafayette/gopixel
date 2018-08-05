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
	force := pixel.Unit(randFloat(0, 2*math.Pi))
	return cell{
		force: force,
	}
}

func (c *cell) calculateDrawPoints(x, y int) {
	pixelPerCellX := screenWidth / numX
	halfX := pixelPerCellX / 2
	pixelPerCellY := screenHeight / numY
	halfY := pixelPerCellY / 2
	shortestHalf := math.Min(float64(halfX), float64(halfY)) * .75
	center := pixel.V(float64(pixelPerCellX*x+halfX), float64(pixelPerCellY*y+halfY))
	c.drawStart = center.Sub(c.force.Scaled(shortestHalf))
	c.drawEnd = center.Add(c.force.Scaled(shortestHalf))
}

// NewField ...
func NewField() Field {
	f := Field{
		cells: [numX][numY]cell{},
		color: pixel.RGB(.5, .5, .5),
	}
	f.randomizeFlow()
	f.averageFlow()
	return f
}

func (f *Field) randomizeFlow() {
	for x := 0; x < numX; x++ {
		for y := 0; y < numY; y++ {
			f.cells[x][y] = newCell(x, y)
		}
	}
}

func (f *Field) averageFlow() {

	for x := 0; x < numX; x++ {
		for y := 0; y < numY; y++ {
			neighbors := make([]pixel.Vec, 8, 8)
			// left
			if x == 0 {
				neighbors[0] = f.cells[x][y].force
			} else {
				neighbors[0] = f.cells[x-1][y].force
			}
			// topleft
			if x == 0 || y == numY-1 {
				neighbors[1] = f.cells[x][y].force
			} else {
				neighbors[1] = f.cells[x-1][y+1].force
			}
			// top
			if y == numY-1 {
				neighbors[2] = f.cells[x][y].force
			} else {
				neighbors[2] = f.cells[x][y+1].force
			}
			// topright
			if x == numX-1 || y == numY-1 {
				neighbors[3] = f.cells[x][y].force
			} else {
				neighbors[3] = f.cells[x+1][y+1].force
			}
			// right
			if x == numX-1 {
				neighbors[4] = f.cells[x][y].force
			} else {
				neighbors[4] = f.cells[x+1][y].force
			}
			// bottomright
			if x == numX-1 || y == 0 {
				neighbors[5] = f.cells[x][y].force
			} else {
				neighbors[5] = f.cells[x+1][y-1].force
			}
			// bottom
			if y == 0 {
				neighbors[6] = f.cells[x][y].force
			} else {
				neighbors[6] = f.cells[x][y-1].force
			}
			// bottomleft
			if x == 0 || y == 0 {
				neighbors[7] = f.cells[x][y].force
			} else {
				neighbors[7] = f.cells[x-1][y-1].force
			}
			sumX := f.cells[x][y].force.X
			sumY := f.cells[x][y].force.Y
			for i := 0; i < len(neighbors); i++ {
				sumX = sumX + neighbors[i].X
				sumY = sumY + neighbors[i].Y
			}
			averageX := sumX / 9
			averageY := sumY / 9
			f.cells[x][y].force = pixel.V(averageX, averageY)
		}
	}
	for x := 0; x < numX; x++ {
		for y := 0; y < numY; y++ {
			f.cells[x][y].calculateDrawPoints(x, y)
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
