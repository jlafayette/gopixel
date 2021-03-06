package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	imageWidth   = 800
	imageHeight  = 800
	initNumSites = 5
)

// DrawMode determines what should be drawn to the screen
type DrawMode int

// Possible DrawModes
const (
	Normal DrawMode = 0
	Debug  DrawMode = 1
)

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

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	c := NewCells(initNumSites, win.Bounds())
	background := color.RGBA{220, 220, 220, 255}
	foreground := colornames.Black
	imd := imdraw.New(nil)
	imd.Color = foreground
	imd.EndShape = imdraw.NoEndShape

	first := true          // switch to determine if it's the first loop.
	displayMode := Debug   // switch that determines what mode to draw.
	dirty := false         // switch that determines if things need to be redrawn.
	relax := false         // determine if cells should relax.
	nSites := initNumSites // track number of sites

	// main loop
	for !win.Closed() {

		// UPDATE
		if win.JustPressed(pixelgl.KeyUp) {
			nSites = nSites + 1
			first = true
		}
		if win.JustPressed(pixelgl.KeyDown) {
			if nSites > 1 {
				nSites = nSites - 1
				first = true
			}
		}
		if win.JustPressed(pixelgl.KeySpace) || win.JustPressed(pixelgl.KeyRight) || first {
			// new voronoi!
			seed := time.Now().UnixNano()
			fmt.Printf("running %v\n", seed)
			rand.Seed(seed)
			c = NewCells(nSites, win.Bounds())
			c.update()
			dirty = true
			first = false
		}
		if win.JustPressed(pixelgl.KeyRight) {
			relax = true
		}
		if relax {
			c.update()
			dirty = true
			if c.relaxed() {
				relax = false
			}
		}
		if win.JustPressed(pixelgl.KeyLeftControl) {
			dirty = true
			displayMode = Normal
		}
		if win.JustReleased(pixelgl.KeyLeftControl) {
			dirty = true
			displayMode = Debug
		}
		frames++
		win.Update()

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | Count: %d", cfg.Title, frames, nSites))
			frames = 0
		default:
		}
		if dirty {
			win.Clear(background)
			imd.Clear()
			switch displayMode {
			case Normal:
				c.draw(imd)
			case Debug:
				c.drawDebug(imd)
			}
			imd.Draw(win)
			dirty = false
		}
	}
}

func main() {
	pixelgl.Run(run)
}
