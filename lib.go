package gowindow

import (
	"fmt"
	"image"
	"log"
	"reflect"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	runtime.LockOSThread()
}

// Framework object
type Framework struct {
	x, y, w, h int
	buildMode  string
	*glfw.Window
	cb Callbacks
}

// Config stores window x,y position and width,height for when framework is started
func Config(x int, y int, w int, h int) Framework {
	return Framework{x: x, y: y, w: w, h: h, buildMode: "Development"}
}

// BuildMode sets Production or Development build
func (f *Framework) BuildMode(mode ...string) string {
	// get
	if mode == nil {
		return f.buildMode
	}

	// set
	switch mode[0] {
	case "Development":
		f.buildMode = mode[0]
	case "Production":
		f.buildMode = mode[0]
	default:
		log.Fatal("unknown build mode")
	}
	return ""

}

// GetBuildMode sets Production or Development build

// ResourcePath gets path where resources are stored
func (f *Framework) ResourcePath() string {
	var path string
	if f.BuildMode() == "Production" {
		path = "/Applications/gwApp.app/Contents/Resources/"
	} else {
		path = "Resources/"
	}
	return path
}

// Start initializes OpenGL/GLFW then runs a render callback on each iteration
// func Init(wc WinConfig, cbUserDefined Callbacks) {
func (f *Framework) Start() {
	dbg.Log("lib.go", "Framework.Init()")
	if f.cb.Render == nil {
		dbg.Log("lib.go", "Framework.Init()", "cb.Render == nil")
		f.cb.Render = onRenderNil
	}
	if f.cb.CursorMove == nil {
		dbg.Log("lib.go", "Framework.Init()", "cb.CursorMove == nil")
		f.cb.CursorMove = onCursorMoveNil
	}
	if f.cb.Key == nil {
		dbg.Log("lib.go", "Framework.Init()", "cb.Key == nil")
		f.cb.Key = onKeyNil
	}
	if f.cb.FPS == nil {
		dbg.Log("lib.go", "Framework.Init()", "cb.FPS == nil")
		f.cb.FPS = onFPSNil
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(f.w, f.h, "", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetPos(f.x, f.y)

	f.initGL()

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		f.cb.Key(&Window{*w}, Key(key), scancode, Action(action), ModifierKey(mods))
	})

	window.SetCursorPos(0.0, 0.0)
	// win := Window{window}
	f.Window = window

	var cursorPrev, cursorCurr Cursor
	fps := GetFPS()

	for !window.ShouldClose() {
		if dbg.Profiling() {
			fps()
		}

		cursorCurr.X, cursorCurr.Y = window.GetCursorPos()
		if cursorCurr.X != cursorPrev.X || cursorCurr.Y != cursorPrev.Y {
			cursorPrev = cursorCurr // xPrev, yPrev = xCurr, yCurr
			f.cb.CursorMove(cursorCurr)
		}

		f.renderGL()

		window.SwapBuffers()
		glfw.PollEvents()

	}

}

// RegisterCallback allows setting user defined callbacks
func (f *Framework) RegisterCallback(cbHandler interface{}) {
	dbg.Log("lib.go", "Framework.RegisterCallback()")
	// f := Callbacks{}

	this := reflect.TypeOf(cbHandler)
	render := reflect.TypeOf(f.cb.Render)
	cursorMove := reflect.TypeOf(f.cb.CursorMove)
	key := reflect.TypeOf(f.cb.Key)
	fps := reflect.TypeOf(f.cb.FPS)

	if this == render {
		fmt.Println("render")
		f.cb.Render, _ = cbHandler.(func() *image.RGBA)
	}

	if this == cursorMove {
		fmt.Println("cursor")
		f.cb.CursorMove = cbHandler.(func(c Cursor))
	}

	if this == key {
		fmt.Println("key")
		f.cb.Key = cbHandler.(func(w *Window, k Key, scancode int, a Action, mk ModifierKey))
	}

	if this == fps {
		fmt.Println("fps")
		f.cb.FPS = cbHandler.(func(fps int))
		f.cb.FPS(1234)
	}
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

// Callbacks holds the callbacks defined in the User Application  ran in the library
type Callbacks struct {
	Render     func() *image.RGBA
	CursorMove func(c Cursor)
	Key        func(w *Window, k Key, scancode int, a Action, mk ModifierKey)
	FPS        func(fps int)
}

func onRenderNil() *image.RGBA {
	dbg.Log("lib.go", "onRenderNil()")
	return BlankImage()
}
func onCursorMoveNil(c Cursor) {
	dbg.Log("lib.go", "onCursorMoveNil")

}
func onKeyNil(w *Window, k Key, scancode int, a Action, mk ModifierKey) {
	dbg.Log("lib.go", "onKeyNil")

}
func onFPSNil(fps int) {
	dbg.Log("lib.go", "onFPSNil")

}
