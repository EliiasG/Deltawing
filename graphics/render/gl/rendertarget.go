package gl

import (
	"github.com/eliiasg/deltawing/graphics/render"
	"github.com/eliiasg/glow/enum"
)

type RenderTarget struct {
	cxt Context
	// Framebuffer object
	Framebuffer any
	// DrawBuffer object, this is either a texture or a renderbuffer (if multisampled)
	DrawBuffer any
	// DepthBuffer object, again eithr a texture or a renderbuffer (if multisampled)
	DepthBuffer   any
	Multisample   bool
	width, height uint16
}

type primaryRenderTarget struct {
	*RenderTarget
	widthFunc, heightFunc func() uint16
}

func GLRenderTarget(target render.RenderTarget) (*RenderTarget, bool) {
	switch t := target.(type) {
	case *RenderTarget:
		return t, true
	case *primaryRenderTarget:
		return t.RenderTarget, true
	}
	return nil, false
}

func (r *Renderer) MakeRenderTarget(width, height uint16, multisample bool) render.RenderTarget {
	// make buffer
	framebuffer := r.cxt.CreateFramebuffer()
	r.cxt.BindFramebuffer(enum.FRAMEBUFFER, framebuffer)
	// add texture
	texture := makeTextureBuffer(r.cxt, multisample)
	// add depth
	depth := makeTextureBuffer(r.cxt, multisample)
	t := &RenderTarget{
		cxt:         r.cxt,
		Framebuffer: framebuffer,
		DrawBuffer:  texture,
		DepthBuffer: depth,
		Multisample: multisample,
	}
	// to init texture and depthbuffer
	t.Resize(width, height)
	return t
}

func makeTextureBuffer(cxt Context, multisample bool) any {
	// make
	var texture any
	if multisample {
		texture = cxt.CreateRenderbuffer()
	} else {
		texture = cxt.CreateTexture()
	}

	return texture
}

func (t *RenderTarget) Free() {
	t.cxt.DeleteFramebuffer(t.Framebuffer)
	if t.Multisample {
		t.cxt.DeleteRenderbuffer(t.DrawBuffer)
		t.cxt.DeleteRenderbuffer(t.DepthBuffer)
	} else {
		t.cxt.DeleteTexture(t.DrawBuffer)
		t.cxt.DeleteTexture(t.DepthBuffer)
	}

}

func (t *RenderTarget) Width() uint16 {
	return t.width
}

func (t *primaryRenderTarget) Width() uint16 {
	return t.widthFunc()
}

func (t *RenderTarget) Height() uint16 {
	return t.height
}

func (t *primaryRenderTarget) Height() uint16 {
	return t.heightFunc()
}

func (t *RenderTarget) Clear(r uint8, g uint8, b uint8) {
	t.cxt.BindFramebuffer(enum.FRAMEBUFFER, t.Framebuffer)
	t.cxt.ClearColor(float32(r)/256, float32(g)/256, float32(b)/256, 1.0)
	t.cxt.Clear(enum.COLOR_BUFFER_BIT | enum.DEPTH_BUFFER_BIT)
}

func (t *RenderTarget) Resize(width, height uint16) {
	if t.Framebuffer == 0 {
		panic("do not resize primary rendertarget")
	}
	t.width = width
	t.height = height
	t.cxt.BindFramebuffer(enum.FRAMEBUFFER, t.Framebuffer)
	if t.Multisample {
		t.resizeMultisample(width, height)
	} else {
		t.resizeNormal(width, height)
	}
}

func (t *RenderTarget) resizeNormal(width, height uint16) {
	// bind texture
	t.cxt.BindTexture(enum.TEXTURE_2D, t.DrawBuffer)
	// init texture
	t.cxt.TexImage2D(enum.TEXTURE_2D, 0, enum.RGB, int32(width), int32(height), 0, enum.RGB, enum.UNSIGNED_BYTE, nil)
	// bind depthbuffer
	t.cxt.BindTexture(enum.TEXTURE_2D, t.DepthBuffer)
	// init depthbuffer
	t.cxt.TexImage2D(enum.TEXTURE_2D, 0, enum.DEPTH_COMPONENT32, int32(width), int32(height), 0, enum.DEPTH_COMPONENT, enum.UNSIGNED_INT, nil)
	// add to framebuffer
	t.cxt.FramebufferTexture2D(enum.FRAMEBUFFER, enum.COLOR_ATTACHMENT0, enum.TEXTURE_2D, t.DrawBuffer, 0)
	t.cxt.FramebufferTexture2D(enum.FRAMEBUFFER, enum.DEPTH_ATTACHMENT, enum.TEXTURE_2D, t.DepthBuffer, 0)
}

func (t *RenderTarget) resizeMultisample(width, height uint16) {
	// bind texture
	t.cxt.BindRenderbuffer(enum.RENDERBUFFER, t.DrawBuffer)
	// init texture
	t.cxt.RenderbufferStorageMultisample(enum.RENDERBUFFER, 4, enum.RGB, int32(width), int32(height))
	// bind depthbuffer
	t.cxt.BindRenderbuffer(enum.RENDERBUFFER, t.DepthBuffer)
	// init depthbuffer
	t.cxt.RenderbufferStorageMultisample(enum.RENDERBUFFER, 4, enum.DEPTH_COMPONENT32, int32(width), int32(height))
	// add to framebuffer
	t.cxt.FramebufferRenderbuffer(enum.FRAMEBUFFER, enum.COLOR_ATTACHMENT0, enum.RENDERBUFFER, t.DrawBuffer)
	t.cxt.FramebufferRenderbuffer(enum.FRAMEBUFFER, enum.DEPTH_ATTACHMENT, enum.RENDERBUFFER, t.DepthBuffer)
}

func (t *RenderTarget) BlitTo(target render.RenderTarget, x, y int32) {
	tar, _ := GLRenderTarget(target)
	if tar.Multisample {
		panic("Do not blit to multisampled target!")
	}
	t.cxt.BindFramebuffer(enum.READ_FRAMEBUFFER, t.Framebuffer)
	t.cxt.BindFramebuffer(enum.DRAW_FRAMEBUFFER, tar.Framebuffer)
	// using target, because it might be a primarytarget
	y = int32(target.Height()) - int32(t.Height()) - y
	t.cxt.BlitFramebuffer(0, 0, int32(t.Width()), int32(t.Height()), int32(x), int32(y), x+int32(t.Width()), y+int32(t.Height()), enum.COLOR_BUFFER_BIT, enum.LINEAR)
}
