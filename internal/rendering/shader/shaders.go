package shader

import (
	_ "embed"

	"github.com/eliiasg/deltawing/graphics/render"
)

const VertexBaseInputAmt = 2

//go:embed vertex.glsl
var VertexBaseSource string

// not called base, sice it should not be modified
//
//go:embed fragment.glsl
var FragmentSource string

func init() {
	// to avoid sahder comp error
	FragmentSource += "\x00"
}

func IsInt(t render.ChannelShaderType) bool {
	return t == render.ShaderInt || t == render.ShaderUnsignedInt
}
