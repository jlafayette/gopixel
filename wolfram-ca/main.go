package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	screenWidth  = 800
	screenHeight = 800
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Wolfram CA",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	// Create an 1 x 1 image
	err = writeColorToPng("out.png", 1, 1, color.RGBA{0, 0, 0, 255})
	if err != nil {
		panic(err)
	}

	// pic, err := loadPicture("out.png")
	// if err != nil {
	// 	panic(err)
	// }
	// sprite := pixel.NewSprite(pic, pic.Bounds())
	// batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	background := pixel.RGB(.9, .9, .9)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	win.Clear(background)

	r := NewRule()
	new := true

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		if new {
			r.update()
		}

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
		if new {
			win.Clear(background)
			r.draw(win)
			new = false
		}
		// batch.Draw(win)
		// batch.Clear()
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
