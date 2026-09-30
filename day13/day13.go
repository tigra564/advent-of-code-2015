package main

import (
	"advent-of-code-2015/common"
	"regexp"
	"strconv"
)

func fromLines(lines []string) (common.Nodes, common.Weights) {
	// Carol would lose 62 happiness units by sitting next to Alice.
	// Carol would gain 60 happiness units by sitting next to Bob.
	re := regexp.MustCompile(`^(.+) would ([a-z]+) ([0-9]+) .* ([A-Za-z]+)\.$`)

	hapUnits := make(common.Weights)
	peopleMap := make(map[string]common.Empty)

	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		from := m[1]
		to := m[4]
		weight, _ := strconv.Atoi(m[3])
		if m[2] == "lose" {
			weight = -weight
		}

		hapUnits[common.MakeEdge(from, to)] += weight
		peopleMap[from] = common.Empty{}
		peopleMap[to] = common.Empty{}
	}

	var people common.Nodes
	for k := range peopleMap {
		people = append(people, k)
	}

	return people, hapUnits
}
