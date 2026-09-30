# Day 9: All in a Single Night

## The Puzzle

Santa needs to visit a list of locations, each exactly once. The order does not matter beyond the total distance travelled. The input gives pairwise distances between every pair of locations: `London to Dublin = 464` and so on. The distances are symmetric: the distance from A to B is the same as from B to A.

Part One asks for the shortest route that visits every location exactly once. Part Two asks for the longest such route. This is the Travelling Salesman Problem, in its path form rather than the more common cycle form.

## Input

Several lines of the form:

```
London to Dublin = 464
London to Belfast = 518
Dublin to Belfast = 141
```

Each line gives a distance between two named locations. The names are arbitrary strings of letters, and the distance is a positive integer. The input is symmetric—the same distance is not repeated in the opposite direction, but we can assume it holds.

The graph is complete and undirected: every location is connected to every other, and the distance does not depend on the direction of travel. Both properties matter. If the graph were not complete, some permutations would correspond to routes that do not exist, and the brute force would need a validity check. If it were directed, the distance from A to B would differ from B to A, and we would need to store two weights per pair instead of one.

There are eight locations in the real input.

## The Approach

The problem is NP-hard in general, and we solve it by brute force. That is only reasonable because the instance is tiny: eight locations, so 8! = 40 320 routes. For each route, we sum the distances between consecutive stops. The number of arithmetic operations is 40 320 × 7 ≈ 280 000, which is nothing.

The practical limit is around twelve to sixteen locations. Beyond that, the factorial grows faster than any polynomial, and brute force stops being viable. My iterator, `Permute16`, is capped at sixteen elements for exactly this reason: there is no point in supporting more, because nobody would use them.

The iterator is hand-written and lives in `common.go`. It is not a general-purpose library—it is limited to sixteen elements, because the "used" bitmask is a `uint16`. For this puzzle and for Day 13 that is more than enough. A comment in the code points to `github.com/etnz/permute` for production use.

The iterator reuses a single buffer across all permutations. This is a deliberate trade: it makes the generator faster—one allocation instead of 40 320—but it means the caller must not hold on to the yielded slice after the next iteration. In `FindBestHamiltonian` this is not a problem, because each permutation is consumed immediately, and the only place where the slice is kept is `slices.Clone(p)`, which makes a copy.

The main loop is where the search actually happens:

```go
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
```

Each permutation is a candidate route. We sum the weights along it, compare with the best result so far using the caller-supplied `cmp`, and keep the winner. The `cycle` flag is a small extension for Day 13: when it is `true`, the route is closed by appending the starting location to the end.

### Part One

The distances are stored in a map keyed by an unordered pair of location names. `MakeEdge` normalises the pair so that `Edge{"London", "Dublin"}` and `Edge{"Dublin", "London"}` are the same key. This works because the graph is undirected.

For each permutation, we sum the edge weights along the path and compare with the best result so far. The comparison is passed in as a function, `cmp`, which lets the same code handle both the minimum (Part One) and the maximum (Part Two):

```go
result, _ := common.FindBestHamiltonian(locs, dists, false, math.MaxInt, func(a, b int) bool {
	return a < b
})
```

The initial value is `math.MaxInt`, so any real path beats it.

See [part1.go](part1.go).

### Part Two

The same call, with two changes: the initial value is `math.MinInt`, and the comparison is reversed:

```go
result, _ := common.FindBestHamiltonian(locs, dists, false, math.MinInt, func(a, b int) bool {
	return a > b
})
```

## Complexity

Let *n* be the number of locations.

- **Time:** *O*(*n*! · *n*). We enumerate *n*! permutations, and for each one we sum *n* − 1 edge weights. For *n* = 8 this is 40 320 × 7 ≈ 280 000 operations. The factorial dominates everything else; the extra *n* factor is a detail.

- **Space:** *O*(*n*) for the current permutation and the best path. The recursion in `Permute16` uses *O*(*n*) stack frames, and the accumulator is a fixed-size array of sixteen integers. As usual, we do not count the input itself.

## Associations

The Traveling Salesman Problem is one of those classical problems that, once you recognize it, starts appearing everywhere. Here are four examples from real research where TSP (or its variants) shows up in unexpected places.

- **Astronomy: scheduling observations.** The problem of ordering telescope targets is formalized as the *Traveling Telescope Problem*, a TSP variant with time windows and time-dependent travel costs. A recent paper models the multi-telescope case as a *multiple Traveling Salesperson Problem*, where several telescopes divide the target list between themselves. See [here](https://www.semanticscholar.org/paper/An-mTSP-model-for-multiple-telescope-short-term-Han-Cater/151b9cdde601b579badb28c9d33b48e4fa5cd524).
- **DNA: genome reconstruction.** In radiation hybrid mapping, a distance matrix between DNA sequences is used to recover their linear order: the Hamming distance between columns gives the edge weights, and a TSP tour recovers the correct ordering. See [here](https://lectures.molgen.mpg.de/Algorithmische_Bioinformatik_WS0405/material/Reinert_Script.pdf).
- **Archaeology: chronological seriation.** Flinders Petrie encoded each grave as a binary vector of pottery types and proposed that ordering the graves so consecutive vectors differ as little as possible would recover their chronology. A 2012 paper traces this history and shows Petrie was solving a weighted Hamiltonian path problem, equivalent to the symmetric TSP. See [here](https://www.zib.de/userpage//groetschel/pubnew/paper/gertzgroe2012.pdf).
- **Embroidery: minimizing thread changes.** A US patent on embroidery thread management formalizes the design sequence as a TSP variant, where the “distance” between two designs is the number of thread colors the machine lacks for the next one. The patent notes that testing all orderings becomes infeasible for many designs, motivating heuristics. See [here](https://patentimages.storage.googleapis.com/38/18/c8/17b8bb77d0869d/US20100017011A1.pdf).

## Notes

For *n* = 8 there is a classical alternative to brute force: the Held–Karp algorithm, which solves the TSP in *O*(*n*² · 2ⁿ) by dynamic programming over subsets of visited locations. For *n* = 8 it performs 8² · 2⁸ ≈ 16 384 operations—fewer than the 8! = 40 320 permutations we enumerate. But the DP table and the backward pass needed to recover the order add enough bookkeeping that the difference disappears in the noise; both methods finish before you can read this sentence. I chose brute force because the code is an order of magnitude simpler. Held–Karp becomes genuinely necessary somewhere around *n* ≈ 15–20, where 20! ≈ 2.4 × 10¹⁸ permutations cease to be tractable while 20² · 2²⁰ ≈ 4 × 10⁸ DP operations remain feasible. See [Held–Karp algorithm](https://en.wikipedia.org/wiki/Held%E2%80%93Karp_algorithm).

`Permute16` is capped at sixteen elements for a related reason: 16! ≈ 2.09 × 10¹³, and enumerating that many permutations takes hours rather than microseconds. The cap is a reminder that brute force is a tool for small inputs only, and that "small" here means roughly a dozen elements at most. For Day 9 it is more than sufficient.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~20 ms |
| 2 | ~15 ms |

Both parts are dominated by process startup and input parsing, not by the enumeration itself. The permutation count is small (8! = 40 320), and each permutation costs a handful of additions and map lookups. The difference between the two parts is within noise—two separate binaries, two separate runs.

I also measured an alternative implementation that uses the [`permute`](https://github.com/etnz/permute) package instead of the hand-written `Permute16`. The results were indistinguishable from the table above: the iterator is not the bottleneck.

---

[← Previous: Day 8](day08.md) · [Up: Contents](../README.md) · [Next: Day 10 →](day10.md)
