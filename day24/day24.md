# Day 24: It Hangs in the Balance

## The Puzzle

Santa has a pile of presents, each with a weight, and he wants to divide them into groups of equal total weight—one for the front of the sleigh, one for the back, and so on. Each present goes into exactly one group, and every group must weigh the same.

Part One asks for the ideal configuration for three groups. The configuration is judged first by the number of presents in the first group (fewer is better), then by the product of their weights, called the quantum entanglement (lower is better). Part Two asks the same question for four groups.

## Input

About thirty lines, one per present:

```
1
2
3
7
11
13
...
```

Each line is a single positive integer: the weight of one present. The presents are distinct even when their weights coincide—two presents of weight 5 are two different presents, and swapping them produces the same group by weight but a different selection. This matters for counting configurations, though the puzzle asks only for the best one.

The number of presents in the real input is around thirty, and the weights are small—mostly below a hundred. The total weight is therefore a few thousand, which is the parameter that shapes the algorithm: the target weight for each group is the total divided by the number of groups, and it is this value that bounds the dynamic program.

The input does not specify the number of groups. Part One fixes it at 3, Part Two at 4, and the same program is run twice with a different argument. The total weight must be divisible by the number of groups, or no valid division exists; for the real input, it is divisible by both 3 and 4.

## The Approach

By Day 24 I was tired of brute force, even though it had served well, and I wanted to show a dynamic-programming solution. There are two catches. First, the optimisation function is not trivial—it is not a simple minimum or maximum, but a lexicographic comparison of two criteria. Second, and more importantly, I solved only half the problem: I found the optimal set of presents for the passenger compartment, and I simply ignored whether the remaining presents could be balanced across the other compartments. It worked anyway, and the reason is worth a note. As soon as I saw that the answer was accepted, I moved to Part Two, and there the answer matched as well.

### Part One and Part Two

Both parts share the same `solve` function, which takes the number of compartments as a parameter. Part One passes 3, Part Two passes 4:

```go
func solve(numCompartments int) int {
	...
	totalWeight := sum / numCompartments
	...
}
```

The target weight per compartment is the total weight divided by the number of compartments. The puzzle guarantees this is an integer for the given input; if it were not, no valid division would exist, and the algorithm would return garbage.

The core is a dynamic program over weights. `bestScores[weight]` holds the best first-group configuration for a group of exactly that weight, where "best" is defined by the puzzle's criterion:

```go
type scoreT struct {
	num int
	qe int

	packsUsed uint64
}
```

The `num` field is the number of presents in the group, `qe` is their product (the quantum entanglement), and `packsUsed` is a bitmask recording which presents were used. The bitmask is not part of the puzzle's criterion—it is a technical necessity: without it, the DP could pick the same present twice, once for a smaller weight and again for a larger one.

The comparison function implements the puzzle's ordering:

```go
func cmpSocres(a, b scoreT) int {
	diffNum := a.num - b.num
	if 0 != diffNum {
		return diffNum
	}
	return a.qe - b.qe
}
```

Fewer presents wins; on a tie, the smaller product wins. The `packsUsed` field is deliberately not compared—two configurations with the same count and the same product are equivalent as far as the puzzle is concerned, and comparing the bitmask would introduce an arbitrary tie-breaker.

The DP loop fills the table from weight 1 up to the target:

```go
for weight := 1; weight <= totalWeight; weight++ {
	minScore := scoreT{math.MaxInt, math.MaxInt, 0}
	for n, p := range packages {
		if p > weight {
			continue
		}
		newScore := bestScores[weight - p]
		if 0 != newScore.packsUsed & (1 << n) {
			continue
		}
		newScore.num++
		newScore.qe *= p
		newScore.packsUsed |= 1 << n

		if 0 > cmpSocres(newScore, minScore) {
			minScore = newScore
		}
	}
	bestScores[weight] = minScore
}
```

For each weight, the algorithm tries every present that could be the last one added, and looks up the best configuration for the remaining weight. If that configuration already uses the present in question, the candidate is rejected—this is what prevents double-counting. Otherwise, the candidate's count is incremented, its product is multiplied by the present's weight, and its bitmask is updated. The best candidate wins.

The answer is the product stored at the target weight:

```go
return bestScores[totalWeight].qe
```

See [part1.go](part1.go) and [part2.go](part2.go).

## Complexity

Let *n* be the number of presents, *W* the target weight per compartment, and *k* the number of compartments.

- **Time:** *O*(*W* · *n*) for the dynamic program. The outer loop runs from 1 to *W*, and the inner loop tries every present. Each candidate is checked against the current best in constant time, since the comparison is a few integer operations. For the real input, *W* is a few hundred and *n* is about thirty, so the DP does on the order of ten thousand operations. The answer is read off the table at index *W*, which is *O*(1).

  The number of compartments *k* does not appear in the time bound, and this is the source of the bug: the algorithm computes the optimal first group but never checks whether the remaining presents can actually be divided into the required number of groups. It is an *O*(*W* · *n*) algorithm for a problem whose full solution would require an additional feasibility check—one that is exponential in the general case, as we will see in `Notes`.

