package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func main() {
	lines := common.ReadInputLines()

	// Vixen can fly 19 km/s for 7 seconds, but then must rest for 124 seconds.
	// Poor Vixen, nobody cares about its name. And the winner remains unsung as well.

	const raceTime = 2503
	result := 0
	for _, line := range lines {
		nums := common.ExtractNumbers(line)
		speed, flightDuration, restDuration := nums[0], nums[1], nums[2]

		dist := (raceTime / (flightDuration + restDuration)) * flightDuration * speed
		dist += min(raceTime % (flightDuration + restDuration), flightDuration) * speed
		result = max(result, dist)
	}

	fmt.Println("Result:", result)
}
