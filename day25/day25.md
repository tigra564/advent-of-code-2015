# Day 25: Let It Snow

## The Puzzle

Santa wants to start the snow machine, but it needs a code. The code is hidden in a grid of numbers that fills diagonally. The first cell (row 1, column 1) contains 20151125, and each subsequent cell—moving down the diagonals—is computed by multiplying the previous value by 252533 and taking the remainder modulo 33554393.

## Input

A single line with two numbers: the row and column of the code. Here is my input, as an example of the format:

```
To continue, please consult the code grid in the manual.  Enter the code at row 3010, column 3019.
```

The numbers are extracted with `common.ExtractNumbers`, which pulls them out in order. There is nothing else in the input—no grid, no sequence, no parameters. The row and column are the only data.Part One asks for the code at a specific row and column, given in the input.

Part Two is not a puzzle: it asks you to collect the remaining stars from the previous days and put all fifty into the machine. There is no algorithm here, only a completion condition. Solving Part One and having the other days done is what fills the machine, and the snow begins to fall.

## The Approach

### Part One

The grid fills diagonally, and the first thing to work out is which cell comes at which position in the sequence. If we number the diagonals from zero—which is convenient here, because we ultimately want the number of steps from the start—then the diagonal of a cell `(row, col)` is `diag = (row − 1) + (col − 1)`. The first cell is `(1, 1)`, on diagonal zero; the next diagonal is `(2, 1)` and `(1, 2)`; the one after that is `(3, 1)`, `(2, 2)`, `(1, 3)`, and so on. The number of cells on all diagonals before `diag` is the triangular number `diag · (diag + 1) / 2`. Within its own diagonal, the cell is at position `col`, counting from the top-right. So the index of the cell is:

```go
func codeNum(row, col int) int {
	diag := (row - 1) + (col - 1)
	return diag * (diag+1) / 2 + col
}
```

For my input, `row = 3010` and `col = 3019`, so `diag = 6027` and `codeNum = 6027 · 6028 / 2 + 3019 = 18 168 397`. That is the number of times we need to advance the sequence from the starting value to reach the target cell.

The sequence itself is a linear congruential generator:

```go
func next(in int) int {
	return (in * 252533) % 33554393
}
```

The multiplier 252533, the modulus 33554393, and the starting value 20151125 are all given in the puzzle statement. The generator is a simple one—no addition, no bit mixing—but it is enough to produce a long sequence of distinct values before repeating.

The main loop advances the sequence the required number of times:

```go
const start = 20151125

coords := common.ExtractNumbers(lines[0])
n := codeNum(coords[0], coords[1])

result := start
for range n - 1 {
	result = next(result)
}
```

The loop runs `n − 1` times, not `n`, because the starting value is the code for cell 1, and we need `n − 1` multiplications to reach cell `n`.

See [part1.go](part1.go).

### Part Two

Part Two is not an algorithm. It asks you to collect the remaining 49 stars from the other days and add them to the one already in the machine. If you have solved the other days, this part is already done; if you have not, no amount of coding will help. It is the event's closing ceremony, dressed as a puzzle.

## Complexity

Let *n* be the index of the target cell in the diagonal ordering.

- **Time:** *O*(*n*) for Part One. The loop runs `n − 1` times, and each iteration is one multiplication and one modulo—constant work. For my input, *n* = 18 168 397, which is about 18 million iterations, and the program runs in tens of milliseconds. There is no faster way to get the *n*-th value of a linear congruential generator without using the generator's algebraic structure; the sequence is not amenable to closed-form evaluation in the general case, though for some parameters it can be jumped ahead using modular exponentiation.

- **Space:** *O*(1). A single integer holds the current value, and the coordinates are parsed once into two ints.

## Associations

- **Linear congruential generators.** The sequence in this puzzle is a linear congruential generator: `xₙ = a · xₙ₋₁ mod m`, with `a = 252533` and `m = 33554393`. LCGs are the oldest and simplest family of pseudorandom number generators, and they are still used where speed matters more than quality—in embedded systems, in simulations, and in the `rand()` of many C libraries. The parameters matter: a poor choice of `a` and `m` produces a sequence with short periods or visible patterns, and the generator fails statistical tests. The parameters here are fine for a puzzle, but a real application would need a generator with better guarantees. See [Linear congruential generator](https://en.wikipedia.org/wiki/Linear_congruential_generator).

- **The Cantor pairing function.** The diagonal numbering in this puzzle is a bijection between pairs of positive integers and positive integers: each cell `(row, col)` gets a unique index. The same idea appears in the Cantor pairing function, which maps two naturals to one, and in the more general pairing functions used in computability theory to encode tuples as single numbers. The diagonal enumeration is a classic example, and it shows up in proofs of countability—the rationals are countable because they can be arranged in a grid and traversed diagonally. See [Pairing function](https://en.wikipedia.org/wiki/Pairing_function).

- **Snowflakes and snow machines.** Real snowflakes grow as they fall, and the journey takes anywhere from 15 minutes to a few hours, depending on the height of the cloud and the conditions inside it. A snowflake that spends longer in the air grows larger, which is why the biggest crystals form in still, cold weather. Snow machines exist, but they produce snow locally—on a ski slope, not across a continent. The Snow4Ever 200, for example, makes 200 cubic metres of technical snow per day at temperatures above freezing. Santa's machine is a different order of magnitude: it does not make snow, it makes it *snow*. See [Snowflake](https://en.wikipedia.org/wiki/Snowflake).

## Notes

The loop in Part One runs about 18 million times, and each iteration is one multiplication and one modulo. That is fast enough—the program finishes in tens of milliseconds—but it is not the fastest way to compute the *n*-th value of a linear congruential generator. Since the generator has no additive constant, it is a pure multiplicative LCG, and its *n*-th value can be computed directly:

> *xₙ* = *aⁿ* · *x*₀ mod *m*

The exponentiation *aⁿ* mod *m* takes *O*(log *n*) multiplications by fast exponentiation—about 25 for our *n*, since 18 million fits in 25 bits. That is a million times fewer operations than the naive loop. I did not use it because the loop is clear, the input is small enough that it does not matter, and the fast version would require explaining modular exponentiation, which has nothing to do with the diagonal numbering that is the actual content of the puzzle. But it is the right tool if the input ever grew.

The Chinese remainder theorem is a tempting next step, because the target index 18 168 397 is composite, and the modulus 33 554 393 factors into primes. The CRT lets you reconstruct a number from its residues modulo pairwise coprime moduli, which sounds like it should help. It does not, at least not directly: the CRT solves a system of congruences for a single unknown, while we have a single recurrence and need to jump ahead in it. The two problems look similar—both involve moduli and residues—but they are not the same. The tool for jumping ahead in an LCG is exponentiation, not the CRT.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~81 ms |

Part One runs about 18 million iterations, each one a multiplication and a modulo. That is the whole cost of the program—there is no parsing to speak of, and the diagonal index is computed in a few arithmetic operations. The runtime is dominated by the loop, and it scales linearly with the index of the target cell. For a cell twice as far into the grid, the program would take twice as long.

There is no benchmark for Part Two: it is not an algorithm, and there is nothing to measure.

---

[← Previous: Day 24](day24.md) · [Up: Contents](../README.md)
