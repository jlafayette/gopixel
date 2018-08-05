package main

import (
	"fmt"
	"image/color"
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
		Title:  "Flow-Field",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	// win.SetSmooth(true)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	imd := imdraw.New(nil)
	imd.Color = color.RGBA{0, 0, 0, 255}
	imd.EndShape = imdraw.RoundEndShape
	background := pixel.RGB(.9, .9, .9)
	win.Clear(background)

	f := NewField()
	v := NewVehicle(pixel.V(200, screenHeight/2), &f)
	// newField := true

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		f.update()
		v.update()

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS %d", cfg.Title, frames))
			frames = 0
		default:
		}
		// if newField {
		// 	win.Clear(background)
		// 	f.draw()
		// 	newField = false
		// }
		win.Clear(background)
		imd.Clear()
		v.draw(imd)
		f.draw(imd)
		imd.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
