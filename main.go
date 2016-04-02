// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3.1 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const size = 512
const windowWidth = size
const windowHeight = size

// ********************************************************************
func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// ********************************************************************
func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err = gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	img, _ := getRGBAFromImage("image.png")
	texture, err := newTexture(img)
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
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	stride := int32(8)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, stride, gl.PtrOffset(0))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		gl.UseProgram(program)

		gl.BindVertexArray(vao)

		// Load the texture
		rgbaFile, _ := getRGBAFromImage("image.png")
		rgba := randomizeRGBA(rgbaFile)
		texture, err = newTexture(rgba)
		// texture, err = newTexture(rgbaFile)

		if err != nil {
			log.Fatalln(err)
		}

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 2*3)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

// ********************************************************************
func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

// ********************************************************************
func compileShader(source string, shaderType uint32) (uint32, error) {

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// ********************************************************************
func getRGBAFromImage(filename string) (*image.RGBA, error) {

	file := filename
	imgFile, err := os.Open(file)
	if err != nil {
		// return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
		return nil, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, err
}

// ********************************************************************
func randomizeRGBA(rgba *image.RGBA) *image.RGBA {

	size := 512
	// blue := color.RGBA{0, 0, 255, 255}
	// draw.Draw(rgba, rgba.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	for i := 0; i < 200000; i++ {

		x := rand.Intn(size)
		y := rand.Intn(size)

		r := uint8(rand.Intn(255))
		g := uint8(rand.Intn(255))
		b := uint8(rand.Intn(255))

		rgba.Set(x, y, color.RGBA{r, g, b, 255})
	}

	// for x := 0; x < size; x++ {
	// 	for y := 0; y < size; y++ {
	// 		r := uint8(rand.Intn(255))
	// 		g := uint8(rand.Intn(255))
	// 		b := uint8(rand.Intn(255))
	// 		rgba.Set(x, y, color.RGBA{r, g, b, 255})
	// 	}
	// }

	return rgba
}

// ********************************************************************
func newTexture(rgba *image.RGBA) (uint32, error) {

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

// ********************************************************************

var vertexShader = `
#version 330
in vec2 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main() {
fragTexCoord = vertTexCoord;
gl_Position = vec4((vert.x - 0.5) * 2.0, (vert.y - 0.5) * 2.0, 0, 1.0);
}
` + "\x00"

var fragmentShader = `
#version 330
uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Front
	-1.0, -1.0, //0.0, 1.0, 0.0,
	1.0, -1.0, //0.0, 0.0, 0.0,
	-1.0, 1.0, //0.0, 1.0, 1.0,
	1.0, -1.0, //0.0, 0.0, 0.0,
	1.0, 1.0, //0.0, 0.0, 1.0,
	-1.0, 1.0, //0.0, 1.0, 1.0,

	// -1.0, -1.0, 1.0, 1.0, 0.0,
	// 1.0, -1.0, 1.0, 0.0, 0.0,
	// -1.0, 1.0, 1.0, 1.0, 1.0,
	// 1.0, -1.0, 1.0, 0.0, 0.0,
	// 1.0, 1.0, 1.0, 0.0, 1.0,
	// -1.0, 1.0, 1.0, 1.0, 1.0,
}
