package opengl

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type renderTarget struct {
	framebufferID uint32
	textureID     uint32
	depthID       uint32
	width, height uint16
}

func (r *renderer) MakeRenderTarget(width, height uint16) render.RenderTarget {
	// make buffer
	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	// add texture
	texture := makeTextureBuffer(width, height)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)
	// add depth
	depth := makeDepthBuffer(width, height)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, depth, 0)
	t := &renderTarget{
		framebufferID: framebuffer,
		textureID:     texture,
		depthID:       depth,
	}
	// to init texture and depthbuffer
	t.Resize(width, height)
	return t
}

func makeDepthBuffer(width, height uint16) uint32 {
	// make
	var texture uint32
	gl.GenTextures(1, &texture)
	return texture
}

func makeTextureBuffer(width, height uint16) uint32 {
	// make
	var texture uint32
	gl.GenTextures(1, &texture)
	return texture
}

func (t *renderTarget) Free() {
	gl.DeleteFramebuffers(1, &t.framebufferID)
	gl.DeleteTextures(1, &t.textureID)
	gl.DeleteTextures(1, &t.depthID)
}

func (t *renderTarget) Width() uint16 {
	if t.framebufferID != 0 {
		return t.width
	} else {
		panic("do not get height of primary rendertarget")
	}
}

func (t *renderTarget) Height() uint16 {
	if t.framebufferID != 0 {
		return t.height
	} else {
		panic("do not get height of primary rendertarget")
	}
}

func (t *renderTarget) Clear(r uint8, g uint8, b uint8) {
	gl.ClearColor(float32(r)/256, float32(g)/256, float32(b)/256, 1.0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.framebufferID)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (t *renderTarget) Resize(width, height uint16) {
	if t.framebufferID == 0 {
		panic("do not resize primary rendertarget")
	}
	t.width = width
	t.height = height
	// bind texture
	gl.BindTexture(gl.TEXTURE_2D, t.textureID)
	// init texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	// bind depthbuffer
	gl.BindTexture(gl.TEXTURE_2D, t.framebufferID)
	// init depthbuffer
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32, int32(width), int32(height), 0, gl.DEPTH, gl.UNSIGNED_INT, nil)
}

func (t *renderTarget) BlitTo(target render.RenderTarget, x, y uint16) {
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, t.framebufferID)
	// funny convertion, but should only be called with renderTarget from same renderer
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, target.(*renderTarget).framebufferID)
	gl.BlitFramebuffer(0, 0, int32(t.width), int32(t.height), int32(x), int32(y), int32(x+t.width), int32(y+t.height), gl.COLOR_BUFFER_BIT, gl.LINEAR)
}

func (t *renderTarget) DrawTo(target render.RenderTarget, x uint16, y uint16, width uint16, height uint16, pivotX uint16, pivotY uint16, rotation float32, shader render.FragmentShader) {
	panic("not implemented") // TODO: Implement
}
