package viewer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	octree "gemilang/src/packages/octree"
)

// SaveVoxelOBJ exports each leaf voxel as a box mesh in Wavefront OBJ format.
func SaveVoxelOBJ(path string, leaves []*octree.Octree) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("cannot create output directory: %w", err)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	if _, err := fmt.Fprintln(w, "# OctoVox voxelized output"); err != nil {
		return err
	}

	vertexOffset := 1
	for i, leaf := range leaves {
		x0 := leaf.Min.X
		y0 := leaf.Min.Y
		z0 := leaf.Min.Z
		x1 := leaf.Max.X
		y1 := leaf.Max.Y
		z1 := leaf.Max.Z

		if _, err := fmt.Fprintf(w, "o voxel_%d\n", i+1); err != nil {
			return err
		}

		verts := [8][3]float64{
			{x0, y0, z0},
			{x1, y0, z0},
			{x1, y1, z0},
			{x0, y1, z0},
			{x0, y0, z1},
			{x1, y0, z1},
			{x1, y1, z1},
			{x0, y1, z1},
		}

		for _, v := range verts {
			if _, err := fmt.Fprintf(w, "v %.6f %.6f %.6f\n", v[0], v[1], v[2]); err != nil {
				return err
			}
		}

		faces := [12][3]int{
			{1, 2, 3}, {1, 3, 4},
			{5, 6, 7}, {5, 7, 8},
			{1, 5, 8}, {1, 8, 4},
			{2, 6, 7}, {2, 7, 3},
			{4, 3, 7}, {4, 7, 8},
			{1, 2, 6}, {1, 6, 5},
		}

		for _, fIdx := range faces {
			a := vertexOffset + fIdx[0] - 1
			b := vertexOffset + fIdx[1] - 1
			c := vertexOffset + fIdx[2] - 1
			if _, err := fmt.Fprintf(w, "f %d %d %d\n", a, b, c); err != nil {
				return err
			}
		}

		vertexOffset += 8
	}

	return nil
}
