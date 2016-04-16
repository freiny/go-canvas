package gowindow

import "time"

var debug = false

// SetDebug turn debugging mode on or off
func SetDebug() {
	debug = true
}

var frameToggle = false

// Toggle alternates between returning true and false each frame
func Toggle() bool {
	frameToggle = !frameToggle
	return frameToggle
}

// GetTime returns current time in nano seconds
func GetTime() int64 {
	return time.Now().UTC().UnixNano()
}
