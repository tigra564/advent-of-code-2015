package main

import (
	"advent-of-code-2015/common"
	"fmt"
	"regexp"
)

func main() {
	lines := common.ReadInputLines()

	re3Vowels := regexp.MustCompile(`[aeiou].*[aeiou].*[aeiou]`)
	// This could do the job, but the "regex" package does not support backreferencing
	//reHasCouple	:= regexp.MustCompile(`([a-z])\1`)
	reBanned := regexp.MustCompile(`(ab|cd|pq|xy)`)

	hasCouple := func(s string) bool {
		for i := 0; i < len(s)-1; i++ {
			if s[i] == s[i+1] {
				return true
			}
		}
		return false
	}

	result := 0

	for _, line := range lines {
		if re3Vowels.MatchString(line) && hasCouple(line) && !reBanned.MatchString(line) {
			result++
		}
	}

	fmt.Println("Result:", result)
}
