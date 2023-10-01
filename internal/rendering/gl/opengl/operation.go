package opengl

import (
	"github.com/eliiasg/deltawing/graphics/render"
	g "github.com/eliiasg/deltawing/internal/rendering/gl"
	"github.com/eliiasg/deltawing/internal/rendering/shader"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type operation struct {
	vaoID         uint32
	idx           uint32
	amount        uint32
	proc          *procedure
	uniformParams map[int32]any
	idxStart      int32
	idxAmt        int32
}

func (r *renderer) MakeOperation(proc render.Procedure) render.Operation {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	return &operation{vao, shader.VertexBaseInputAmt, 0, proc.(*procedure), make(map[int32]any), 0, 0}
}

func (o *operation) Free() {
	gl.DeleteVertexArrays(1, &o.vaoID)
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
	default:
		return 0
	}
}

func (o *operation) SetChannelValue(channel render.Channel, data any) {
	uniform := gl.GetUniformLocation(o.proc.progID, gl.Str(shader.ChannelName(channel)))
	if !g.AssertType(shader.ChannelType(channel), data) {
		// Maybe bad?
		panic("Invalid type")
	}
	o.uniformParams[uniform] = data
}

func (o *operation) DrawTo(target render.RenderTarget) {
	o.bind(target)
	o.initShader(target.Width(), target.Height())
	gl.DrawElementsInstancedBaseVertex(gl.TRIANGLES, o.idxAmt, gl.UNSIGNED_INT, nil, int32(o.amount), o.idxStart)
}

func (o *operation) bind(target render.RenderTarget) {
	gl.UseProgram(o.proc.progID)
	gl.BindVertexArray(o.vaoID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, target.(*renderTarget).framebufferID)
}

func (o *operation) initShader(width, height uint16) {
	for location, value := range o.uniformParams {
		setUniform(location, value)
	}
	setUniform(o.proc.screenSizeLocation, [2]float32{float32(width), float32(height)})
}

// very smart to have a function for every type, i just love OpenGL
func setUniform(location int32, data any) {
	switch v := data.(type) {
	case int32:
		gl.Uniform1i(location, v)
	case uint32:
		gl.Uniform1ui(location, v)
	case float32:
		gl.Uniform1f(location, v)
	case float64:
		gl.Uniform1d(location, v)
	case [2]int32:
		gl.Uniform2i(location, v[0], v[1])
	case [2]uint32:
		gl.Uniform2ui(location, v[0], v[1])
	case [2]float32:
		gl.Uniform2f(location, v[0], v[1])
	case [2]float64:
		gl.Uniform2d(location, v[0], v[1])
	case [3]int32:
		gl.Uniform3i(location, v[0], v[1], v[2])
	case [3]uint32:
		gl.Uniform3ui(location, v[0], v[1], v[2])
	case [3]float32:
		gl.Uniform3f(location, v[0], v[1], v[2])
	case [3]float64:
		gl.Uniform3d(location, v[0], v[1], v[2])
	case [4]int32:
		gl.Uniform4i(location, v[0], v[1], v[2], v[3])
	case [4]uint32:
		gl.Uniform4ui(location, v[0], v[1], v[2], v[3])
	case [4]float32:
		gl.Uniform4f(location, v[0], v[1], v[2], v[3])
	case [4]float64:
		gl.Uniform4d(location, v[0], v[1], v[2], v[3])
	default:
		// type is checked when added
		panic("This should never happen")
	}
}

func (o *operation) SetSprite(buffer render.SpriteBuffer, id uint32) {
	buf := buffer.(*spriteBuffer)
	// tell operation what sprite to draw
	o.idxStart = int32(buf.idxPositions[id])
	o.idxAmt = int32(buf.idxPositions[id+1]) - o.idxStart
	// setup vao
	gl.BindVertexArray(o.vaoID)
	gl.BindBuffer(gl.ARRAY_BUFFER, buf.vertsID)
	// to store on VAO
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.indsID)
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
