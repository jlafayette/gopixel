package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	imageWidth  = 600
	imageHeight = 600
	nSites      = 25
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
// For easier pixel calulations, each cell center is always an integer
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
		newcells.cells[i] = NewCell(x, y, i)
	}
	return newcells
}

// randomColor generates a random color
func randomColor() color.NRGBA {
	return color.NRGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
}

func findClosestSeed(c *Cells, x, y int) (int16, pixel.Vec) {
	v := pixel.V(float64(x), float64(y))
	var closestSeedIndex int16
	var minDistance int
	var currentDistance int
	minDistance = distance(c.boundsMaxX, c.boundsMaxY)
	for i := 0; i < len(c.cells); i++ {
		currentDistance = distance(c.cells[i].cx-x, c.cells[i].cy-y)
		if currentDistance <= minDistance {
			closestSeedIndex = int16(i)
			minDistance = currentDistance
		}
	}
	return closestSeedIndex, v
}

func (c *Cells) generateVoronoi() {
	width := c.boundsMaxX
	height := c.boundsMaxY
	master := make([]int16, width*height)

	// top right corner
	x := width - 1
	y := height - 1
	closestSeedIndex, v := findClosestSeed(c, x, y)
	c.cells[closestSeedIndex].addPoint(v)
	master[x+y*width] = closestSeedIndex

	// horizontal edges and left corners
	var leftIndex int16
	for _, y := range []int{0, height - 1} {
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
		btIndex = -1
		for y := 0; y < height-1; y++ {
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

// uses dot to efficently calcuate distance between two points.
func distance(x, y int) int {
	return x*x + y*y
}

func (c *Cells) draw(tgt pixel.Target) {
	for _, cell := range c.cells {
		cell.draw(tgt)
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Tile-Fill",
		Bounds: pixel.R(0, 0, imageWidth, imageHeight),
		// Bounds: pixel.R(0, 0, 600, 400),
		VSync: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// imd := imdraw.New(nil)
	// imd.Color = colornames.Whitesmoke
	// imd.EndShape = imdraw.NoEndShape

	rand.Seed(time.Now().Unix())
	// rand.Seed(99)
	c := NewCells(nSites, win.Bounds())

	pic := generateVoronoi(sitesFromCells(c))
	sprite := pixel.NewSprite(pic, win.Bounds())
	// Move to main loop later ... testing voronoi
	win.Clear(colornames.Gray)

	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(1, -1))
	sprite.Draw(win, mat)

	// imd.Clear()
	c.generateVoronoi()
	c.draw(win)
	// imd.Draw(win)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	// main loop
	for !win.Closed() {
		// UPDATE

		// DRAW
		win.Update()

		// framerate
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
