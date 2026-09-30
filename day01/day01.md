# Day 1: Not Quite Lisp

## The Puzzle

Santa is given a string of `(` and `)` characters. Each `(` moves him one floor up, each `)` moves him one floor down. Part One asks for the final floor after processing the entire string. Part Two asks for the position of the first character that takes him to the basement (floor `-1`).

## Input

A single line of opening and closing parentheses. The exact length is not specified, but it is large enough that the answer is not obvious from a quick glance.

The example from the puzzle statement is short and easy to follow by hand:

```
(()(()(
```

Processing it character by character gives floors `1, 2, 1, 2, 3, 2, 3`. The final floor is `3`. The first time the floor drops to `-1` does not happen here—this example is for Part One only.

The real input is read from a file (or from standard input) through `common.ReadInputLines`. The `common` package holds three functions used across all days: `ReadLinesFromFile`, `ReadLinesFromStdin`, and `ReadInputLines`, which picks one of the two based on command-line arguments. I will not describe the package separately; each function appears in the day where it is first used.

## The Approach

There is no clever algorithm here. The problem is a one-dimensional walk: the floor is a counter, and each character either increments or decrements it. Both parts are solved in a single pass over the string.

### Part One

We keep a single counter, start it at zero, and walk the string. For each character, we increment the counter on an opening parenthesis and decrement it on a closing one. When the walk is done, the counter holds the answer.

```go
result := 0
for _, c := range lines[0] {
	switch c {
	case ')':
		result--
	case '(':
		result++
	}
}
```

See [part1.go](part1.go).

### Part Two

The walk is the same; only the stopping condition changes. We need the position of the first character that takes the floor below zero. This means we cannot simply process the whole string—we have to check after every step and stop as soon as the condition is met.

```go
level := 0
result := 0
for p, c := range lines[0] {
	switch c {
	case ')':
		level--
	case '(':
		level++
	}
	if level < 0 {
		result = p + 1
		break
	}
}
```

See [part2.go](part2.go).

## Complexity

- **Time:** *O*(*n*), where *n* is the length of the input. One pass for Part One, and one early-exit pass for Part Two.
- **Space:** *O*(1). A single integer counter. As a rule, we do not count the space taken by the input itself.

## Associations

- **The ballot problem.** If we read an opening parenthesis as a vote for candidate A and a closing one as a vote for candidate B, the sequence of characters is a path in the classical ballot problem. The question "when does the walk first go below zero" is a variant of the question "when does candidate B first lead". See [Ballot problem](https://en.wikipedia.org/wiki/Ballot_problem).

- **Catalan numbers.** The number of strings that stay at or above zero after every prefix is the *n*-th Catalan number. Our input does not stay non-negative—otherwise Part Two would have no answer—but the counting problem is the same. See [Catalan number](https://en.wikipedia.org/wiki/Catalan_number).

## Notes

There is nothing tricky in this puzzle. The only thing worth a note is the order of operations in Part Two: the check must happen after the step, not before. If we check first, we look at the floor before the character has been processed, and the answer is off by one.

I do not include a benchmark for this day. The input is short enough that any reasonable implementation runs in microseconds; a benchmark would measure process startup overhead more than the algorithm. Starting with Day 2, the benchmarks become meaningful, and I include them from there on.

---

[Up: Contents](../README.md) · [Next: Day 2 →](day02.md)
