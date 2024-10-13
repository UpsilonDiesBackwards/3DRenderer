package rendering

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/3DRenderer/tools"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"strings"
	"time"
)

type Renderer struct {
	Window  *Window
	Objects map[string]*RenderableObject
	Camera  *Camera
	Shader  *Shader

	project  mgl32.Mat4
	lastTime time.Time
}

func NewRenderer(window *Window) *Renderer {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	winWidth := window.FramebufferSize()[0]
	winHeight := window.FramebufferSize()[1]
	gl.Viewport(0, 0, int32(winWidth), int32(winHeight))

	gl.ClearColor(0.52, 0.80, 0.96, 1.0)

	projection := mgl32.Perspective(mgl32.DegToRad(60), float32(winWidth)/float32(winHeight), 0.1, 2000.0)

	shader, err := NewShader("res/shaders/shader.vert", "res/shaders/shader.frag")
	if err != nil {
		fmt.Println("Error initializing OpenGL shader: ", err)
	}

	return &Renderer{
		Window:  window,
		Objects: make(map[string]*RenderableObject),
		Camera: &Camera{
			Position:    mgl64.Vec3{0, 0, 0},
			Up:          mgl64.Vec3{0, 1, 0},
			WorldUp:     mgl64.Vec3{0, 1, 0},
			Yaw:         -90,
			Pitch:       0,
			Speed:       12,
			Sensitivity: 0.075,
			Fov:         60,
		},
		Shader:   shader,
		project:  projection,
		lastTime: time.Now(),
	}
}

func (r *Renderer) NewObject(filePath, mtlPath, name string) {
	if mtlPath == "" {
		mtlPath = strings.Replace(filePath, ".obj", ".mtl", 1)
	}

	model := tools.CreateNewOBJ(filePath, mtlPath)

	renderableObject := NewRenderableObject(model, mtlPath)

	r.AddNewObject(renderableObject, name)
}

func (r *Renderer) AddNewObject(object *RenderableObject, name string) {
	r.Objects[name] = object
}

func (r *Renderer) GetObject(name string) *RenderableObject {
	if r.Objects[name] == nil {
		fmt.Println("Object not found: ", name)
	}

	return r.Objects[name]
}

func (r *Renderer) Draw() {
	r.CalculateDeltaTime()

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	r.Shader.Use()

	r.Shader.SetMat4ByName("projection", r.project)
	r.Shader.SetMat4ByName("view", r.Camera.GetTransform())

	for _, object := range r.Objects {
		object.Draw(r.Shader)
	}

	r.Window.SwapBuffers()
	glfw.PollEvents()
}

func (r *Renderer) CalculateDeltaTime() {
	currentTime := time.Now()
	currentTime.Sub(r.lastTime).Seconds()
	r.lastTime = currentTime
}
