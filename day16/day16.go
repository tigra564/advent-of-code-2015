package main

import (
	"advent-of-code-2015/common"
	"regexp"
	"strconv"
)

var (
	reSue = regexp.MustCompile(`Sue ([0-9]+): `)
	reMap = regexp.MustCompile(`([a-z]+): ([0-9]+)`)
)

// From the puzzle statement: the MFCSAM analysis of the real Aunt Sue.
var traitsMFCSAM = map[string]int{
	"children":    3,
	"cats":        7,
	"samoyeds":    2,
	"pomeranians": 3,
	"akitas":      0,
	"vizslas":     0,
	"goldfish":    5,
	"trees":       3,
	"cars":        2,
	"perfumes":    1,
}

func readSue(s string) (int, map[string]int) {
	m := reSue.FindStringSubmatch(s)
	numSue, _ := strconv.Atoi(m[1])
	traitsSue := readMapStrInt(s[len(m[0]):])

	return numSue, traitsSue
}

func readMapStrInt(s string) map[string]int {
	matches := reMap.FindAllStringSubmatch(s, -1)

	mapping := make(map[string]int)
	for _, m := range matches {
		mapping[m[1]], _ = strconv.Atoi(m[2])
	}
	return mapping
}

func findSue(validate func(traitsMFCSAM, traitsSue map[string]int) bool) int {
	lines := common.ReadInputLines()
	for _, line := range lines {
		numSue, traitsSue := readSue(line)
		if validate(traitsMFCSAM, traitsSue) {
			return numSue
		}
	}

	return 0
}
