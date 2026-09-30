# Day 5: Doesn't He Have Intern-Elves For This?

## The Puzzle

Santa needs help figuring out which strings in a list are "nice" and which are "naughty". The list is the input, one string per line. The rules for niceness differ between the two parts.

Part One declares a string nice if it contains at least three vowels, at least one letter that appears twice in a row, and none of the forbidden substrings `ab`, `cd`, `pq`, `xy`.

Part Two changes the rules: a string is nice if it contains a pair of adjacent letters that appears at least twice without overlapping, and at least one letter that repeats with exactly one character between the two occurrences.

## Input

About a thousand strings, one per line. In the real input each string is exactly 16 lowercase letters long, but the puzzle does not state any length restriction.

The example from the puzzle statement is short:

```
ugknbfddgicrmopn
aaa
jchzalrnumimnmhp
haegwjzuvuyypxyu
dvszwmarrgswjxmb
```

The first two lines are nice by the Part One rules; the last three are naughty, each violating a different rule.

## The Approach

Both parts are line-by-line validation. For each string we check a handful of properties and count how many strings pass all of them.

The interesting design question is how to express the checks. Some of them are naturally expressed as regular expressions, others are not. I use regex where it fits and hand-written loops where it does not.

### Part One

Three conditions: at least three vowels, at least one doubled letter, no forbidden substrings.

The vowels and the forbidden substrings are easy to express as regexes:

```go
re3Vowels := regexp.MustCompile(`[aeiou].*[aeiou].*[aeiou]`)
reBanned := regexp.MustCompile(`(ab|cd|pq|xy)`)
```

The doubled letter is a different story. It would be natural to write `([a-z])\1`, but Go's `regexp` package implements RE2, which does not support backreferences. This is a deliberate design choice in RE2: it keeps matching linear in the input size, at the cost of expressive power. So we write the check by hand:

```go
hasCouple := func(s string) bool {
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			return true
		}
	}
	return false
}
```

The three conditions are then combined with short-circuit evaluation:

```go
if re3Vowels.MatchString(line) && hasCouple(line) && !reBanned.MatchString(line) {
	result++
}
```

See [part1.go](part1.go).

### Part Two

The rules change, and neither is easy to express as a regex.

The first requires a pair of adjacent letters that appears at least twice without overlapping. We keep a map from each pair to the position of its first occurrence. When we see the same pair again, we check that the earlier position is not adjacent to the current one—otherwise the two occurrences would share a letter and would not count as two separate pairs.

```go
seen := make(map[[2]byte]int)

for c := 0; c < len(s)-1; c++ {
	seq := [2]byte{s[c], s[c+1]}
	if pos, ok := seen[seq]; ok {
		if pos < c-1 {
			return true
		}
	} else {
		seen[seq] = c
	}
}
```

The second rule requires a letter that repeats with exactly one character between the two occurrences. That is a simple loop with a step of three positions:


```go
for c := 0; c < len(s)-2; c++ {
	if s[c] == s[c+2] {
		return true
	}
}
```

Both functions are used in the same way as in Part One:

```go
if hasDoublePair(line) && hasEmbraced(line) {
	result++
}
```

See [part2.go](part2.go).

## Complexity

Let *n* be the number of strings and *m* the length of a string.

- **Time:** *O*(*n* · *m*). Part One uses two regexes and one loop over the string; the regexes are linear in *m* because RE2 guarantees it, and the loop is linear as well. Part Two uses a map keyed by a two-byte array, which is *O*(1) per lookup on average, plus a loop over the string. So each string is processed in time linear in its length.

- **Space:** *O*(*m*) for Part Two, where the map can hold up to *m* - 1 distinct pairs. Part One uses a fixed number of regex states—determined by the regex specification, not by the input—so we count it as *O*(1).

## Associations

- **Early static analysis tools.** Checking each string against a set of rules, with no memory of the previous lines, is how the first static analysis tools worked. Lint, created by Stephen C. Johnson at Bell Labs in 1977, scanned source code for local patterns without building any model of the program. See [Lint](https://en.wikipedia.org/wiki/Lint_(software)). The first security tools followed the same recipe: ITS4 and RATS searched for dangerous function calls and left the rest to manual review. See [ITS4](https://homes.cs.washington.edu/~yoshi/papers/ACSAC00/its4.pdf) and [RATS](https://github.com/andrew-d/rough-auditing-tool-for-security). Modern SAST tools parse source into a syntax tree and trace data flow across function boundaries; our solution is at the level of the early tools—line-by-line pattern matching, no context, no model.

- **Why static analysis can never be perfect.** The reason early tools were "rough" is not laziness: Rice's theorem says that no general algorithm can decide non-trivial semantic properties of programs—whether a function dereferences null, releases all file handles, or has any other interesting behavior. Every static analyzer must trade off between false positives and false negatives. Our solution is firmly in the heuristic camp: it checks local patterns and says nothing about the program as a whole. See [Rice's theorem](https://en.wikipedia.org/wiki/Rice%27s_theorem).

## Notes

The regex versus hand-written loop split is worth a comment. For the vowels and forbidden substrings, the regex reads almost like the puzzle statement, and there is no reason to write the loop by hand. For the doubled letter, the regex would be `([a-z])\1`, which is exactly what the puzzle describes—but RE2 does not support backreferences.

It is tempting to ask whether the rules could be combined into a single regex. In principle they can: the language of nice strings is regular. In practice, the result is a monster. The doubled-letter rule alone, without backreferences, must be written as `aa|bb|cc|...|zz`—twenty-six alternatives. The first rule of Part Two, which requires a pair of letters to appear twice without overlapping, needs all 26² = 676 pairs enumerated explicitly: `aa.*aa|ab.*ab|...|zz.*zz`. And that is just one rule. Combining two such rules into a single regex requires intersecting two languages, which multiplies the size again.

This is the lesson: regular expressions are a powerful tool, and we use them liberally throughout the event. But they are not a universal tool. When the language is a product of several independent constraints, the regex becomes a machine-generated monster, and the hand-written code stays readable.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~4 ms |
| 2 | ~4 ms |

Both parts are fast because the input is small—about a thousand short strings. Each string is scanned a constant number of times, one pass per rule, so the total work is linear in the input size. There is no point in optimising further.

---

[← Previous: Day 4](day04.md) · [Up: Contents](../README.md) · [Next: Day 6 →](day06.md)
