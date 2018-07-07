package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func branch(imd *imdraw.IMDraw, posV, lengthV pixel.Vec, m pixel.Matrix, angle, offsetAngle, thickness float64) {

	imd.Push(posV)
	m = m.Moved(lengthV)
	m = m.Rotated(posV, angle)
	posV = m.Project(pixel.ZV)
	imd.Push(posV)
	imd.Line(thickness)
	if lengthV.Len() > 4 {
		lengthV = lengthV.Scaled(.67)
		if thickness > 1 {
			thickness--
		}
		// left branch
		branch(
			imd,
			posV,
			lengthV,
			m,
			angle+offsetAngle,
			offsetAngle,
			thickness,
		)
		// right branch
		branch(
			imd,
			posV,
			lengthV,
			m,
			angle-offsetAngle,
			offsetAngle,
			thickness,
		)
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "IMDraw",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	imd.Color = colornames.Whitesmoke
	imd.EndShape = imdraw.SharpEndShape

	var (
		frames = 0
		second = time.Tick(time.Second)
		length = pixel.V(0, 250)
		angle  = math.Pi / 4
	)

	// main loop
	for !win.Closed() {
		// UPDATE

		// w longer
		// s shorter
		if win.Pressed(pixelgl.KeyW) {
			length = length.Add(pixel.V(0, 1))
		} else if win.Pressed(pixelgl.KeyS) {
			length = length.Add(pixel.V(0, -1))
		}
		// a smaller angle
		// d larger angle
		if win.Pressed(pixelgl.KeyA) {
			angle = angle + 0.01
		} else if win.Pressed(pixelgl.KeyD) {
			angle = angle - 0.01
		}

		// DRAW
		win.Clear(colornames.Grey)
		imd.Clear()

		root := pixel.V(win.Bounds().Center().X, 0)
		mat := pixel.IM.Moved(root)
		branch(imd, root, length, mat, 0, angle, 12)

		imd.Draw(win)
		win.Update()

		// framerate
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
