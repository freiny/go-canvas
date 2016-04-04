# go-canvas
Render images and pixel data to the screen with Go/Golang

### Using go-canvas

<pre><code>
package main

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	gc "github.com/freiny/go-canvas"
)

const width = 512
const height = 512

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	gc.Init(width, height, render)
}

func render() *image.RGBA {

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	for i := 0; i < 200000; i++ {

		x := rand.Intn(width)
		y := rand.Intn(height)

		r := uint8(rand.Intn(255))
		g := uint8(rand.Intn(255))
		b := uint8(rand.Intn(255))
		//
		rgba.Set(x, y, color.RGBA{r, g, b, 255})
	}
	return rgba

}
</code></pre>
