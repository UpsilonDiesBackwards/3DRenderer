package common

type Material struct {
	Name       string
	TextureID  uint32
	DiffuseMap string
}

type ObjectPrimitive struct {
	Vertices []float32
	Indices  []uint32
	Normals  []float32
	UVs      []float32

	Textures map[string]uint32
	Material *Material
}
