package glfw

import (
	"github.com/eliiasg/deltawing/input"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type keyboard struct {
	win         *glfw.Window
	keyPressed  func(input.Key)
	keyReleased func(input.Key)
	keyHeld     func(input.Key)
	charTyped   func(rune)
}

func makeKeyboard(win *glfw.Window) *keyboard {
	kb := &keyboard{win: win}
	win.SetKeyCallback(kb.keyCallback)
	win.SetCharCallback(kb.charCallback)
	return kb
}

func (k *keyboard) SetKeyPressedHandler(handler func(input.Key)) {
	k.keyPressed = handler
}

func (k *keyboard) SetKeyReleasedHandler(handler func(input.Key)) {
	k.keyReleased = handler
}

func (k *keyboard) SetKeyHeldHandler(handler func(input.Key)) {
	k.keyHeld = handler
}

func (k *keyboard) SetKeyTypedHandler(handler func(rune)) {
	k.charTyped = handler
}

func (k *keyboard) keyCallback(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		if k.keyPressed != nil {
			k.keyPressed(input.Key(key))
		}
	case glfw.Release:
		if k.keyReleased != nil {
			k.keyReleased(input.Key(key))
		}
	case glfw.Repeat:
		if k.keyHeld != nil {
			k.keyHeld(input.Key(key))
		}
	}
}

func (k *keyboard) charCallback(_ *glfw.Window, char rune) {
	if k.charTyped != nil {
		k.charTyped(char)
	}
}
