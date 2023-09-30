package shader

import _ "embed"

const VertexBaseInputAmt = 2

//go:embed vertex.glsl
var VertexBaseSource string

// not called base, sice it should not be modified
//go:embed fragment.glsl
var FragmentSource string
