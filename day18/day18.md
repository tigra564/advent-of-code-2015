# Day 18: Like a GIF For Your Yard

## The Puzzle

You have a grid of lights, each either on (`#`) or off (`.`), and a set of rules from Conway's Game of Life: a lit cell stays lit if it has two or three lit neighbours, an unlit cell becomes lit if it has exactly three, and everything else turns off or stays off.

Part One asks how many lights are on after 100 steps of this rule, starting from the given configuration.

Part Two adds a twist: the four corner lights are stuck on and can never be turned off. They count as lit for their neighbours' counts, and they remain lit at every step, regardless of the rule. Part Two asks for the number of lit lights after 100 steps under this constraint.

## Input

A 100 × 100 grid of `.` and `#` characters, one row per line. Each character represents one light: `.` is off, `#` is on.

The grid is square and its size is fixed at 100 × 100 by the puzzle statement, so there is no need to parse dimensions from the input. The starting configuration is arbitrary—the puzzle does not guarantee any particular density or pattern—but for the real input it is sparse enough that the number of lit lights stays in the hundreds throughout the simulation.

## The Approach

Both parts simulate the same rule for 100 steps. The only difference is the corner constraint in Part Two, and that difference is small enough that the same frame-advance function serves both.

The grid is stored as a slice of strings, one string per row. This representation is convenient for two reasons: a row is a contiguous chunk of memory, so scanning it is cache-friendly, and the border trick works naturally—each row can be padded with a `.` on the left and right, and an extra row of dots can be added at the top and bottom.

The border is the first design decision. Instead of bounds-checking every neighbour access, `getGrid` pads the input with a one-cell frame of `.`:

```go
width := len(lines[0])
border := strings.Repeat(".", width+2)

newLines := make([]string, 0, len(lines)+2)
newLines = append(newLines, border)
for _, line := range lines {
	newLines = append(newLines, "."+line+".")
}
newLines = append(newLines, border)
```

With the border in place, the neighbour-counting loop is uniform: every interior cell has eight neighbours in all directions, and no index ever goes out of range. The cost is a grid that is two rows and two columns larger than the input, and a border that must be ignored when counting lit lights. Both are trivial.

The eight neighbour offsets are stored in a package-level slice:

```go
var dirs = [][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{ 0, -1},          { 0, 1},
	{ 1, -1}, { 1, 0}, { 1, 1},
}
```

The order of the offsets does not matter—they are all visited for every cell—but writing them in a grid pattern makes it easy to verify that all eight are present and that none is duplicated.

`nextFrame` builds the next generation one row at a time. For each interior cell, it counts the lit neighbours and applies the Game of Life rule:

```go
switch in[i][j] {
case '#':
	if 2 == neighborsOn || 3 == neighborsOn {
		cell = '#'
	}
case '.':
	if 3 == neighborsOn {
		cell = '#'
	}
}
```

A lit cell survives on two or three neighbours; an unlit cell is born on exactly three. Everything else defaults to `.`, which is why the switch does not have a `default` branch—the variable `cell` is initialised to `.` before the switch, and only the surviving and newborn cases overwrite it.

The new frame is built with a `strings.Builder`, preallocated to the row width. The builder writes the left border, then the computed cells, then the right border. The top and bottom border rows are copied verbatim from the input—they are all dots and never change.

### Part One

Part One runs the simulation for 100 steps and counts the lit lights:

```go
grid := getGrid()

const nIters = 100
for range nIters {
	grid = nextFrame(grid)
}

result := countLit(grid)
```

`countLit` sums the `#` characters in every row with `strings.Count`. The border rows and columns are all `.`, so they contribute nothing and require no special handling.

See [part1.go](part1.go).

### Part Two

Part Two needs the four corners to stay lit. There are two ways to enforce this: after every frame, overwrite the corners with `#`, or modify the rule so that corner cells always produce `#`. The first is simpler and does not touch `nextFrame`, so the frame-advance function stays the same in both parts.

`setStuckLights` writes `#` into the four corner cells of the interior grid—the border cells are left alone, because they are not part of the puzzle:

```go
func setStuckLights(grid []string) {
	right := len(grid[0]) - 2
	bottom := len(grid) - 2

	for _, i := range []int{1, bottom} {
		b := []byte(grid[i])
		b[1] = '#'
		b[right] = '#'
		grid[i] = string(b)
	}
}
```

