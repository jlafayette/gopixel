package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	screenWidth  = 800
	screenHeight = 800
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Flocking",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	// win.SetSmooth(true)

	var (
		frames     = 0
		prevFrames = 0
		second     = time.Tick(time.Second)
	)

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(0, 0, 0)
	imd.EndShape = imdraw.RoundEndShape
	background := pixel.RGB(.9, .9, .9)
	win.Clear(background)

	var boids []Boid
	for i := 0; i < 40; i++ {
		boids = append(boids, NewBoid(pixel.V(randFloat(10, screenWidth-10), randFloat(10, screenHeight-10))))
	}

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			boids = append(boids, NewBoid(win.MousePosition()))
		}
		for i := 0; i < len(boids); i++ {
			boids[i].update()
		}

		// DRAW
		win.SetTitle(fmt.Sprintf("%s | FPS %d | Vehicles %d", cfg.Title, prevFrames, len(boids)))
		select {
		case <-second:
			prevFrames = frames
			frames = 0
		default:
		}
		win.Clear(background)
		imd.Clear()
		for i := 0; i < len(boids); i++ {
			boids[i].draw(imd)
		}
		imd.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
