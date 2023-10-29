package gl

import (
	"github.com/eliiasg/deltawing/graphics/render"
)

type Renderer struct {
	Context Context
	primary *primaryRenderTarget
}

func (r *Renderer) PrimaryRenderTarget() render.RenderTarget {
	return r.primary
}

// should be called after gl and GLFW is initialized
// assumes primary rendertarget is set up properly
func NewRenderer(winWdith, winHeight func() uint16) render.Renderer {
	return &Renderer{
		primary: &primaryRenderTarget{&renderTarget{0, 0, 0, 0, 0, false}, winWdith, winHeight},
	}
}
