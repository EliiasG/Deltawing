package opengl

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/glow/v3.3-core/gl"
)

type renderTarget struct {
	framebufferID uint32
	textureID     uint32
	depthID       uint32
	width, height uint16
	multisample   bool
}

type primaryRenderTarget struct {
	*renderTarget
	widthFunc, heightFunc func() uint16
}

func (r *renderer) MakeRenderTarget(width, height uint16, multisample bool) render.RenderTarget {
	// make buffer
	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	// add texture
	texture := makeTextureBuffer(width, height, multisample)
	// add depth
	depth := makeTextureBuffer(width, height, multisample)
	t := &renderTarget{
		framebufferID: framebuffer,
		textureID:     texture,
		depthID:       depth,
		multisample:   multisample,
	}
	// to init texture and depthbuffer
	t.Resize(width, height)
	return t
}

func makeTextureBuffer(width, height uint16, multisample bool) uint32 {
	// make
	var texture uint32
	if multisample {
		gl.GenRenderbuffers(1, &texture)
	} else {
		gl.GenTextures(1, &texture)
	}

	return texture
}

func (t *renderTarget) Free() {
	gl.DeleteFramebuffers(1, &t.framebufferID)
	if t.multisample {
		gl.DeleteRenderbuffers(1, &t.textureID)
		gl.DeleteRenderbuffers(1, &t.depthID)
	} else {
		gl.DeleteTextures(1, &t.textureID)
		gl.DeleteTextures(1, &t.depthID)
	}

}

func (t *renderTarget) Width() uint16 {
	return t.width
}

func (t *primaryRenderTarget) Width() uint16 {
	return t.widthFunc()
}

func (t *renderTarget) Height() uint16 {
	return t.height
}

func (t *primaryRenderTarget) Height() uint16 {
	return t.heightFunc()
}

func (t *renderTarget) Clear(r uint8, g uint8, b uint8) {
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.framebufferID)
	gl.ClearColor(float32(r)/256, float32(g)/256, float32(b)/256, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (t *renderTarget) Resize(width, height uint16) {
	if t.framebufferID == 0 {
		panic("do not resize primary rendertarget")
	}
	t.width = width
	t.height = height
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.framebufferID)
	if t.multisample {
		t.resizeMultisample(width, height)
	} else {
		t.resizeNormal(width, height)
	}
}

func (t *renderTarget) resizeNormal(width, height uint16) {
	// bind texture
	gl.BindTexture(gl.TEXTURE_2D, t.textureID)
	// init texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	// bind depthbuffer
	gl.BindTexture(gl.TEXTURE_2D, t.depthID)
	// init depthbuffer
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT32, int32(width), int32(height), 0, gl.DEPTH_COMPONENT, gl.UNSIGNED_INT, nil)
	// add to framebuffer
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, t.textureID, 0)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, t.depthID, 0)
}

func (t *renderTarget) resizeMultisample(width, height uint16) {
	// bind texture
	gl.BindRenderbuffer(gl.RENDERBUFFER, t.textureID)
	// init texture
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, 4, gl.RGB, int32(width), int32(height))
	// bind depthbuffer
	gl.BindRenderbuffer(gl.RENDERBUFFER, t.depthID)
	// init depthbuffer
	gl.RenderbufferStorageMultisample(gl.RENDERBUFFER, 4, gl.DEPTH_COMPONENT32, int32(width), int32(height))
	// add to framebuffer
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.RENDERBUFFER, t.textureID)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, t.depthID)
}

func getRenderTarget(target render.RenderTarget) *renderTarget {
	switch t := target.(type) {
	case *renderTarget:
		return t
	case *primaryRenderTarget:
		return t.renderTarget
	}
	return nil
}

func (t *renderTarget) BlitTo(target render.RenderTarget, x, y int32) {
	tar := getRenderTarget(target)
	if tar.multisample {
		panic("Do not blit to multisampled target!")
	}
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, t.framebufferID)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, tar.framebufferID)
	// using target, because it might be a primarytarget
	y = int32(target.Height()) - int32(t.Height()) - y
	gl.BlitFramebuffer(0, 0, int32(t.Width()), int32(t.Height()), int32(x), int32(y), x+int32(t.Width()), y+int32(t.Height()), gl.COLOR_BUFFER_BIT, gl.LINEAR)
}
