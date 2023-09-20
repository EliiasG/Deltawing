package vec

import "github.com/eliiasg/deltawing/math"

type Vec2[T math.Signed] struct {
	X T
	Y T
}

func MakeVec2[T math.Signed](x, y T) Vec2[T] {
	return Vec2[T]{
		X: x,
		Y: y,
	}
}

func Add[T math.Signed](v1, v2 Vec2[T]) Vec2[T] {
	return Vec2[T]{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}
