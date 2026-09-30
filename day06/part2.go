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
			*point += 2
		case "turn off":
			if *point > 0 {
				*point--
			}
		case "turn on":
			*point++
		}
	})

	fmt.Println("Result:", result)
}
