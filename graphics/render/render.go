package render

import (
	"github.com/eliiasg/deltawing/graphics/vecsprite"
)

type ChannelInputType uint8

type ChannelShaderType uint8

// Based on OpenGL types: https://www.khronos.org/opengl/wiki/OpenGL_Type
const (
	// int8
	InputByte ChannelInputType = iota
	InputUnsignedByte
	// int16
	InputShort
	InputUnsignedShort
	// int32
	InputInt
	InputUnsignedInt
	// float32
	InputFloat
	// float64
	InputDouble
)

// Based on GLSL types: https://www.khronos.org/opengl/wiki/Data_Type_(GLSL)
const (
	ShaderInt ChannelShaderType = iota
	ShaderUnsignedInt
	ShaderFloat
	// Double missing since it is not avalibe in GLSL ES 300
)

func IsInt(t ChannelShaderType) bool {
	return t == ShaderInt || t == ShaderUnsignedInt
}

// Describes a channel
// Channels are information used to transform sprites
// Channels are either initalized to zero, per operation, or per sprite
// Channels can then be read and modified by functions

// represents a GLSL Type, amount is used to represent vecs
type ShaderType struct {
	Type ChannelShaderType
	// must be 1, 2, 3 or 4
	Amount uint8
}

type InputType struct {
	Type ChannelInputType
	// must be 1, 2, 3 or 4
	Amount uint8
}

func SizeOf(typ InputType) uint8 {
	var r uint8
	switch typ.Type {
	case InputByte, InputUnsignedByte:
		r = 1
	case InputShort, InputUnsignedShort:
		r = 2
	case InputInt, InputUnsignedInt, InputFloat:
		r = 4
	case InputDouble:
		r = 8
	default:
		return 0
	}
	return r * typ.Amount
}

// Modifies and reads from channels, screenSize and aColor exist as variables, aColor.a is undefined
type Function struct {
	Parameters []ShaderType
	Source     string
	Name       string
}

func Type(t ChannelShaderType, amt uint8) ShaderType {
	return ShaderType{t, amt}
}

func Input(t ChannelInputType, amt uint8) InputType {
	return InputType{t, amt}
}

func NewFunction(source, name string, params ...ShaderType) *Function {
	return &Function{params, source, name}
}

type RendererObject interface {
	Free()
}

// Used to tarnsform sprites
// Can (and should when possible) be used for multiple different sprites
// Must be uint because of how go interfaces work, if you wish to combine types they should be converted to uints, this can be done with buffers.AddTo
type DataBuffer interface {
	RendererObject
	// Very smart to not support generic methods
	SetData8(data []uint8)
	SetData16(data []uint16)
	SetData32(data []uint32)
	SetData64(data []uint64)
	// Describes layout of buffer
	// If the data changes often then it should likely only contain one attribute
	// For static data it might be worth combining multiple attributes into a single buffer
	SetLayout(layout ...InputType)
	// TODO maybe add method to replace only parts of data?
}

// A buffer used to store sprites in the renderer
// Use as few as possible, while still keeping as many sprites out of memory as possible
type SpriteBuffer interface {
	RendererObject
	// bit of a hack, this is used to have RendererObjects not be interchangable
	// without this you could pass a SpriteBuffer to a function expecting a Procedure
	spriteBuffer()
}

type SpriteBufferIdentifier struct{}

func (s SpriteBufferIdentifier) spriteBuffer() {
	panic("should never be called")
}

type SpriteBufferBuilder interface {
	// adds a sprite to the buffer and returns the ID
	AddSprite(sprite *vecsprite.VecSprite) uint32
	// static specifies whether the buffer is optimized to not be reallocated
	MakeBuffer(static bool) SpriteBuffer
	Reallocate(buffer SpriteBuffer)
	Clear()
}

// An image/screen/texture/framebuffer/whatever that can be drawn to.
type RenderTarget interface {
	RendererObject
	Width() uint16
	Height() uint16
	Clear(r, g, b uint8)
	Resize(width, height uint16)
	// Draw on other RenderTarget using bliting
	BlitTo(target RenderTarget, x, y int32)
	// Draw on other RenderTarget with given shader,position, size, rotation and pivot, pivot is realative to given size
	// Disabled for now, i'll need a proper FragmentShader system sometime
	//DrawTo(target RenderTarget, x, y int32, width, height, pivotX, pivotY uint16, rotation float32, shader FragmentShader)
}

