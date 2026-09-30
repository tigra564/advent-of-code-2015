package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"sort"
)

func main() {
	lines := common.ReadInputLines()

	result := 0
	for _, l := range lines {
		dims := common.ExtractNumbers(l)
		sort.Ints(dims)
		result += 2*(dims[0]+dims[1]) +
			dims[0]*dims[1]*dims[2]
	}

	fmt.Println("Result:", result)
}
