package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"math"
)

func main() {
	lines := common.ReadInputLines()
	locs, dists := fromLines(lines)

	result, _ := common.FindBestHamiltonian(locs, dists, false, math.MaxInt, func(a, b int) bool {
		return a < b
	})

	fmt.Println("Result:", result)
}