// Describes how to transform a sprite from the given data
type Procedure interface {
	RendererObject
	// the hack again
	procedure()
}

type ProcedureIdentifier struct{}

func (s ProcedureIdentifier) procedure() {
	panic("should never be called")
}

type Channel interface {
	// the hack once again
	channel()
}

type ChannelIdentifier struct{}

func (s ChannelIdentifier) channel() {
	panic("should never be called")
}

type ProcedureBuilder interface {
	// A channel that is neither initialized per operation, nor per sprite
	// 'expression' specifies the default value, as a GLSL expression
	// Description is Shadertype because the input data is not realavent to the shader
	AddIntermediateChannel(shaderType ShaderType, expression string) Channel

	// A channel initialized per sprite, this is called an attribute for the drawn sprite
	// These channels may only be read from
	AddAttributeChannel(shaderType ShaderType) Channel

	// A channel initialized per operation
	// These channels may only be read from
	AddOperationChannel(shaderType ShaderType) Channel

	// Adds a function, keep in mind that order matters
	CallFunction(function *Function, channels ...Channel) error

	// Sets the channel to use for the position, must be 2 floats - final position is (0, 0) in top left and (width-1, height-1) in bottom right
	SetPositionChannel(channel Channel) error

	// Set the channel to use for the layer, must be an uint
	SetLayerChannel(channel Channel) error

	// Set the channel tp use for the color, must be 3 ints - if not set ivec4(aColor.rgb, 255) variable be used
	SetColorChannel(channel Channel) error

	// Use following methods for scaling and rotation, must be 2 floats per channel
	// Before translation, every vertex in a sprite will be recalculated with the following formula: (XAxis * x + YAxis * y) where x and y is the original position
	SetXAxisChannel(channel Channel) error
	SetYAxisChannel(channel Channel) error

	// Used to "compile" the Procedure
	Finish() (Procedure, error)
}

// Describes how to draw sprites, this is more specfic than the precedure, because the procedure only says how to use data, and this says what data to use
type Operation interface {
	RendererObject
	// Supply an attribute for the procedure, this should be called as many times as the procedure has attributes
	// Offset says where in the DataBuffer to start, and bufferIndex says what data from the DataBufferLayout to use
	SetInstanceAttribute(channel Channel, buffer DataBuffer, offset uint32, bufferIndex uint16)

	// Set a OperationChannel returned by ProcedureBuilder.AddOperationChannel()
	SetChannelValue(channel Channel, data any)

	// Set sprite given buffer and index returned by SpriteBufferBuilder.AddSprite()
	SetSprite(buffer SpriteBuffer, id uint32)

	// Set the amount of sprites to draw, if this is longer than the avalible buffers the result is undefined
	SetAmount(amount uint32)

	// Runs the operation and reads the buffers
	DrawTo(target RenderTarget)
}

type FragmentShader interface {
	// the hack once again
	fragmentShader()
}

type FragmentShaderIdentifier struct{}

func (s FragmentShaderIdentifier) fragmentShader() {
	panic("should never be called")
}

type Renderer interface {
	// if static is true buffer is optimized to be only written to once
	MakeDataBuffer(static bool) DataBuffer
	MakeSpriteBufferBuilder() SpriteBufferBuilder
	MakeRenderTarget(width, height uint16, multisample bool) RenderTarget
	MakeProcedureBuilder() ProcedureBuilder
	MakeOperation(procedure Procedure) Operation
	// Only allows simple shaders for small effects, because the input data is not modifiable
	// Expects a function that with the following parameters:
	// in original: sampler2d
	// in uv:       vec2
	// out color:   vec4
	// Disabled for now, seems unnessecary, and if i want a way to make effects on RenderTargets it should be designed better
	// MakeFragmentShader(source string) FragmentShader

	PrimaryRenderTarget() RenderTarget
}
