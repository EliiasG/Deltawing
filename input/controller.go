package input

// GLFW does not seem to support controller rumble, not sure about web
type Controller interface {
	GetAxis(controllerID uint16, axis ControllerAxis)
	GetButton(controllerID uint16, button ControllerButton)
}

type ControllerAxis int32

const (
	AxisLeftX ControllerAxis = iota
	AxisLeftY
	AxisRightX
	AxisRightY
	AxisLeftTrigger
	AxisRightTrigger
)

type ControllerButton int32

const (
	ControllerA ControllerButton = iota
	ControllerB
	ControllerX
	ControllerY
	ControllerLeftBumper
	ControllerRightBumper
	// not sure about some of these
	ControllerBack
	ControllerStart
	ControllerGuide
	ControllerLeftThumb
	ControllerRightThumb
	ControllerDPadUp
	ControllerDPadDown
	ControllerDPadLeft
	ControllerDPadRight
)
