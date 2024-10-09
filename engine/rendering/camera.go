package rendering

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type Camera struct {
	Position    mgl64.Vec3
	Front       mgl64.Vec3
	Up          mgl64.Vec3
	Right       mgl64.Vec3
	WorldUp     mgl64.Vec3
	Yaw         float64
	Pitch       float64
	Speed       float64
	Sensitivity float64
	Fov         float32
}

func (c *Camera) UpdateDirection(dx, dy float64) {
	c.Yaw += dx * c.Sensitivity
	c.Pitch -= dy * c.Sensitivity

	if c.Pitch > 89 {
		c.Pitch = 89
	} else if c.Pitch < -89 {
		c.Pitch = -89
	}
	c.Yaw = math.Mod(c.Yaw, 360)
	c.UpdateVec()
}

func (c *Camera) UpdateVec() {
	c.Front = mgl64.Vec3{
		math.Cos(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch)), // X component
		math.Sin(mgl64.DegToRad(c.Pitch)),                                   // Y component
		math.Sin(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch)), // Z component
	}.Normalize()

	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

func (c *Camera) GetTransform() mgl32.Mat4 {
	cameraTarget := c.Position.Add(c.Front)
	return mgl32.LookAt(
		float32(c.Position.X()), float32(c.Position.Y()), float32(c.Position.Z()),
		float32(cameraTarget.X()), float32(cameraTarget.Y()), float32(cameraTarget.Z()),
		float32(c.Up.X()), float32(c.Up.Y()), float32(c.Up.Z()),
	)
}

func (c *Camera) GetFov() float32 {
	return c.Fov
}
