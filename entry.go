package gowindow

import (
	"fmt"
	"image"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type cbRender func() *image.RGBA
type cbCursorMove func(float64, float64)

// type cbKeyGLFW func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
type cbKey func(w *Window)
type cbFPS func(fps int)

// Callbacks holds the callbacks defined in the User Application  ran in the library
type Callbacks struct {
	Render     cbRender
	CursorMove cbCursorMove
	Key        cbKey
	FPS        cbFPS
}

var cb = Callbacks{}

// Config holds global data (e.g. window dimensions, cursor location)
type Config struct {
	Width  int
	Height int
	X      int
	Y      int
}

// Action hides the glfw from the User
type Action glfw.Action

// Window wraps glfw window
type Window struct {
	glfw.Window
}

// var onKey func()
// func cbKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	fmt.Println(w, key, scancode, action, mods)
// 	onKey()
// }
// w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey

func init() {
	runtime.LockOSThread()
}

// Init initializes OpenGL/GLFW then runs a render callback on each iteration of
// the library's Render Loop. Allows render function to be defined externally
// inside a user's application.
func Init(config Config, cbUserDefined Callbacks) {
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
	window, err := glfw.CreateWindow(config.Width, config.Height, "", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetPos(config.X, config.Y)

	if err = gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	texture, err := newTexture(cb.Render())
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(screenVertices)*4, gl.Ptr(screenVertices), gl.STATIC_DRAW)

	stride := int32(8)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, stride, gl.PtrOffset(0))

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// var x *GWWindow = &w
		x := Window{*w}
		// fmt.Println(x)
		cb.Key(&x)
	})

	window.SetCursorPos(0.0, 0.0)

	var xPrev, yPrev, xCurr, yCurr float64
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

		xCurr, yCurr = window.GetCursorPos()
		if xCurr != xPrev || yCurr != yPrev {
			xPrev, yPrev = xCurr, yCurr
			cb.CursorMove(xCurr, yCurr)
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.BindVertexArray(vao)

		texture, err = newTexture(cb.Render())

		if err != nil {
			log.Fatalln(err)
		}

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 2*3)

		window.SwapBuffers()
		glfw.PollEvents()

	}
}
