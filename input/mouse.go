package input

type Mouse interface {
	SetMoveHandler(handler func(float64, float64))
	SetClickHandler(handler func(MouseButton))
	SetReleaseHandler(handler func(MouseButton))
	// positive values mean upwards scroll
	SetScrollHandler(handler func(float64))
	CursorLocked() bool
	// Locks and hides the cursor, this will make the arguments of the mouse handler be the amount the mouse moved
	LockCursor()
	UnlockCursor()
	// unlocks cursor
	SetCursorStyle(style CursorStyle)
}

type CursorStyle int32

const (
	CursorArrow CursorStyle = iota
	CursorIBeam
	CursorCrossHair
	CursorHand
	CursorHResize
	CursorVResize
	CursorHidden
)

type MouseButton int32

const (
	MousePrimary MouseButton = iota
	MouseSecondary
	MouseMiddle
	MouseExtra1
	MouseExtra2
)
