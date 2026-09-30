package main

import (
	"advent-of-code-2015/common"
	"fmt"
)

func next(passwd []byte) {
	// We won't get more than 8 bytes anyway, so no overflow at the first byte
	for i := len(passwd) - 1; i >= 0; i-- {
		c := &passwd[i]
		if *c++; *c <= 'z' {
			break
		}
		*c = 'a'
	}
}

func has3Straight(passwd []byte) bool {
	for i := 0; i < len(passwd)-2; i++ {
		if passwd[i]+1 == passwd[i+1] && passwd[i+1]+1 == passwd[i+2] {
			return true
		}
	}
	return false
}

func hasNoIllegals(passwd []byte) bool {
	for i := 0; i < len(passwd); i++ {
		if c := passwd[i]; c == 'i' || c == 'o' || c == 'l' {
			return false
		}
	}
	return true
}

func hasDifferentPairs(passwd []byte) bool {
	i := 0
	for ; i < len(passwd)-1; i++ {
		if passwd[i+1] == passwd[i] {
			break
		}
	}
	for j := i + 1; j < len(passwd)-1; j++ {
		if passwd[j+1] == passwd[j] && passwd[j] != passwd[i] {
			return true
		}
	}
	return false
}

func isValid(passwd []byte) bool {
	res1 := has3Straight(passwd)
	res2 := hasNoIllegals(passwd)
	res3 := hasDifferentPairs(passwd)

	return res1 && res2 && res3
}

func main() {
	lines := common.ReadInputLines()

	passwd := []byte(lines[0])
	next(passwd)
	for !isValid(passwd) {
		next(passwd)
	}

	result := string(passwd)

	fmt.Println("Result:", result)
}
