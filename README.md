# go-window

A multi-platform library for easy access to window and input in Go/Golang

### Usage

<pre><code>
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"

	gc "github.com/freiny/go-canvas"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var config = gc.Config{}
var cb = gc.Callbacks{}

func main() {
	config.Width = 512
	config.Height = 512
	config.X = 0
	config.Y = 0

	cb.Render = onRender
	cb.Key = onKey
	cb.CursorMove = onCursorMove

	rand.Seed(time.Now().UTC().UnixNano())
	gc.Init(config, cb)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == 1 {
		fmt.Print(string(key), " DOWN")

	}

	if action == 0 {
		fmt.Print(string(key), " UP")
	}
}

var xCursor, yCursor float64

func onCursorMove(xPos float64, yPos float64) {
	xCursor, yCursor = xPos, yPos
	// fmt.Println("CURSOR: ", xCursor, yCursor)
}

func onRender() *image.RGBA {
	xCur, yCur := int(xCursor), int(yCursor)

	rgba := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	for i := 0; i < 9000; i++ {

		x := rand.Intn(config.Width)
		y := rand.Intn(config.Height)

		r := uint8(rand.Intn(255))
		g := uint8(rand.Intn(255))
		b := uint8(rand.Intn(255))

		rgba.Set(x, y, color.RGBA{r, g, b, 255})
	}

	var d = 20
	for x := 0; x < d; x++ {
		for y := 0; y < d; y++ {
			rgba.Set(xCur+x-d, config.Height-yCur+y-d, color.RGBA{255, 0, 0, 255})
		}
	}

	return rgba

}
</code></pre>
