//go:build wasm
// +build wasm

package webgl

import (
	"syscall/js"
	"unsafe"

	"github.com/eliiasg/deltawing/graphics/render/gl"
	"github.com/eliiasg/glow/enum"
)

type context struct {
	dataSize                       int
	dataView                       js.Value
	dataBuffer                     js.Value
	createBuffer                   js.Value
	createFramebuffer              js.Value
	createProgram                  js.Value
	createRenderbuffer             js.Value
	createShader                   js.Value
	createTexture                  js.Value
	createVertexArray              js.Value
	deleteBuffer                   js.Value
	deleteFramebuffer              js.Value
	deleteProgram                  js.Value
	deleteRenderbuffer             js.Value
	deleteShader                   js.Value
	deleteTexture                  js.Value
	deleteVertexArray              js.Value
	bindBuffer                     js.Value
	bindFramebuffer                js.Value
	bindRenderbuffer               js.Value
	bindTexture                    js.Value
	bindVertexArray                js.Value
	attachShader                   js.Value
	blitFramebuffer                js.Value
	bufferData                     js.Value
	clear                          js.Value
	clearColor                     js.Value
	clearDepth                     js.Value
	compileShader                  js.Value
	drawElementsInstanced          js.Value
	enableVertexAttribArray        js.Value
	framebufferRenderbuffer        js.Value
	framebufferTexture2D           js.Value
	getProgramInfoLog              js.Value
	getProgramParameter            js.Value
	getShaderInfoLog               js.Value
	getUniformLocation             js.Value
	getShaderParameter             js.Value
	linkProgram                    js.Value
	renderbufferStorageMultisample js.Value
	shaderSource                   js.Value
	texImage2D                     js.Value
	useProgram                     js.Value
	vertexAttribDivisor            js.Value
	vertexAttribIPointer           js.Value
	vertexAttribPointer            js.Value
	viewport                       js.Value
	uniform1i                      js.Value
	uniform2i                      js.Value
	uniform3i                      js.Value
	uniform4i                      js.Value
	uniform1ui                     js.Value
	uniform2ui                     js.Value
	uniform3ui                     js.Value
	uniform4ui                     js.Value
	uniform1f                      js.Value
	uniform2f                      js.Value
	uniform3f                      js.Value
	uniform4f                      js.Value
}

func getFunction(target js.Value, name string) js.Value {
	return target.Get(name).Call("bind", target)
}

