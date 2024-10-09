package rendering

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	window      *glfw.Window
	aspectRatio float32
}

func NewWindow(title string) (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	} else {
		fmt.Println("Initialised GLFW")
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Maximized, glfw.True)

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	win, err := glfw.CreateWindow(mode.Width, mode.Height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	win.MakeContextCurrent()

	glfw.SwapInterval(0)

	win.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	return &Window{window: win, aspectRatio: float32(mode.Width) / float32(mode.Height)}, nil
}

func (w *Window) ShouldClose() bool {
	return w.window.ShouldClose()
}

func (w *Window) SwapBuffers() {
	w.window.SwapBuffers()
}

func (w *Window) PollEvents() {
	glfw.PollEvents()
}

func (w *Window) SetKeyCallback(callback glfw.KeyCallback) {
	w.window.SetKeyCallback(callback)
}

func (w *Window) SetCursorPosCallback(posCallback glfw.CursorPosCallback) {
	w.window.SetCursorPosCallback(posCallback)
}

func (w *Window) SetSizeCallback(callback glfw.SizeCallback) {
	w.window.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		w.aspectRatio = float32(width) / float32(height)
	})
}

func (w *Window) AspectRatio() float32 {
	return w.aspectRatio
}

func (w *Window) DisplaySize() [2]float32 {
	width, height := w.window.GetSize()
	return [2]float32{float32(width), float32(height)}
}

func (w *Window) FramebufferSize() [2]float32 {
	width, height := w.window.GetFramebufferSize()
	return [2]float32{float32(width), float32(height)}
}
