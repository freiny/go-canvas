package gowindow

import (
	"fmt"
	"image"
	"image/draw"
	"os"
)

// RGBAFromImage returns an RGBA pointer from a .png file
func RGBAFromImage(filename string) (*image.RGBA, error) {

	file := filename
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, err
}
