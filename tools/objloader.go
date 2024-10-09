package tools

import (
	"bufio"
	"github.com/UpsilonDiesBackwards/3DRenderer/engine/common"
	"os"
	"strconv"
	"strings"
)

type Vertex struct {
	Position [3]float32
	UV       [2]float32
	Normal   [3]float32
}

func CreateNewOBJ(modelFPath, mtlFPath string) *common.ObjectPrimitive {
	objPrimitive := &common.ObjectPrimitive{}
	objPrimitive.Vertices, objPrimitive.Normals, objPrimitive.UVs, objPrimitive.Indices = loadOBJFromFile(modelFPath)

	return objPrimitive
}

func loadOBJFromFile(filePath string) (vertices, normals, textureCoords []float32, indices []uint32) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var positions [][]float32
	var uvs [][]float32
	var normalsList [][]float32
	var vertexMap = make(map[Vertex]uint32)
	var uniqueVertices []float32
	var uniqueNormals []float32
	var uniqueUVs []float32
	var newIndices []uint32
	var nextIndex uint32 = 0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(line, "v "):
			positions = append(positions, parseVec3(line))
		case strings.HasPrefix(line, "vn "):
			normalsList = append(normalsList, parseVec3(line))
		case strings.HasPrefix(line, "vt "):
			uvs = append(uvs, parseVec2(line))
		case strings.HasPrefix(line, "f "):
			faceVertices := parseFace(line)

			for _, fv := range faceVertices {
				vertex := Vertex{
					Position: [3]float32{
						positions[fv.Position-1][0],
						positions[fv.Position-1][1],
						positions[fv.Position-1][2],
					},
					UV: [2]float32{
						uvs[fv.UV-1][0],
						uvs[fv.UV-1][1],
					},
					Normal: [3]float32{
						normalsList[fv.Normal-1][0],
						normalsList[fv.Normal-1][1],
						normalsList[fv.Normal-1][2],
					},
				}

				if index, found := vertexMap[vertex]; found {
					newIndices = append(newIndices, index)
				} else {
					uniqueVertices = append(uniqueVertices, vertex.Position[:]...)
					uniqueUVs = append(uniqueUVs, vertex.UV[:]...)
					uniqueNormals = append(uniqueNormals, vertex.Normal[:]...)
					vertexMap[vertex] = nextIndex
					newIndices = append(newIndices, nextIndex)
					nextIndex++
				}
			}
		}
	}

	return uniqueVertices, uniqueNormals, uniqueUVs, newIndices
}

func parseVec3(line string) []float32 {
	parts := strings.Fields(line[2:])
	x, _ := strconv.ParseFloat(parts[0], 32)
	y, _ := strconv.ParseFloat(parts[1], 32)
	z, _ := strconv.ParseFloat(parts[2], 32)
	return []float32{float32(x), float32(y), float32(z)}
}

func parseVec2(line string) []float32 {
	parts := strings.Fields(line[3:])
	u, _ := strconv.ParseFloat(parts[0], 32)
	v, _ := strconv.ParseFloat(parts[1], 32)
	return []float32{float32(u), float32(v)}
}

type FaceVertex struct {
	Position int
	UV       int
	Normal   int
}

func parseFace(line string) []FaceVertex {
	parts := strings.Fields(line[2:])
	vertices := make([]FaceVertex, len(parts))

	for i, part := range parts {
		indices := strings.Split(part, "/")
		posIdx, _ := strconv.Atoi(indices[0])
		uvIdx, _ := strconv.Atoi(indices[1])
		normalIdx, _ := strconv.Atoi(indices[2])

		vertices[i] = FaceVertex{
			Position: posIdx,
			UV:       uvIdx,
			Normal:   normalIdx,
		}
	}

	return vertices
}