- **Space:** *O*(*W*) for the `bestScores` table, plus *O*(*n*) for the packages. Each entry in the table is a small struct: two integers and a 64-bit mask. For *W* in the hundreds, the table is a few kilobytes.

## Associations

Classical knapsack packing is a dull affair. What matters here is stability—the sleigh must fly without tipping over. The same problem appears throughout engineering, from helicopter rotors to spacecraft to ecosystems, and the examples below are all about balance rather than packing.

- **Gyroscopes.** A gyroscope stays upright because angular momentum resists changes to its axis. The effect is used everywhere from aircraft instruments to smartphone sensors, and it is the simplest illustration of stabilisation: a system that maintains its orientation without active correction. Our sleigh has no gyroscope, but the idea is the same—find a configuration that does not tip. See [Gyroscope](https://en.wikipedia.org/wiki/Gyroscope).

- **The N-1 rocket and the KORD system.** The Soviet N-1 moon rocket failed four times, and the failures were not in the engines alone. The KORD system was designed to shut down engines that showed signs of trouble, maintaining thrust symmetry. On the first launch, electrical interference caused KORD to falsely shut down engine 12, then engine 24 for symmetry, and finally all 28 engines at T+68.7 seconds. The system meant to stabilise the rocket destroyed it—worse, by removing engines it removed the very degrees of freedom that could have compensated for the failure. A gimballed thrust vector would have allowed the remaining engines to correct for the loss; shutting down the opposite engine only made the asymmetry permanent. Stability, it turns out, is not the same as symmetry. See [N-1 rocket](https://en.wikipedia.org/wiki/N1_(rocket)).

- **The lynx and hare cycle.** The ten-year oscillation of Canada lynx and snowshoe hare populations is the classic example of predator-prey dynamics. The populations rise and fall in a stable pattern, kept in balance by predation, food supply, and climate. It is a system where stability emerges from the interaction of opposing forces—not from any single component being stable on its own. Our sleigh is the same: each compartment is only balanced if the others are. See [Lotka–Volterra equations](https://en.wikipedia.org/wiki/Lotka%E2%80%93Volterra_equations).

## Notes

- **Why the DP is correct.** The dynamic program finds the best first-compartment configuration because the problem has **optimal substructure**: the best way to reach weight *w* is the best way to reach *w − p* plus the present *p*. The comparison between configurations is a **lexicographic order** on `(num, qe)`: fewer presents wins, and on a tie the smaller product wins. This order is **monotone** under adding a present—both the count and the product can only grow, never shrink—which is what makes the greedy choice inside the DP safe. The **principle of optimality** then guarantees that the table entry for the target weight is the global optimum, not just a local one. The DP is not a heuristic; it is exact, and it would remain exact for any input where the target weight is an integer.

- **Why solving half the problem was enough.** The DP chooses the best first compartment but never checks whether the remaining presents can be split into the required number of equal-weight groups. On this input, that check is unnecessary—and it is worth seeing why. The helper program `count-dissections.go` enumerates every subset of the input whose weight equals the target, and then checks whether the remainder can be partitioned into the remaining compartments. For Part One, it finds 216 871 subsets of weight 508, and **all** of them leave a splittable remainder. For Part Two, it finds 54 966 subsets of weight 381, and again **all** of them can be completed. The input is constructed so that any valid first compartment extends to a full solution, which is why the DP's answer was accepted on the first try. This is a property of the input, not of the algorithm: a different set of weights could easily have subsets of the right sum that cannot be completed, and the DP would return an answer for a configuration that does not actually exist. The lesson is that a solution can be correct for an input without being correct for the problem.

- **Primes and Goldbach.** The input is a set of primes—specifically, the primes from 3 to 113 with 2, 7, and 61 omitted—plus a leading 1. The total is 1524, which is divisible by both 3 and 4, giving target weights of 508 and 381. The puzzle is therefore a partition of a set of primes into equal-sum groups, which is thematically close to Goldbach's conjecture: both are about expressing numbers as sums of primes. The connection is not technical—Goldbach speaks of representing a single number, while we are partitioning a set—but it is a pleasant coincidence that the input for a puzzle about balance happens to be built from the atoms of additive number theory. The ternary version of Goldbach's conjecture (every odd number greater than 5 is the sum of three primes) was proved by Harald Helfgott in 2013; the binary version (every even number greater than 2 is the sum of two primes) remains open, though it has been verified up to 4 × 10¹⁸. Neither version is used by the algorithm, but they are the kind of thing one thinks about when staring at a list of primes that must be divided into equal parts.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~2 ms |
| 2 | ~1.5 ms |

Both parts are dominated by process startup. The DP itself is a few thousand iterations over a table of a few hundred entries, and each iteration does a constant amount of work. The difference between the two parts is within run-to-run variation, and there is nothing here that would benefit from optimisation.

---

[← Previous: Day 23](day23.md) · [Up: Contents](../README.md) · [Next: Day 25 →](day25.md)
