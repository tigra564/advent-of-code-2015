package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

type deerState int
const (
	FLYING deerState = iota
	RESTING
)

type reindeer struct {
	speed int
	flightDuration int
	restDuration int

	state deerState
	timer int
	distance int
	score int
}

func main() {
	lines := common.ReadInputLines()

	const raceTime = 2503
	contestants := []reindeer{}
	for _, line := range lines {
		var deer reindeer
		nums := common.ExtractNumbers(line)

		deer.speed = nums[0]
		deer.flightDuration = nums[1]
		deer.restDuration = nums[2]
		deer.state = FLYING
		deer.timer = deer.flightDuration

		contestants = append(contestants, deer)
	}

	for range raceTime {
		maxDist := 0
		for i := range contestants {
			c := &contestants[i]
			switch c.state {
			case FLYING:
				c.distance += c.speed // * 1 second
				if c.timer--; c.timer == 0 {
					c.state = RESTING
					c.timer = c.restDuration
				}
			case RESTING:
				if c.timer--; c.timer == 0 {
					c.state = FLYING
					c.timer = c.flightDuration
				}
			}
			maxDist = max(maxDist, c.distance)
		}

		for i := range contestants {
			c := &contestants[i]
			if c.distance == maxDist {
				c.score++
			}
		}
	}

	result := 0
	for _, c := range contestants {
		result = max(result, c.score)
	}
	fmt.Println("Result:", result)
}
