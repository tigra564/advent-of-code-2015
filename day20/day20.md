# Day 20: Infinite Elves and Infinite Houses

## The Puzzle

An infinite supply of elves delivers presents to an infinite row of houses, numbered 1, 2, 3, ... Each elf has a number, and elf number *n* visits every *n*-th house—house *n*, house 2*n*, house 3*n*, and so on—delivering a fixed number of presents each time.

Part One asks for the lowest house number that receives at least the target number of presents, given that each elf delivers 10 presents per visit.

Part Two changes the rules: each elf delivers 11 presents per visit, and stops after 50 houses. Part Two asks for the lowest house number that reaches the target under these rules.

## Input

A single line with one integer: the target number of presents. For the real input it is in the tens of millions, which is the parameter that drives the algorithm: we need to find the smallest house whose present count reaches this target.

## The Approach

Both parts compute the number of presents delivered to each house, then find the first house that reaches the target. The difference is in the rules: Part One counts every elf whose number divides the house, and Part Two stops each elf after 50 visits.

The naive way to compute the present counts would be to iterate over houses and, for each house, find all its divisors. That would be *O*(*H*√*H*) or worse, where *H* is the number of houses. The better way—the one used in both parts—is to invert the loop: iterate over elves, and for each elf, add its contribution to every house it visits. This is the classic sieve pattern, and it turns the divisor sum into a single pass over the multiples of each number.

### Part One

Part One works in units of "tens of presents" to keep the numbers small. The target is divided by 10, rounding up:

```go
maxPresents, _ := strconv.Atoi(lines[0])
// mtp — maximum of tens of presents. Here +9 is in case the input is not divisible by 10
mtp := (maxPresents + 9) / 10
```

The array size is not the target itself, but an upper bound on the answer. The bound comes from two facts about the sum-of-divisors function *σ*(*n*).

The first is that *σ*(*n*) ≥ *n*, with equality only for *n* = 1. A house *h* therefore receives at least *h* tens of presents, so *h* ≤ *mtp* is a valid bound—but a weak one.

