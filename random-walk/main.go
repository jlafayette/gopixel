package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	screenWidth  = 800
	screenHeight = 800
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Random-Walk",
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
	imd.Color = pixel.RGB(.1, .1, .1)
	imd.EndShape = imdraw.RoundEndShape
	background := pixel.RGB(.9, .9, .9)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	win.Clear(background)

	w := NewWalker(100)

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		w.update()
		win.Update()

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		imd.Clear()
		w.draw(imd)
		imd.Draw(win)
	}
}

func main() {
	pixelgl.Run(run)
}
