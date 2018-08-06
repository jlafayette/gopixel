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
		Title:  "Follow Path",
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

	path := NewPath()
	var vehicles []Vehicle
	vehicles = append(vehicles, NewVehicle(pixel.V(screenWidth/2, screenHeight/2), &path))

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			vehicles = append(vehicles, NewVehicle(win.MousePosition(), &path))
		}

		// path.update()
		for i := 0; i < len(vehicles); i++ {
			vehicles[i].update()
		}

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS %d | Vehicles %d", cfg.Title, frames, len(vehicles)))
			frames = 0
		default:
		}
		win.Clear(background)
		imd.Clear()
		// path.draw(imd)
		for i := 0; i < len(vehicles); i++ {
			vehicles[i].draw(imd)
		}
		imd.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
