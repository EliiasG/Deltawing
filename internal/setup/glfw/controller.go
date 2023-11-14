//go:build cgo
// +build cgo

package glfw

import (
	"github.com/eliiasg/deltawing/input"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type controller struct{}

func makeController() controller {
	return controller{}
}

// nil if no controller
func (c controller) GetState(id uint16) input.ControllerState {
	joy := glfw.Joystick(id)
	if !joy.IsGamepad() {
		return nil
	}
	return controllerState{joy.GetGamepadState()}
}

type controllerState struct {
	*glfw.GamepadState
}

func (c controllerState) GetAxis(axis input.ControllerAxis) float32 {
	if axis < 0 || axis > 5 {
		panic("invalid axis")
	}
	return c.Axes[axis]
}

func (c controllerState) GetButton(button input.ControllerButton) bool {
	if button < 0 || button > 14 {
		panic("invalid button")
	}
	return c.Buttons[button] == glfw.Press
}
