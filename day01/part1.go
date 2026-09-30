package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	result := 0
	for _, c := range lines[0] {
		switch c {
		case ')':
			result--
		case '(':
			result++
		}
	}

	fmt.Println("Result:", result)
}
