//go:build cgo
// +build cgo

package glfw

import (
	"runtime"

	"github.com/eliiasg/deltawing/desktop/program"
	"github.com/eliiasg/deltawing/graphics/render"
	g "github.com/eliiasg/deltawing/graphics/render/gl"
	"github.com/eliiasg/deltawing/internal/rendering/opengl"
	"github.com/eliiasg/glow/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type glfwProgram struct {
	window   *window
	renderer render.Renderer
}

func (p *glfwProgram) Renderer() render.Renderer {
	return p.renderer
}

func (p *glfwProgram) Window() program.Window {
	return p.window
}

func (p *glfwProgram) Terminate() {
	glfw.Terminate()
}

func (p *glfwProgram) Time() float64 {
	return glfw.GetTime()
}

func NewProgram(width, height uint16, name string) program.Program {
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
	gl.Enable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)
	gl.DepthFunc(gl.GREATER)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearDepth(0)

	return &glfwProgram{win, g.NewRenderer(
		func() uint16 {
			width, _ := win.WindowSize()
			return uint16(width)
		},
		func() uint16 {
			_, height := win.WindowSize()
			return uint16(height)
		},
		opengl.MakeContext(),
		"#version 330 core",
		false,
	)}
}
