package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"strconv"
)

func main() {
	lines := common.ReadInputLines()

	maxPresents, _ := strconv.Atoi(lines[0])
	// mtp — maximum of tens of presents. Here +9 is in case the input is not divisible by 10
	mtp := (maxPresents + 9) / 10

	// We have a 1.0 probability of meeting an abundant number
	// Hence we can put upperBoundHouse := mtp / 2

	// But we can do better since 9! is highly abundant and quite close to our input
	const (
		fact9 = 362880
		sigma_fact9 = 1481040
	)
	upperBoundHouse := fact9 * (mtp / sigma_fact9 + 1)

	// Here +1 is to accomodate the house #0, which is non-existent, but indices start with 1
	presents := make([]int, upperBoundHouse + 1, upperBoundHouse + 1)

	for i := 1; i <= upperBoundHouse; i++ {
		for j := i; j <= upperBoundHouse; j += i {
			presents[j] += i
			if (presents[j] >= mtp) {
				upperBoundHouse = j
			}
		}
	}

	result := upperBoundHouse
	fmt.Println("Result:", result)
}
