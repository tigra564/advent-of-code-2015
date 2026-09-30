package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func hasDoublePair(s string) bool {
	// We put into "seen" the positions of the first occurence of a pair of letters
	seen := make(map[[2]byte]int)

	for c := 0; c < len(s)-1; c++ {
		seq := [2]byte{s[c], s[c+1]}
		if pos, ok := seen[seq]; ok {
			// We've already met such pair earlier, let's check for overlapping
			if pos < c-1 {
				return true
			}
		} else {
			seen[seq] = c
		}
	}

	return false
}

func hasEmbraced(s string) bool {
	for c := 0; c < len(s)-2; c++ {
		if s[c] == s[c+2] {
			return true
		}
	}
	return false
}

func main() {
	lines := common.ReadInputLines()

	result := 0

	for _, line := range lines {
		if hasDoublePair(line) && hasEmbraced(line) {
			result++
		}
	}

	fmt.Println("Result:", result)
}
