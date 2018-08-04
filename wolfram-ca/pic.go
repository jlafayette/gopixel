package main

import (
	"image"
	"os"

	"image/color"
	"image/png"

	"github.com/faiface/pixel"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func writePngFile(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err = png.Encode(f, img); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

func writeColorToPng(path string, w, h int, fill color.RGBA) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, fill)
		}
	}
	return writePngFile(path, img)
}