func MakeContext(g js.Value) gl.Context {
	return &context{
		createBuffer:                   getFunction(g, "createBuffer"),
		createFramebuffer:              getFunction(g, "createFramebuffer"),
		createProgram:                  getFunction(g, "createProgram"),
		createRenderbuffer:             getFunction(g, "createRenderbuffer"),
		createShader:                   getFunction(g, "createShader"),
		createTexture:                  getFunction(g, "createTexture"),
		createVertexArray:              getFunction(g, "createVertexArray"),
		deleteBuffer:                   getFunction(g, "deleteBuffer"),
		deleteFramebuffer:              getFunction(g, "deleteFramebuffer"),
		deleteProgram:                  getFunction(g, "deleteProgram"),
		deleteRenderbuffer:             getFunction(g, "deleteRenderbuffer"),
		deleteShader:                   getFunction(g, "deleteShader"),
		deleteTexture:                  getFunction(g, "deleteTexture"),
		deleteVertexArray:              getFunction(g, "deleteVertexArray"),
		bindBuffer:                     getFunction(g, "bindBuffer"),
		bindFramebuffer:                getFunction(g, "bindFramebuffer"),
		bindRenderbuffer:               getFunction(g, "bindRenderbuffer"),
		bindTexture:                    getFunction(g, "bindTexture"),
		bindVertexArray:                getFunction(g, "bindVertexArray"),
		attachShader:                   getFunction(g, "attachShader"),
		blitFramebuffer:                getFunction(g, "blitFramebuffer"),
		bufferData:                     getFunction(g, "bufferData"),
		clear:                          getFunction(g, "clear"),
		clearColor:                     getFunction(g, "clearColor"),
		clearDepth:                     getFunction(g, "clearDepth"),
		compileShader:                  getFunction(g, "compileShader"),
		drawElementsInstanced:          getFunction(g, "drawElementsInstanced"),
		enableVertexAttribArray:        getFunction(g, "enableVertexAttribArray"),
		framebufferRenderbuffer:        getFunction(g, "framebufferRenderbuffer"),
		framebufferTexture2D:           getFunction(g, "framebufferTexture2D"),
		getProgramInfoLog:              getFunction(g, "getProgramInfoLog"),
		getProgramParameter:            getFunction(g, "getProgramParameter"),
		getShaderInfoLog:               getFunction(g, "getShaderInfoLog"),
		getUniformLocation:             getFunction(g, "getUniformLocation"),
		getShaderParameter:             getFunction(g, "getShaderParameter"),
		linkProgram:                    getFunction(g, "linkProgram"),
		renderbufferStorageMultisample: getFunction(g, "renderbufferStorageMultisample"),
		shaderSource:                   getFunction(g, "shaderSource"),
		texImage2D:                     getFunction(g, "texImage2D"),
		useProgram:                     getFunction(g, "useProgram"),
		vertexAttribDivisor:            getFunction(g, "vertexAttribDivisor"),
		vertexAttribIPointer:           getFunction(g, "vertexAttribIPointer"),
		vertexAttribPointer:            getFunction(g, "vertexAttribPointer"),
		viewport:                       getFunction(g, "viewport"),
		uniform1i:                      getFunction(g, "uniform1i"),
		uniform2i:                      getFunction(g, "uniform2i"),
		uniform3i:                      getFunction(g, "uniform3i"),
		uniform4i:                      getFunction(g, "uniform4i"),
		uniform1ui:                     getFunction(g, "uniform1ui"),
		uniform2ui:                     getFunction(g, "uniform2ui"),
		uniform3ui:                     getFunction(g, "uniform3ui"),
		uniform4ui:                     getFunction(g, "uniform4ui"),
		uniform1f:                      getFunction(g, "uniform1f"),
		uniform2f:                      getFunction(g, "uniform2f"),
		uniform3f:                      getFunction(g, "uniform3f"),
		uniform4f:                      getFunction(g, "uniform4f"),
	}
}

func dataViewSetUint[T uint8 | uint16 | uint32 | uint64](view js.Value, val T, idx int) {
	switch v := any(val).(type) {
	case uint8:
		view.Call("setUint8", idx, v, true)
	case uint16:
		view.Call("setUint16", idx*2, v, true)
	case uint32:
		view.Call("setUint32", idx*4, v, true)
	case uint64:
		ints := *(*[2]uint32)(unsafe.Pointer(&v))
		view.Call("setUint32", idx*8, ints[0], true)
		view.Call("setUint32", idx*8+4, ints[1], true)
	}
}

// easier to have a generic type
func makeArrBuf[T uint8 | uint16 | uint32 | uint64](c *context, data []T) js.Value {
	// generic size hack
	var v T
	size := int(unsafe.Sizeof(v))
	byteSize := size * len(data)
	if byteSize > c.dataSize {
		c.dataBuffer = js.Global().Get("ArrayBuffer").New(byteSize)
		c.dataView = js.Global().Get("DataView").New(c.dataBuffer)
		c.dataSize = byteSize
	}

	for i, elem := range data {
		dataViewSetUint(c.dataView, elem, i)
	}
	return c.dataView
}

func (c *context) jsData(data any) js.Value {
	if data == nil {
		return js.Null()
	}
	switch v := data.(type) {
	case []uint8:
		return makeArrBuf(c, v)
	case []uint16:
		return makeArrBuf(c, v)
	case []uint32:
		return makeArrBuf(c, v)
	case []uint64:
		return makeArrBuf(c, v)
	default:
		// FIXME maybe add support for directly making typedarrays
		panic("WebGL2 implementation only supports uint[8,16,32,64] slices")
	}
}

