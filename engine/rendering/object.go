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

	Material          map[string]*common.Material
	AlbedoTextures    []uint32
	NormalTextures    []uint32
	SpecularTextures  []uint32
	RoughnessTextures []uint32

	ModelMatrix mgl32.Mat4
}

func NewRenderableObject(obj *common.ObjectPrimitive, mtlPath string) *RenderableObject {
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

	var albedoTextures, normalTextures, specularTextures, roughnessTextures []uint32
	materials, err := tools.ParseMTL(mtlPath)
	if err != nil {

		fmt.Println("Failed to parse mtl file: ", err)
	} else {
		fmt.Println("Parsed materials from file: ", materials)

		for name, material := range materials {
			if material.DiffuseMap != "" {
				tex, err := tools.LoadTexture(material.DiffuseMap)
				if err != nil {
					fmt.Println("Failed to load texture: ", name, " : ", err)
					tex = tools.CreateWhiteTexture()
				}
				albedoTextures = append(albedoTextures, tex)
			} else {
				tex := tools.CreateWhiteTexture()
				albedoTextures = append(albedoTextures, tex)
				fmt.Println("(A) No texture path for material: ", material)
			}

			if material.NormalMap != "" {
				tex, err := tools.LoadTexture(material.NormalMap)
				if err != nil {
					fmt.Println("Failed to load texture: ", name, " : ", err)
					tex = tools.CreateWhiteTexture()
				}
				normalTextures = append(normalTextures, tex)
			} else {
				tex := tools.CreateWhiteTexture()
				normalTextures = append(normalTextures, tex)
				fmt.Println("(N) No texture path for material: ", material)
			}

			if material.SpecularMap != "" {
				tex, err := tools.LoadTexture(material.SpecularMap)
				if err != nil {
					fmt.Println("Failed to load texture: ", name, " : ", err)
					tex = tools.CreateWhiteTexture()
				}
				specularTextures = append(specularTextures, tex)
			} else {
				tex := tools.CreateWhiteTexture()
				specularTextures = append(specularTextures, tex)
				fmt.Println("(S) No texture path for material: ", material)
			}

			if material.RoughnessMap != "" {
				tex, err := tools.LoadTexture(material.RoughnessMap)
				if err != nil {
					fmt.Println("Failed to load texture: ", name, " : ", err)
					tex = tools.CreateWhiteTexture()
				}
				roughnessTextures = append(roughnessTextures, tex)
			} else {
				tex := tools.CreateWhiteTexture()
				roughnessTextures = append(roughnessTextures, tex)
				fmt.Println("(R) No texture path for material: ", material)
			}

		}
	}

	return &RenderableObject{
		VAO:               vao,
		VBO:               vbo,
		EBO:               ebo,
		Vertices:          obj.Vertices,
		Normals:           obj.Normals,
		TexCoords:         obj.UVs,
		Indices:           obj.Indices,
		ModelMatrix:       mgl32.Ident4(),
		Material:          materials,
		AlbedoTextures:    albedoTextures,
		NormalTextures:    normalTextures,
		SpecularTextures:  specularTextures,
		RoughnessTextures: roughnessTextures,
	}
}

func (obj *RenderableObject) Draw(shader *Shader) {
	gl.BindVertexArray(obj.VAO)

	shader.SetMat4ByName("model", obj.ModelMatrix)

	for i, texture := range obj.AlbedoTextures {
		gl.ActiveTexture(uint32(gl.TEXTURE0 + i))
		gl.BindTexture(gl.TEXTURE_2D, texture)
		shader.SetInt(fmt.Sprintf("texture%d", i), int(int32(i)))
	}

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
		tex := tools.CreateColorMaterial(R, G, B, A)
		obj.AlbedoTextures = append(obj.AlbedoTextures, tex)
	}
}

func CombineVertices(obj *common.ObjectPrimitive) []float32 {
	combinedVertices := make([]float32, 0, len(obj.Vertices)+len(obj.UVs)+len(obj.Normals))
	for i := 0; i < len(obj.Vertices)/3; i++ {
		combinedVertices = append(combinedVertices, obj.Vertices[i*3], obj.Vertices[i*3+1], obj.Vertices[i*3+2])

		if i < len(obj.UVs)/2 {
			combinedVertices = append(combinedVertices, obj.UVs[i*2], obj.UVs[i*2+1])
		} else {
			combinedVertices = append(combinedVertices, 0.0, 0.0) // Default UV
		}

		if i < len(obj.Normals)/3 {
			combinedVertices = append(combinedVertices, obj.Normals[i*3], obj.Normals[i*3+1], obj.Normals[i*3+2])
		} else {
			combinedVertices = append(combinedVertices, 0.0, 0.0, 0.0) // Default normal
		}
	}
	return combinedVertices
}
