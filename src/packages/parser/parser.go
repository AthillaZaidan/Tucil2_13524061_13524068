package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Vec3 struct {
	X, Y, Z float64
}

type Face struct {
	V1, V2, V3 int
}

func ParseOBJ(path string) ([]Vec3, []Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer file.Close()

	verts := []Vec3{}
	faces := []Face{}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)

		if len(parts) == 0 || parts[0] == "#" {
			continue
		}

		switch parts[0] {
		case "v":
			if len(parts) < 4 {
				return nil, nil, fmt.Errorf("line %d: 'v' needs 3 coords", lineNum)
			}
			x, err1 := strconv.ParseFloat(parts[1], 64)
			y, err2 := strconv.ParseFloat(parts[2], 64)
			z, err3 := strconv.ParseFloat(parts[3], 64)
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, nil, fmt.Errorf("line %d: invalid vertex coords", lineNum)
			}
			verts = append(verts, Vec3{x, y, z})

		case "f":
			if len(parts) < 4 {
				return nil, nil, fmt.Errorf("line %d: 'f' needs 3 indices", lineNum)
			}
			i1, err1 := strconv.Atoi(parts[1])
			i2, err2 := strconv.Atoi(parts[2])
			i3, err3 := strconv.Atoi(parts[3])
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, nil, fmt.Errorf("line %d: invalid face indices", lineNum)
			}
			if i1 < 1 || i2 < 1 || i3 < 1 {
				return nil, nil, fmt.Errorf("line %d: face index must be >= 1", lineNum)
			}
			faces = append(faces, Face{i1 - 1, i2 - 1, i3 - 1})

		default:
			continue
		}
	}

	// cek index face tidak melebihi jumlah vertex
	for i, f := range faces {
		if f.V1 >= len(verts) || f.V2 >= len(verts) || f.V3 >= len(verts) {
			return nil, nil, fmt.Errorf("face %d references vertex out of range", i+1)
		}
	}

	if len(verts) == 0 {
		return nil, nil, fmt.Errorf("no vertices found in file")
	}
	if len(faces) == 0 {
		return nil, nil, fmt.Errorf("no faces found in file")
	}

	return verts, faces, nil
}