The strings are immutable in Go, so changing a single character means converting the row to a `[]byte`, mutating it, and converting back. This is a small inefficiency, but the corners are only four cells and the conversion happens twice per frame, so it does not matter.

The main loop calls `setStuckLights` both before the simulation and after every step:

```go
grid := getGrid()
setStuckLights(grid)

const nIters = 100
for range nIters {
	grid = nextFrame(grid)
	setStuckLights(grid)
}
```

The call before the loop is necessary because `nextFrame` reads the corners of the current grid when computing the next one—if the corners are not lit at the start, the first frame will be computed as if they were off.

See [part2.go](part2.go).

## Complexity

Let *W* and *H* be the width and height of the grid, and *T* the number of steps.

- **Time:** *O*(*T* · *W* · *H*) for both parts. Each step visits every interior cell once and, for each cell, loops over the eight neighbour offsets. The constant factor is small—eight array accesses and a few comparisons per cell—which is why 100 steps over a 100 × 100 grid run in a fraction of a second. Part Two adds a constant amount of work per step for `setStuckLights`, which is negligible.

- **Space:** *O*(*W* · *H*) for the two grids—the current frame and the next one. Both are stored as slices of strings, and each string is a contiguous block of bytes. The `strings.Builder` in `nextFrame` allocates one string per row per step, which is *O*(*W*) per row and *O*(*W* · *H*) per frame. In principle the next frame could overwrite the current one in place, but the rule reads neighbours from the current frame while writing the next, so a second buffer is needed.

## Associations

- **Conway's Game of Life.** The rules in this puzzle are exactly the rules of Conway's Game of Life, published by John Horton Conway in 1970 in *Scientific American*. It is the most famous cellular automaton: a grid of cells, each with two states, updated synchronously according to the number of live neighbours. Despite the simplicity of the rule, the Game of Life is Turing-complete—it can simulate any computation—and it exhibits a rich zoo of stable, oscillating, and travelling patterns. See [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life).

- **Synchronous versus asynchronous update.** A key feature of this puzzle, and of the Game of Life, is that all cells update at the same time: the new state of a cell depends on the old states of its neighbours, not on their new states. This is why a second buffer is needed—writing into the same grid while reading from it would mix generations. The same distinction appears in other cellular automata, in synchronous digital circuits, and in double-buffered graphics, where the frame being drawn is not the frame being displayed. See [Cellular automaton](https://en.wikipedia.org/wiki/Cellular_automaton).

## Notes

The simulation allocates a fresh slice of strings on every step. Each `nextFrame` call creates a `strings.Builder` per row, builds the row, and appends it to a new slice. For 100 steps over a 100 × 100 grid, that is 10 000 row allocations—small enough to be invisible at this size, but not free. A double-buffered version would allocate two grids once and alternate between them, writing into the second while reading from the first. The rule reads neighbours from the current frame and writes to the next, so two buffers are genuinely needed—in-place update would corrupt the neighbour counts. But the buffers themselves could be reused across steps, and the `strings.Builder` could be replaced with direct writes into a preallocated `[]byte` per row. I did not do this because the puzzle is small and the code is clearer without the buffer management, but it is the natural optimisation if the grid or the step count ever grew.

The choice to keep `nextFrame` unchanged in Part Two, rather than adding a flag or a special case for the corners, is deliberate. The corner rule is a property of the puzzle, not of the Game of Life, and it belongs in the caller. `setStuckLights` is a small function that does exactly one thing, and it is called at the two points where it is needed: once before the loop, so the initial frame has lit corners, and once after every step, so the corners stay lit. The alternative—passing a `stuckCorners bool` into `nextFrame` and branching inside the cell loop—would put the rule in the hot path and complicate the function for the sake of avoiding two function calls per step.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~24 ms |
| 2 | ~24 ms |

Both parts are dominated by the simulation: 100 steps over a 100 × 100 grid, with eight neighbour lookups per cell. Part Two does the same work plus four corner writes per step, which is negligible. The two parts are effectively the same speed.

---

[← Previous: Day 17](day17.md) · [Up: Contents](../README.md) · [Next: Day 19 →](day19.md)
