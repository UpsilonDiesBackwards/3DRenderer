package main

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/3DRenderer/engine/io"
	"github.com/UpsilonDiesBackwards/3DRenderer/engine/rendering"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
	"time"
)

var DeltaTime float64

func main() {
	runtime.LockOSThread()
	defer glfw.Terminate()

	window, err := rendering.NewWindow("Test")
	if err != nil {
		fmt.Println(err)
	}

	rend := rendering.NewRenderer(window)

	rend.NewObject("res/models/2b-dec.obj",
		"res/textures/UVChecker.png", "char")

	rend.GetObject("char").SetPosition(mgl32.Vec3{0, -1, -4})
	rend.GetObject("char").SetScale(mgl32.Vec3{1, 1, 1})

	var previousTime = time.Now()
	for !window.ShouldClose() {
		DeltaTime = CalculateDeltaTime(previousTime)
		previousTime = time.Now()

		io.InputRunner(window, DeltaTime, rend.Camera, 6)

		rend.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func CalculateDeltaTime(previousTime time.Time) float64 {
	currentTime := time.Now()
	deltaTime := currentTime.Sub(previousTime).Seconds()
	return deltaTime
}
