package gowindow

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"
	"os"
)

// BlankImage returns a blank RGBA image
func BlankImage() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, 1, 1))
}

// ClearImage clears input image.RGBA to specified color
func (f Framework) ClearImage(rgba *image.RGBA, c color.RGBA) *image.RGBA {
	point := rgba.Bounds().Size()
	w := point.X
	h := point.Y

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			rgba.Set(x, y, c)
		}
	}

	return rgba
}

// GetImagePart returns an RGBA pointer from a partial .png file
func (f Framework) GetImagePart(filename string, point image.Point) *image.RGBA {

	file := filename
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatal(fmt.Errorf("file %q not found on disk: %v", file, err))
		return nil
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
		// log.Fatal(fmt.Errorf("error decoding file %q: %v", imgFile, err))
		return nil
	}

	// rgba := image.NewRGBA(img.Bounds())
	rgba := image.NewRGBA(image.Rect(0, 0, point.X, point.Y))
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatal(fmt.Errorf("unsupported stride"))
		return nil
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}

// RandomImage clears input image.RGBA to specified color
func (f Framework) RandomImage(p image.Point) *image.RGBA {

	w, h := p.X, p.Y
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r := uint8(rand.Intn(255))
			g := uint8(rand.Intn(255))
			b := uint8(rand.Intn(255))
			rgba.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	return rgba
}

// var d = 20
// for x := 0; x < d; x++ {
// 	for y := 0; y < d; y++ {
// 		rgba.Set(xCur+x-d, cfg.Height-yCur+y-d, color.RGBA{255, 0, 0, 255})
// 	}
// }

// GetImage returns an RGBA pointer from a .png file
func (f Framework) GetImage(filename string) *image.RGBA {

	file := filename
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatal(fmt.Errorf("texture %q not found on disk: %v", file, err))
		return nil
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(fmt.Errorf("error decoding file %q: %v", imgFile, err))
		return nil
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatal("unsupported stride")
		return nil
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}
