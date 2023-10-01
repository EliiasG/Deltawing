package deltawing

import (
	"github.com/eliiasg/deltawing/desktop"
	"github.com/eliiasg/deltawing/internal/setup/glfw"
)

func NewProgram(width, height uint16, name string) desktop.Program {
	return glfw.NewProgram(width, height, name)
}
