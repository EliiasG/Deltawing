package input

// GLFW does not seem to support controller rumble, not sure about web
type Controller interface {
	// nil if no controller
	GetState(id uint16) ControllerState
}

type ControllerState interface {
	GetAxis(axis ControllerAxis) float32
	GetButton(button ControllerButton) bool
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
	ControllerDPadRight
	ControllerDPadDown
	ControllerDPadLeft
)
