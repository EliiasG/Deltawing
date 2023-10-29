package gl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/graphics/render/gl/shader"
	"github.com/eliiasg/deltawing/internal/rendering/shader_sources"
	"github.com/eliiasg/glow/v3.3-core/gl"
)

type procedureBuilder struct {
	sb *shader.ShaderBuilder
}

func (r *Renderer) MakeProcedureBuilder() render.ProcedureBuilder {
	source := shader.ShaderSource{
		SourceCode:     shader_sources.VertexBaseSource,
		LayoutStartPos: 2,
		Version:        "#version 330 core",
		Variables: []shader.Variable{
			{Name: "pos", Type: render.Type(render.ShaderFloat, 2), DefaultValue: ""},
			{Name: "layer", Type: render.Type(render.ShaderUnsignedInt, 1), DefaultValue: ""},
			{Name: "color", Type: render.Type(render.ShaderInt, 3), DefaultValue: "aColor.rgb"},
			{Name: "xAxis", Type: render.Type(render.ShaderFloat, 2), DefaultValue: "vec2(1, 0)"},
			{Name: "yAxis", Type: render.Type(render.ShaderFloat, 2), DefaultValue: "vec2(0, 1)"},
		},
	}
	return &procedureBuilder{shader.NewShaderBuilder(source)}
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
	vertSource, attribTypes, err := p.sb.Finish()
	if err != nil {
		return nil, err
	}
	return compileProgram(vertSource, attribTypes)
}

type procedure struct {
	render.ProcedureIdentifier
	progID             uint32
	screenSizeLocation int32
	attribChannels     map[render.Channel]shader.AttribChannelInfo
}

func (p *procedure) Free() {
	gl.DeleteProgram(p.progID)
}

// will be called by the ShaderBuilder
func compileProgram(vertSource string, attribTypes map[render.Channel]shader.AttribChannelInfo) (render.Procedure, error) {
	// vertex shader
	vert, err := compileShader(gl.VERTEX_SHADER, vertSource)
	if err != nil {
		return nil, err
	}
	// fragment shader
	// FIXME fragment shader is not generated at runtime, so maybe only compile once
	frag, err := compileShader(gl.FRAGMENT_SHADER, shader_sources.FragmentSource)
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
	// get size uniform location
	sizeLoc := gl.GetUniformLocation(prog, gl.Str("screenSize\x00"))
	return &procedure{progID: prog, screenSizeLocation: sizeLoc, attribChannels: attribTypes}, nil
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
	gl.CompileShader(shader)
	// check for error
	err := getShaderError(shader)
	if err != "" {
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
