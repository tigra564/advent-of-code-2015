package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	result := countNumbers(lines[0])

	fmt.Println("Result:", result)
}
