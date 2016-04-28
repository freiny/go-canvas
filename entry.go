package gowindow

// Init initializes OpenGL/GLFW then runs a render callback on each iteration
// func Init(wc WinConfig, cbUserDefined Callbacks) {
// func Init(wc WinConfig) {
//
// 	if cb.Render == nil {
// 		cb.Render = onRenderNil
// 	}
// 	if cb.CursorMove == nil {
// 		cb.CursorMove = onCursorMoveNil
// 	}
// 	if cb.Key == nil {
// 		cb.Key = onKeyNil
// 	}
// 	if cb.FPS == nil {
// 		cb.FPS = onFPSNil
// 	}
//
// 	if err := glfw.Init(); err != nil {
// 		log.Fatalln("failed to initialize glfw:", err)
// 	}
// 	defer glfw.Terminate()
//
// 	glfw.WindowHint(glfw.Resizable, glfw.False)
// 	glfw.WindowHint(glfw.ContextVersionMajor, 4)
// 	glfw.WindowHint(glfw.ContextVersionMinor, 1)
// 	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
// 	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
// 	window, err := glfw.CreateWindow(wc.W, wc.H, "", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	window.MakeContextCurrent()
//
// 	window.SetPos(wc.X, wc.Y)
//
// 	initGL()
//
// 	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 		cb.Key(&Window{*w}, Key(key), scancode, Action(action), ModifierKey(mods))
// 	})
//
// 	window.SetCursorPos(0.0, 0.0)
//
// 	var cursorPrev, cursorCurr Cursor
// 	fps := GetFPS()
//
// 	for !window.ShouldClose() {
// 		if debug {
// 			fps()
// 		}
//
// 		cursorCurr.X, cursorCurr.Y = window.GetCursorPos()
// 		if cursorCurr.X != cursorPrev.X || cursorCurr.Y != cursorPrev.Y {
// 			cursorPrev = cursorCurr // xPrev, yPrev = xCurr, yCurr
// 			cb.CursorMove(cursorCurr)
// 		}
//
// 		renderGL()
//
// 		window.SwapBuffers()
// 		glfw.PollEvents()
//
// 	}
// }
