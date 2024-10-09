package io

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/3DRenderer/engine/rendering"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var ViewportTransform mgl32.Mat4
var u = &UserInput{}

var isZooming = false

func InputRunner(win *rendering.Window, deltaTime float64, camera *rendering.Camera, zoomMult float64) error {
	var adjCSpeed float64

	if isZooming {
		adjCSpeed = deltaTime * float64(camera.Speed*zoomMult)
	} else {
		adjCSpeed = deltaTime * float64(camera.Speed)
	}

	if ActionState[VP_FORW] {
		camera.Position = camera.Position.Add(camera.Front.Mul(adjCSpeed))
	}
	if ActionState[VP_BACK] {
		camera.Position = camera.Position.Sub(camera.Front.Mul(adjCSpeed))
	}
	if ActionState[VP_LEFT] {
		camera.Position = camera.Position.Sub(camera.Front.Cross(camera.Up).Mul(adjCSpeed))
	}
	if ActionState[VP_RGHT] {
		camera.Position = camera.Position.Add(camera.Front.Cross(camera.Up).Mul(adjCSpeed))
	}
	if ActionState[VP_UP] {
		camera.Position = camera.Position.Add(camera.Up.Mul(adjCSpeed))
	}
	if ActionState[VP_DOWN] {
		camera.Position = camera.Position.Sub(camera.Up.Mul(adjCSpeed))
	}

	if ActionState[ED_ZOOM] {
		isZooming = !isZooming
	}

	if ActionState[ED_QUIT] {
		fmt.Println("Exiting!")
		glfw.Terminate()
	}

	camera.UpdateDirection(u.CursorChange().X(), u.CursorChange().Y())
	u.CheckpointCursorChange()
	ViewportTransform = camera.GetTransform()
	InputManager(win, u)

	return nil
}

func GetUserInput() *UserInput {
	return u
}
