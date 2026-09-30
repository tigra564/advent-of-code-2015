package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"math"
)

func main() {
	lines := common.ReadInputLines()
	people, hapUnits := fromLines(lines)

	result, _ := common.FindBestHamiltonian(people, hapUnits, true, math.MinInt, func(a, b int) bool {
		return a > b
	})

	fmt.Println("Result:", result)
}
