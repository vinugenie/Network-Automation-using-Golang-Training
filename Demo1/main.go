package main

import "fmt"

var grid [3][3]int

func main() {
	fmt.Printf("This is our first program...")
	grid[1][0], grid[1][1], grid[1][2] = 8, 10, 6
	fmt.Printf("%v", grid)
	nodes := [...]string{"XE-1", "XE-2", "NX-3"}
	nodes[len(nodes)-1] = "XR-4"

	// nodes1 := append(nodes, "XR-4")
	fmt.Printf("%v", nodes)

}
