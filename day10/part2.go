package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()
	result := lookAndSayCount(lines[0], 50)

	fmt.Println("Result:", result)
}
