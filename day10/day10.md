# Day 10: Elves Look, Elves Say

## The Puzzle

The elves are playing a game: given a sequence of digits, read it aloud. "1" becomes "one 1", or "11". "11" becomes "two 1s", or "21". "21" becomes "one 2, one 1", or "1211". And so on.

Part One asks for the length of the sequence after forty iterations of this process, starting from the input line. Part Two asks for the length after fifty iterations.

## Input

A single line of digits, like `1113122113`. The sequence grows quickly, so the length after fifty iterations is around seven million characters—large enough that we cannot print it, only measure it.

## The Approach

The process is a straightforward string transformation. For each iteration, we scan the current string, count runs of identical digits, and write each run as `count` followed by the digit. There is no clever shortcut: the result of iteration *n* depends entirely on iteration *n* − 1, and the sequence does not stabilize or cycle.

The function `lookAndSay` does one iteration:

```go
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
```

The inner loop counts how many times the current digit repeats. As soon as the next digit differs—or the string ends—we return to the outer loop, reset the counter to 1, and continue scanning from the current position. The outer loop has no increment of its own: the index `i` is advanced by the inner loop and by the initial `i++`, so the outer loop simply picks up where the previous run left off.

The `strings.Builder` is preallocated with `len(s) * 2`, which is a reasonable guess: the sequence grows, but not explosively per iteration, and doubling the buffer is enough to avoid most reallocations. The actual growth factor is closer to 1.3, so the buffer is slightly oversized, which is fine—one allocation is cheaper than several.

For Part One and Part Two, the only difference is the number of iterations. The function `lookAndSayCount` wraps this in a loop:

```go
func lookAndSayCount(s string, num int) int {
	seq := s
	for i := 0; i < num; i++ {
		seq = lookAndSay(seq)
	}

	return len(seq)
}
```

### Part One

Forty iterations, then `len(seq)`:

```go
result := lookAndSayCount(lines[0], 40)
```

See [part1.go](part1.go).

### Part Two

The same call, with fifty iterations:

```go
result := lookAndSayCount(lines[0], 50)
```

See [part2.go](part2.go).

## Complexity

Let *L* be the length of the sequence after the final iteration, and *n* the number of iterations.

- **Time:** *O*(*L*). Each iteration scans the entire string once, and the total work is the sum of the lengths over all iterations. Since the sequence grows geometrically, the last iteration dominates: the sum is *O*(*L*), not *O*(*n* · *L*). For Part Two, *L* ≈ 7 × 10⁶, which matches the measured runtime.

- **Space:** *O*(*L*). We keep the current sequence and build the next one in the same iteration; the old one is discarded when the new one is assigned. The peak usage is roughly twice the length of the current sequence, which is still *O*(*L*).

## Associations

- **The look-and-say sequence.** This puzzle is a direct reference to the look-and-say sequence, studied by John Conway in the 1980s. Conway proved that the sequence grows at a fixed rate (the Conway constant) and that it decomposes into a finite set of "atoms" that never interact with each other after a certain point. See [Look-and-say sequence](https://en.wikipedia.org/wiki/Look-and-say_sequence).

- **The cosmological theorem.** Conway's analysis of the look-and-say sequence led him to a result he called the "cosmological theorem": no matter what sequence you start from (with one exception), after finitely many iterations it decomposes into a fixed set of 92 elements that evolve independently. The name comes from the fact that Conway saw an analogy with the evolution of the universe—initial chaos, then stable structure. See [Cosmological theorem](https://en.wikipedia.org/wiki/Look-and-say_sequence#Cosmological_decay).

- **Rewriting systems.** One step of look-and-say is a parallel rewriting rule: every run of *k* identical digits *d* is replaced by "*k*" followed by *d*. This is not a context-free grammar—the rules act on runs, not on single symbols—but a semi-Thue system, a set of rewriting rules applied simultaneously to the whole string. Iterating the same rule produces the sequence. See [Semi-Thue system](https://en.wikipedia.org/wiki/Semi-Thue_system).

## Notes

The growth rate of the sequence is the key property worth remembering. Each iteration multiplies the length by roughly 1.303577—the Conway constant—so the length after *n* iterations is approximately *L* · 1.303577ⁿ, where *L* is the length of the input. For the given input, ten additional iterations in Part Two multiply the length by about 13.9, which explains why Part Two takes an order of magnitude longer than Part One. The sequence never stabilizes, never cycles, and never repeats: every iteration produces genuinely new digits, and the only way to get the answer is to run the process.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~17 ms |
| 2 | ~150 ms |

Part One runs forty iterations, Part Two fifty. Ten extra iterations multiply the sequence length by roughly 1.303577¹⁰ ≈ 13.9, and the measured ratio (150 / 17 ≈ 8.8) is a bit lower. The gap comes from the smaller sequence in Part One spending proportionally more time in allocation relative to scanning. Either way, the final iteration dominates the total cost: both parts are bounded by the length of the last sequence, not by the per-iteration overhead.

---

[← Previous: Day 9](day09.md) · [Up: Contents](../README.md) · [Next: Day 11 →](day11.md)
