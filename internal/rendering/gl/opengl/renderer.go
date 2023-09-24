package opengl

import (
	"github.com/eliiasg/deltawing/graphics/render"
)

type Renderer struct {
}

func (r *Renderer) MakeRenderTarget() render.RenderTarget {
	panic("not implemented") // TODO: Implement
}

func (r *Renderer) MakeProcedureBuilder() render.ProcedureBuilder {
	panic("not implemented") // TODO: Implement
}

func (r *Renderer) MakeOperation(precedure render.Procedure) render.Opteration {
	panic("not implemented") // TODO: Implement
}

func (r *Renderer) MakeFragmentShader(source string) render.FragmentShader {
	panic("not implemented") // TODO: Implement
}

func (r *Renderer) PrimaryRenderTarget() render.RenderTarget {
	panic("not implemented") // TODO: Implement
}
