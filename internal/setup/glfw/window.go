//go:build cgo
// +build cgo

package glfw

import (
	"github.com/eliiasg/deltawing/input"
	"github.com/eliiasg/glow/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type window struct {
	glfwWin     *glfw.Window
	sizeHandler func(uint16, uint16)
	keyboard    *keyboard
}

func makeWindow(width, height uint16, name string) *window {
	// init glfw window
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)
	// for layer prececion
	glfw.WindowHint(glfw.DepthBits, 32)
	glfwWin, e := glfw.CreateWindow(int(width), int(height), name, nil, nil)
	if e != nil {
		panic("Window failed to init with following error: " + e.Error())
	}
	// init window
	win := &window{glfwWin, nil, makeKeyboard(glfwWin)}
	glfwWin.SetSizeCallback(win.sizeCallback)
	return win
}

func (w *window) SetSize(width uint16, height uint16) {
	w.glfwWin.SetSize(int(width), int(height))
}

func (w *window) SetMaximized(maximized bool) {
	if maximized {
		w.glfwWin.Maximize()
	} else {
		w.glfwWin.Restore()
	}
}

func (w *window) SetFullScreen(fullscreen bool) {
	if fullscreen {
		mon := getCurrentMonitor(w.glfwWin)
		// mode for some reason just means information about monitor
		mode := mon.GetVideoMode()
		w.glfwWin.SetMonitor(mon, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		w.glfwWin.Restore()
	}
}

func (w *window) SetSizeChanged(handler func(uint16, uint16)) {
	w.sizeHandler = handler
}

func (w *window) sizeCallback(_ *glfw.Window, width, height int) {
	if w.sizeHandler != nil {
		w.sizeHandler(uint16(width), uint16(height))
	}
}

func (w *window) WindowSize() (uint16, uint16) {
	width, height := w.glfwWin.GetSize()
	return uint16(width), uint16(height)
}

func (w *window) UpdateView() {
	w.glfwWin.SwapBuffers()
	glfw.PollEvents()
}

func (w *window) ShouldClose() bool {
	return w.glfwWin.ShouldClose()
}

func getCurrentMonitor(win *glfw.Window) *glfw.Monitor {
	winX, winY := win.GetPos()
	mons := glfw.GetMonitors()
	for _, mon := range mons {
		// monitor info
		mode := mon.GetVideoMode()
		monW, monH := mode.Width, mode.Height
		monX, monY := mon.GetPos()
		// amazing
		if winX >= monX && winY >= monY && winX < monX+monW && winY < monY+monH {
			return mon
		}
	}
	return mons[0]
}

func (w *window) Keyboard() input.Keyboard {
	return w.keyboard
}

func (w *window) Mouse() input.Mouse {
	panic("not implemented") // TODO: Implement
}

func (w *window) Controller() input.Controller {
	panic("not implemented") // TODO: Implement
}
