package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	regs := registers{}
	regs.gp[1] = 1
	ops := compile(lines)
	run(&regs, ops, false)

	result := regs.gp[2]
	fmt.Println("Result:", result)
}
