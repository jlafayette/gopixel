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
	var vehicles []Vehicle
	vehicles = append(vehicles, NewVehicle(pixel.V(screenWidth/2, screenHeight/2), &f))
	debug := true

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			vehicles = append(vehicles, NewVehicle(win.MousePosition(), &f))
			if len(vehicles) > 82 {
				debug = false
			}
		}
		if win.JustPressed(pixelgl.MouseButtonRight) {
			for i := 0; i < 39; i++ {
				vehicles = append(vehicles, NewVehicle(pixel.V(randFloat(1, screenWidth-1), randFloat(1, screenHeight-1)), &f))
			}
			if len(vehicles) > 82 {
				debug = false
			}
		}

		f.update()
		for i := 0; i < len(vehicles); i++ {
			vehicles[i].update(win.Bounds())
		}

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS %d | Vehicles %d | debug: %v", cfg.Title, frames, len(vehicles), debug))
			frames = 0
		default:
		}
		win.Clear(background)
		imd.Clear()
		f.draw(imd)
		for i := 0; i < len(vehicles); i++ {
			vehicles[i].draw(imd, debug)
		}
		imd.Draw(win)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
