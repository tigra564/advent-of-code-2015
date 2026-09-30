package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()
	circ := parseCircuit(lines)

	result := circ["a"].eval(circ)

	fmt.Println("Result:", result)
}
