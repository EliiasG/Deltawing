//go:build cgo
// +build cgo

package glfw

import (
	"github.com/eliiasg/deltawing/input"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type mouse struct {
	win            *glfw.Window
	moveHandler    func(float64, float64)
	clickHandler   func(input.MouseButton)
	releaseHandler func(input.MouseButton)
	scrollHandler  func(float64)
	style          input.CursorStyle
	locked         bool
	oldX, oldY     float64

	cursorIBeam     *glfw.Cursor
	cursorCrossHair *glfw.Cursor
	cursorHand      *glfw.Cursor
	cursorHResize   *glfw.Cursor
	cursorVResize   *glfw.Cursor
}

func makeMouse(win *glfw.Window) *mouse {
	m := &mouse{win: win}
	win.SetCursorPosCallback(m.moveCallback)
	win.SetMouseButtonCallback(m.buttonCallback)
	win.SetScrollCallback(m.scrollCallback)
	m.cursorIBeam = glfw.CreateStandardCursor(glfw.IBeamCursor)
	m.cursorCrossHair = glfw.CreateStandardCursor(glfw.CrosshairCursor)
	m.cursorHand = glfw.CreateStandardCursor(glfw.HandCursor)
	m.cursorHResize = glfw.CreateStandardCursor(glfw.HResizeCursor)
	m.cursorVResize = glfw.CreateStandardCursor(glfw.VResizeCursor)
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

func (m *mouse) LockCursor() {
	m.locked = true
	m.win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func (m *mouse) UnlockCursor() {
	m.SetCursorStyle(m.style)
}

func (m *mouse) SetCursorStyle(style input.CursorStyle) {
	m.locked = false
	m.style = style
	if style == input.CursorHidden {
		m.win.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		return
	}
	m.win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	switch style {
	case input.CursorArrow:
		m.win.SetCursor(nil)
	case input.CursorIBeam:
		m.win.SetCursor(m.cursorIBeam)
	case input.CursorCrossHair:
		m.win.SetCursor(m.cursorCrossHair)
	case input.CursorHand:
		m.win.SetCursor(m.cursorHand)
	case input.CursorHResize:
		m.win.SetCursor(m.cursorHResize)
	case input.CursorVResize:
		m.win.SetCursor(m.cursorVResize)
	}
}

func (m *mouse) CursorLocked() bool {
	return m.locked
}

func (m *mouse) moveCallback(_ *glfw.Window, x, y float64) {
	if m.moveHandler == nil {
		return
	}
	var rx, ry float64
	if m.locked {
		rx = x - m.oldX
		ry = y - m.oldY
	} else {
		rx, ry = x, y
	}
	m.oldX, m.oldY = x, y
	m.moveHandler(rx, ry)
}

func (m *mouse) buttonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	// unsupported button
	if button > 4 {
		return
	}
	// else should only fire if action is release
	if action == glfw.Press {
		if m.clickHandler != nil {
			// should map directly
			m.clickHandler(input.MouseButton(button))
		}
	} else if m.releaseHandler != nil {
		m.releaseHandler(input.MouseButton(button))
	}
}

func (m *mouse) scrollCallback(_ *glfw.Window, xoffset, yoffset float64) {
	// horizontal scroll not supported
	m.scrollHandler(yoffset)
}
