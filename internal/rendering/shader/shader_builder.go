package shader

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	r "github.com/eliiasg/deltawing/graphics/render"
)

type glslChannel struct {
	r.ChannelIdentifier
	id      uint16
	varType r.ShaderType
}

type funcCall struct {
	fun *r.Function
	// only the id of the channels are required
	params []uint16
}

type interChannel struct {
	*glslChannel
	expr string
}

func getGLChannel(channel r.Channel) *glslChannel {
	switch c := channel.(type) {
	case *glslChannel:
		return c
	case *interChannel:
		return c.glslChannel
	}
	return nil
}

func ChannelName(channel r.Channel) string {
	return getVarName(getGLChannel(channel).id)
}

func ChannelType(channel r.Channel) r.ShaderType {
	return getGLChannel(channel).varType
}

func getVarName(id uint16) string {
	// Beatiful
	// The only way for this to fail should be if another renderer is used to supply the channel
	return "c" + strconv.Itoa(int(id))
}

// just collects the data, shader will be composed at end
type shaderBuilder struct {
	baseSource  string
	version     string
	done        func(string, []r.ShaderType) (r.Procedure, error)
	chanID      uint16
	interChans  []*interChannel
	attribChans []*glslChannel
	operChans   []*glslChannel
	calls       []funcCall
	// start position of layout
	startPos uint8
	// IDs for output channels, only storing ID as it is the only requirement for getting the var name
	posID   uint16
	layerID uint16
	xAxisID uint16
	yAxisID uint16
	colorID uint16
}

/*
The baseSource is a shader that will be added parts to.
The following keywords will be replaced:
'<version>' version and possiblr precision calls, should only be used once
'<attributes>' instance data, should only be used once
'<uniforms>' uniforms, should only be used once
'<functions>' function declarations, should only be used once
'<variables>' variables, should only be used once
'<calls>' function calls, should only be used once
'<pos>' output position variable, this should not be modified
'<layer>' output layer variable, this should not be modified
'<xAxis>' output xAxis variable, this should not be modified
'<yAxis>' output xAxis variable, this should not be modified
*/
func NewShaderBuilder(baseSource string, layoutStartPos uint8, version string, done func(string, []r.ShaderType) (r.Procedure, error)) r.ProcedureBuilder {
	return &shaderBuilder{
		baseSource: baseSource,
		done:       done,
		// start index of channels, 0 is used unset channels
		chanID:      1,
		interChans:  make([]*interChannel, 0),
		attribChans: make([]*glslChannel, 0),
		operChans:   make([]*glslChannel, 0),
		calls:       make([]funcCall, 0),
		startPos:    layoutStartPos,
		version:     version,
	}
}

func (s *shaderBuilder) makeChannel(varType r.ShaderType) *glslChannel {
	channel := &glslChannel{
		id:      s.chanID,
		varType: varType,
	}
	s.chanID++
	return channel
}

func (s *shaderBuilder) AddIntermediateChannel(shaderType r.ShaderType, expression string) r.Channel {
	channel := s.makeChannel(shaderType)
	s.interChans = append(s.interChans, &interChannel{
		channel,
		expression,
	})
	return channel
}

func (s *shaderBuilder) AddAttributeChannel(shaderType r.ShaderType) r.Channel {
	channel := s.makeChannel(shaderType)
	s.attribChans = append(s.attribChans, channel)
	return channel
}

func (s *shaderBuilder) AddOperationChannel(shaderType r.ShaderType) r.Channel {
	channel := s.makeChannel(shaderType)
	s.operChans = append(s.operChans, channel)
	return channel
}

func (s *shaderBuilder) CallFunction(function *r.Function, channels ...r.Channel) error {
	if len(function.Parameters) != len(channels) {
		return errors.New("amount of parameters given must match amount of parameters expected")
	}

	call := funcCall{
		fun:    function,
		params: make([]uint16, 0),
	}

	// set params of call, and check for wrong parameters
	for i, channel := range channels {
		glChan := getGLChannel(channel)
		param := function.Parameters[i]
		if glChan.varType != param {
			typ := glChan.varType
			return errors.New(fmt.Sprintf("Expected %v, %v, but got %v, %v", param.Type, param.Amount, typ.Type, typ.Amount))
		}
		call.params = append(call.params, glChan.id)
	}

	s.calls = append(s.calls, call)

	return nil
}

// hacky helper method for the output setters
func setOutputChannel(channel r.Channel, typ r.ShaderType) (uint16, error) {
	glChan := getGLChannel(channel)
	if glChan.varType != typ {
		return 0, errors.New("Invalid output type")
	}
	return glChan.id, nil
}

func (s *shaderBuilder) SetPositionChannel(channel r.Channel) (e error) {
	s.posID, e = setOutputChannel(channel, r.Type(r.ShaderFloat, 2))
	return
}

func (s *shaderBuilder) SetLayerChannel(channel r.Channel) (e error) {
	s.layerID, e = setOutputChannel(channel, r.Type(r.ShaderUnsignedInt, 1))
	return
}

func (s *shaderBuilder) SetXAxisChannel(channel r.Channel) (e error) {
	s.xAxisID, e = setOutputChannel(channel, r.Type(r.ShaderFloat, 2))
	return
}

func (s *shaderBuilder) SetYAxisChannel(channel r.Channel) (e error) {
	s.yAxisID, e = setOutputChannel(channel, r.Type(r.ShaderFloat, 2))
	return
}

