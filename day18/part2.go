package main

import (
	"fmt"
)

func setStuckLights(grid []string) {
	right := len(grid[0]) - 2
	bottom := len(grid) - 2

	for _, i := range []int{1, bottom} {
		b := []byte(grid[i])
		b[1] = '#'
		b[right] = '#'
		grid[i] = string(b)
	}
}

func main() {
	grid := getGrid()
	setStuckLights(grid)

	const nIters = 100
	for range nIters {
		grid = nextFrame(grid)
		setStuckLights(grid)
	}

	result := countLit(grid)
	fmt.Println("Result:", result)
}
