package web

import "github.com/eliiasg/deltawing/input"

type keyboard struct{}

func (k *keyboard) SetKeyPressedHandler(handler func(input.Key)) {
	// TODO URGENT!
	// FIXME not panicing for testing
}

func (k *keyboard) SetKeyReleasedHandler(handler func(input.Key)) {
	panic("not implemented") // TODO: Implement
}

func (k *keyboard) SetKeyHeldHandler(handler func(input.Key)) {
	panic("not implemented") // TODO: Implement
}

func (k *keyboard) SetKeyTypedHandler(handler func(rune)) {
	panic("not implemented") // TODO: Implement
}
