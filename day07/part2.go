package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()
	circ := parseCircuit(lines)

	sig := circ["a"].eval(circ)

	for _, v := range circ {
		v.isSet = false
	}
	circ["b"].conn = sig

	result := circ["a"].eval(circ)

	fmt.Println("Result:", result)
}
