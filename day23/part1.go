package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	regs := registers{}
	ops := compile(lines)
	run(&regs, ops, false)

	result := regs.gp[2]
	fmt.Println("Result:", result)
}
