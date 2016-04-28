package gowindow

import "fmt"

var dbg = Debugger{
	isRunning:   true,
	isLogging:   true,
	isProfiling: true,
}

// Debugger general purpose logging and profiling
type Debugger struct {
	isRunning   bool
	isLogging   bool
	isProfiling bool
}

// Init set debugger runnng on start up
func (d Debugger) Init() Debugger {
	d.Log("debug.go", "Debugger.Init()")
	d.isRunning = true
	d.isLogging = true
	d.isProfiling = true
	return d
}

// Log prints to terminal
func (d *Debugger) Log(i ...interface{}) {
	if d.Logging() {
		fmt.Print("[LOG] ")
		fmt.Println(i...)
	}
}

// Running checks if debugger is on
func (d *Debugger) Running() bool {
	return d.isRunning
}

// Logging checks if logging is on
func (d *Debugger) Logging() bool {
	return d.isRunning && d.isLogging
}

// Profiling checks if profiling is on
func (d *Debugger) Profiling() bool {
	return d.isRunning && d.isProfiling
}

// // SetDebug turn debugging mode on or off
// func SetDebug() {
// 	debug = true
// }