The second fact is stronger. Abundant numbers are those with *σ*(*n*) > 2*n*, and their natural density is about 0.25 [1](https://en.wikipedia.org/wiki/Abundant_number). A quarter of all houses receive more than twice their number in tens of presents, so if *mtp* is large, the answer is at most about *mtp* / 2.

The code uses a third bound, tighter still. The factorial 9! = 362880 is not highly abundant—there is a smaller number with a larger divisor sum—but its divisor sum is *σ*(9!) = 1481040 [2](https://en.wikipedia.org/wiki/Highly_abundant_number), and *σ*(*k*·*n*) ≥ *σ*(*n*)·*k* for any *k* [1](https://en.wikipedia.org/wiki/Abundant_number). So for any multiple *k*·9!, the divisor sum is at least *k*·1481040. Setting this above *mtp* gives a bound on *k*, and hence on the answer:

```go
const (
	fact9 = 362880
	sigma_fact9 = 1481040
)
upperBoundHouse := fact9 * (mtp / sigma_fact9 + 1)
```

The array is allocated after the bound is known, so its size is exactly the bound and never needs to grow:

```go
// Here +1 is to accomodate the house #0, which is non-existent, but indices start with 1
presents := make([]int, upperBoundHouse + 1, upperBoundHouse + 1)
```

The sieve fills the array. For each elf *i*, it adds *i* to every multiple of *i*. The addition is followed by a check: if any house reaches the target, the upper bound is reduced to that house's index, and the loops stop early:

```go
for i := 1; i <= upperBoundHouse; i++ {
	for j := i; j <= upperBoundHouse; j += i {
		presents[j] += i
		if (presents[j] >= mtp) {
			upperBoundHouse = j
		}
	}
}
```

This shrink-as-you-go approach is the real optimisation. The sieve does not need to run to the initial upper bound if the answer is found earlier—and the answer is usually found much earlier. The upper bound starts as an estimate from *σ*(9!), but as soon as a house is confirmed to reach the target, the bound tightens, and the remaining iterations shrink accordingly.

The final answer is whatever `upperBoundHouse` holds when the loops finish: the smallest house whose present count reached the target.

See [part1.go](part1.go).

### Part Two

Part Two changes two things: the multiplier is 11 instead of 10, and each elf stops after 50 visits. The second change is the important one—it truncates the divisor sum, so each house receives contributions only from its first 50 divisors.

The truncation also breaks the connection with the sum-of-divisors function *σ*(*n*), which is what made the bounds in Part One possible. The present count per house is no longer *σ*(*h*) but a partial sum of its divisors, and the theory of abundant and highly abundant numbers does not apply. So Part Two starts with a much weaker bound—the target itself—and relies entirely on the shrink-as-you-go optimisation to cut the work down.

The code works in units of "elevens of presents" for the same reason Part One works in tens:

```go
maxPresents, _ := strconv.Atoi(lines[0])
// mep — maximum "elevens" of presents
mep := (maxPresents + 10) / 11
presents := make([]int, mep + 1, mep + 1)
```

The `+10` is the round-up for base 11: `(x + 10) / 11` is `ceil(x / 11)`. The array is sized to `mep`, which is a valid upper bound because the truncated divisor sum is at most the full one, and the full one is at most *mep* for the answer.

The sieve now has a counter that breaks the inner loop after 50 iterations:

```go
for i := 1; i <= upperBoundHouse; i++ {
	k := 0
	for j := i; j <= upperBoundHouse; j += i {
		if k++; k > 50 {
			break
		}
		// We don't multiply by eleven here, we have already divided the target properly
		presents[j] += i
		if (presents[j] >= mep) {
			upperBoundHouse = j
		}
	}
}
```

The counter `k` is incremented and checked before the addition, so the first 50 houses get their present, and the 51st does not. The multiplication by 11 is folded into the target: since both sides are divided by 11, the comparison `presents[j] >= mep` is equivalent to comparing the real present count with the real target.

The shrink-as-you-go trick from Part One is still here: `upperBoundHouse` is reduced whenever a house reaches the target. In Part Two this matters even more, because the initial bound `mep` is the target itself, not a fraction of it. Even with this weaker starting point, the dynamic bound brings the runtime down to a fraction of what a full sieve to `mep` would cost.

See [part2.go](part2.go).

## Complexity

Let *T* be the target number of presents, and *N* the number of houses actually processed by the sieve.

- **Time:** *O*(*N* log *N*) for both parts. The sieve visits every multiple of every elf up to *N*, which is the harmonic sum: *N* + *N*/2 + *N*/3 + ... = *N* · *H*<sub>*N*</sub> ≈ *N* ln *N*, where *H*<sub>*N*</sub> is the *N*-th harmonic number. In Part Two, each elf stops after 50 visits, so the inner loop is capped. The cap of 50 is the same order of magnitude as log₂ *N* for the given input, so the asymptotic bound is unchanged—the cap only trims the constant factor.

  The value of *N* is not the initial upper bound but the final `upperBoundHouse` when the loops stop. In both parts, this is the answer itself: the sieve stops as soon as the first house reaches the target. So the actual work is proportional to the answer, not to the initial bound, which is what makes the shrink-as-you-go optimisation effective.

- **Space:** *O*(*N*) for the `presents` array. In Part One, *N* is the answer's upper bound computed from *σ*(9!); in Part Two, it is the target divided by 11. Both are linear in the input size.

## Associations

- **From Eratosthenes to the quadratic sieve.** The sieve in this puzzle is a direct descendant of the Sieve of Eratosthenes, the 2200-year-old algorithm for finding primes by crossing out multiples. The same technique—walk through a range and eliminate or mark numbers at regular intervals—is the engine behind the Quadratic Sieve (QS), one of the two fastest known algorithms for factoring large integers. QS factored RSA-129 in 1994, a 129-digit challenge that had stood for 17 years and was the first public break of an RSA number. The connection to our puzzle is not the problem but the method: our elf sieve marks every *n*-th house, and QS marks every *p*-th value. See [Quadratic sieve](https://en.wikipedia.org/wiki/Quadratic_sieve).

## Notes

The search for a ready-made answer is worth a note, because it illustrates a common temptation. Part One computes *σ*(*n*) for every house *n* up to the answer, and OEIS has a page for that sequence—A000203, the sum-of-divisors function. My first instinct was to look for a precomputed list: surely someone had already tabulated the values up to a few million, and I could just read off the first one that crosses the target. That is not how OEIS works. The page provides formulas, code, and references, but not a bulk data file of the sequence up to arbitrary limits. There are bounds on the *index* at which the sequence is tabulated, and those bounds are far below what this puzzle needs. So I wrote the sieve, which took twenty lines and runs in milliseconds.

The lesson is not that OEIS is unhelpful—it is an extraordinary resource, and the pages for A000203 and A005101 gave me the two facts I actually used: the density of abundant numbers and the value of *σ*(9!). But there is a difference between knowing a sequence and having its values in hand. For this puzzle, the sieve was the right tool, and the OEIS visit was a detour that happened to pay off in a different way.

The bounds in Part One are a small study in what "good enough" means. The trivial bound *h* ≤ *mtp* is correct but useless: sieving up to it would process about 2.9 million houses instead of the roughly 665 000 that actually matter. The array would be four times larger, but the work is quadratic in the bound—two nested loops—so the real cost is closer to (2.9 / 0.665)² ≈ 19 times more iterations. The abundant-number bound *h* ≤ *mtp* / 2 is better and rests on a real theorem, but it still leaves a factor of two on the table. The *σ*(9!) bound is sharper—it uses the specific magnitude of the target—and it is the one that made the sieve fast enough to finish in tens of milliseconds. None of these bounds is "the right answer"; they are successive refinements, each justified by a different fact about divisors.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~44 ms |
| 2 | ~63 ms |

Part Two is slower despite doing less work per elf—each elf stops after 50 visits, so the inner loop is capped. The reason is the initial bound: Part One starts from a tight estimate based on *σ*(9!), while Part Two starts from the target itself, which is several times larger. The shrink-as-you-go optimisation narrows the gap, but it takes time to kick in, and the extra iterations at the start are not recovered. The measured ratio (63 / 44 ≈ 1.4) reflects that difference.

---

[← Previous: Day 19](day19.md) · [Up: Contents](../README.md) · [Next: Day 21 →](day21.md)
