package gowindow

import (
	"fmt"
	"time"
)

var frameToggle = false

// Cursor for storing the current cursor position
type Cursor struct {
	X float64
	Y float64
}

// SetFrameToggle alternates between returning true and false each frame
func SetFrameToggle() bool {
	frameToggle = !frameToggle
	return frameToggle
}

// GetTime returns current time in nano seconds
func GetTime() int64 {
	return time.Now().UTC().UnixNano()

}

func delay(n int) {
	for j := 0; j < n*1000; j++ {
		fmt.Print("")
	}
}

// GetFPS fps counter in closure
// usage:
// fps := GetFPS()
// RenderLoop() {
//   fps()
// }
func GetFPS() func() {
	isReady := false
	frames := 0
	var prev, curr int64
	return func() {
		if isReady {
			fmt.Println(frames)
			frames = 0
			prev = GetTime()
			isReady = false
		}

		curr = GetTime()

		if curr-prev >= 1000000000 {
			isReady = true
		}
		frames++
	}
}
