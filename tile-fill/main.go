package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
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

// draw the dots
func (d *Cells) drawdots(imd *imdraw.IMDraw) {
	for _, cell := range d.cells {
		imd.Push(cell.center)
		imd.Circle(5, 0)
	}
}

func (d *Cells) generateVoronoi() {
	// evaluate each pixel
	var v pixel.Vec
	var minDistance float64
	var currentDistance float64
	var closestCellIndex int
	leftIndex := -1
	bttmIndexes := make([]int, d.boundsMaxX+1)
	// start at lower left, process the whole row, then go up one and and continue
	for y := 0; y < d.boundsMaxY+1; y++ {
		for x := 0; x < d.boundsMaxX+1; x++ {

			closestCellIndex = -1
			if x == d.boundsMaxX {
				// far right case
				v = pixel.V(float64(x-1), float64(y))
			} else if y == d.boundsMaxY {
				// top case
				v = pixel.V(float64(x), float64(y-1))
			} else {
				v = pixel.V(float64(x), float64(y))
			}
			// find closest dot
			minDistance = d.bounds.Size().Len()
			for i, cell := range d.cells {
				currentDistance = v.Sub(cell.center).Len()
				if currentDistance < minDistance {
					closestCellIndex = i
					minDistance = currentDistance
				}
			}

			// idSet stores a map of the different cell indexes that are around the current
			// pixel being evaluated. Pixels being evaluated are the current pixel, the
			// pixel to the left, the pixel down, the pixel down and to the left. If 3 or
			// more of them are different, it's a vertex.
			idSet := make(map[int]bool)
			if x == 0 {
				leftIndex = -1
			}
			if y == 0 {
				bttmIndexes[x] = -1
			}
			// if x == 0 || x == d.boundsMaxX || y == 0 || y == d.boundsMaxY {
			// 	idSet[-1] = true
			// }
			idSet[closestCellIndex] = true
			idSet[leftIndex] = true
			if y > 0 {
				idSet[bttmIndexes[x]] = true
				if x > 0 {
					idSet[bttmIndexes[x-1]] = true
				}
			}
			if len(idSet) > 2 {
				for k := range idSet {
					// debug circle
					if k >= 0 {
						d.cells[k].addPoint(v)
						// d.cells[k].imd.Push(v)
						// d.cells[k].imd.Circle(float64(k), 1)
					}
				}
			}

			leftIndex = closestCellIndex
			bttmIndexes[x] = closestCellIndex
		}
	}
}

// distance between two vectors. This is the same as v1.Sub(v2).Len()
// TODO: Test these to find most efficient option
func distance(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
	// return v1.Sub(v2).Len()
}

func (d *Cells) draw(tgt pixel.Target) {
	for _, cell := range d.cells {
		cell.draw(tgt)
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Tile-Fill",
		Bounds: pixel.R(0, 0, 1024, 768),
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

	// rand.Seed(time.Now().Unix())
	rand.Seed(99)
	d := NewCells(25, win.Bounds())

	// Move to main loop later ... testing voronoi
	win.Clear(colornames.Gray)
	// imd.Clear()
	// d.drawdots(imd)
	d.generateVoronoi()
	d.draw(win)
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
