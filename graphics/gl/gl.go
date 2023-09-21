package gl

//TODO wip file
import (
	"runtime"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}
