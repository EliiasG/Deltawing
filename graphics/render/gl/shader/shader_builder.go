package shader

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	r "github.com/eliiasg/deltawing/graphics/render"
)

type Channel struct {
	r.ChannelIdentifier
	id      uint16
	varType r.ShaderType
}

func (c *Channel) Name() string {
	return "c" + strconv.Itoa(int(c.id))
}

func (c *Channel) ShaderType() r.ShaderType {
	return c.varType
}

type funcCall struct {
	fun    *r.Function
	params []*Channel
}

type interChannel struct {
	*Channel
	expr string
}

func GLChannel(channel r.Channel) *Channel {
	switch c := channel.(type) {
	case *Channel:
		return c
	case *interChannel:
		return c.Channel
	}
	return nil
}

// just collects the data, shader will be composed at end
type ShaderBuilder struct {
	baseSource  string
	version     string
	chanID      uint16
	interChans  []*interChannel
	attribChans []*Channel
	operChans   []*Channel
	calls       []funcCall
	// start position of layout
	startPos uint8
	// storing channel names, since vars can have defaault values
	// maybe a bit hacky
	shaderVars map[string]*varValue
}

// A variable existing in a sahder
type Variable struct {
	// Variable name in shader (without <>)
	Name string
	// Type of variable
	Type r.ShaderType
	// Default value of variable, leave empty to make setting required
	DefaultValue string
}

type varValue struct {
	typ r.ShaderType
	val string
}

type ShaderSource struct {
	/*
		The SourceCode is a shader that will be added parts to.
		The following keywords will be replaced:
		'<version>' version and possible precision calls, should only be used once
		'<attributes>' instance data, should only be used once
		'<uniforms>' uniforms, should only be used once
		'<functions>' function declarations, should only be used once
		'<variables>' variables, should only be used once
		'<calls>' function calls, should only be used once
		'<varname>' for every variable added - can be used multiple times, but should not be written to
	*/
	SourceCode string
	// Start position of layout for generated attributes, in case attributes are predefined in the shader
	LayoutStartPos uint8
	// String to replace '<version>' in the shader with
	Version string
	// Variables used in the shader, should be without <> - keys are varname, values are default values leave value empty to make setting it required
	Variables []Variable
}

func NewShaderBuilder(source ShaderSource) *ShaderBuilder {
	shaderVars := make(map[string]*varValue)
	for _, variable := range source.Variables {
		if strings.ContainsAny(variable.Name, "<>") {
			panic("Error making ShaderBuilder, variable names must be without <> : " + variable.Name)
		}
		if variable.Name == "" {
			panic("Cannot make variable with empty name")
		}
		shaderVars[variable.Name] = &varValue{variable.Type, variable.DefaultValue}
	}
	return &ShaderBuilder{
		// start index of channels, 0 is used unset channels
		chanID:      1,
		interChans:  make([]*interChannel, 0),
		attribChans: make([]*Channel, 0),
		operChans:   make([]*Channel, 0),
		calls:       make([]funcCall, 0),
		startPos:    source.LayoutStartPos,
		version:     source.Version,
		baseSource:  source.SourceCode,
		shaderVars:  shaderVars,
	}
}

func (s *ShaderBuilder) makeChannel(varType r.ShaderType) *Channel {
	channel := &Channel{
		id:      s.chanID,
		varType: varType,
	}
	s.chanID++
	return channel
}

func (s *ShaderBuilder) AddIntermediateChannel(shaderType r.ShaderType, expression string) r.Channel {
	channel := s.makeChannel(shaderType)
	s.interChans = append(s.interChans, &interChannel{
		channel,
		expression,
	})
	return channel
}

func (s *ShaderBuilder) AddAttributeChannel(shaderType r.ShaderType) r.Channel {
	channel := s.makeChannel(shaderType)
	s.attribChans = append(s.attribChans, channel)
	return channel
}

func (s *ShaderBuilder) AddOperationChannel(shaderType r.ShaderType) r.Channel {
	channel := s.makeChannel(shaderType)
	s.operChans = append(s.operChans, channel)
	return channel
}

