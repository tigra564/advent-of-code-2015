package main

import (
	"advent-of-code-2015/common"
	"math"
	"strconv"
)

type scoreT struct {
	num int
	qe int

	packsUsed uint64
}

func cmpSocres(a, b scoreT) int {
	diffNum := a.num - b.num
	if 0 != diffNum {
		return diffNum
	}
	return a.qe - b.qe
}

func solve(numCompartments int) int {
	lines := common.ReadInputLines()

	packages := make([]int, 0, len(lines))
	sum := 0
	for _, line := range lines {
		pack, _ := strconv.Atoi(line)
		packages = append(packages, pack)
		sum += pack
	}

	totalWeight := sum / numCompartments

	bestScores := make([]scoreT, totalWeight + 1, totalWeight + 1)
	bestScores[0] = scoreT{0, 1, 0}

	for weight := 1; weight <= totalWeight; weight++ {
		minScore := scoreT{math.MaxInt, math.MaxInt, 0}
		for n, p := range packages {
			if p > weight {
				continue
			}
			newScore := bestScores[weight - p]
			if 0 != newScore.packsUsed & (1 << n) {
				continue
			}
			newScore.num++
			newScore.qe *= p
			newScore.packsUsed |= 1 << n

			if 0 > cmpSocres(newScore, minScore) {
				minScore = newScore
			}
		}
		bestScores[weight] = minScore
	}

	return bestScores[totalWeight].qe
}
