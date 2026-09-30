# Day 6: Probably a Fire Hazard

## The Puzzle

Santa has a grid of a million lights (1000 × 1000) and a list of instructions. Each instruction tells him to turn some rectangular region on, off, or toggle it. The regions are given as two opposite corners.

Part One treats the lights as binary: on or off. Part Two changes the semantics: each light has a brightness level, and the three actions mean something different.

## Input

About 300 instructions, one per line. Each line has one of three forms:

```
turn on 0,0 through 999,999
turn off 499,499 through 500,500
toggle 0,0 through 999,0
```

Each instruction gives two opposite corners of a rectangle. The puzzle does not specify which corner comes first, nor which of the two numbers in a pair is *x* and which is *y*, nor which direction is up. None of this matters: the grid is square, and the question is how many lights are on (or how bright they are), not where they are or how they are oriented.

Because the choice is free, I picked the convention that makes the code simplest. I read the instruction as `(left, bottom)` through `(right, top)` and iterate from the first coordinate to the second in each pair. The real input happens to be ordered this way. If it had not been, I would have needed to sort the pair first—but the puzzle does not require this, and I did not add the sort.

Both corners are included in the region.

## The Approach

Both parts are the same simulation over a grid, with a different operation for each action. The only structural decision is how to share the code between the two parts.

I put the common simulation into `day06.go` and pass the per-light operation as a callback. The callback receives a pointer to the light and the action name, and mutates the light in place. Part One and Part Two differ only in what that callback does.

The grid is a fixed `[1000][1000]uint16`. The type `uint16` is chosen for Part Two: a light can be toggled many times, and with about 300 instructions, the maximum brightness is well above 255. A `byte` would overflow. For Part One the values are only 0 and 1, but I use the same type in both parts and let the callback decide what the values mean.

### Part One

The callback treats the light as a boolean:

```go
func(action string, point *uint16) {
	switch action {
	case "toggle":
		*point ^= 1
	case "turn off":
		*point = 0
	case "turn on":
		*point = 1
	}
}
```

`toggle` uses XOR, which flips between 0 and 1 regardless of the current value. `turn off` and `turn on` set the value directly.

See [part1.go](part1.go).

### Part Two

The same callback, but the semantics change: a light now has a brightness, and the three actions mean "add 2", "subtract 1 but not below 0", and "add 1". The switch is the only difference:

```go
func(action string, point *uint16) {
	switch action {
	case "toggle":
		*point += 2
	case "turn off":
		if *point > 0 {
			*point--
		}
	case "turn on":
		*point++
	}
}
```

The `if *point > 0` guard prevents underflow. Without it, `turn off` on a dark light would wrap around to a huge number, and the answer would be nonsense. Go does not panic on unsigned underflow; it wraps silently.

See [part2.go](part2.go).

### Parsing

The instruction is parsed with a single regex:

```go
re := regexp.MustCompile(`^(.*) ([0-9]+),([0-9]+) through ([0-9]+),([0-9]+)$`)
```

The first group is the action, which can be `turn on`, `turn off`, or `toggle`. The other four are the coordinates.

The two coordinate pairs are separated by the word `through`, and the numbers inside each pair are separated by commas. This structure is what makes the regex unambiguous: the separator between the action and the first number is a space, and `turn on` itself contains a space, so the regex engine simply takes the last space before the first digit as the boundary. The commas and the word `through` pin down the rest.

## Complexity

Let *n* be the number of instructions and *A* the area of a rectangle. In the worst case *A* = 10⁶, the whole grid.

- **Time:** *O*(*n* · *A*) worst case, *O*(*n* · *M*(*a*)) on average, where *M*(*a*) is the mean rectangle area. If each rectangle is produced by taking two independent coordinates on each axis and sorting them into min and max, the expected width is one third of the axis length: for two independent uniform points *a* and *b* on an interval of length *L*, the expectation of |*a* − *b*| is *L*/3. The same holds for the height, so *M*(*a*) ≈ *L*² / 9 ≈ 1.1 × 10⁵. For the real input, the rectangles are small enough that the whole simulation runs in tens of milliseconds. The final count over the grid is a separate *O*(10⁶) pass, which is negligible compared to the simulation.

- **Space:** *O*(1) beyond the grid itself. The grid is a fixed `[1000][1000]uint16`, which is 2 MB. This is the only large allocation in the program.

## Associations

- **The symmetry of the square.** The puzzle fixes the grid size but not the orientation of the coordinates. It does not matter which axis we call *x* and which one we call *y*, nor which direction is up. The answer depends only on the rectangles, not on how they sit in the plane. The symmetry group of the square has eight elements—four rotations and four reflections—and the puzzle is invariant under all of them. This is a small example of a general principle: when the question is symmetric, the coordinate system can be chosen for convenience. See [Symmetry group of the square](https://en.wikipedia.org/wiki/Dihedral_group).

- **Range update and point query.** The straightforward solution touches every cell in every rectangle. A 2D Fenwick tree (or a 2D segment tree with lazy propagation) can apply each range update in *O*(log² *N*) and query a point in the same time. For Part One, where the operations are linear (XOR, assignment), this would work. For Part Two, the `turn off` operation clamps at zero—a non-linear operation that a Fenwick tree cannot represent without knowing the current value of each cell. So the second part resists the clever solution and rewards the simple one. See [Fenwick tree](https://en.wikipedia.org/wiki/Fenwick_tree).

## Notes

The callback design deserves a comment. Part One and Part Two differ only in what happens to a single light. The simulation itself—parse the instruction, walk the rectangle, apply the operation, count at the end—is identical. I could have duplicated it in each part; the result would have been a few dozen lines of near-identical code, differing in one `switch`.

The callback is not just a way to avoid duplication. It is the right shape for this problem because the varying piece is genuinely small and self-contained: one function that receives a light and an action, and mutates the light. Everything else is fixed. When the variation is this narrow, a callback expresses it directly, and the rest of the code stays unchanged.

The rule of thumb: use a callback when the varying piece is a single, well-defined operation, and the surrounding logic is identical. Do not use it when the varying piece is large, or when it needs to affect the control flow of the caller. In those cases, duplication or a more explicit abstraction—an interface, a generic function, a strategy object—will be clearer.

The Go compiler often inlines small callbacks, so the cost of the extra indirection is usually zero. Even when it does not, the call is cheap compared to the array access inside the loop.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~72 ms |
| 2 | ~76 ms |

Both parts do the same amount of work: they touch the same lights in the same order. The difference is in the callback, which is a few arithmetic operations per light. The 4 ms gap is mostly noise—two separate binaries, two separate runs, and the callbacks differ by a branch or two. For practical purposes the two parts are the same speed.

---

[← Previous: Day 5](day05.md) · [Up: Contents](../README.md) · [Next: Day 7 →](day07.md)
