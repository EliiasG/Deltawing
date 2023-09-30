package opengl

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/deltawing/internal/rendering/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type operation struct {
	vaoID  uint32
	idx    uint32
	amount uint32
}

func (r *renderer) MakeOperation(procedure render.Procedure) render.Operation {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	return &operation{vao, shader.VertexBaseInputAmt, 0}
}

func (o *operation) Free() {
	panic("not implemented") // TODO: Implement
}

func (o *operation) AddInstanceAttribute(buffer render.DataBuffer, offset uint32, index uint16) {
	buf := buffer.(*dataBuffer)
	// calculate offset
	var off uintptr
	for i := uint32(0); i < offset; i++ {
		off += uintptr(render.SizeOf(buf.layout[i]))
	}
	// binding
	gl.BindVertexArray(o.vaoID)
	gl.BindBuffer(gl.ARRAY_BUFFER, buf.id)
	// setup
	typ := buf.layout[index]
	gl.VertexAttribPointerWithOffset(o.idx, int32(render.SizeOf(typ)), glType(typ.Type), false, int32(buf.layoutSize), off)
	o.idx++
}

// very exiting function
// maybe using a map would look better
func glType(typ render.ChannelInputType) uint32 {
	switch typ {
	case render.InputByte:
		return gl.BYTE
	case render.InputUnsignedByte:
		return gl.UNSIGNED_BYTE
	case render.InputShort:
		return gl.SHORT
	case render.InputUnsignedShort:
		return gl.UNSIGNED_SHORT
	case render.InputInt:
		return gl.INT
	case render.InputUnsignedInt:
		return gl.UNSIGNED_INT
	case render.InputFloat:
		return gl.FLOAT
	case render.InputDouble:
		return gl.DOUBLE
	}
}

// Set a OperationChannel returned by ProcedureBuilder.AddOperationChannel()
func (o *operation) SetChannelValue(channel render.Channel, data any) {
	panic("not implemented") // TODO: Implement
}

// Runs the operation and reads the buffers
func (o *operation) DrawTo(target render.RenderTarget) {
	panic("not implemented") // TODO: Implement
}

func (o *operation) SetSprite(buffer render.SpriteBuffer, id uint32) {
	gl.BindVertexArray(o.vaoID)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.(*spriteBuffer).vertsID)
	// position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 12, 0)
	// color / layer
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 4, gl.UNSIGNED_BYTE, false, 12, 8)
}

func (o *operation) SetAmount(amount uint32) {
	o.amount = amount
}
