package main

import (
	"advent-of-code-2015/common"
	"cmp"
	"fmt"
	"math"
	"slices"
	"strconv"
)

func main() {
	lines := common.ReadInputLines()

	const nLiters = 150

	containers := []int{}
	for _, line := range lines {
		cont, _ := strconv.Atoi(line)
		containers = append(containers, cont)
	}
	slices.SortFunc(containers, func(a, b int) int {
		return cmp.Compare(b, a)
	})

	nCont := len(containers)
	minUsed := math.MaxInt
	result := 0
	for mask := range uint32(1) << nCont {
		sum := 0
		contsUsed := 0
		for i := range nCont {
			if 0 != mask&(uint32(1)<<i) {
				sum += containers[i]
				if sum > nLiters {
					break
				}
				contsUsed++
			}
		}
		if sum == nLiters {
			switch {
			case contsUsed < minUsed:
				minUsed = contsUsed
				result = 1
			case contsUsed == minUsed:
				result++
			}
		}
	}

	fmt.Println("Result:", result)
}
