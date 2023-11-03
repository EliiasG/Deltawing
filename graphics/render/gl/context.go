package gl

type Context interface {
	// Using any for gl tpyes as OpenGL represents them at ints while webgl has types
	// Not using generics, since the same code should work for all Contexts
	// It would technically be possible to make this type safe, but it would require many weird workarounds
	// Neither OpenGL nor JavaScript is type safe anyway...
	// Functions names are based on WebGL2

	CreateBuffer() any
	CreateFramebuffer() any
	CreateProgram() any
	CreateRenderbuffer() any
	CreateShader(xtype uint32) any
	CreateTexture() any
	CreateVertexArray() any

	DeleteBuffer(buffer any)
	DeleteFramebuffer(framebuffer any)
	DeleteProgram(progarm any)
	DeleteRenderbuffer(renderbuffer any)
	DeleteShader(shader any)
	DeleteTexture(texture any)
	DeleteVertexArray(vertexArray any)

	BindBuffer(target uint32, buffer any)
	BindFramebuffer(target uint32, framebuffer any)
	BindRenderbuffer(target uint32, renderbuffer any)
	BindTexture(target uint32, texture any)
	BindVertexArray(array any)

	AttachShader(program any, shader any)
	// WARNING: might override bound TEXTURE_2D in webgl
	BlitFramebuffer(srcX0 int32, srcY0 int32, srcX1 int32, srcY1 int32, dstX0 int32, dstY0 int32, dstX1 int32, dstY1 int32, mask uint32, filter uint32)
	BufferData(target uint32, data any, usage uint32)
	Clear(mask uint32)
	ClearColor(r, g, b, a float32)
	CompileShader(shader any)
	DrawElementsInstanced(mode uint32, count int32, xtype uint32, indexOffset uintptr, instancecount int32)
	EnableVertexAttribArray(index uint32)
	FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer any)
	FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture any, level int32)
	GetProgramInfoLog(program any) string
	GetProgramParameter(program any, pname uint32) int32
	GetShaderInfoLog(shader any) string
	GetShaderParameter(shader any, pname uint32) int32
	GetUniformLocation(program any, name string) any
	LinkProgram(program any)
	RenderbufferStorageMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32)
	ShaderSource(shader any, source string)
	TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels any)
	UseProgram(program any)
	VertexAttribDivisor(index uint32, divisor uint32)
	VertexAttribIPointer(index uint32, size int32, xtype uint32, stride int32, offset uintptr)
	VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uintptr)
	Viewport(x int32, y int32, width int32, height int32)

	Uniform1i(location any, v0 int32)
	Uniform2i(location any, v0, v1 int32)
	Uniform3i(location any, v0, v1, v2 int32)
	Uniform4i(location any, v0, v1, v2, v3 int32)
	Uniform1ui(location any, v0 uint32)
	Uniform2ui(location any, v0, v1 uint32)
	Uniform3ui(location any, v0, v1, v2 uint32)
	Uniform4ui(location any, v0, v1, v2, v3 uint32)
	Uniform1f(location any, v0 float32)
	Uniform2f(location any, v0, v1 float32)
	Uniform3f(location any, v0, v1, v2 float32)
	Uniform4f(location any, v0, v1, v2, v3 float32)

	// gl.Enable left to platform specific init functions
}
