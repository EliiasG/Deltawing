package opengl

import (
	"github.com/eliiasg/deltawing/graphics/render"
)

type renderer struct {
	primary *renderTarget
}

func (r *renderer) PrimaryRenderTarget() render.RenderTarget {
	return r.primary
}

// should be called after gl and GLFW is initialized
// assumes primary rendertarget is set up properly
func NewRenderer() render.Renderer {
	return &renderer{&renderTarget{0, 0, 0, 0, 0}}
}
