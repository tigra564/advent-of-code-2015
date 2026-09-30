package main

import (
	"advent-of-code-2015/common"
	"regexp"
	"strconv"
)

func fromLines(lines []string) (common.Nodes, common.Weights) {
	re := regexp.MustCompile(`^(.+) to (.+) = ([0-9]+)`)

	dists := make(common.Weights)
	locsMap := make(map[string]common.Empty)

	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		from := m[1]
		to := m[2]
		dist, _ := strconv.Atoi(m[3])

		dists[common.MakeEdge(from, to)] = dist
		locsMap[from] = common.Empty{}
		locsMap[to] = common.Empty{}
	}

	var locs common.Nodes
	for k := range locsMap {
		locs = append(locs, k)
	}

	return locs, dists
}
