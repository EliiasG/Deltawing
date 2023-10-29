package opengl

import (
	"strings"
	"unsafe"

	g "github.com/eliiasg/deltawing/graphics/render/gl"
	"github.com/eliiasg/glow/v3.3-core/gl"
)

type context struct{}

func MakeContext() g.Context {
	return context{}
}

func glObj(obj any) uint32 {
	if obj == nil {
		return 0
	}
	r := obj.(uint32)
	return r
}

func glLocation(obj any) int32 {
	r := obj.(int32)
	return r
}

func glPtr(slice any) (unsafe.Pointer, int) {
	// couldnt find a way to do this in a unsafe/smart way
	switch s := slice.(type) {
	case nil:
		return nil, 0
	case []uint8:
		return gl.Ptr(s), len(s)
	case []uint16:
		return gl.Ptr(s), len(s) * 2
	case []uint32:
		return gl.Ptr(s), len(s) * 4
	case []uint64:
		return gl.Ptr(s), len(s) * 8
	default:
		panic("OpenGL implementation only supports uint[8,16,32,64] slices")
	}
}

func genUsing(f func(int32, *uint32)) uint32 {
	var ref uint32
	f(1, &ref)
	return ref
}

func delUsing(f func(int32, *uint32), obj any) {
	conv := glObj(obj)
	f(1, &conv)
}

func (c context) CreateBuffer() any {
	return genUsing(gl.GenBuffers)
}

func (c context) CreateFramebuffer() any {
	return genUsing(gl.GenFramebuffers)
}

func (c context) CreateProgram() any {
	return gl.CreateProgram()
}

func (c context) CreateRenderbuffer() any {
	return genUsing(gl.GenRenderbuffers)
}

func (c context) CreateShader(xtype uint32) any {
	return gl.CreateShader(xtype)
}

func (c context) CreateTexture() any {
	return genUsing(gl.GenTextures)
}

func (c context) CreateVertexArray() any {
	return genUsing(gl.GenVertexArrays)
}

func (c context) DeleteBuffer(buffer any) {
	delUsing(gl.DeleteBuffers, buffer)
}

func (c context) DeleteFramebuffer(framebuffer any) {
	delUsing(gl.DeleteFramebuffers, framebuffer)
}

func (c context) DeleteProgram(progarm any) {
	gl.DeleteProgram(glObj(progarm))
}

func (c context) DeleteRenderbuffer(renderbuffer any) {
	delUsing(gl.DeleteRenderbuffers, renderbuffer)
}

func (c context) DeleteShader(shader any) {
	gl.DeleteShader(glObj(shader))
}

func (c context) DeleteTexture(texture any) {
	delUsing(gl.DeleteTextures, texture)
}

func (c context) DeleteVertexArray(vertexArray any) {
	delUsing(gl.DeleteVertexArrays, vertexArray)
}

func (c context) BindBuffer(target uint32, buffer any) {
	gl.BindBuffer(target, glObj(buffer))
}

func (c context) BindFramebuffer(target uint32, framebuffer any) {
	gl.BindFramebuffer(target, glObj(framebuffer))
}

func (c context) BindRenderbuffer(target uint32, renderbuffer any) {
	gl.BindRenderbuffer(target, glObj(renderbuffer))
}

func (c context) BindTexture(target uint32, texture any) {
	gl.BindTexture(target, glObj(texture))
}

func (c context) BindVertexArray(array any) {
	gl.BindVertexArray(glObj(array))
}

func (c context) AttachShader(program any, shader any) {
	gl.AttachShader(glObj(program), glObj(shader))
}

func (c context) BlitFramebuffer(srcX0 int32, srcY0 int32, srcX1 int32, srcY1 int32, dstX0 int32, dstY0 int32, dstX1 int32, dstY1 int32, mask uint32, filter uint32) {
	gl.BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter)
}

func (c context) BufferData(target uint32, data any, usage uint32) {
	slice, size := glPtr(data)
	gl.BufferData(target, size, slice, usage)
}

func (c context) Clear(mask uint32) {
	gl.Clear(mask)
}

func (c context) ClearColor(r float32, g float32, b float32, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (c context) ClearDepth(depth float64) {
	gl.ClearDepth(depth)
}

func (c context) CompileShader(shader any) {
	gl.CompileShader(glObj(shader))
}

func (c context) DrawElementsInstanced(mode uint32, count int32, xtype uint32, indexOffset uintptr, instancecount int32) {
	gl.DrawElementsInstancedWithOffset(mode, count, xtype, indexOffset, instancecount)
}

func (c context) EnableVertexAttribArray(index uint32) {
	gl.EnableVertexAttribArray(index)
}

func (c context) FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer any) {
	gl.FramebufferRenderbuffer(target, attachment, renderbuffertarget, glObj(renderbuffer))
}

