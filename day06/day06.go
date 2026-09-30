package main

import (
	"log"
	"regexp"
	"strconv"
)

func LightGrid(lines []string, apply func(action string, point *uint16)) int {
	re := regexp.MustCompile(`^(.*) ([0-9]+),([0-9]+) through ([0-9]+),([0-9]+)$`)

	const (
		width, height = 1000, 1000
	)
	// Input data contains 300 instructions. A byte could overflow in Part 2
	var grid [width][height]uint16

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			log.Fatalf("Invalid input format: %q", line)
		}

		action := matches[1]
		left, _ := strconv.Atoi(matches[2])
		bottom, _ := strconv.Atoi(matches[3])
		right, _ := strconv.Atoi(matches[4])
		top, _ := strconv.Atoi(matches[5])

		for x := left; x <= right; x++ {
			for y := bottom; y <= top; y++ {
				apply(action, &grid[x][y])
			}
		}
	}

	result := 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			result += int(grid[x][y])
		}
	}

	return result
}
