package desktop

import (
	"github.com/eliiasg/deltawing/graphics/render"
)

type Program interface {
	GetRenderer() render.Renderer
	GetWindow() Window
	Terminate()
}

type Window interface {
	SetSize(width, height uint16)
	SetMaximized(maximized bool)
	SetFullScreen(fullscreen bool)
	SetSizeShanged(handler func(uint16, uint16))
	WindowSize() (uint16, uint16)
	UpdateView()
}
