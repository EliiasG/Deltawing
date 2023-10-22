package desktop

import (
	"github.com/eliiasg/deltawing/desktop/program"
	"github.com/eliiasg/deltawing/internal/setup/glfw"
)

func NewProgram(width, height uint16, name string) program.Program {
	return glfw.NewProgram(width, height, name)
}