func (c context) FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture any, level int32) {
	gl.FramebufferTexture2D(target, attachment, textarget, glObj(texture), level)
}

func (c context) GetProgramInfoLog(program any) string {
	prog := glObj(program)
	// get length
	var length int32
	gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &length)
	// get log
	log := strings.Repeat("\x00", int(length+1))
	gl.GetProgramInfoLog(prog, length, nil, gl.Str(log))
	return log
}

func (c context) GetProgramParameter(program any, pname uint32) int32 {
	var res int32
	gl.GetProgramiv(glObj(program), pname, &res)
	return res
}

func (c context) GetShaderInfoLog(shader any) string {
	shad := glObj(shader)
	// get length
	var length int32
	gl.GetShaderiv(shad, gl.INFO_LOG_LENGTH, &length)
	// get log
	log := strings.Repeat("\x00", int(length+1))
	gl.GetShaderInfoLog(shad, length, nil, gl.Str(log))
	return log
}

func (c context) GetShaderParameter(shader any, pname uint32) int32 {
	var res int32
	gl.GetShaderiv(glObj(shader), pname, &res)
	return res
}

func (c context) GetUniformLocation(program any, name string) any {
	return gl.GetUniformLocation(glObj(program), gl.Str(name))
}

func (c context) LinkProgram(program any) {
	gl.LinkProgram(glObj(program))
}

func (c context) RenderbufferStorageMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32) {
	gl.RenderbufferStorageMultisample(target, samples, internalformat, width, height)
}

func (c context) ShaderSource(shader any, source string) {
	cSource, free := gl.Strs(source)
	gl.ShaderSource(glObj(shader), 1, cSource, nil)
	free()
}

func (c context) TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels any) {
	pix, _ := glPtr(pixels)
	gl.TexImage2D(target, level, internalformat, width, height, border, format, xtype, pix)
}

func (c context) UseProgram(program any) {
	gl.UseProgram(glObj(program))
}

func (c context) VertexAttribDivisor(index uint32, divisor uint32) {
	gl.VertexAttribDivisor(index, divisor)
}

func (c context) VertexAttribIPointer(index uint32, size int32, xtype uint32, stride int32, offset uintptr) {
	gl.VertexAttribIPointerWithOffset(index, size, xtype, stride, offset)
}

func (c context) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uintptr) {
	gl.VertexAttribPointerWithOffset(index, size, xtype, normalized, stride, offset)
}

func (c context) Viewport(x int32, y int32, width int32, height int32) {
	gl.Viewport(x, y, width, height)
}

func (c context) Uniform1i(location any, v0 int32) {
	gl.Uniform1i(glLocation(location), v0)
}

func (c context) Uniform2i(location any, v0 int32, v1 int32) {
	gl.Uniform2i(glLocation(location), v0, v1)
}

func (c context) Uniform3i(location any, v0 int32, v1 int32, v2 int32) {
	gl.Uniform3i(glLocation(location), v0, v1, v2)
}

func (c context) Uniform4i(location any, v0 int32, v1 int32, v2 int32, v3 int32) {
	gl.Uniform4i(glLocation(location), v0, v1, v2, v3)
}

func (c context) Uniform1ui(location any, v0 uint32) {
	gl.Uniform1ui(glLocation(location), v0)
}

func (c context) Uniform2ui(location any, v0 uint32, v1 uint32) {
	gl.Uniform2ui(glLocation(location), v0, v1)
}

func (c context) Uniform3ui(location any, v0 uint32, v1 uint32, v2 uint32) {
	gl.Uniform3ui(glLocation(location), v0, v1, v2)
}

func (c context) Uniform4ui(location any, v0 uint32, v1 uint32, v2 uint32, v3 uint32) {
	gl.Uniform4ui(glLocation(location), v0, v1, v2, v3)
}

func (c context) Uniform1f(location any, v0 float32) {
	gl.Uniform1f(glLocation(location), v0)
}

func (c context) Uniform2f(location any, v0 float32, v1 float32) {
	gl.Uniform2f(glLocation(location), v0, v1)
}

func (c context) Uniform3f(location any, v0 float32, v1 float32, v2 float32) {
	gl.Uniform3f(glLocation(location), v0, v1, v2)
}

func (c context) Uniform4f(location any, v0 float32, v1 float32, v2 float32, v3 float32) {
	gl.Uniform4f(glLocation(location), v0, v1, v2, v3)
}
