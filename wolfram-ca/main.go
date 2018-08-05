package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	screenWidth  = 1800
	screenHeight = 960
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
	// win.SetSmooth(true)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	background := pixel.RGB(.9, .9, .9)
	win.Clear(background)

	var ruleNum uint8 = 30
	var step uint8 = 1
	r := NewRule(ruleNum, 4)
	new := true

	// main loop
	for !win.Closed() {

		// UPDATE
		frames++
		if win.Pressed(pixelgl.KeyLeftControl) {
			step = 8
		} else if win.Pressed(pixelgl.KeyLeftShift) {
			step = 16
		} else if win.Pressed(pixelgl.KeyLeftAlt) {
			step = 32
		} else {
			step = 1
		}

		if win.JustPressed(pixelgl.KeyRight) {
			ruleNum = ruleNum + step
			r = NewRule(ruleNum, 4)
			new = true
		}
		if win.JustPressed(pixelgl.KeyLeft) {
			ruleNum = ruleNum - step
			r = NewRule(ruleNum, 4)
			new = true
		}
		if new {
			r.update()
		}

		// DRAW
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | Rule %d   %08b", cfg.Title, ruleNum, byte(ruleNum)))
			frames = 0
		default:
		}
		if new {
			win.Clear(background)
			r.draw(win)
			new = false
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
