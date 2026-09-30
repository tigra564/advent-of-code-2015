# Day 17: No Such Thing as Too Much

## The Puzzle

You have a collection of containers, each with a capacity in liters, and a fridge that holds exactly 150 liters. You need to choose some subset of the containers to fill, such that their capacities sum to exactly 150.

Part One asks for the number of distinct subsets that sum to 150.

Part Two asks for the number of subsets of the *minimum possible size* that sum to 150—that is, among all valid fillings, how many use the fewest containers.

## Input

About twenty lines, one per container:

```
43
3
4
10
21
```

Each line is a single positive integer: the capacity of one container in liters. The containers are distinct even when their capacities are equal, so two containers with the same size count as two different choices. This matters for the answer: the puzzle asks for the number of subsets, not the number of capacity multisets.

The number of containers in the real input is twenty. This is what makes the brute-force enumeration possible: 2²⁰ = 1 048 576 subsets, each summing at most twenty values. The subset count grows exponentially with the number of containers, so the puzzle is calibrated to be just large enough to be tedious by hand and just small enough to be tractable by machine.

## The Approach

The counting version of the subset-sum problem has a classical dynamic-programming solution that runs in *O*(*n* · *S*) time, where *n* is the number of containers and *S* is the target volume. For *n* = 20 and *S* = 150, that would be about 3000 operations—three orders of magnitude fewer than the 2²⁰ ≈ 1 million subsets we enumerate. I chose the enumeration anyway, for two reasons. First, it is simpler: a bitmask and a nested loop, with no table to size or indices to juggle. Second, the problem is small enough that the difference is invisible: both approaches finish in milliseconds, and the enumeration's early exit on overshoot narrows the gap further. If the number of containers were larger—say, thirty or forty—the dynamic-programming approach would become necessary, and the code would look very different.

Both parts enumerate subsets of the container list. Twenty containers give 2²⁰ = 1 048 576 subsets, which is a small enough space to walk through exhaustively. The only difference between the parts is what is counted: Part One counts all subsets summing to 150, Part Two counts only those of minimum size.

The enumeration uses a bitmask. Each subset corresponds to an integer from 0 to 2ⁿ − 1, where bit *i* is set if container *i* is included:

```go
for mask := range uint32(1) << nCont {
	sum := 0
	for i := range nCont {
		if 0 != mask&(uint32(1)<<i) {
			sum += containers[i]
			if sum > nLiters {
				break
			}
		}
	}
	if sum == nLiters {
		result++
	}
}
```

The inner loop walks the bits of the mask from least significant to most significant. For each set bit, it adds the corresponding container's capacity to the running sum. If the sum exceeds 150, the loop breaks early—there is no way to recover by adding more containers, since capacities are positive.

This early exit is the reason the containers are sorted in descending order before the enumeration:

```go
slices.SortFunc(containers, func(a, b int) int {
	return cmp.Compare(b, a)
})
```

With large containers first, the sum crosses 150 sooner, and the inner loop terminates after fewer iterations. If the containers were sorted ascending, the sum would creep up slowly and the break would rarely trigger. The saving is not asymptotic—the worst case is still *O*(*n*) per mask—but it is real, and it costs one line of code. The measurements bear this out: the optimised version runs roughly two and a half times faster than the version without sorting and early exit. That is a substantial gain for a change that takes a minute to write.

### Part One

Part One counts every subset whose sum is exactly 150:

```go
if sum == nLiters {
	result++
}
```

The early break on overshoot means the comparison at the end is only reached for subsets that did not exceed the limit, so `sum == nLiters` is equivalent to `sum >= nLiters` in this context.

See [part1.go](part1.go).

### Part Two

Part Two needs the minimum number of containers used, and then the count of subsets achieving that minimum. The mask loop is the same, but the inner loop also tracks how many containers were included:

```go
contsUsed := 0
for i := range nCont {
	if 0 != mask&(uint32(1)<<i) {
		sum += containers[i]
		if sum > nLiters {
			break
		}
		contsUsed++
	}
}
```

The `contsUsed` increment is placed after the break check, so it only counts containers that did not push the sum over the limit. This is correct because a subset that overshoots is discarded anyway.

When a valid subset is found, the code compares its size to the current minimum and updates accordingly:

```go
switch {
case contsUsed < minUsed:
	minUsed = contsUsed
	result = 1
case contsUsed == minUsed:
	result++
}
```

If a smaller subset appears, the previous count is discarded and the counter resets to 1. If an equally small subset appears, the counter is incremented. This is the standard pattern for "count the number of optimal solutions" when the optimum is discovered during the scan.

See [part2.go](part2.go).

## Complexity

Let *n* be the number of containers and *S* the target volume.

- **Time:** *O*(2ⁿ · *n*) in the worst case, but substantially less in practice. The outer loop runs 2ⁿ times, and the inner loop walks up to *n* bits. The early exit cuts the inner loop short whenever the running sum exceeds *S*, which happens frequently because the containers are sorted in descending order. The dynamic-programming alternative would be *O*(*n* · *S*), which for this input is about 3000 operations against roughly 1 million—but the constant factors and the early exit make the enumeration competitive at this size.

- **Space:** *O*(*n*) for the container slice and the sort. The enumeration itself uses only a few scalars: the mask, the running sum, and the counters.

## Associations

- **Knapsack cryptography.** The subset-sum problem at the heart of this puzzle was the basis of the first public-key cryptosystem, published by Ralph Merkle and Martin Hellman in 1978. The idea was to disguise a superincreasing sequence—one where each element exceeds the sum of all previous elements—by multiplying it by a secret factor modulo a large number. For a superincreasing sequence, the subset sum is solved in a single pass by a greedy algorithm, so the legitimate receiver could decrypt quickly; an attacker would face a general subset-sum instance, which is NP-complete. The system was broken by Adi Shamir in 1984, and most subsequent knapsack cryptosystems have also been broken. The history is a useful reminder that a problem being NP-complete in general does not guarantee that a particular family of instances is hard. See [Merkle–Hellman knapsack cryptosystem](https://en.wikipedia.org/wiki/Merkle%E2%80%93Hellman_knapsack_cryptosystem).

## Notes

The dynamic-programming solution deserves a closer look, even though I did not use it. The recurrence for counting subsets with a given sum is `A[i][j] = A[i-1][j] + A[i-1][j - w_i]`, where `A[i][j]` is the number of ways to reach sum `j` using the first `i` containers, and `w_i` is the capacity of the `i`-th container. Part One is a direct application: `A[n][S]` is the answer. Part Two is subtler, because it asks for the count of *minimum-size* subsets, not all subsets. The DP would need a second dimension tracking the number of containers used, or a two-pass approach: first compute the minimum size, then count subsets achieving it. Both are doable, and both are more code than the enumeration. The enumeration's simplicity is the reason it won here, not its speed.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Optimised | Unoptimised |
|------|-----------|-------------|
| 1 | ~28 ms | ~76 ms |
| 2 | ~30 ms | ~77 ms |

The "unoptimised" column is the same code with the sort removed and the early exit on overshoot disabled. Both changes are three lines each, and together they cut the runtime by a factor of about 2.6. This is the kind of optimisation that pays for itself immediately: the change is trivial to write, trivial to verify, and the gain is large enough to matter on a problem of this size.

The two parts are almost identical in runtime, as expected: they enumerate the same subsets with the same early exit, and Part Two only adds a counter and a comparison per valid subset.

---

[← Previous: Day 16](day16.md) · [Up: Contents](../README.md) · [Next: Day 18 →](day18.md)
