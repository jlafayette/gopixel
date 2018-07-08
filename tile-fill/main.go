package main

import (
	"fmt"
	"image/color"
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

func (d *Cells) drawdots(imd *imdraw.IMDraw) {
	for i := range d.dots {
		imd.Push(d.dots[i])
		imd.Circle(5, 0)
	}
}

func (d *Cells) generateVoronoi() {
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
	imd.Clear()
	imd.Color = colornames.Whitesmoke
	d.drawdots(imd)
	// d.generateVoronoi(...)
	imd.Draw(win)

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
