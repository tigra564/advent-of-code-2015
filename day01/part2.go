package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	level := 0
	result := 0
	for p, c := range lines[0] {
		switch c {
		case ')':
			level--
		case '(':
			level++
		}
		if level < 0 {
			result = p + 1
			break
		}
	}

	fmt.Println("Result:", result)
}
