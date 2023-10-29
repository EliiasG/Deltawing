package shader_sources

import (
	_ "embed"
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
	// tecnically not required, since ShaderBuilder adds end automatically, but seems nice to do it here
	VertexBaseSource += "\x00"
}
