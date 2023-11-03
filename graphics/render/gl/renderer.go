package gl

import (
	"github.com/eliiasg/deltawing/graphics/render"
)

type Renderer struct {
	cxt     Context
	primary *primaryRenderTarget
	// retured by PrimaryRendertarget, used by webgl since cannot blit multisampled to real primary
	primaryOverride render.RenderTarget
	version         string
}

// doing it like this since some types might be extended (like primaryRenderTarget)
func GLRenderer(r render.Renderer) (*Renderer, bool) {
	res, ok := r.(*Renderer)
	return res, ok
}

// renderer types do not have this method, since you can just get the context here and save it
func (r *Renderer) Context() Context {
	return r.cxt
}

func (r *Renderer) GLSLVersion() string {
	return r.version
}

func (r *Renderer) PrimaryRenderTarget() render.RenderTarget {
	return r.primaryOverride
}

func (r *Renderer) RealPrimaryRenderTarget() render.RenderTarget {
	return r.primary
}

// should be called after gl and GLFW is initialized
// assumes primary rendertarget is set up properly
func NewRenderer(winWdith, winHeight func() uint16, cxt Context, version string, overrideTarget bool) *Renderer {
	rend := &Renderer{
		primary: &primaryRenderTarget{&RenderTarget{cxt, nil, nil, nil, false, 0, 0}, winWdith, winHeight},
		cxt:     cxt,
		version: version,
	}
	if overrideTarget {
		rend.primaryOverride = rend.MakeRenderTarget(1, 1, false)
	} else {
		rend.primaryOverride = rend.primary
	}
	return rend
}
