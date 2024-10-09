package rendering

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/3DRenderer/engine/common"
	"github.com/UpsilonDiesBackwards/3DRenderer/tools"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type RenderableObject struct {
	VAO       uint32
	VBO       uint32
	EBO       uint32
	Vertices  []float32
	Normals   []float32
	TexCoords []float32
	Indices   []uint32

	Material    *common.Material
	Texture     uint32
	ModelMatrix mgl32.Mat4
}

func NewRenderableObject(obj *common.ObjectPrimitive, textPath string) *RenderableObject {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	combinedVertices := CombineVertices(obj)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(combinedVertices)*4, gl.Ptr(combinedVertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(obj.Indices)*4, gl.Ptr(obj.Indices), gl.STATIC_DRAW)

	stride := int32(32)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, stride, gl.PtrOffset(18))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	var texture uint32
	var err error
	if textPath == "" {
		texture = tools.CreateWhiteTexture()
		fmt.Println("Could not find texture path for object. Using white texture.")
	} else {
		texture, err = tools.LoadTexture(textPath)
		if err != nil {
			fmt.Println("Failed to load texture from", textPath)
			return nil
		}
	}
	return &RenderableObject{
		VAO:         vao,
		VBO:         vbo,
		EBO:         ebo,
		Vertices:    obj.Vertices,
		Normals:     obj.Normals,
		TexCoords:   obj.UVs,
		Indices:     obj.Indices,
		ModelMatrix: mgl32.Ident4(),
		Texture:     texture,
	}
}

func (obj *RenderableObject) Draw(shader *Shader) {
	gl.BindVertexArray(obj.VAO)

	shader.SetMat4ByName("model", obj.ModelMatrix)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, obj.Texture)
	shader.SetInt("texture0", 0)

	gl.DrawElements(gl.TRIANGLES, int32(len(obj.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindVertexArray(0)
}

func (obj *RenderableObject) SetPosition(position mgl32.Vec3) {
	if obj != nil {
		obj.ModelMatrix = mgl32.Translate3D(position.X(), position.Y(), position.Z())
	}
}

func (obj *RenderableObject) SetRotation(rotation mgl32.Quat) {
	if obj != nil {
		rotationMatrix := rotation.Mat4()
		obj.ModelMatrix = obj.ModelMatrix.Mul4(rotationMatrix)
	}
}

func (obj *RenderableObject) SetScale(scale mgl32.Vec3) {
	if obj != nil {
		scaleMatrix := mgl32.Scale3D(scale.X(), scale.Y(), scale.Z())
		obj.ModelMatrix = obj.ModelMatrix.Mul4(scaleMatrix)
	}
}

func (obj *RenderableObject) SetColor(R, G, B, A uint8) {
	if obj != nil {
		obj.Texture = tools.CreateColorMaterial(R, G, B, A)
	}
}

func CombineVertices(obj *common.ObjectPrimitive) []float32 {
	combinedVertices := make([]float32, 0, len(obj.Vertices)+len(obj.UVs)+len(obj.Normals))
	for i := 0; i < len(obj.Vertices)/3; i++ {
		// Add vertex position (3 floats)
		combinedVertices = append(combinedVertices, obj.Vertices[i*3], obj.Vertices[i*3+1], obj.Vertices[i*3+2])

		// Add UV coordinates (2 floats), assuming UVs are available
		if i < len(obj.UVs)/2 {
			combinedVertices = append(combinedVertices, obj.UVs[i*2], obj.UVs[i*2+1])
		} else {
			combinedVertices = append(combinedVertices, 0.0, 0.0) // Default UV
		}

		// Add normals (3 floats), assuming normals are available
		if i < len(obj.Normals)/3 {
			combinedVertices = append(combinedVertices, obj.Normals[i*3], obj.Normals[i*3+1], obj.Normals[i*3+2])
		} else {
			combinedVertices = append(combinedVertices, 0.0, 0.0, 0.0) // Default normal
		}
	}
	return combinedVertices
}
