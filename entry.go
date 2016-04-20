package gowindow

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

var debug = false
var frameToggle = false
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
		// x := Window{*w}
		// cb.Key(&x)
		cb.Key(&Window{*w}, Key(key), scancode, Action(action), ModifierKey(mods))
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
