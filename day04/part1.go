package main

import (
	"advent-of-code-2015/common"
	"crypto/md5"
	"fmt"
	"strconv"
)

func main() {
	lines := common.ReadInputLines()

	result := 1
	prefix := lines[0]

	for {
		sample := []byte(prefix + strconv.Itoa(result))
		digest := md5.Sum(sample)
		if 0 == digest[0] && 0 == digest[1] && 0 == digest[2]&0b11110000 {
			break
		}
		result++
	}

	fmt.Println("Result:", result)
}
