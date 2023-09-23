package shader

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	r "github.com/eliiasg/deltawing/graphics/render"
)

type glslChannel struct {
	r.ChannelIdentifyer
	id      uint16
	varType r.ShaderType
}

type funcCall struct {
	fun *r.Function
	// only the id of the channels are required
	params []uint16
}

type interChannel struct {
	glslChannel
	expr string
}

func GetChannelName(channel r.Channel) string {
	return getVarName(channel.(glslChannel).id)
}

func getVarName(id uint16) string {
	// Beatiful
	// The only way for this to fail should be if another renderer is used to supply the channel
	return "c" + strconv.Itoa(int(id))
}

// just collects the data, shader will be composed at end
type shaderBuilder struct {
	baseSource  string
	done        func(string) (r.Procedure, error)
	chanID      uint16
	interChans  []interChannel
	attribChans []glslChannel
	operChans   []glslChannel
	calls       []funcCall
	// IDs for output channels, only storing ID as it is the only requirement for getting the var name
	posID   uint16
	layerID uint16
	xAxisID uint16
	yAxisID uint16
}

/*
The baseSource is a shader that will be added parts to.
The following keywords will be replaced:
'<attributes>' instance data, should only be used once
'<uniforms>' uniforms, should only be used once
'<functions>' function declarations, should only be used once
'<variabless>' variables, should only be used once
'<calls>' function calls, should only be used once
'<pos>' output position variable, this should not be modified
'<layer>' output layer variable, this should not be modified
'<xAxis>' output xAxis variable, this should not be modified
'<yAxis>' output xAxis variable, this should not be modified
*/
func NewShaderBuilder(baseSource string, done func(string) (r.Procedure, error)) r.ProcedureBuilder {
	return &shaderBuilder{
		baseSource:  baseSource,
		done:        done,
		chanID:      0,
		interChans:  make([]interChannel, 0),
		attribChans: make([]glslChannel, 0),
		operChans:   make([]glslChannel, 0),
		calls:       make([]funcCall, 0),
	}
}

func (s *shaderBuilder) makeChannel(varType r.ShaderType) glslChannel {
	channel := glslChannel{
		id:      s.chanID,
		varType: varType,
	}
	s.chanID++
	return channel
}

func (s *shaderBuilder) AddIntermediateChannel(shaderType r.ShaderType, expression string) r.Channel {
	channel := s.makeChannel(shaderType)
	s.interChans = append(s.interChans, interChannel{
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
		glChan := channel.(glslChannel)
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
	glChan := channel.(glslChannel)
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
	s.layerID, e = setOutputChannel(channel, r.Type(r.ShaderInt, 1))
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

var typeMap = map[r.ChannelShaderType][2]string{
	r.ShaderDouble:      {"double", "dvec"},
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
		sb.WriteString(fmt.Sprintf("layout(location=%v) %v %v;\n", i, getGLSLTypeName(channel.varType), GetChannelName(channel)))
	}
	return sb.String()
}

func (s *shaderBuilder) makeUniformSection() string {
	var sb strings.Builder
	for _, channel := range s.operChans {
		sb.WriteString(fmt.Sprintf("uniform %v %v;\n", getGLSLTypeName(channel.varType), GetChannelName(channel)))
	}
	return sb.String()
}

func (s *shaderBuilder) getFunctions() []*r.Function {
	funcSet := make(map[*r.Function]bool, 0)
	for _, call := range s.calls {
		funcSet[call.fun] = true
	}
	keys := make([]*r.Function, 0, len(funcSet))
	for f := range funcSet {
		keys = append(keys, f)
	}
	return keys
}

func (s *shaderBuilder) makeDeclarationSection() string {

}

func (s *shaderBuilder) Finish() (r.Procedure, error) {
	panic("not implemented") // TODO: Implement
}
