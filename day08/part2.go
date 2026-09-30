package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"strconv"
)

func main() {
	lines := common.ReadInputLines()

	result := 0

	for _, line := range lines {
		quoted := strconv.Quote(line)
		result += len(quoted) - len(line)
	}

	fmt.Println("Result:", result)
}
