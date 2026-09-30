package common

import (
	"bufio"
	"iter"
	"fmt"
	"log"
	"math/bits"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func ReadLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", filename, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filename, err)
	}

	return lines, nil
}

func ReadLinesFromStdin() ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func ReadInputLines() []string {
	var (
		lines []string
		err   error
	)

	if len(os.Args) > 1 {
		lines, err = ReadLinesFromFile(os.Args[1])
	} else {
		lines, err = ReadLinesFromStdin()
	}

	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	return lines
}

var reNum = regexp.MustCompile(`-?[0-9]+`)

func ExtractNumbers(s string) []int {
	var nums []int
	for _, m := range reNum.FindAllString(s, -1) {
		n, _ := strconv.Atoi(m)
		nums = append(nums, n)
	}
	return nums
}

func permuteIndices16Impl(n int, used uint16, accum *[16]int, yield func([]int) bool) {
	count := bits.OnesCount16(used)
	if n == count {
		yield(accum[:n])
		return
	}
	for i := range n {
		if 0 == used&(1<<i) {
			accum[count] = i
			permuteIndices16Impl(n, used|(1<<i), accum, yield)
		}
	}
}

func permuteIndices16(n int) iter.Seq[[]int] {
	var used uint16
	var accum [16]int

	return func(yield func([]int) bool) {
		permuteIndices16Impl(min(n, 16), used, &accum, yield)
	}
}

// In production use something more rigorous, for example: "github.com/etnz/permute"
func Permute16[S ~[]E, E any](slice S) iter.Seq[S] {
	return func(yield func(S) bool) {
		val := make(S, len(slice), len(slice))
		for p := range permuteIndices16(len(slice)) {
			for n, i := range p {
				val[n] = slice[i]
			}
			if !yield(val) {
				return
			}
		}
	}
}

type Empty struct{}

type Edge struct {
	From, To string
}

func MakeEdge(a, b string) Edge {
	if a > b {
		a, b = b, a
	}
	return Edge{a, b}
}

type Weights map[Edge]int
type Nodes []string

func FindBestHamiltonian(nodes Nodes, weights Weights, cycle bool, initBest int, cmp func(int, int) bool) (int, []string) {
	var bestPath []string
	bestTotalWeight := initBest

	for p := range Permute16(nodes) {
		if cycle {
			p = append(slices.Clone(p), p[0])
		}

		totalWeight := 0
		for i := 0; i < len(p)-1; i++ {
			totalWeight += weights[MakeEdge(p[i], p[i+1])]
		}
		if cmp(totalWeight, bestTotalWeight) {
			bestPath = slices.Clone(p)
			bestTotalWeight = totalWeight
		}
	}

	return bestTotalWeight, bestPath
}