func (s *ShaderBuilder) CallFunction(function *r.Function, channels ...r.Channel) error {
	if len(function.Parameters) != len(channels) {
		return errors.New("amount of parameters given must match amount of parameters expected")
	}

	call := funcCall{
		fun:    function,
		params: make([]*Channel, 0),
	}

	// set params of call, and check for wrong parameters
	for i, channel := range channels {
		glChan := GLChannel(channel)
		param := function.Parameters[i]
		if glChan.varType != param {
			typ := glChan.varType
			return errors.New(fmt.Sprintf("Expected %v, %v, but got %v, %v", param.Type, param.Amount, typ.Type, typ.Amount))
		}
		call.params = append(call.params, glChan)
	}

	s.calls = append(s.calls, call)

	return nil
}

func (s *ShaderBuilder) SetOutputChannel(varName string, channel r.Channel) error {
	glChan := GLChannel(channel)
	if glChan.ShaderType() != s.shaderVars[varName].typ {
		return errors.New("Invalid type for " + varName)
	}
	s.shaderVars[varName].val = glChan.Name()
	return nil
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

func (s *ShaderBuilder) makeAtrribSection() string {
	var sb strings.Builder
	for i, channel := range s.attribChans {
		sb.WriteString(fmt.Sprintf("layout(location=%v) in %v %v;\n", i+int(s.startPos), getGLSLTypeName(channel.varType), channel.Name()))
	}
	return sb.String()
}

func (s *ShaderBuilder) makeUniformSection() string {
	var sb strings.Builder
	for _, channel := range s.operChans {
		sb.WriteString(fmt.Sprintf("uniform %v %v;\n", getGLSLTypeName(channel.varType), channel.Name()))
	}
	return sb.String()
}

func (s *ShaderBuilder) getFunctions() []*r.Function {
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

func (s *ShaderBuilder) makeDeclarationSection() string {
	var sb strings.Builder
	for _, function := range s.getFunctions() {
		sb.WriteString(function.Source + "\n")
	}
	return sb.String()
}

func (s *ShaderBuilder) makeVariablesSection() string {
	var sb strings.Builder
	for _, variable := range s.interChans {
		// '<type> <name>'
		sb.WriteString(fmt.Sprintf("%v %v", getGLSLTypeName(variable.varType), variable.Name()))
		if variable.expr != "" {
			// '= <expression>'
			sb.WriteString("= " + variable.expr)
		}
		sb.WriteString(";\n")
	}
	return sb.String()
}

func (s *ShaderBuilder) makeCallsSection() string {
	var sb strings.Builder
	// calls
	for _, call := range s.calls {
		sb.WriteString(call.fun.Name + "(")
		// parameters
		for i, param := range call.params {
			r := param.Name()
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
func (s *ShaderBuilder) composeSections(shader string, oldnew *[]string) error {
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
			return errors.New(fmt.Sprintf("Expected %v once, found it %v times", section[0], count))
		}
		// should only contain section[0] once
		*oldnew = append(*oldnew, section[0], section[1])
	}

	return nil
}

func (s *ShaderBuilder) composeVars(oldnew *[]string) error {
	for name, value := range s.shaderVars {
		if value.val == "" {
			return errors.New("Variable '<" + name + ">' must be set")
		}
		*oldnew = append(*oldnew, "<"+name+">", value.val)
	}

	return nil
}

func (s *ShaderBuilder) Finish() (string, map[r.Channel]AttribChannelInfo, []string, error) {
	oldnew := make([]string, 0)
	// sections
	err := s.composeSections(s.baseSource, &oldnew)
	if err != nil {
		return "", nil, nil, err
	}
	// vars
	err = s.composeVars(&oldnew)
	if err != nil {
		return "", nil, nil, err
	}
	shader := strings.NewReplacer(oldnew...).Replace(s.baseSource)
	// make sure to properly end shader with escape char
	if shader[len(shader)-1] != '\x00' {
		shader += "\x00"
	}

	return shader, s.getAttribTypes(), s.getUniformNames(), nil
}

type AttribChannelInfo struct {
	Type  r.ShaderType
	Index uint32
}

func (s *ShaderBuilder) getAttribTypes() map[r.Channel]AttribChannelInfo {
	res := make(map[r.Channel]AttribChannelInfo)
	for i, channel := range s.attribChans {
		res[channel] = AttribChannelInfo{channel.varType, uint32(i) + uint32(s.startPos)}
	}
	return res
}

func (s *ShaderBuilder) getUniformNames() []string {
	res := make([]string, 0, len(s.operChans))
	for _, channel := range s.operChans {
		res = append(res, channel.Name())
	}
	return res
}
