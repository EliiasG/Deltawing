package opengl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/internal/rendering/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func (r *renderer) MakeProcedureBuilder() render.ProcedureBuilder {
	return shader.NewShaderBuilder(shader.VertexBaseSource, 2, compileProgram)
}

type procedure struct {
	render.ProcedureIdentifier
	id uint32
}

func (p *procedure) Free() {
	gl.DeleteProgram(p.id)
}

// will be called by the ShaderBuilder
func compileProgram(vertSource string) (render.Procedure, error) {
	// vertex shader
	vert, err := compileShader(gl.VERTEX_SHADER, vertSource)
	if err != nil {
		return nil, err
	}
	// fragment shader
	// FIXME fragment shader is not generated, so maybe only compile once
	frag, err := compileShader(gl.FRAGMENT_SHADER, shader.FragmentSource)
	if err != nil {
		return nil, err
	}
	// program
	prog, err := createProgram(vert, frag)
	if err != nil {
		return nil, err
	}
	// delete shaders (program is not deleted)
	gl.DeleteShader(vert)
	gl.DeleteShader(frag)
	return &procedure{id: prog}, nil
}

func createProgram(vertShader, fragShader uint32) (uint32, error) {
	// make and link
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertShader)
	gl.AttachShader(prog, fragShader)
	gl.LinkProgram(prog)

	// ckeck for error
	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		// get log length
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)

		// get error
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(log))

		return 0, errors.New(fmt.Sprintf("Program linking failed:\n%v", log))
	}
	return prog, nil
}

func compileShader(typ uint32, source string) (uint32, error) {
	// make and compile shader
	shader := gl.CreateShader(typ)
	cSource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cSource, nil)
	free()
	// check for error
	err := getShaderError(shader)
	if err != "" {
		gl.DeleteShader(shader)
		return 0, errors.New(fmt.Sprintf("Failed to compile shader:\n%v\nWith error:\n%v", source, err))
	}
	return shader, nil
}

func getShaderError(shader uint32) string {
	var status int32
	// get status
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		// get log length
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		// get log
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return log
	}
	return ""
}

func (r *renderer) MakeFragmentShader(source string) render.FragmentShader {
	panic("not implemented") // TODO: Implement
}
