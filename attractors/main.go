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
	screenWidth  = 800
	screenHeight = 800
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Attractors",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		// VSync:  true,
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

	background := color.RGBA{220, 220, 220, 255}
	foreground := colornames.Black
	imd := imdraw.New(nil)
	imd.Color = foreground
	imd.EndShape = imdraw.RoundEndShape
	imd.Intensity = .1
	win.Clear(background)

	a1 := NewAttractor(pixel.V(screenWidth/2, screenHeight/2))
	var particles []Particle
	for i := 0; i < 50; i++ {
		p := NewParticle(pixel.V(screenWidth/2-(150-float64(i)), screenHeight/2), pixel.V(0, 0+(float64(i)*.04)))
		particles = append(particles, p)
	}
	engine := NewEngine([]Attractor{a1}, particles)

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		win.Update()
		engine.update()

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		// win.Clear(background)
		imd.Clear()
		engine.draw(imd)
		imd.Draw(win)
	}
}

func main() {
	pixelgl.Run(run)
}
