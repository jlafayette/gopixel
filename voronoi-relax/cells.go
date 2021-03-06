package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Cells ...
type Cells struct {
	bounds     pixel.Rect // bounds to fill with cells
	boundsMinX int        // lower bound X as int
	boundsMaxX int        // upper bound X as int
	boundsMinY int        // lower bound Y as int
	boundsMaxY int        // upper bound Y as int
	cells      []Cell
}

// NewCells returns a new Cells object with given number of cells
// For easier pixel calulations, each cell seed has an int for x
// and y.
func NewCells(n int, bounds pixel.Rect) Cells {
	newcells := Cells{
		bounds:     bounds,
		boundsMinX: int(bounds.Min.X),
		boundsMaxX: int(bounds.Max.X),
		boundsMinY: int(bounds.Min.Y),
		boundsMaxY: int(bounds.Max.Y),
		cells:      make([]Cell, n),
	}
	for i := 0; i < n; i++ {
		x := rand.Intn(newcells.boundsMaxX-newcells.boundsMinX) + newcells.boundsMinX
		y := rand.Intn(newcells.boundsMaxY-newcells.boundsMinY) + newcells.boundsMinY
		newcells.cells[i] = NewCell(x, y)
	}
	return newcells
}

// NewCentroidTestCells generates a particular set of cells to compare against a centroid example
func NewCentroidTestCells(bounds pixel.Rect) Cells {
	newcells := Cells{
		bounds:     bounds,
		boundsMinX: int(bounds.Min.X),
		boundsMaxX: int(bounds.Max.X),
		boundsMinY: int(bounds.Min.Y),
		boundsMaxY: int(bounds.Max.Y),
		cells: []Cell{
			NewCell(106, 272),
			NewCell(134, 207),
			NewCell(190, 287),
			NewCell(198, 254),
			NewCell(198, 153),
		},
	}
	return newcells
}

func (c *Cells) randomize() {
	for i := 0; i < len(c.cells); i++ {
		x := rand.Intn(c.boundsMaxX-c.boundsMinX) + c.boundsMinX
		y := rand.Intn(c.boundsMaxY-c.boundsMinY) + c.boundsMinY
		c.cells[i].reset(x, y)
	}
}

func (c *Cells) relaxed() bool {
	largestDistance := 0.0
	for i := 0; i < len(c.cells); i++ {
		if c.cells[i].lastRelaxDistance > largestDistance {
			largestDistance = c.cells[i].lastRelaxDistance
		}
	}
	if largestDistance <= 1.3 {
		return true
	}
	return false
}

func (c *Cells) generateVoronoi() {
	for i := 0; i < len(c.cells); i++ {
		c.cells[i].points = nil
	}
	width := c.boundsMaxX
	height := c.boundsMaxY
	master := make([]int16, width*height)

	// right corners
	for _, y := range []int{0, height - 1} {
		x := width - 1
		//y := height - 1
		closestSeedIndex, v := findClosestSeed(c, x, y)
		c.cells[closestSeedIndex].addPoint(v)
		master[x+y*width] = closestSeedIndex
	}

	// horizontal edges and left corners
	for _, y := range []int{0, height - 1} {
		var leftIndex int16
		leftIndex = -1
		for x := 0; x < width-1; x++ {
			closestSeedIndex, v := findClosestSeed(c, x, y)
			master[x+y*width] = closestSeedIndex
			if closestSeedIndex != leftIndex {
				c.cells[closestSeedIndex].addPoint(v)
				if leftIndex >= 0 {
					c.cells[leftIndex].addPoint(v)
				}
			}
			leftIndex = closestSeedIndex
		}
	}
	// vertical edges
	for _, x := range []int{0, width - 1} {
		var btIndex int16
		btIndex = master[x+0*width] // y is 0 for starting case
		for y := 1; y < height-1; y++ {
			closestSeedIndex, v := findClosestSeed(c, x, y)
			master[x+y*width] = closestSeedIndex
			if closestSeedIndex != btIndex {
				c.cells[closestSeedIndex].addPoint(v)
				if btIndex >= 0 {
					c.cells[btIndex].addPoint(v)
				}
			}
			btIndex = closestSeedIndex
		}
	}
	// middle points
	for y := 1; y < height; y++ {
		for x := 1; x < width; x++ {
			closestSeedIndex, v := findClosestSeed(c, x, y)
			master[x+y*width] = closestSeedIndex
			// idSet stores a map of the different cell indexes that are around the current
			// pixel being evaluated. Pixels being evaluated are the current pixel, the
			// pixel to the left, the pixel down, the pixel down and to the left. If 3 or
			// more of them are different, it's a vertex.
			idSet := make(map[int16]bool)
			idSet[closestSeedIndex] = true
			idSet[master[(x-1)+y*width]] = true
			idSet[master[x+(y-1)*width]] = true
			idSet[master[(x-1)+(y-1)*width]] = true
			if len(idSet) > 2 {
				for k := range idSet {
					if k >= 0 {
						c.cells[k].addPoint(v)
					}
				}
			}
		}
	}
}

func findClosestSeed(c *Cells, x, y int) (int16, pixel.Vec) {
	v := pixel.V(float64(x), float64(y))
	var closestSeedIndex int16
	var minDistance int
	var currentDistance int
	minDistance = distance(c.boundsMaxX, c.boundsMaxY)
	for i := 0; i < len(c.cells); i++ {
		currentDistance = distance(c.cells[i].seedX-x, c.cells[i].seedY-y)
		if currentDistance <= minDistance {
			closestSeedIndex = int16(i)
			minDistance = currentDistance
		}
	}
	return closestSeedIndex, v
}

// uses dot to efficently calcuate distance between two points.
func distance(x, y int) int {
	return x*x + y*y
}

func (c *Cells) update() {
	c.generateVoronoi()
	for i := 0; i < len(c.cells); i++ {
		c.cells[i].update()
	}
}

func (c *Cells) draw(imd *imdraw.IMDraw) {
	for i := 0; i < len(c.cells); i++ {
		c.cells[i].draw(imd)
	}
}

func (c *Cells) drawDebug(imd *imdraw.IMDraw) {
	for i := 0; i < len(c.cells); i++ {
		// imd.Color = randomColor(10, 150)
		c.cells[i].drawDebug(imd)
	}
}
