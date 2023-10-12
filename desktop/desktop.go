package desktop

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/input"
)

type Program interface {
	Renderer() render.Renderer
	Window() Window
	Terminate()
}

type Window interface {
	SetSize(width, height uint16)
	SetMaximized(maximized bool)
	SetFullScreen(fullscreen bool)
	SetSizeChanged(handler func(uint16, uint16))
	ShouldClose() bool
	WindowSize() (uint16, uint16)
	UpdateView()

	// Input is local to window
	Keyboard() input.Keyboard
	Mouse() input.Mouse
	Controller() input.Controller
}
