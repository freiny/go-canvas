package gowindow

import (
	"fmt"
	"image"
	"reflect"

	"github.com/go-gl/glfw/v3.1/glfw"
)

var wc = WinConfig{}
var cb = Callbacks{}

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

// Callbacks holds the callbacks defined in the User Application  ran in the library
type Callbacks struct {
	Render     func() *image.RGBA
	CursorMove func(c Cursor)
	Key        func(w *Window, k Key, scancode int, a Action, mk ModifierKey)
	FPS        func(fps int)
}

func onRenderNil() *image.RGBA                                          { return BlankImage() }
func onCursorMoveNil(c Cursor)                                          {}
func onKeyNil(w *Window, k Key, scancode int, a Action, mk ModifierKey) {}
func onFPSNil(fps int)                                                  {}

// RegisterCallback allows setting user defined callbacks
func RegisterCallback(cbName string, cbHandler interface{}) {
	f := Callbacks{}

	this := reflect.TypeOf(cbHandler)
	render := reflect.TypeOf(f.Render)
	cursorMove := reflect.TypeOf(f.CursorMove)
	key := reflect.TypeOf(f.Key)
	fps := reflect.TypeOf(f.FPS)

	if this == render {
		fmt.Println("render")
		cb.Render, _ = cbHandler.(func() *image.RGBA)
	}

	if this == cursorMove {
		fmt.Println("cursor")
		cb.CursorMove = cbHandler.(func(c Cursor))
	}

	if this == key {
		fmt.Println("key")
		cb.Key = cbHandler.(func(w *Window, k Key, scancode int, a Action, mk ModifierKey))
	}

	if this == fps {
		fmt.Println("fps")
		cb.FPS = cbHandler.(func(fps int))
	}
}
