package main

import (
	"fmt"
)

func main() {
	v, f, err := ParseOBJ("../test/line.obj")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(v), "verts,", len(f), "faces")
}
