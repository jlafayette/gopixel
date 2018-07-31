package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 1200
	screenHeight = 1000
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Attractors",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	imd := imdraw.New(nil)
	imd.Color = color.RGBA{220, 220, 220, 255}
	imd.EndShape = imdraw.RoundEndShape
	background := colornames.Black
	win.Clear(background)

	engine := NewEngine(basic())
	new := false

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		if win.JustPressed(pixelgl.Key1) {
			engine = NewEngine(basic())
			new = true
		}
		if win.JustPressed(pixelgl.Key2) {
			engine = NewEngine(gasGiant())
			new = true
		}
		if win.JustPressed(pixelgl.Key3) {
			engine = NewEngine(random())
			new = true
		}
		if win.JustPressed(pixelgl.Key4) {
			engine = NewEngine(gravityPaths())
			new = true
		}
		engine.update()
		win.Update()

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		if new {
			win.Clear(background)
			new = false
		}
		imd.Clear()
		engine.draw(imd)
		imd.Draw(win)
	}
}

func main() {
	pixelgl.Run(run)
}
