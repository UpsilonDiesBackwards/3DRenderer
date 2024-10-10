package tools

import (
	"bufio"
	"github.com/UpsilonDiesBackwards/3DRenderer/engine/common"
	"os"
	"strings"
)

func ParseMTL(mtlFPath string) (map[string]*common.Material, error) {
	materials := make(map[string]*common.Material)

	file, err := os.Open(mtlFPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var currentMaterial *common.Material

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(line, "newmtl "):
			materialName := line[7:]
			currentMaterial = &common.Material{
				Name: materialName,
			}
			materials[materialName] = currentMaterial
		case strings.HasPrefix(line, "map_Kd "):
			if currentMaterial != nil {
				currentMaterial.DiffuseMap = line[7:]
			}
		}
	}

	return materials, scanner.Err()
}
