package main

import (
	"fmt"
)

func main() {
	grid := getGrid()

	const nIters = 100
	for range nIters {
		grid = nextFrame(grid)
	}

	result := countLit(grid)
	fmt.Println("Result:", result)
}
