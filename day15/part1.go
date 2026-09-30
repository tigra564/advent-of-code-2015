package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	var ingredients [][]int
	for _, line := range lines {
		ingredients = append(ingredients, common.ExtractNumbers(line))
	}

	const nSpoons = 100
	result := 0

	getMaxScore := func(parts []int) {
		caloriesPos := len(ingredients[0]) - 1
		sums := make([]int, caloriesPos)
		for i, coeff := range parts {
			for j, val := range ingredients[i][:caloriesPos] {
				sums[j] += coeff * val
			}
		}

		product := 1
		for _, sum := range sums {
			if sum < 0 {
				product = 0
				break
			}
			product *= sum
		}

		result = max(result, product)
	}

	iteratePartitions(nSpoons, len(ingredients), []int{}, getMaxScore)

	fmt.Println("Result:", result)
}
