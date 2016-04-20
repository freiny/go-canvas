package gowindow

import (
	"image"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
)

type cbRender func() *image.RGBA
type cbCursorMove func(float64, float64)
type cbKey func(w *Window, k Key, scancode int, action Action, mods ModifierKey)
type cbFPS func(fps int)

// Callbacks holds the callbacks defined in the User Application  ran in the library
type Callbacks struct {
	Render     cbRender
	CursorMove cbCursorMove
	Key        cbKey
	FPS        cbFPS
}

// WinConfig holds global data (e.g. window dimensions, cursor location)
type WinConfig struct {
	W int
	H int
	X int
	Y int
}

// Window wraps glfw.Window
type Window struct {
	glfw.Window
}

// Key replaces glfw.Key
type Key int

// Action replaces glfw.Action
type Action int

// ModifierKey replaces glfw.ModifierKey
type ModifierKey int

// SetDebug turn debugging mode on or off
func SetDebug() {
	debug = true
}

// Toggle alternates between returning true and false each frame
func Toggle() bool {
	frameToggle = !frameToggle
	return frameToggle
}

// GetTime returns current time in nano seconds
func GetTime() int64 {
	return time.Now().UTC().UnixNano()
}
