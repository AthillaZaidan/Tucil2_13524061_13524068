package main

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
		return nil, nil, fmt.Errorf("Cannot open %s: %w", path, err)
	}
	defer file.Close()

	verts := []Vec3{}
	faces := []Face{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "v":
			if len(parts) < 4 {
				continue
			}
			x, _ := strconv.ParseFloat(parts[1], 64)
			y, _ := strconv.ParseFloat(parts[2], 64)
			z, _ := strconv.ParseFloat(parts[3], 64)
			verts = append(verts, Vec3{x, y, z})
		case "f":
			if len(parts) < 4 {
				continue
			}
			i1, _ := strconv.Atoi(parts[1])
			i2, _ := strconv.Atoi(parts[2])
			i3, _ := strconv.Atoi(parts[3])
			faces = append(faces, Face{i1 - 1, i2 - 1, i3 - 1})
		}
	}
	return verts, faces, nil
}
