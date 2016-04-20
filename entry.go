package gowindow

import (
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
)

var wc = WinConfig{}
var cb = Callbacks{}

func init() {
	runtime.LockOSThread()
}

// Init initializes OpenGL/GLFW then runs a render callback on each iteration
func Init(wc WinConfig, cbUserDefined Callbacks) {
	cb = cbUserDefined

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(wc.W, wc.H, "", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetPos(wc.X, wc.Y)

	initGL()

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		cb.Key(&Window{*w}, Key(key), scancode, Action(action), ModifierKey(mods))
	})

	window.SetCursorPos(0.0, 0.0)

	var cursorPrev, cursorCurr Cursor
	var tPrev, tCurr, tDelta int64
	var fps int

	for !window.ShouldClose() {
		if debug {

			tCurr = GetTime()
			tDelta = tCurr - tPrev
			if tDelta > 1000000000 {
				cb.FPS(fps)
				tPrev = tCurr
				fps = 0

			} else {
				fps++
			}

		}

		cursorCurr.X, cursorCurr.Y = window.GetCursorPos()
		if cursorCurr.X != cursorPrev.X || cursorCurr.Y != cursorPrev.Y {
			cursorPrev = cursorCurr
			// xPrev, yPrev = xCurr, yCurr
			cb.CursorMove(cursorCurr)
		}

		renderGL()

		window.SwapBuffers()
		glfw.PollEvents()

	}
}
