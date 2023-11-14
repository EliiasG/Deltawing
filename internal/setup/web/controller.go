//go:build wasm
// +build wasm

package web

import (
	"syscall/js"

	"github.com/eliiasg/deltawing/input"
)

type controller struct{}

func makeController() *controller {
	return &controller{}
}

// nil if no controller
func (c *controller) GetState(id uint16) input.ControllerState {
	jsCont := getController(int(id))
	if jsCont.IsNull() {
		return nil
	}
	return getState(jsCont)
}

// using index as js button index and value as result button
var buttonMap = [17]input.ControllerButton{
	input.ControllerA,
	input.ControllerB,
	input.ControllerX,
	input.ControllerY,
	input.ControllerLeftBumper,
	input.ControllerRightBumper,
	-1,
	-1,
	input.ControllerBack,
	input.ControllerStart,
	input.ControllerLeftThumb,
	input.ControllerRightThumb,
	input.ControllerDPadUp,
	input.ControllerDPadDown,
	input.ControllerDPadLeft,
	input.ControllerDPadRight,
	input.ControllerGuide,
}

func getState(controller js.Value) *controllerState {
	res := new(controllerState)
	// buttons
	jsButtons := controller.Get("buttons")
	res.buttons = make([]bool, 15)
	for i, e := range buttonMap {
		if e == -1 {
			continue
		}
		res.buttons[e] = jsButtons.Index(i).Get("pressed").Bool()
	}
	// axes
	jsAxes := controller.Get("axes")
	res.axes = []float32{
		float32(jsAxes.Index(0).Float()),
		float32(jsAxes.Index(1).Float()),
		float32(jsAxes.Index(2).Float()),
		float32(jsAxes.Index(3).Float()),
		float32(2*jsButtons.Index(6).Get("value").Float() - 1),
		float32(2*jsButtons.Index(7).Get("value").Float() - 1),
	}
	return res
}

func getController(index int) js.Value {
	gamepads := js.Global().Get("window").Get("navigator").Call("getGamepads")
	return gamepads.Index(index)
}

type controllerState struct {
	buttons []bool
	axes    []float32
}

func (c *controllerState) GetAxis(axis input.ControllerAxis) float32 {
	if axis < 0 || axis > 5 {
		panic("invalid axis")
	}
	return c.axes[axis]
}

func (c *controllerState) GetButton(button input.ControllerButton) bool {
	if button < 0 || button > 14 {
		panic("invalid button")
	}
	return c.buttons[button]
}
