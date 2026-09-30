package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"strings"
)

func getGrid() []string {
	lines := common.ReadInputLines()

	width := len(lines[0])
	border := strings.Repeat(".", width+2)

	newLines := make([]string, 0, len(lines)+2)
	newLines = append(newLines, border)
	for _, line := range lines {
		newLines = append(newLines, "."+line+".")
	}
	newLines = append(newLines, border)

	return newLines
}

var dirs = [][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{ 0, -1},          { 0, 1},
	{ 1, -1}, { 1, 0}, { 1, 1},
}

func nextFrame(in []string) []string {
	width := len(in[0])
	height := len(in)
	border := in[0]
	out := make([]string, 0, height)
	out = append(out, border)

	for i := 1; i < height - 1; i++ {
		var b strings.Builder
		b.Grow(width)
		b.WriteByte('.')
		for j := 1; j < width - 1; j++ {
			neighborsOn := 0
			for _, d := range dirs {
				if in[i + d[0]][j + d[1]] == '#' {
					neighborsOn++
				}
			}
			var cell byte = '.'
			switch in[i][j] {
			case '#':
				if 2 == neighborsOn || 3 == neighborsOn {
					cell = '#'
				}
			case '.':
				if 3 == neighborsOn {
					cell = '#'
				}
			}
			b.WriteByte(cell)
		}
		b.WriteByte('.')
		out = append(out, b.String())
	}
	out = append(out, border)

	return out
}

func countLit(grid []string) int {
	result := 0
	for _, line := range grid {
		result += strings.Count(line, "#")
	}

	return result
}

func printGrid(grid []string) {
	for _, l := range grid {
		fmt.Println(l)
	}
	fmt.Println()
}
