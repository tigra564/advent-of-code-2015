package main

import (
	"advent-of-code-2015/common"
	"cmp"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type production struct {
	from, to string

	numNonTerms, numAtoms int
}

var reAtom = regexp.MustCompile(`[A-Z][a-z]?`)

func fillProductions(prods []production) {
	nonTerms := make(map[string]common.Empty)

	for _, prod := range prods {
		nonTerms[prod.from] = common.Empty{}
	}

	for i := range prods {
		prod := &prods[i]

		matches := reAtom.FindAllString(prod.to, -1)
		prod.numAtoms = len(matches)

		prod.numNonTerms = 0
		for _, m := range matches {
			if _, ok := nonTerms[m]; ok {
				prod.numNonTerms++
			}
		}
	}
}

func sortProductions(prods []production) {
	fillProductions(prods)

	slices.SortFunc(prods, func(a, b production) int {
		if c := cmp.Compare(b.numNonTerms, a.numNonTerms); c != 0 {
			return c
		}
		if c := cmp.Compare(b.numAtoms, a.numAtoms); c != 0 {
			return c
		}
		return cmp.Compare(b.to, a.to)
	})
}

func main() {
	lines := common.ReadInputLines()

	prods := make([]production, 0, len(lines)-2)
	for _, line := range lines[:len(lines)-2] {
		s := strings.Split(line, " => ")
		prods = append(prods, production{s[0], s[1], 0, 0})
	}
	sortProductions(prods)

	molecule := lines[len(lines)-1]

	changed := true
	result := 0
	for changed {
		changed = false
		for _, prod := range prods {
			numSubst := strings.Count(molecule, prod.to)
			if 0 == numSubst {
				continue
			}
			molecule = strings.ReplaceAll(molecule, prod.to, prod.from)
			result += numSubst
			changed = true
			break
		}
	}

	if "e" != molecule {
		panic(`Algorithm ended with "` + molecule + `", not with just "e"`)
	}

	fmt.Println("Result:", result)
}
