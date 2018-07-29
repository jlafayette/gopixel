package main

// Copied from https://rosettacode.org/wiki/Voronoi_diagram#Go with modifications
// so it works with pixel instead of writing to an image.
// Copyright (c) 2018 rosettacode.org
// License (GNU Free Documentation License 1.2) http://www.gnu.org/licenses/fdl-1.2.html

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"

	"github.com/faiface/pixel"
)

func generateVoronoi(sx, sy []int) pixel.Picture {
	// generate a random color for each site
	sc := make([]color.NRGBA, initNumSites)
	for i := range sx {
		sc[i] = color.NRGBA{
			uint8(rand.Intn(156) + 100),
			uint8(rand.Intn(156) + 100),
			uint8(rand.Intn(156) + 100), 255}
	}

	// generate diagram by coloring each pixel with color of nearest site
	img := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			dMin := dot(imageWidth, imageHeight)
			var sMin int
			for s := 0; s < initNumSites; s++ {
				if d := dot(sx[s]-x, sy[s]-y); d < dMin {
					sMin = s
					dMin = d
				}
			}
			img.SetNRGBA(x, y, sc[sMin])
		}
	}
	// mark each cell center with a black box
	black := image.NewUniform(color.Black)
	for s := 0; s < initNumSites; s++ {
		draw.Draw(img, image.Rect(sx[s]-2, sy[s]-2, sx[s]+2, sy[s]+2),
			black, image.ZP, draw.Src)
	}
	pic := pixel.PictureDataFromImage(img)
	return pic
}

func dot(x, y int) int {
	return x*x + y*y
}

func sitesFromCells(c Cells) (sx, sy []int) {
	sx = make([]int, len(c.cells))
	sy = make([]int, len(c.cells))
	for i, cell := range c.cells {
		sx[i] = cell.seedX
		sy[i] = cell.seedY
	}
	return
}
