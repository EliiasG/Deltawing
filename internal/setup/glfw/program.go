package glfw

import (
	"runtime"

	"github.com/eliiasg/deltawing/desktop"
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/internal/rendering/gl/opengl"
	"github.com/eliiasg/glow/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type program struct {
	window   *window
	renderer render.Renderer
}

func (p *program) GetRenderer() render.Renderer {
	return p.renderer
}

func (p *program) GetWindow() desktop.Window {
	return p.window
}

func (p *program) Terminate() {
	glfw.Terminate()
}

func NewProgram(width, height uint16, name string) desktop.Program {
	// requird for OpenGL bindings (and i think also GLFW)
	runtime.LockOSThread()
	e := glfw.Init()
	if e != nil {
		panic("GLFW failed to init with following error: " + e.Error())
	}
	// window setup
	win := makeWindow(width, height, name)
	win.glfwWin.MakeContextCurrent()
	// OpenGL init must be called after MakeContextCurrent
	e = gl.Init()
	if e != nil {
		panic("OpenGL failed to init with following error: " + e.Error())
	}
	// OpenGL setup
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.MULTISAMPLE)
	gl.Disable(gl.CULL_FACE)
	gl.DepthFunc(gl.GREATER)
	gl.ClearDepth(0)

	return &program{win, opengl.NewRenderer(
		func() uint16 {
			width, _ := win.WindowSize()
			return uint16(width)
		},
		func() uint16 {
			_, height := win.WindowSize()
			return uint16(height)
		},
	)}
}
