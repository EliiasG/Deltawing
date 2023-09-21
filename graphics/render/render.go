package render

import "github.com/eliiasg/deltawing/graphics/vecsprite"

type ChannelInputType uint8

type ChannelShaderType uint8

// Based on OpenGL types: https://www.khronos.org/opengl/wiki/OpenGL_Type
const (
	// int8
	InputByte ChannelInputType = iota
	InputUnsignedbyte
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
	ShaderDouble
)

// Describes a channel
// Channels are information used to transform sprites
// Channels are either initalized to zero, per operation, or per sprite
// Channels can then be read and modified by functions
type ChannelDescription struct {
	InputType  ChannelInputType
	Shadertype ChannelShaderType
	// must be 1, 2, 3 or 4
	Amount uint8
}

// represents a GLSL Type, amount is used to represent vecs
type FunctionParameter struct {
	Type ChannelShaderType
	// must be 1, 2, 3 or 4
	Amount uint8
}

// Modifies and reads from channels
type Function struct {
	Parameters []FunctionParameter
	Source     string
}

func Param(t ChannelShaderType, amt uint8) FunctionParameter {
	return FunctionParameter{t, amt}
}

func NewFunction(source string, params []FunctionParameter) *Function {
	return &Function{params, source}
}

func Channel(inType ChannelInputType, shaderType ChannelShaderType, amount uint8) ChannelDescription {
	return ChannelDescription{inType, shaderType, amount}
}

type RendererObject interface {
	Free()
}

// Describes layout of buffer
// If the data changes often then it should likely only contain one attribute
// For more static data it might be worth combining multiple attributes into a single buffer
type DataBufferLayout *[]ChannelInputType

// Used to tarnsform sprites
// Can (and should when possible) be used for multiple sprites
type DataBuffer interface {
	RendererObject
	SetData(*[]any)
	SetLayout(DataBufferLayout)
	// TODO maybe add method to replace only parts of data?
}

// A buffer used to store sprites in the renderer
// Use as few as possible, while still keeping as many sprites out of memory as possible
type SpriteBuffer RendererObject

type SpriteBufferBuilder interface {
	AddSprite(sprite vecsprite.VecSprite) uint32
	Finish() SpriteBuffer
}

// An image/screen/texture/framebuffer/whatever that can be drawn to.
type RenderTarget interface {
	RendererObject
	Width() uint16
	Height() uint16

	// Draw on other RenderTarget with given position, size, rotation and pivot, pivot is realative to given size
	DrawTo(target RenderTarget, x, y, width, height, pivotX, pivotY uint16, rotation float32)
}

// Describes how to transform a sprite from the given data
type Procedure RendererObject

type ChannelID uint16

type ProcedureBuilder interface {
	// A channel that is neither initialized per operation, nor per sprite
	// 'expression' specifies the default value, as a GLSL expression
	// Description is FunctionParameter because it does not have a input type
	AddIntermediateChannel(description FunctionParameter, expression string) ChannelID

	// A channel initialized per sprite, this is called an attribute for the drawn sprite
	AddAttributeChannel(description ChannelDescription) ChannelID

	// A channel initialized per operation
	AddOperationChannel(description ChannelDescription) ChannelID

	// Adds a function, keep in mind that order matters
	AddFunction(function Function, channels ...ChannelID)

	// Sets the channel to use for the position, must be 2 ints
	SetPositionChannel(channel ChannelID)

	// Set the channel to use for the layer, must be an int
	SetLayerChannel(channel ChannelID)

	// Use following methods for scaling and rotation
	// Before translation, every vertex in a sprite will be recalculated with the following formula: (XAxis * x + YAxis * y) where x and y is the original position
	SetXAxisChannel(channel ChannelID)
	SetYAxisChannel(channel ChannelID)

	// Used to "compile" the Procedure
	Finish() Procedure
}

// Describes how to draw sprites, this is more specfic than the precedure, because the procedure only says how to use data, and this says what data to use
type Opteration interface {
	RendererObject
	// Supply the next attribute for the procedure, this should be called as many times as the procedure has attributes
	// The attribute added to the ProcedureBuilder first will be supplied first
	// Offset says where in the DataBuffer to start, and the index says what data from the DataBuffer to use
	AddInstanceAttribute(DataBuffer, offset uint32, index uint16)

	// Set a OperationChannel, the channel parameter should be a id returned by the ProcedureBuilder
	SetChannelValue(channel ChannelID, data any)

	// Runs the operation and reads the buffers
	DrawTo(target RenderTarget)
}

type Renderer interface {
	MakeDataBuffer() DataBuffer
	MakeSpriteBufferBuilder() SpriteBufferBuilder
	MakeRenderTarget() RenderTarget
	MakeProcedureBuilder() ProcedureBuilder
	MakeOperation(Procedure) Opteration
	PrimaryRenderTarget() RenderTarget
}
