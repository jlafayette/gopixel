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
	bounds     pixel.Rect      // bounds to fill with dots
	boundsMinX int             // lower bound X as int
	boundsMaxX int             // upper bound X as int
	boundsMinY int             // lower bound Y as int
	boundsMaxY int             // upper bound Y as int
	dots       []pixel.Vec     // vector for each dot, these are cell centers
	colors     []color.NRGBA   // color for each dot / cell
	imds       []imdraw.IMDraw // polygon drawer for each cell
	dx         []int           // dot x coordinates as int
	dy         []int           // dot y coordinates as int
}

// NewCells returns a new Cells object with given number of cells
// For easier pixel calulations, each dot is always an integer
func NewCells(n int, bounds pixel.Rect) Cells {
	d := Cells{
		bounds:     bounds,
		boundsMinX: int(bounds.Min.X),
		boundsMaxX: int(bounds.Max.X),
		boundsMinY: int(bounds.Min.Y),
		boundsMaxY: int(bounds.Max.Y),
		dots:       make([]pixel.Vec, n),
		dx:         make([]int, n),
		dy:         make([]int, n),
	}
	for i := range d.dots {
		x := rand.Intn(d.boundsMaxX-d.boundsMinX) + d.boundsMinX
		y := rand.Intn(d.boundsMaxY-d.boundsMinY) + d.boundsMinY
		d.dots[i] = pixel.V(float64(x), float64(y))
		d.computeColors()
		d.dx[i] = x
		d.dy[i] = y
	}
	return d
}

// computeColors generates a random color for each cell / dot
func (d *Cells) computeColors() {
	d.colors = make([]color.NRGBA, len(d.dots))
	for i := range d.dots {
		d.colors[i] = color.NRGBA{
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			255,
		}
	}
}

// draw the dots
func (d *Cells) drawdots(imd *imdraw.IMDraw) {
	for i := range d.dots {
		imd.Push(d.dots[i])
		imd.Circle(5, 0)
	}
}

func (d *Cells) generateVoronoi() {
	// generate a polygon drawer for each cell to track them
	d.imds = make([]imdraw.IMDraw, len(d.dots))
	for i := range d.imds {
		d.imds[i] = *imdraw.New(nil)
		d.imds[i].Color = d.colors[i]
		d.imds[i].EndShape = imdraw.NoEndShape
	}

	// evaluate each pixel
	var v pixel.Vec
	var minDistance float64
	var currentDistance float64
	var closestDotIndex int
	leftIndex := -1
	bttmIndexes := make([]int, d.boundsMaxX+1)
	// start at lower left, process the whole row, then go up one and and continue
	for y := 0; y < d.boundsMaxY+1; y++ {
		for x := 0; x < d.boundsMaxX+1; x++ {

			// idSet stores a map of the different cell indexes that are around the current
			// pixel being evaluated. Pixels being evaluated are the current pixel, the
			// pixel to the left, the pixel down, the pixel down and to the left.
			idSet := make(map[int]bool)
			closestDotIndex = -1

			if x == d.boundsMaxX {
				// far right case
				idSet[-1] = true
				v = pixel.V(float64(x-1), float64(y))
			} else if y == d.boundsMaxY {
				// top case
				idSet[-1] = true
				v = pixel.V(float64(x), float64(y-1))
			} else {
				v = pixel.V(float64(x), float64(y))
			}
			// find closest dot
			minDistance = d.bounds.Size().Len()
			for i, dotV := range d.dots {
				currentDistance = v.Sub(dotV).Len()
				if currentDistance < minDistance {
					closestDotIndex = i
					minDistance = currentDistance
				}
			}

			// Evaluate the bottom left corner of the current pixel, there are 4 pixels,
			// if 3 of them are different, it's a meeting point
			idSet[closestDotIndex] = true
			idSet[leftIndex] = true
			idSet[bttmIndexes[x]] = true
			if x > 0 {
				idSet[bttmIndexes[x-1]] = true
			}
			if len(idSet) > 2 {
				for k := range idSet {
					// debug circle
					if k >= 0 {
						d.imds[k].Push(v)
						d.imds[k].Circle(float64(k), 1)
					}
				}

				// TODO: handle top case, handle far right case

				// draw all the lines... diagonals are broken
				// if numCells >= 2 {
				// 	// add to drawer for each cell
				// 	for _, index := range cellsToDraw {

				// 		// debug circle
				// 		d.imds[index].Push(v)
				// 		d.imds[index].Circle(1, 1)

				// 	}
			}

			leftIndex = closestDotIndex
			bttmIndexes[x] = closestDotIndex
		}
	}

	for i, imd := range d.imds {
		// imd.Circle(5, 1)
		imd.Push(d.dots[i])
		imd.Circle(5, 0)
	}
}

// distance between two vectors. This is the same as v1.Sub(v2).Len()
// TODO: Test these to find most efficient option
func distance(v1, v2 pixel.Vec) float64 {
	return math.Sqrt(math.Pow(v1.X-v2.X, 2) + math.Pow(v1.Y-v2.Y, 2))
	// return v1.Sub(v2).Len()
}

func (d *Cells) drawVoronoi(tgt pixel.Target) {
	for _, imd := range d.imds {
		imd.Draw(tgt)
	}
}

func dot(x, y int) int {
	return (x * x) + (y * y)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title: "Tile-Fill",
		// Bounds: pixel.R(0, 0, 1024, 768),
		Bounds: pixel.R(0, 0, 200, 200),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	imd.Color = colornames.Whitesmoke
	imd.EndShape = imdraw.NoEndShape

	rand.Seed(time.Now().Unix())
	d := NewCells(10, win.Bounds())

	// Move to main loop later ... testing voronoi
	win.Clear(colornames.Gray)
	// imd.Clear()
	// d.drawdots(imd)
	d.generateVoronoi()
	d.drawVoronoi(win)
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
