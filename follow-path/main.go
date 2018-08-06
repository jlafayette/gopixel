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
	screenWidth  = 1600
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
		frames     = 0
		prevFrames = 0
		second     = time.Tick(time.Second)
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
	var removeIndex int
	var remove bool
	var spawnDelay int

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		spawnDelay++
		if spawnDelay > 10 {
			vehicles = append(vehicles, NewVehicle(pixel.V(1, randFloat(10, screenHeight-10)), &path))
			spawnDelay = 0
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			vehicles = append(vehicles, NewVehicle(win.MousePosition(), &path))
		}
		if win.Pressed(pixelgl.MouseButtonRight) {
			vehicles = append(vehicles, NewVehicle(win.MousePosition(), &path))
		}
		if win.JustPressed(pixelgl.KeyLeftControl) {
			path.points = randomPoints()
		}
		for i := 0; i < len(vehicles); i++ {
			vehicles[i].update()
			if !win.Bounds().Contains(vehicles[i].pos) {
				removeIndex = i
				remove = true
			}
		}
		if remove {
			// Delete without preserving order
			vehicles[removeIndex] = vehicles[len(vehicles)-1]
			vehicles[len(vehicles)-1] = Vehicle{}
			vehicles = vehicles[:len(vehicles)-1]
			remove = false
		}
		// DRAW
		win.SetTitle(fmt.Sprintf("%s | FPS %d | Vehicles %d", cfg.Title, prevFrames, len(vehicles)))
		select {
		case <-second:
			prevFrames = frames
			frames = 0
		default:
		}
		win.Clear(background)
		imd.Clear()
		path.draw(imd)
		// path.drawClosest(win.MousePosition(), imd) // debug
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
