package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var packages []int
var W int
var n int

// sum returns the sum of weights for the given mask.
func sum(mask uint64) int {
	s := 0
	for i := 0; i < n; i++ {
		if mask&(1<<i) != 0 {
			s += packages[i]
		}
	}
	return s
}

// numBits counts the set bits in mask.
func numBits(mask uint64) int {
	cnt := 0
	for mask != 0 {
		mask &= mask - 1
		cnt++
	}
	return cnt
}

// qe returns the product of weights for the given mask.
func qe(mask uint64) int {
	p := 1
	for i := 0; i < n; i++ {
		if mask&(1<<i) != 0 {
			p *= packages[i]
		}
	}
	return p
}

// canSplit reports whether mask can be partitioned into k groups of weight W.
func canSplit(mask uint64, k int) bool {
	if k == 1 {
		return sum(mask) == W
	}
	// Enumerate only subsets containing the lowest set bit of mask
	// to avoid considering the same partition in different orders.
	firstBit := mask & -mask
	sub := mask
	for sub != 0 {
		if sub&firstBit != 0 && sum(sub) == W {
			if canSplit(mask&^sub, k-1) {
				return true
			}
		}
		sub = (sub - 1) & mask
	}
	return false
}

func main() {
	data, _ := os.ReadFile("input.txt")
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		v, _ := strconv.Atoi(strings.TrimSpace(line))
		packages = append(packages, v)
	}
	n = len(packages)

	total := 0
	for _, p := range packages {
		total += p
	}

	for _, k := range []int{3, 4} {
		if total%k != 0 {
			continue
		}
		W = total / k
		fmt.Printf("=== k = %d, W = %d ===\n", k, W)

		full := uint64(1)<<n - 1

		// Enumerate all subsets with sum W.
		countT := 0
		countValid := 0
		bestNum := n + 1
		bestQe := 0
		var bestT uint64

		for mask := uint64(1); mask < full; mask++ {
			if sum(mask) != W {
				continue
			}
			countT++
			rest := full &^ mask
			if canSplit(rest, k-1) {
				countValid++
				num := numBits(mask)
				q := qe(mask)
				if num < bestNum || (num == bestNum && q < bestQe) {
					bestNum = num
					bestQe = q
					bestT = mask
				}
			}
			if countT%1000000 == 0 {
				fmt.Printf("  ... enumerated %d T\n", countT)
			}
		}

		fmt.Printf("Total T with sum W: %d\n", countT)
		fmt.Printf("Of them leaving a splittable remainder: %d\n", countValid)
		fmt.Printf("Best T: num = %d, qe = %d\n", bestNum, bestQe)
		fmt.Printf("Best T: ")
		for i := 0; i < n; i++ {
			if bestT&(1<<i) != 0 {
				fmt.Printf("%d ", packages[i])
			}
		}
		fmt.Println()
	}
}
