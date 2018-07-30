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

	background := colornames.Black
	foreground := color.RGBA{220, 220, 220, 255}
	imd := imdraw.New(nil)
	imd.Color = foreground
	imd.EndShape = imdraw.RoundEndShape
	win.Clear(background)

	var particles []Particle
	a1 := NewParticle(pixel.V(screenWidth/2, screenHeight/2), pixel.V(0, 0), 5000)
	a1.color = colornames.White
	particles = append(particles, a1)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	for i := 0; i < 5; i++ {
		p := NewOrbiter(a1)
		particles = append(particles, p)
	}
	engine := NewEngine(particles)

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