// This is not auto-generated...
// I should get a life

func glEnum(val js.Value) int32 {
	if val.Type() != js.TypeBoolean {
		return int32(val.Int())
	}
	if val.Bool() {
		return enum.TRUE
	}
	return enum.FALSE
}

func (c *context) CreateBuffer() any {
	return c.createBuffer.Invoke()
}

func (c *context) CreateFramebuffer() any {
	return c.createFramebuffer.Invoke()
}

func (c *context) CreateProgram() any {
	return c.createProgram.Invoke()
}

func (c *context) CreateRenderbuffer() any {
	return c.createRenderbuffer.Invoke()
}

func (c *context) CreateShader(xtype uint32) any {
	return c.createShader.Invoke(xtype)
}

func (c *context) CreateTexture() any {
	return c.createTexture.Invoke()
}

func (c *context) CreateVertexArray() any {
	return c.createVertexArray.Invoke()
}

func (c *context) DeleteBuffer(buffer any) {
	c.deleteBuffer.Invoke(buffer)
}

func (c *context) DeleteFramebuffer(framebuffer any) {
	c.deleteFramebuffer.Invoke(framebuffer)
}

func (c *context) DeleteProgram(progarm any) {
	c.deleteProgram.Invoke(progarm)
}

func (c *context) DeleteRenderbuffer(renderbuffer any) {
	c.deleteRenderbuffer.Invoke(renderbuffer)
}

func (c *context) DeleteShader(shader any) {
	c.deleteShader.Invoke(shader)
}

func (c *context) DeleteTexture(texture any) {
	c.deleteTexture.Invoke(texture)
}

func (c *context) DeleteVertexArray(vertexArray any) {
	c.deleteVertexArray.Invoke(vertexArray)
}

func (c *context) BindBuffer(target uint32, buffer any) {
	c.bindBuffer.Invoke(target, buffer)
}

func (c *context) BindFramebuffer(target uint32, framebuffer any) {
	c.bindFramebuffer.Invoke(target, framebuffer)
}

func (c *context) BindRenderbuffer(target uint32, renderbuffer any) {
	c.bindRenderbuffer.Invoke(target, renderbuffer)
}

func (c *context) BindTexture(target uint32, texture any) {
	c.bindTexture.Invoke(target, texture)
}

func (c *context) BindVertexArray(array any) {
	c.bindVertexArray.Invoke(array)
}

func (c *context) AttachShader(program any, shader any) {
	c.attachShader.Invoke(program, shader)
}

func (c *context) BlitFramebuffer(srcX0 int32, srcY0 int32, srcX1 int32, srcY1 int32, dstX0 int32, dstY0 int32, dstX1 int32, dstY1 int32, mask uint32, filter uint32) {
	c.blitFramebuffer.Invoke(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter)
}

func (c *context) BufferData(target uint32, data any, usage uint32) {
	c.bufferData.Invoke(target, c.jsData(data), usage)
}

func (c *context) Clear(mask uint32) {
	c.clear.Invoke(mask)
}

func (c *context) ClearColor(r float32, g float32, b float32, a float32) {
	c.clearColor.Invoke(r, g, b, a)
}

func (c *context) CompileShader(shader any) {
	c.compileShader.Invoke(shader)
}

func (c *context) DrawElementsInstanced(mode uint32, count int32, xtype uint32, indexOffset uintptr, instancecount int32) {
	c.drawElementsInstanced.Invoke(mode, count, xtype, indexOffset, instancecount)
}

func (c *context) EnableVertexAttribArray(index uint32) {
	c.enableVertexAttribArray.Invoke(index)
}

func (c *context) FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer any) {
	c.framebufferRenderbuffer.Invoke(target, attachment, renderbuffertarget, renderbuffer)
}

