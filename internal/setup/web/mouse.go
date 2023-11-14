package web

import (
	"syscall/js"

	"github.com/eliiasg/deltawing/input"
)

type mouse struct {
	moveHandler    func(float64, float64)
	clickHandler   func(input.MouseButton)
	releaseHandler func(input.MouseButton)
	scrollHandler  func(float64)
}

func makeMouse() *mouse {
	m := &mouse{}
	doc := js.Global().Get("document")
	doc.Call("addEventListener", "mousemove", js.FuncOf(m.jsMouseMoved), false)
	doc.Get("body").Call("addEventListener", "mousedown", js.FuncOf(m.jsMouseDown))
	doc.Get("body").Call("addEventListener", "mouseup", js.FuncOf(m.jsMouseUp))
	doc.Get("body").Call("addEventListener", "wheel", js.FuncOf(m.jsWheel))
	return m
}

func (m *mouse) SetMoveHandler(handler func(float64, float64)) {
	m.moveHandler = handler
}

func (m *mouse) SetClickHandler(handler func(input.MouseButton)) {
	m.clickHandler = handler
}

func (m *mouse) SetReleaseHandler(handler func(input.MouseButton)) {
	m.releaseHandler = handler
}

func (m *mouse) SetScrollHandler(handler func(float64)) {
	m.scrollHandler = handler
}

func (m *mouse) CursorLocked() bool {
	doc := js.Global().Get("document")
	return doc.Get("pointerLockElement").Equal(doc.Get("body"))
}

func (m *mouse) LockCursor() {
	js.Global().Get("document").Get("body").Call("requestPointerLock")
}

func (m *mouse) UnlockCursor() {
	js.Global().Get("document").Call("exitPointerLock")
}

func (m *mouse) SetCursorStyle(style input.CursorStyle) {
	m.UnlockCursor()
	var name string
	switch style {
	case input.CursorArrow:
		name = "default"
	case input.CursorIBeam:
		name = "text"
	case input.CursorCrossHair:
		name = "crosshair"
	case input.CursorHand:
		name = "pointer"
	case input.CursorHResize:
		name = "ew-resize"
	case input.CursorVResize:
		name = "ns-resize"
	case input.CursorHidden:
		name = "none"
	}
	js.Global().Get("document").Get("body").Get("style").Set("cursor", name)
}

func (m *mouse) jsMouseMoved(_ js.Value, args []js.Value) any {
	jsEvent := args[0]
	if m.moveHandler == nil {
		return nil
	}
	if m.CursorLocked() {
		m.moveHandler(jsEvent.Get("movementX").Float(), jsEvent.Get("movementY").Float())
	} else {
		m.moveHandler(jsEvent.Get("clientX").Float(), jsEvent.Get("clientY").Float())
	}
	return nil
}

func getButton(e js.Value) input.MouseButton {
	r := e.Get("button").Int()
	if r > 4 {
		return -1
	}
	// order is the same except for reversed right and middle
	if r == 1 {
		return input.MouseMiddle
	}
	if r == 2 {
		return input.MouseSecondary
	}
	return input.MouseButton(r)
}

func (m *mouse) jsMouseDown(_ js.Value, args []js.Value) any {
	jsEvent := args[0]
	if m.clickHandler != nil {
		m.clickHandler(getButton(jsEvent))
	}
	return nil
}

func (m *mouse) jsMouseUp(_ js.Value, args []js.Value) any {
	jsEvent := args[0]
	if m.releaseHandler != nil {
		m.releaseHandler(getButton(jsEvent))
	}
	return nil
}

func (m *mouse) jsWheel(_ js.Value, args []js.Value) any {
	jsEvent := args[0]
	v := jsEvent.Get("deltaY").Float()
	// get sign of value since browsers use different units
	// also reversed, since it should do the same as glfw impl
	if v > 0 {
		v = -1
	} else if v < 0 {
		v = 1
	}
	if m.scrollHandler != nil {
		m.scrollHandler(v)
	}
	return nil
}
