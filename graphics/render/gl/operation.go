package gl

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/graphics/render/gl/shader"
	"github.com/eliiasg/deltawing/graphics/render/gl/util"
	"github.com/eliiasg/glow/enum"
)

type Operation struct {
	cxt Context
	// Vao Object
	Vao any
	// Amount of instances to draw
	InstanceAmt uint32
	// Procedure to use for drawing
	Proc *Procedure
	// Parameters for uniforms
	UniformParams map[string]any
	// Start index of sprite in
	SpriteIdxStart int32
	// Amount of indices in sprite
	SpriteIdxAmt int32
}

func GLOperation(o render.Operation) (*Operation, bool) {
	res, ok := o.(*Operation)
	return res, ok
}

func (r *Renderer) MakeOperation(proc render.Procedure) render.Operation {
	return &Operation{r.cxt, r.cxt.CreateVertexArray(), 0, proc.(*Procedure), make(map[string]any), 0, 0}
}

func (o *Operation) Free() {
	o.cxt.DeleteVertexArray(o.Vao)
}

func (o *Operation) SetInstanceAttribute(channel render.Channel, buffer render.DataBuffer, offset uint32, index uint16) {
	channelInfo := o.Proc.AttribChannels[channel]
	buf := buffer.(*dataBuffer)
	if len(buf.Layout) == 0 {
		panic("missing buffer layout")
	}
	// calculate offset
	off := uintptr(offset) * uintptr(buf.LayoutSize)
	for i := uint16(0); i < index; i++ {
		off += uintptr(render.SizeOf(buf.Layout[i]))
	}
	// binding
	o.cxt.BindVertexArray(o.Vao)
	o.cxt.BindBuffer(enum.ARRAY_BUFFER, buf.Buffer)
	// setup

	typ := buf.Layout[index]
	// layout index
	o.cxt.EnableVertexAttribArray(channelInfo.Index)
	// OpenGL is more annoying than i thought, amazing!
	// IPointer must be used if its an int to int for some reason, thought that was what the normalized param was for
	if render.IsInt(channelInfo.Type.Type) {
		o.cxt.VertexAttribIPointer(channelInfo.Index, int32(typ.Amount), glType(typ.Type), int32(buf.LayoutSize), off)
	} else {
		o.cxt.VertexAttribPointer(channelInfo.Index, int32(typ.Amount), glType(typ.Type), false, int32(buf.LayoutSize), off)
	}
	o.cxt.VertexAttribDivisor(channelInfo.Index, 1)

}

// very exiting function
// maybe using a map would look better
func glType(typ render.ChannelInputType) uint32 {
	switch typ {
	case render.InputByte:
		return enum.BYTE
	case render.InputUnsignedByte:
		return enum.UNSIGNED_BYTE
	case render.InputShort:
		return enum.SHORT
	case render.InputUnsignedShort:
		return enum.UNSIGNED_SHORT
	case render.InputInt:
		return enum.INT
	case render.InputUnsignedInt:
		return enum.UNSIGNED_INT
	case render.InputFloat:
		return enum.FLOAT
	case render.InputDouble:
		return enum.DOUBLE
	default:
		return 0
	}
}

func (o *Operation) SetChannelValue(channel render.Channel, data any) {
	glChan := shader.GLChannel(channel)
	//uniform := o.cxt.GetUniformLocation(o.Proc.Prog, glChan.Name())
	if !util.AssertType(glChan.ShaderType(), data) {
		// Maybe bad?
		panic("Unable to set channel value: Invalid type")
	}
	// set param
	o.UniformParams[glChan.Name()] = data
}

func (o *Operation) DrawTo(target render.RenderTarget) {
	// tell OpenGl how big target is, i don't really understand why this would be required
	o.cxt.Viewport(0, 0, int32(target.Width()), int32(target.Height()))
	o.bind(target)
	o.initShader(target.Width(), target.Height())
	// o.spriteIdxStart is *4, because the argument is in bytes, but type is 32bit
	o.cxt.DrawElementsInstanced(enum.TRIANGLES, o.SpriteIdxAmt, enum.UNSIGNED_INT, uintptr(o.SpriteIdxStart*4), int32(o.InstanceAmt))
}

func (o *Operation) bind(target render.RenderTarget) {
	o.cxt.UseProgram(o.Proc.Prog)
	o.cxt.BindVertexArray(o.Vao)
	tar, _ := GLRenderTarget(target)
	o.cxt.BindFramebuffer(enum.FRAMEBUFFER, tar.Framebuffer)
}

func (o *Operation) initShader(width, height uint16) {
	for name, param := range o.UniformParams {
		setUniform(o.cxt, o.Proc.UniformLocations[name], param)
	}
	setUniform(o.cxt, o.Proc.ScreenSizeLocation, [2]int32{int32(width), int32(height)})
}

// very smart to have a function for every type, i just love OpenGL
func setUniform(cxt Context, location any, data any) {
	switch v := data.(type) {
	case int32:
		cxt.Uniform1i(location, v)
	case uint32:
		cxt.Uniform1ui(location, v)
	case float32:
		cxt.Uniform1f(location, v)
	case [2]int32:
		cxt.Uniform2i(location, v[0], v[1])
	case [2]uint32:
		cxt.Uniform2ui(location, v[0], v[1])
	case [2]float32:
		cxt.Uniform2f(location, v[0], v[1])
	case [3]int32:
		cxt.Uniform3i(location, v[0], v[1], v[2])
	case [3]uint32:
		cxt.Uniform3ui(location, v[0], v[1], v[2])
	case [3]float32:
		cxt.Uniform3f(location, v[0], v[1], v[2])
	case [4]int32:
		cxt.Uniform4i(location, v[0], v[1], v[2], v[3])
	case [4]uint32:
		cxt.Uniform4ui(location, v[0], v[1], v[2], v[3])
	case [4]float32:
		cxt.Uniform4f(location, v[0], v[1], v[2], v[3])
	default:
		// type is checked when added
		panic("This should never happen")
	}
}

func (o *Operation) SetSprite(buffer render.SpriteBuffer, id uint32) {
	buf := buffer.(*SpriteBuffer)
	// tell operation what sprite to draw
	o.SpriteIdxStart = int32(buf.IdxPositions[id])
	o.SpriteIdxAmt = int32(buf.IdxPositions[id+1]) - o.SpriteIdxStart
	// setup vao
	o.cxt.BindVertexArray(o.Vao)
	o.cxt.BindBuffer(enum.ARRAY_BUFFER, buf.Verts)
	// to store on VAO
	o.cxt.BindBuffer(enum.ELEMENT_ARRAY_BUFFER, buf.Inds)
	// position
	o.cxt.EnableVertexAttribArray(0)
	o.cxt.VertexAttribPointer(0, 2, enum.FLOAT, false, 12, 0)
	// color / layer
	o.cxt.EnableVertexAttribArray(1)
	o.cxt.VertexAttribIPointer(1, 4, enum.UNSIGNED_BYTE, 12, 8)
}

func (o *Operation) SetAmount(amount uint32) {
	o.InstanceAmt = amount
}
