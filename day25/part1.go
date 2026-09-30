package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func codeNum(row, col int) int {
	diag := (row - 1) + (col - 1)
	return diag * (diag+1) / 2 + col
}

func next(in int) int {
	return (in * 252533) % 33554393
}

func main() {
	lines := common.ReadInputLines()

	const start = 20151125

	coords := common.ExtractNumbers(lines[0])
	n := codeNum(coords[0], coords[1])

	result := start
	for range n - 1 {
		result = next(result)
	}

	fmt.Println("Result:", result)
}
