package glfw

import "github.com/go-gl/glfw/v3.3/glfw"

type window struct {
	win *glfw.Window
}

func (w *window) SetSize(width uint16, height uint16) {
	w.win.SetSize(int(width), int(height))
}

func (w *window) SetMaximized(maximized bool) {
	w.win.Maximize()
}

func (w *window) SetFullScreen(fullscreen bool) {
	panic("not implementd") // TODO: Implement
}

func (w *window) SetSizeShanged(handler func(uint16, uint16)) {
	panic("not implemented") // TODO: Implement
}

func (w *window) WindowSize() (uint16, uint16) {
	width, height := w.win.GetSize()
	return uint16(width), uint16(height)
}

func (w *window) UpdateView() {
	w.win.SwapBuffers()
	glfw.PollEvents()
}

func (w *window) ShouldClose() bool {
	return w.win.ShouldClose()
}
