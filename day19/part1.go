package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"strings"
)

type production struct {
	from, to string
}

func main() {
	lines := common.ReadInputLines()

	prods := make([]production, 0, len(lines)-2)
	for _, line := range lines[:len(lines)-2] {
		s := strings.Split(line, " => ")
		prods = append(prods, production{s[0], s[1]})
	}
	
	molecule := lines[len(lines)-1]

	newMolecules := make(map[string]common.Empty)
	for _, prod := range prods {
		offset := 0
		for {
			idx := strings.Index(molecule[offset:], prod.from)
			if idx == -1 {
				break
			}
			offset += idx
			nm := molecule[:offset] + prod.to + molecule[offset + len(prod.from):]
			offset++
			newMolecules[nm] = common.Empty{}
		}
	}

	result := len(newMolecules)
	fmt.Println("Result:", result)
}
