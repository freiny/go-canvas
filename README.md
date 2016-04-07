# go-canvas
Render images and pixel data to the screen with Go/Golang

### Using go-canvas

<pre><code>
package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	gc "github.com/freiny/go-canvas"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const width = 512
const height = 512

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	gc.Init(width, height, render, onKey, onCursorMove)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == 1 {
		fmt.Print(string(key), " DOWN")

	}

	if action == 0 {
		fmt.Print(string(key), " UP")
	}
}

var x, y float64

func onCursorMove(xPos float64, yPos float64) {
	x, y = xPos, yPos
	fmt.Println("CURSOR: ", x, y)
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