func (s *shaderBuilder) SetColorChannel(channel r.Channel) (e error) {
	s.colorID, e = setOutputChannel(channel, r.Type(r.ShaderInt, 3))
	return
}

var typeMap = map[r.ChannelShaderType][2]string{
	r.ShaderFloat:       {"float", "vec"},
	r.ShaderInt:         {"int", "ivec"},
	r.ShaderUnsignedInt: {"uint", "uvec"},
}

func getGLSLTypeName(typ r.ShaderType) string {
	if typ.Amount == 1 {
		return typeMap[typ.Type][0]
	}
	return typeMap[typ.Type][1] + strconv.Itoa(int(typ.Amount))
}

func (s *shaderBuilder) makeAtrribSection() string {
	var sb strings.Builder
	for i, channel := range s.attribChans {
		sb.WriteString(fmt.Sprintf("layout(location=%v) in %v %v;\n", i+int(s.startPos), getGLSLTypeName(channel.varType), ChannelName(channel)))
	}
	return sb.String()
}

func (s *shaderBuilder) makeUniformSection() string {
	var sb strings.Builder
	for _, channel := range s.operChans {
		sb.WriteString(fmt.Sprintf("uniform %v %v;\n", getGLSLTypeName(channel.varType), ChannelName(channel)))
	}
	return sb.String()
}

func (s *shaderBuilder) getFunctions() []*r.Function {
	// make set of functions
	funcSet := make(map[*r.Function]bool, 0)
	for _, call := range s.calls {
		funcSet[call.fun] = true
	}
	// convert set to list
	keys := make([]*r.Function, 0, len(funcSet))
	for f := range funcSet {
		keys = append(keys, f)
	}
	return keys
}

func (s *shaderBuilder) makeDeclarationSection() string {
	var sb strings.Builder
	for _, function := range s.getFunctions() {
		sb.WriteString(function.Source + "\n")
	}
	return sb.String()
}

func (s *shaderBuilder) makeVariablesSection() string {
	var sb strings.Builder
	for _, variable := range s.interChans {
		// '<type> <name>'
		sb.WriteString(fmt.Sprintf("%v %v", getGLSLTypeName(variable.varType), ChannelName(variable)))
		if variable.expr != "" {
			// '= <expression>'
			sb.WriteString("= " + variable.expr)
		}
		sb.WriteString(";\n")
	}
	return sb.String()
}

func (s *shaderBuilder) makeCallsSection() string {
	var sb strings.Builder
	// calls
	for _, call := range s.calls {
		sb.WriteString(call.fun.Name + "(")
		// parameters
		for i, param := range call.params {
			r := getVarName(param)
			if i < len(call.params)-1 {
				r += ", "
			}
			sb.WriteString(r)
		}
		sb.WriteString(");\n")
	}
	return sb.String()
}

// takes a shader to then replace section keywords with generated code
// it expects a shader because other edits might be made to the baseSource before this
func (s *shaderBuilder) composeSections(shader string) (string, error) {
	sections := [][2]string{
		{"<version>", s.version},
		{"<attributes>", s.makeAtrribSection()},
		{"<uniforms>", s.makeUniformSection()},
		{"<functions>", s.makeDeclarationSection()},
		{"<variables>", s.makeVariablesSection()},
		{"<calls>", s.makeCallsSection()},
	}

	for _, section := range sections {
		if count := strings.Count(shader, section[0]); count != 1 {
			return "", errors.New(fmt.Sprintf("Expected %v once, found it %v times", section[0], count))
		}
		// should only contain section[0] once
		shader = strings.Replace(shader, section[0], section[1], 1)
	}

	return shader, nil
}

func (s *shaderBuilder) composeVars(shader string) (string, error) {
	// assert builder
	if s.layerID == 0 || s.posID == 0 {
		return "", errors.New("Both layer and position must be set")
	}

	// get x and y axis values
	xAxis, yAxis := "vec2(1, 0)", "vec2(0, 1)"

	if s.xAxisID != 0 {
		xAxis = getVarName(s.xAxisID)
	}
	if s.yAxisID != 0 {
		yAxis = getVarName(s.yAxisID)
	}

	// get color value
	color := "aColor.rgb"
	if s.colorID != 0 {
		color = getVarName(s.colorID)
	}

	// compose
	shader = strings.ReplaceAll(shader, "<xAxis>", xAxis)
	shader = strings.ReplaceAll(shader, "<yAxis>", yAxis)
	shader = strings.ReplaceAll(shader, "<pos>", getVarName(s.posID))
	shader = strings.ReplaceAll(shader, "<layer>", getVarName(s.layerID))
	shader = strings.ReplaceAll(shader, "<color>", color)

	return shader, nil
}

func (s *shaderBuilder) composeShader() (string, error) {
	// sections
	shader, err := s.composeSections(s.baseSource)
	if err != nil {
		return "", err
	}
	// vars
	shader, err = s.composeVars(shader)
	if err != nil {
		return "", err
	}
	// make sure to properly end shader with escape char
	if shader[len(shader)-1] != '\x00' {
		shader += "\x00"
	}

	return shader, nil
}

func (s *shaderBuilder) getAttribTypes() []r.ShaderType {
	res := make([]r.ShaderType, 0, len(s.attribChans)+2)
	for _, channel := range s.attribChans {
		res = append(res, channel.varType)
	}
	return res
}

func (s *shaderBuilder) Finish() (r.Procedure, error) {
	shader, err := s.composeShader()
	if err != nil {
		return nil, err
	}
	return s.done(shader, s.getAttribTypes())
}
