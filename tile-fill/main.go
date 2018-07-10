package main

import (
	"fmt"
	"image/color"
	"math"
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

func (d *Cells) generateVoronoi() {
	// evaluate each pixel
	var v pixel.Vec
	var minDistance float64
	var currentDistance float64
	var closestCellIndex int
	leftIndex := -1
	bttmIndexes := make([]int, d.boundsMaxX)
	nextBtIndexes := make([]int, d.boundsMaxX)
	for i := 0; i < d.boundsMaxX; i++ {
		nextBtIndexes[i] = -1
	}
	// start at lower left, process the whole row, then go up one and and continue
	for y := 0; y < d.boundsMaxY; y++ {
		for i := 0; i < d.boundsMaxX; i++ {
			bttmIndexes[i] = nextBtIndexes[i]
		}
		for x := 0; x < d.boundsMaxX; x++ {

			closestCellIndex = -1
			if x == d.boundsMaxX {
				// far right case
				v = pixel.V(float64(x), float64(y))
			} else if y == d.boundsMaxY {
				// top case
				v = pixel.V(float64(x), float64(y))
			} else {
				v = pixel.V(float64(x), float64(y))
			}
			// find closest cell center
			minDistance = d.bounds.Size().Len()
			for i, cell := range d.cells {
				currentDistance = v.Sub(cell.center).Len()
				if currentDistance <= minDistance {
					closestCellIndex = i
					minDistance = currentDistance
				}
			}

			// 0 bt lf corner
			// 1    lf edge
			// 2 tp lf corner
			// 3 tp    edge
			// 4 tp rt corner
			// 5    rt edge
			// 6 bt rt corner
			// 7 bt    edge
			// 8 middle
			var edgeCornerStatus int
			if x == 0 {
				if y == 0 {
					edgeCornerStatus = 0
				} else if y >= d.boundsMaxY-1 {
					edgeCornerStatus = 2
				} else {
					edgeCornerStatus = 1
				}
			} else if x >= d.boundsMaxX-1 {
				if y == 0 {
					edgeCornerStatus = 6
				} else if y >= d.boundsMaxY-1 {
					edgeCornerStatus = 4
				} else {
					edgeCornerStatus = 5
				}
			} else {
				if y == 0 {
					edgeCornerStatus = 7
				} else if y >= d.boundsMaxY-1 {
					edgeCornerStatus = 3
				} else {
					edgeCornerStatus = 8
				}
			}

			switch edgeCornerStatus {
			case 0: // 0 bt lf corner
				d.cells[closestCellIndex].addPoint(v)
			case 1: // 1    lf edge
				if closestCellIndex != bttmIndexes[x] {
					d.cells[closestCellIndex].addPoint(v)
					d.cells[bttmIndexes[x]].addPoint(v)
				}
			case 2: // 2 tp lf corner
				d.cells[closestCellIndex].addPoint(v)
			case 3: // 3 tp    edge
				if closestCellIndex != leftIndex {
					d.cells[closestCellIndex].addPoint(v)
					d.cells[leftIndex].addPoint(v)
				}
			case 4: // 4 tp rt corner
				d.cells[closestCellIndex].addPoint(v)
			case 5: // 5    rt edge
				if closestCellIndex != bttmIndexes[x] {
					d.cells[closestCellIndex].addPoint(v)
					d.cells[bttmIndexes[x]].addPoint(v)
				}
			case 6: // 6 bt rt corner
				d.cells[closestCellIndex].addPoint(v)
			case 7: // 7 bt    edge
				if closestCellIndex != leftIndex {
					d.cells[closestCellIndex].addPoint(v)
					d.cells[leftIndex].addPoint(v)
				}
			default: // 8 middle
				// idSet stores a map of the different cell indexes that are around the current
				// pixel being evaluated. Pixels being evaluated are the current pixel, the
				// pixel to the left, the pixel down, the pixel down and to the left. If 3 or
				// more of them are different, it's a vertex.
				idSet := make(map[int]bool)
				idSet[closestCellIndex] = true
				idSet[leftIndex] = true
				idSet[bttmIndexes[x]] = true
				idSet[bttmIndexes[x-1]] = true
				if len(idSet) > 2 {
					for k := range idSet {
						if k >= 0 {
							d.cells[k].addPoint(v)
						}
					}
				}
			}
			leftIndex = closestCellIndex
			nextBtIndexes[x] = closestCellIndex
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
