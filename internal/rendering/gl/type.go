package gl

import "github.com/eliiasg/deltawing/graphics/render"

// Welcome to graphics programming, it's super fun
func AssertType(typ render.ShaderType, val any) bool {
	switch val.(type) {
	case int32:
		return checkType(typ, render.ShaderInt, 1)
	case uint32:
		return checkType(typ, render.ShaderUnsignedInt, 1)
	case float32:
		return checkType(typ, render.ShaderFloat, 1)
	case [2]int32:
		return checkType(typ, render.ShaderInt, 2)
	case [2]uint32:
		return checkType(typ, render.ShaderUnsignedInt, 2)
	case [2]float32:
		return checkType(typ, render.ShaderFloat, 2)
	case [3]int32:
		return checkType(typ, render.ShaderInt, 3)
	case [3]uint32:
		return checkType(typ, render.ShaderUnsignedInt, 3)
	case [3]float32:
		return checkType(typ, render.ShaderFloat, 3)
	case [4]int32:
		return checkType(typ, render.ShaderInt, 4)
	case [4]uint32:
		return checkType(typ, render.ShaderUnsignedInt, 4)
	case [4]float32:
		return checkType(typ, render.ShaderFloat, 4)
	default:
		return false
	}
}

func checkType(typ render.ShaderType, typTyp render.ChannelShaderType, amt uint8) bool {
	return typ.Type == typTyp && typ.Amount == amt
}
