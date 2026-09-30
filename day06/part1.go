package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	result := LightGrid(lines, func(action string, point *uint16) {
		switch action {
		case "toggle":
			*point ^= 1
		case "turn off":
			*point = 0
		case "turn on":
			*point = 1
		}
	})

	fmt.Println("Result:", result)
}
