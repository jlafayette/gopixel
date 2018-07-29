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
	imageWidth  = 300 //800
	imageHeight = 300 //800
	nSites      = 50
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

	// to display a reference image in the background
	pic, err := loadPicture("300px-LloydsMethod1.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	// c := NewCells(nSites, win.Bounds())
	c := NewCentroidTestCells(win.Bounds())
	background := color.RGBA{220, 220, 220, 255}
	foreground := colornames.Black
	imd := imdraw.New(nil)
	imd.Color = foreground
	imd.EndShape = imdraw.NoEndShape

	first := true        // switch to determine if it's the first loop.
	displayMode := Debug // switch that determines what mode to draw.
	dirty := false       // switch that determines if things need to be redrawn.

	// main loop
	for !win.Closed() {

		// UPDATE
		if win.JustPressed(pixelgl.KeySpace) || first {
			// new voronoi!
			seed := time.Now().UnixNano()
			fmt.Printf("running %v\n", seed)
			rand.Seed(seed)
			// c.randomize()
			c.generateVoronoi()
			c.update()
			dirty = true
			first = false
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
		// framerate
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
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
				sprite.Draw(win, mat) // background reference
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
