package main

import (
	"strconv"
	"strings"
)

func lookAndSay(s string) string {
	var b strings.Builder
	b.Grow(len(s) * 2)

	for i := 0; i < len(s); {
		c := s[i]
		count := 1
		i++

		for i < len(s) && s[i] == c {
			count++
			i++
		}

		b.WriteString(strconv.Itoa(count))
		b.WriteByte(c)
	}

	return b.String()
}

func lookAndSayCount(s string, num int) int {
	seq := s
	for i := 0; i < num; i++ {
		seq = lookAndSay(seq)
	}

	return len(seq)
}
