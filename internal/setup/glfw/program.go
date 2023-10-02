package glfw

import (
	"runtime"

	"github.com/eliiasg/deltawing/desktop"
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/internal/rendering/gl/opengl"
	"github.com/go-gl/gl/v3.3-core/gl"
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
	runtime.LockOSThread()
	e := glfw.Init()
	if e != nil {
		panic("GLFW failed to init with following error: " + e.Error())
	}
	// init window
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)
	// for layer prececion
	glfw.WindowHint(glfw.DepthBits, 32)
	win, e := glfw.CreateWindow(int(width), int(height), name, nil, nil)
	if e != nil {
		panic("Window failed to init with following error: " + e.Error())
	}
	win.MakeContextCurrent()
	e = gl.Init()
	if e != nil {
		panic("OpenGL failed to init with following error: " + e.Error())
	}
	return &program{&window{win}, opengl.NewRenderer(
		func() uint16 {
			width, _ := win.GetSize()
			return uint16(width)
		},
		func() uint16 {
			_, height := win.GetSize()
			return uint16(height)
		},
	)}
}
