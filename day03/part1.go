package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	type house struct {
		x, y int
	}

	line := lines[0]
	loc := house{}
	visits := make(map[house]int)
	visits[loc] = 1

	for _, m := range line {
		switch m {
		case '^':
			loc.y++
		case 'v':
			loc.y--
		case '>':
			loc.x++
		case '<':
			loc.x--
		}
		visits[loc]++
	}

	result := len(visits)
	fmt.Println("Result:", result)
}
