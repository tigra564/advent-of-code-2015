# Day 3: Perfectly Spherical Houses in a Vacuum

## The Puzzle

Santa is delivering presents in an infinite two-dimensional grid of houses. Each character in the input is a direction: `^`, `v`, `>`, or `<`. Santa starts at the origin and moves one house per character.

Part One asks how many houses receive at least one present if Santa does all the walking himself. Part Two introduces Robo-Santa: the two of them alternate moves, and we ask how many houses receive at least one present in total.

## Input

A single line of direction characters, one of `^`, `v`, `>`, `<`:

```
^>v<^>v<
```

The grid is infinite, so there are no bounds to check. The line length is not specified; the real input is roughly 8 KiB, which is enough that we cannot trace the walk by hand.

The walk is read from the same input stream as in previous days. No parsing is needed—each character is itself a direction.

## The Approach

Both parts track positions on the grid and count how many distinct positions are visited. Let us use a map with an integer value as a counter: it lets us count visits per house, even though the answer only needs the number of distinct houses. The counter costs barely nothing and keeps the code uniform across both parts.

A house is a pair of coordinates. Let us define a small `house` type inside `main` because it is used only here:

```go
type house struct {
	x, y int
}
```

### Part One

We start at the origin, mark it as visited, and walk the string. For each direction we update the current position and increment its visit counter in the map. When the walk is done, `len(visits)` is the number of distinct houses.

```go
loc := house{}
visits := make(map[house]int)
visits[loc] = 1

for _, m := range line {
	switch m {
	case '^':
		loc.y++
	case 'v':
		loc.y--
	case '>':
		loc.x++
	case '<':
		loc.x--
	}
	visits[loc]++
}

result := len(visits)
```

See [part1.go](part1.go).

### Part Two

The change is small but not trivial. Now there are two walkers, Santa and Robo-Santa, and they take turns. The first character is Santa's, the second is Robo-Santa's, and so on. They share the same grid and the same set of visited houses.

We keep two positions and select the active one by the parity of the index. Using a pointer avoids duplicating the `switch` for each walker:

```go
var loc *house
if 0 == i%2 {
	loc = &s_loc
} else {
	loc = &r_loc
}
```

The rest of the loop is the same as in Part One. The initial positions of both walkers are the origin; both must be counted, even before any move is made.

See [part2.go](part2.go).

## Complexity

Let *n* be the length of the input.

- **Time:** *O*(*n*) expected. Each move is a map insertion or lookup, which is *O*(1) on average. In the worst case a hash map degrades to *O*(*n*) per operation, but Go's map is a hash table with randomised hashing, and adversarial inputs are not a concern for Advent of Code. So in practice—and in expectation—the walk is linear in the number of moves.

- **Space:** *O*(*n*). The map holds at most *n* + 1 distinct houses in Part One and *n* + 1 in Part Two (both walkers start at the origin and can visit at most one new house each).

## Associations

- **Random walks.** Santa's walk is a two-dimensional random walk: at each step, one of four directions, chosen uniformly if the input were random. The puzzle does not ask for probabilities, but the same walk appears in the theory of random walks. See [Random walk](https://en.wikipedia.org/wiki/Random_walk).

- **Pólya's theorem.** A two-dimensional random walk returns to its starting point with probability 1, while a three-dimensional one does not—"a drunk man will find his way home, but a drunk bird may get lost forever". Santa's grid is two-dimensional, so this is the same walk. See [Pólya's theorem](https://en.wikipedia.org/wiki/P%C3%B3lya%27s_theorem).

## Notes

Two details are worth pointing out.

First, the `house` type is declared inside `main`. In Go this is legal: a type can be local to a function. I do it because the type is not used anywhere else, and keeping it close to where it is used makes the code easier to read. If another day needed a grid position without an explicitly defined grid, I would lift it into a shared package—but that day has not come yet, and I prefer not to guess at the shape of a shared type before there is a second use case.

Second, both walkers in Part Two start at the same position. Both are counted as visited, even before either has moved. The answer does not depend on this—it counts distinct houses, not visits—but the code is explicit about it, and it would matter if we ever needed visit counts.

There is no benchmark for this day either. The input is small and the algorithm is linear; a benchmark would mostly measure process startup, not the walk. Benchmarks become meaningful from the next day onward.

---

[← Previous: Day 2](day02.md) · [Up: Contents](../README.md) · [Next: Day 4 →](day04.md)
