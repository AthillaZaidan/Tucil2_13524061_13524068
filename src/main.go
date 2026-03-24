package main

import (
	"fmt"
	"os"
	"time"

	octree "gemilang/src/packages/octree"
	parser "gemilang/src/packages/parser"
	viewer "gemilang/src/packages/viewer"
)

// hitung jumlah node per depth
func CountNodes(node *octree.Octree, depth int, counts map[int]int) {
	if node == nil {
		return
	}
	counts[depth]++
	for _, child := range node.Children {
		CountNodes(child, depth+1, counts)
	}
}

// kumpulkan semua leaf node (= voxel)
func CollectLeaves(node *octree.Octree, leaves *[]*octree.Octree) {
	if node == nil {
		return
	}
	if node.IsLeaf {
		*leaves = append(*leaves, node)
		return
	}
	for _, child := range node.Children {
		CollectLeaves(child, leaves)
	}
}

func main() {
	var filename string
	fmt.Print("Masukkan nama file OBJ (contoh: pumpkin.obj): ")
	fmt.Scan(&filename)
	path := "obj/" + filename

	var maxDepth int
	fmt.Print("Masukkan kedalaman maksimum: ")
	_, err := fmt.Scan(&maxDepth)
	if err != nil || maxDepth < 1 {
		fmt.Println("Error: max-depth harus angka positif")
		os.Exit(1)
	}

	// 1. parse
	verts, faces, err := parser.ParseOBJ(path)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// 2. bounding box + root
	bbMin, bbMax := octree.BoundingBox(verts)
	root := &octree.Octree{Min: bbMin, Max: bbMax}

	// 3. build octree + ukur waktu
	start := time.Now()
	pruned := map[int]int{}
	octree.Build(root, verts, faces, 0, maxDepth, pruned)
	elapsed := time.Since(start)

	// 4. kumpulkan hasil
	var leaves []*octree.Octree
	CollectLeaves(root, &leaves)
	voxelCount := len(leaves)

	nodeCounts := map[int]int{}
	CountNodes(root, 0, nodeCounts)

	// 5. print statistik
	fmt.Println("=== Hasil Voxelization ===")
	fmt.Printf("Banyak voxel   : %d\n", voxelCount)
	fmt.Printf("Banyak vertex  : %d\n", voxelCount*8)
	fmt.Printf("Banyak face    : %d\n", voxelCount*12)
	fmt.Printf("Kedalaman max  : %d\n", maxDepth)
	fmt.Printf("Waktu eksekusi : %v\n", elapsed)

	fmt.Println("\nStatistik node per depth:")
	for d := 0; d <= maxDepth; d++ {
		fmt.Printf("  %d : %d\n", d, nodeCounts[d])
	}
	fmt.Println("\nStatistik Node yang tidak perlu ditelusuri:")
	for i := 1; i <= maxDepth; i++ {
		fmt.Printf("%d : %d\n", i, pruned[i])
	}
	// 6. tulis output
	outputPath := "output.obj"
	err = writeOBJ(leaves, outputPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("\nOutput disimpan di: %s\n", outputPath)

	var showViewer string
	fmt.Print("\nTampilkan 3D viewer? (y/n): ")
	fmt.Scan(&showViewer)
	if showViewer == "y" || showViewer == "Y" {
		viewer.Launch(leaves, bbMin, bbMax)
	}
}

func writeOBJ(leaves []*octree.Octree, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot create output: %w", err)
	}
	defer f.Close()

	vertOffset := 1
	for _, leaf := range leaves {
		mn := leaf.Min
		mx := leaf.Max

		fmt.Fprintf(f, "v %f %f %f\n", mn.X, mn.Y, mn.Z)
		fmt.Fprintf(f, "v %f %f %f\n", mx.X, mn.Y, mn.Z)
		fmt.Fprintf(f, "v %f %f %f\n", mx.X, mn.Y, mx.Z)
		fmt.Fprintf(f, "v %f %f %f\n", mn.X, mn.Y, mx.Z)
		fmt.Fprintf(f, "v %f %f %f\n", mn.X, mx.Y, mn.Z)
		fmt.Fprintf(f, "v %f %f %f\n", mx.X, mx.Y, mn.Z)
		fmt.Fprintf(f, "v %f %f %f\n", mx.X, mx.Y, mx.Z)
		fmt.Fprintf(f, "v %f %f %f\n", mn.X, mx.Y, mx.Z)

		o := vertOffset
		fmt.Fprintf(f, "f %d %d %d\n", o, o+1, o+2)
		fmt.Fprintf(f, "f %d %d %d\n", o, o+2, o+3)
		fmt.Fprintf(f, "f %d %d %d\n", o+4, o+6, o+5)
		fmt.Fprintf(f, "f %d %d %d\n", o+4, o+7, o+6)
		fmt.Fprintf(f, "f %d %d %d\n", o, o+4, o+5)
		fmt.Fprintf(f, "f %d %d %d\n", o, o+5, o+1)
		fmt.Fprintf(f, "f %d %d %d\n", o+2, o+6, o+7)
		fmt.Fprintf(f, "f %d %d %d\n", o+2, o+7, o+3)
		fmt.Fprintf(f, "f %d %d %d\n", o, o+3, o+7)
		fmt.Fprintf(f, "f %d %d %d\n", o, o+7, o+4)
		fmt.Fprintf(f, "f %d %d %d\n", o+1, o+5, o+6)
		fmt.Fprintf(f, "f %d %d %d\n", o+1, o+6, o+2)

		vertOffset += 8
	}
	return nil
}
