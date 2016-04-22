package gowindow

import "time"

var debug = false
var frameToggle = false

// Cursor for storing the current cursor position
type Cursor struct {
	X float64
	Y float64
}

// SetDebug turn debugging mode on or off
func SetDebug() {
	debug = true
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
