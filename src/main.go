package main

import (
	"fmt"
)

func main() {
	verts, faces, err := ParseOBJ("../test/cow.obj")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	bbMin, bbMax := BoundingBox(verts)
	fmt.Println("min:", bbMin)
	fmt.Println("max:", bbMax)
	fmt.Println(len(verts), "verts,", len(faces), "faces")
}
