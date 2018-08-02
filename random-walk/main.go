package main

import (
	"fmt"
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
		Title:  "Random-Walk",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	// win.SetSmooth(true)

	pic, err := loadPicture("1-black-pixel.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	batch := pixel.NewBatch(&pixel.TrianglesData{}, pic)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	background := pixel.RGB(.9, .9, .9)
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	win.Clear(background)

	w := NewWalker(10, sprite)
	new := true

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		if win.JustPressed(pixelgl.KeyLeftControl) {
			w = NewWalker(10, sprite)
			new = true
		}
		w.update()

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
		w.draw(batch)
		batch.Draw(win)
		batch.Clear()

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
