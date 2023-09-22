package shader

import (
	"strconv"

	r "github.com/eliiasg/deltawing/graphics/render"
)

type glslChannel struct {
	id      uint16
	varType r.ShaderType
}

type interChannel struct {
	glslChannel
	expr string
}

func GetUniformName(channel r.Channel) string {
	return getVarName(channel)
}

func getVarName(channel r.Channel) string {
	// Beatiful
	// The only way for this to fail should be if another renderer is used to supply the channel
	return "c" + strconv.Itoa(int(channel.(glslChannel).id))
}

type shaderBuilder struct {
	baseSource  string
	done        func(string) (r.Procedure, error)
	chanID      uint16
	interChans  []interChannel
	attribChans []r.Channel
	operChans   []r.Channel
}

func NewShaderBuilder(baseSource string, done func(string) (r.Procedure, error)) r.ProcedureBuilder {
	return &shaderBuilder{
		baseSource:  baseSource,
		done:        done,
		chanID:      0,
		interChans:  make([]interChannel, 0),
		attribChans: make([]r.Channel, 0),
		operChans:   make([]r.Channel, 0),
	}
}

func (s *shaderBuilder) makeChannel(varType r.ShaderType) {

}

func (s *shaderBuilder) AddIntermediateChannel(shaderType r.ShaderType, expression string) r.Channel {
	channel := glslChannel{
		id:      s.chanID,
		varType: shaderType,
	}
	s.interChans = append(s.interChans, interChannel{
		channel,
		expression,
	})

}

func (s *shaderBuilder) AddAttributeChannel(description r.ChannelDescription) r.Channel {
	panic("not implemented") // TODO: Implement
}

func (s *shaderBuilder) AddOperationChannel(description r.ChannelDescription) r.Channel {
	panic("not implemented") // TODO: Implement
}

func (s *shaderBuilder) AddFunction(function r.Function, channels ...r.Channel) error {
	panic("not implemented") // TODO: Implement
}

func (s *shaderBuilder) SetPositionChannel(channel r.Channel) error {
	panic("not implemented") // TODO: Implement
}

func (s *shaderBuilder) SetLayerChannel(channel r.Channel) error {
	panic("not implemented") // TODO: Implement
}

func (s *shaderBuilder) SetXAxisChannel(channel r.Channel) error {
	panic("not implemented") // TODO: Implement
}

func (s *shaderBuilder) SetYAxisChannel(channel r.Channel) error {
	panic("not implemented") // TODO: Implement
}

func (s *shaderBuilder) Finish() (r.Procedure, error) {
	panic("not implemented") // TODO: Implement
}
