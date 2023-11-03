package gl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/graphics/render/gl/shader"
	"github.com/eliiasg/deltawing/internal/rendering/shader_sources"
	"github.com/eliiasg/glow/enum"
)

type procedureBuilder struct {
	cxt     Context
	sb      *shader.ShaderBuilder
	version string
}

func (r *Renderer) MakeProcedureBuilder() render.ProcedureBuilder {
	source := shader.ShaderSource{
		SourceCode:     shader_sources.VertexBaseSource,
		LayoutStartPos: 2,
		Version:        r.version,
		Variables: []shader.Variable{
			{Name: "pos", Type: render.Type(render.ShaderFloat, 2), DefaultValue: ""},
			{Name: "layer", Type: render.Type(render.ShaderUnsignedInt, 1), DefaultValue: ""},
			{Name: "color", Type: render.Type(render.ShaderInt, 3), DefaultValue: "aColor.rgb"},
			{Name: "xAxis", Type: render.Type(render.ShaderFloat, 2), DefaultValue: "vec2(1, 0)"},
			{Name: "yAxis", Type: render.Type(render.ShaderFloat, 2), DefaultValue: "vec2(0, 1)"},
		},
	}
	return &procedureBuilder{r.cxt, shader.NewShaderBuilder(source), r.version}
}

func (p *procedureBuilder) AddAttributeChannel(shaderType render.ShaderType) render.Channel {
	return p.sb.AddAttributeChannel(shaderType)
}

func (p *procedureBuilder) AddIntermediateChannel(shaderType render.ShaderType, expression string) render.Channel {
	return p.sb.AddIntermediateChannel(shaderType, expression)
}

func (p *procedureBuilder) AddOperationChannel(shaderType render.ShaderType) render.Channel {
	return p.sb.AddOperationChannel(shaderType)
}

func (p *procedureBuilder) CallFunction(function *render.Function, channels ...render.Channel) error {
	return p.sb.CallFunction(function, channels...)
}

func (p *procedureBuilder) SetColorChannel(channel render.Channel) error {
	return p.sb.SetOutputChannel("color", channel)
}

func (p *procedureBuilder) SetLayerChannel(channel render.Channel) error {
	return p.sb.SetOutputChannel("layer", channel)
}

func (p *procedureBuilder) SetPositionChannel(channel render.Channel) error {
	return p.sb.SetOutputChannel("pos", channel)
}

func (p *procedureBuilder) SetXAxisChannel(channel render.Channel) error {
	return p.sb.SetOutputChannel("xAxis", channel)
}

func (p *procedureBuilder) SetYAxisChannel(channel render.Channel) error {
	return p.sb.SetOutputChannel("yAxis", channel)
}

func (p *procedureBuilder) Finish() (render.Procedure, error) {
	vertSource, attribTypes, uniformNames, err := p.sb.Finish()
	if err != nil {
		return nil, err
	}
	return compileProgram(p.cxt, p.version, vertSource, attribTypes, uniformNames)
}

type Procedure struct {
	render.ProcedureIdentifier
	cxt Context
	// Shader program object
	Prog any
	// Uniform location of screen size
	ScreenSizeLocation any
	// Attribute channels
	AttribChannels map[render.Channel]shader.AttribChannelInfo
	// Uniform locations
	UniformLocations map[string]any
}

func (p *Procedure) Free() {
	p.cxt.DeleteProgram(p.Prog)
}

func compileProgram(cxt Context, version string, vertSource string, attribTypes map[render.Channel]shader.AttribChannelInfo, uniformNames []string) (render.Procedure, error) {
	// vertex shader
	vert, err := compileShader(cxt, enum.VERTEX_SHADER, vertSource)
	if err != nil {
		return nil, err
	}
	// fragment shader
	// FIXME fragment shader is not changed, so maybe only compile once
	fragSource := strings.Replace(shader_sources.FragmentSource, "<version>", version, 1)
	frag, err := compileShader(cxt, enum.FRAGMENT_SHADER, fragSource)
	if err != nil {
		return nil, err
	}
	// program
	prog, err := createProgram(cxt, vert, frag)
	if err != nil {
		return nil, err
	}
	// delete shaders (program is not deleted)
	cxt.DeleteShader(vert)
	cxt.DeleteShader(frag)
	// get size uniform location
	sizeLoc := cxt.GetUniformLocation(prog, "screenSize")
	// get uniform locations
	uniformLocations := make(map[string]any)
	for _, name := range uniformNames {
		uniformLocations[name] = cxt.GetUniformLocation(prog, name)
	}
	return &Procedure{Prog: prog, ScreenSizeLocation: sizeLoc, AttribChannels: attribTypes, UniformLocations: uniformLocations}, nil
}

func createProgram(cxt Context, vertShader, fragShader any) (any, error) {
	// make and link
	prog := cxt.CreateProgram()
	cxt.AttachShader(prog, vertShader)
	cxt.AttachShader(prog, fragShader)
	cxt.LinkProgram(prog)

	// ckeck for error
	status := cxt.GetProgramParameter(prog, enum.LINK_STATUS)
	if status == enum.FALSE {
		// get error
		log := cxt.GetProgramInfoLog(prog)
		return 0, errors.New(fmt.Sprintf("Program linking failed:\n%v", log))
	}
	return prog, nil
}

func compileShader(cxt Context, typ uint32, source string) (any, error) {
	// make and compile shader
	shader := cxt.CreateShader(typ)
	cxt.ShaderSource(shader, source)
	cxt.CompileShader(shader)
	// check for error
	err := getShaderError(cxt, shader)
	if err != "" {
		return 0, errors.New(fmt.Sprintf("Failed to compile shader:\n%v\nWith error:\n%v", source, err))
	}
	return shader, nil
}

func getShaderError(cxt Context, shader any) string {
	status := cxt.GetShaderParameter(shader, enum.COMPILE_STATUS)
	if status == enum.FALSE {
		return cxt.GetShaderInfoLog(shader)
	}
	return ""
}
