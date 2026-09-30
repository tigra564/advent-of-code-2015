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
		unquoted, _ := strconv.Unquote(line)
		result += len(line) - len(unquoted)
	}

	fmt.Println("Result:", result)
}
