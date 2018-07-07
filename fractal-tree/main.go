package main

import (
	"fmt"
	"math"
	"math/rand"
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
	branchAngle float64
	offsetAngle float64
	randAngle   float64
	randLen     float64
	r           *rand.Rand
}

// NewBranchInfo ...
func NewBranchInfo(posV, lenV pixel.Vec, angle, branchAngle, offsetAngle, randAngle, randLen float64, randSeed int64) BranchInfo {
	return BranchInfo{
		posV:        posV,
		lenV:        lenV,
		m:           pixel.IM.Moved(posV),
		angle:       angle,
		branchAngle: branchAngle,
		offsetAngle: offsetAngle,
		randAngle:   randAngle,
		randLen:     randLen,
		r:           rand.New(rand.NewSource(randSeed)),
	}
}

func branch(imd *imdraw.IMDraw, b BranchInfo) {

	if b.randLen > 0 {
		scaleMult := 1 + ((b.r.Float64() - 0.5) * b.randLen)
		b.lenV = b.lenV.Scaled(scaleMult)
	}

	thickness := b.lenV.Len() / 16
	if thickness > 5 {
		imd.EndShape = imdraw.SharpEndShape
	} else {
		imd.EndShape = imdraw.NoEndShape
	}

	imd.Push(b.posV)
	b.m = b.m.Moved(b.lenV)
	if b.angle != 0 {
		b.m = b.m.Rotated(b.posV, b.angle+((b.r.Float64()-0.5)*b.randAngle))
	}
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
		b.angle = b.angle + b.branchAngle + b.offsetAngle
		branch(imd, b)
		// right branch
		b.angle = b.angle - b.branchAngle*2
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
		frames      = 0
		second      = time.Tick(time.Second)
		length      = pixel.V(0, 250)
		branchAngle = math.Pi / 4
		offsetAngle = 0.0
		randAngle   = 0.0
		randLen     = 0.0
	)
	var randSeed int64
	randSeed = 99

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
			if branchAngle < math.Pi {
				branchAngle = branchAngle + 0.01
			}
		} else if win.Pressed(pixelgl.KeyD) {
			if branchAngle > 0 {
				branchAngle = branchAngle - 0.01
			}
		}
		// f smaller random
		// r larger random
		// n get new random seed
		if win.Pressed(pixelgl.KeyR) {
			if randAngle < math.Pi {
				randAngle = randAngle + 0.05
			}
		} else if win.Pressed(pixelgl.KeyF) {
			if randAngle > 0 {
				randAngle = randAngle - 0.05
			}
		}
		if win.JustPressed(pixelgl.KeyN) {
			randSeed = time.Now().UnixNano()
		}
		// z small length random
		// x larger length random (0 -> 1)
		if win.Pressed(pixelgl.KeyX) {
			if randLen < 1 {
				randLen = randLen + 0.02
			}
		} else if win.Pressed(pixelgl.KeyZ) {
			if randLen > 0 {
				randLen = randLen - 0.02
			}
		}

		// q offsetLeft
		// e offsetRight
		if win.Pressed(pixelgl.KeyQ) {
			if offsetAngle < math.Pi/2 {
				offsetAngle = offsetAngle + 0.01
			}
		} else if win.Pressed(pixelgl.KeyE) {
			if offsetAngle > -math.Pi/2 {
				offsetAngle = offsetAngle - 0.01
			}
		}

		// DRAW

		// test for changes before redrawing
		win.Clear(colornames.Grey)
		imd.Clear()

		root := pixel.V(win.Bounds().Center().X, 0)
		b := NewBranchInfo(root, length, 0, branchAngle, offsetAngle, randAngle, randLen, randSeed)
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
