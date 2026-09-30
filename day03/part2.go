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
	s_loc := house{}
	r_loc := house{}
	visits := make(map[house]int)
	visits[s_loc]++
	visits[r_loc]++

	for i, m := range line {
		var loc *house
		if 0 == i%2 {
			loc = &s_loc
		} else {
			loc = &r_loc
		}

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
		visits[*loc]++
	}

	result := len(visits)
	fmt.Println("Result:", result)
}
