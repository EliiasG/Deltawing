package app

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/input"
)

type App interface {
	Renderer() render.Renderer
	Time() float64
	Keyboard() input.Keyboard
	Mouse() input.Mouse
	Controller() input.Controller
}
