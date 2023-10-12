package input

type Mouse interface {
	SetMouseMoveHandler(handler func(float64, float64))
	SetMouseClickHandler(handler func(MouseButton))
	SetMouseReleaseHandler(handler func(MouseButton))
	SetMouseScrollHandler(handler func(float64))

	SetCursorStyle(CursorStyle)
}

type CursorStyle int32

const (
	Arrow CursorStyle = iota
	IBeam
	CrossHair
	Hand
	HResize
	VResize
)

type MouseButton int32

const (
	MousePrimary MouseButton = iota
	MouseSecondary
	MouseMiddle
	MouseExtra1
	MouseExtra2
)
