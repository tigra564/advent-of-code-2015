package main

import (
	"advent-of-code-2015/common"
)

func countNumbers(str string) int64 {
	result := int64(0)
	for _, n := range common.ExtractNumbers(str) {
		result += int64(n)
	}

	return result
}
