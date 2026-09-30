package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"math"
)

func main() {
	lines := common.ReadInputLines()
	people, hapUnits := fromLines(lines)

	// We simply add "I" to people. Happiness values for edges adjacent to "I" will automatically be set to 0.
	people = append(people, "I")

	result, _ := common.FindBestHamiltonian(people, hapUnits, true, math.MinInt, func(a, b int) bool {
		return a > b
	})

	fmt.Println("Result:", result)
}
