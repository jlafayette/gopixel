package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func branch(imd *imdraw.IMDraw, posV, lengthV pixel.Vec, m pixel.Matrix, angle, thickness float64) {

	imd.Push(posV)
	m = m.Moved(lengthV)
	m = m.Rotated(posV, angle)
	posV = m.Project(pixel.ZV)
	imd.Push(posV)
	imd.Line(thickness)
	if lengthV.Len() > 2 {
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
			angle+math.Pi/4,
			thickness,
		)
		// right branch
		branch(
			imd,
			posV,
			lengthV,
			m,
			angle-math.Pi/4,
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
	length := pixel.V(0, 250)
	root := pixel.V(win.Bounds().Center().X, 0)
	mat := pixel.IM.Moved(root)
	branch(imd, root, length, mat, 0, 12)

	// main loop
	for !win.Closed() {
		// update

		// draw
		win.Clear(colornames.Grey)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
