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

// BranchInfo ...
type BranchInfo struct {
	posV        pixel.Vec
	lenV        pixel.Vec
	m           pixel.Matrix
	angle       float64
	offsetAngle float64
}

// NewBranchInfo ...
func NewBranchInfo(posV, lenV pixel.Vec, angle, offsetAngle float64) BranchInfo {
	return BranchInfo{
		posV:        posV,
		lenV:        lenV,
		m:           pixel.IM.Moved(posV),
		angle:       angle,
		offsetAngle: offsetAngle,
	}
}

func branch(imd *imdraw.IMDraw, b BranchInfo) {

	thickness := b.lenV.Len() / 16
	if thickness > 5 {
		imd.EndShape = imdraw.SharpEndShape
	} else {
		imd.EndShape = imdraw.NoEndShape
	}

	imd.Push(b.posV)
	b.m = b.m.Moved(b.lenV)
	b.m = b.m.Rotated(b.posV, b.angle)
	b.posV = b.m.Project(pixel.ZV)
	imd.Push(b.posV)
	if thickness > 1 {
		imd.Line(thickness)
	} else {
		imd.Line(1)
	}
	if b.lenV.Len() > 4 {
		b.lenV = b.lenV.Scaled(.618)
		// left branch
		b.angle = b.angle + b.offsetAngle
		branch(imd, b)
		// right branch
		b.angle = b.angle - b.offsetAngle*2
		branch(imd, b)
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
	imd.EndShape = imdraw.NoEndShape

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
			if angle < math.Pi {
				angle = angle + 0.01
			}
		} else if win.Pressed(pixelgl.KeyD) {
			if angle > 0 {
				angle = angle - 0.01
			}
		}

		// DRAW
		win.Clear(colornames.Grey)
		imd.Clear()

		root := pixel.V(win.Bounds().Center().X, 0)
		b := NewBranchInfo(root, length, 0, angle)
		branch(imd, b)

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
