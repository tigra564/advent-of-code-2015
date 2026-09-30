package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"strconv"
)

func main() {
	lines := common.ReadInputLines()

	maxPresents, _ := strconv.Atoi(lines[0])
	// mep — maximum "elevens" of presents
	mep := (maxPresents + 10) / 11
	presents := make([]int, mep + 1, mep + 1)

	upperBoundHouse := mep
	for i := 1; i <= upperBoundHouse; i++ {
		k := 0
		for j := i; j <= upperBoundHouse; j += i {
			if k++; k > 50 {
				break
			}
			// We don't multiply by eleven here, we have already divided the target properly
			presents[j] += i
			if (presents[j] >= mep) {
				upperBoundHouse = j
			}
		}
	}

	result := upperBoundHouse
	fmt.Println("Result:", result)
}