func (c *context) FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture any, level int32) {
	c.framebufferTexture2D.Invoke(target, attachment, textarget, texture, level)
}

func (c *context) GetProgramInfoLog(program any) string {
	return c.getProgramInfoLog.Invoke(program).String()
}

func (c *context) GetProgramParameter(program any, pname uint32) int32 {
	return glEnum(c.getProgramParameter.Invoke(program, pname))
}

func (c *context) GetShaderInfoLog(shader any) string {
	return c.getShaderInfoLog.Invoke(shader).String()
}

func (c *context) GetShaderParameter(shader any, pname uint32) int32 {
	return glEnum(c.getShaderParameter.Invoke(shader, pname))
}

func (c *context) GetUniformLocation(shader any, name string) any {
	return c.getUniformLocation.Invoke(shader, name)
}

func (c *context) LinkProgram(program any) {
	c.linkProgram.Invoke(program)
}

func (c *context) RenderbufferStorageMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32) {
	c.renderbufferStorageMultisample.Invoke(target, samples, internalformat, width, height)
}

func (c *context) ShaderSource(shader any, source string) {
	c.shaderSource.Invoke(shader, source)
}

func (c *context) TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels any) {
	c.texImage2D.Invoke(target, level, internalformat, width, height, border, format, xtype, c.jsData(pixels))
}

func (c *context) UseProgram(program any) {
	c.useProgram.Invoke(program)
}

func (c *context) VertexAttribDivisor(index uint32, divisor uint32) {
	c.vertexAttribDivisor.Invoke(index, divisor)
}

func (c *context) VertexAttribIPointer(index uint32, size int32, xtype uint32, stride int32, offset uintptr) {
	c.vertexAttribIPointer.Invoke(index, size, xtype, stride, offset)
}

func (c *context) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uintptr) {
	c.vertexAttribPointer.Invoke(index, size, xtype, normalized, stride, offset)
}

func (c *context) Viewport(x int32, y int32, width int32, height int32) {
	c.viewport.Invoke(x, y, width, height)
}

func (c *context) Uniform1i(location any, v0 int32) {
	c.uniform1i.Invoke(location, v0)
}

func (c *context) Uniform2i(location any, v0 int32, v1 int32) {
	c.uniform2i.Invoke(location, v0, v1)
}

func (c *context) Uniform3i(location any, v0 int32, v1 int32, v2 int32) {
	c.uniform3i.Invoke(location, v0, v1, v2)
}

func (c *context) Uniform4i(location any, v0 int32, v1 int32, v2 int32, v3 int32) {
	c.uniform4i.Invoke(location, v0, v1, v2, v3)
}

func (c *context) Uniform1ui(location any, v0 uint32) {
	c.uniform1ui.Invoke(location, v0)
}

func (c *context) Uniform2ui(location any, v0 uint32, v1 uint32) {
	c.uniform2ui.Invoke(location, v0, v1)
}

func (c *context) Uniform3ui(location any, v0 uint32, v1 uint32, v2 uint32) {
	c.uniform3ui.Invoke(location, v0, v1, v2)
}

func (c *context) Uniform4ui(location any, v0 uint32, v1 uint32, v2 uint32, v3 uint32) {
	c.uniform4ui.Invoke(location, v0, v1, v2, v3)
}

func (c *context) Uniform1f(location any, v0 float32) {
	c.uniform1f.Invoke(location, v0)
}

func (c *context) Uniform2f(location any, v0 float32, v1 float32) {
	c.uniform2f.Invoke(location, v0, v1)
}

func (c *context) Uniform3f(location any, v0 float32, v1 float32, v2 float32) {
	c.uniform3f.Invoke(location, v0, v1, v2)
}

func (c *context) Uniform4f(location any, v0 float32, v1 float32, v2 float32, v3 float32) {
	c.uniform4f.Invoke(location, v0, v1, v2, v3)
}
