# Day 15: Science for Hungry People

## The Puzzle

The elves are baking cookies, and each ingredient contributes four properties: capacity, durability, flavor, and texture. Each property is a signed integer, and the total score of a recipe is the product of the four property sums—clamped at zero if any sum is negative. A recipe uses exactly 100 spoons of ingredients in total, distributed in whole spoons among a fixed set of ingredients.

Part One asks for the highest possible score, given that only the total number of spoons is constrained and the individual amounts are free.

Part Two adds a constraint: the recipe must contain exactly 500 calories. Calories are a fifth property, and unlike the four scoring properties, they are summed without being multiplied.

## Input

Four lines, one per ingredient:

```
Butterscotch: capacity -1, durability -2, flavor 6, texture 3, calories 8
Cinnamon: capacity 2, durability 3, flavor -2, texture -1, calories 3
```

Each line gives a name and five integers. The numbers are extracted with `common.ExtractNumbers`, which pulls out the values in order; the name is discarded. The order matters: the first four numbers are the scoring properties, and the last is calories.

The parameter that drives the algorithm is the number of ingredients: the ways to distribute 100 spoons among *k* ingredients is *C*(100 + *k* − 1, *k* − 1), which grows as a polynomial in 100 of degree *k* − 1. For *k* = 4 the count is *C*(103, 3) = 176 851, which is small enough to enumerate exhaustively.

The calorie constraint in Part Two does not change the structure of the input; it only adds a filter on the recipes that are considered valid.

## The Approach

The search space is small enough that brute force is the right tool. There are four ingredients and 100 spoons to distribute, which gives *C*(103, 3) = 176 851 compositions. For each one, the score is a constant amount of arithmetic: multiply each ingredient's four properties by its spoon count, sum per property, clamp negatives to zero, and take the product. That is about three million operations in total—less than the time it takes the process to start. **Brute force** wins on the trade-off between programmer time and CPU time, which is the principle I follow throughout this event.

**Linear programming** does not apply here. The objective is a product of four linear forms, not a linear combination, and ILP requires the objective and constraints to be linear. One could linearise the product by taking logarithms, but that gives a real-valued solution that would then need to be rounded—and rounding an integer program is a separate problem in itself. Not the right tool.

**Dynamic programming** is possible in principle but not in practice. The natural state is `(ingredient index, spoons remaining, four property sums)`, and the property sums can range from −100 to 100 in each coordinate. The state space is on the order of 200⁴ × 4 × 100, many orders of magnitude larger than the 176 851 compositions we enumerate. The DP would only pay off if the constraints on the properties were tight enough to prune most states, and here they are not.

**Branch and bound** would speed up the enumeration by pruning branches whose optimistic upper bound cannot beat the current best. For four ingredients it would help by a constant factor, but it would not change the asymptotics—the search remains over all compositions. The technique becomes genuinely necessary in Day 22, where the state space is much larger and the rules are more complex; we will look at it there.

**Local search and simulated annealing** are the tools of choice when the number of ingredients grows large enough that exhaustive enumeration becomes impossible. They find good solutions quickly, but they do not guarantee optimality—and for this puzzle we want the exact maximum. Making annealing exact would require additional machinery, which is not worth it for a space this small.

### Part One

The enumeration is a recursive partition of 100 spoons into four parts. The `iteratePartitions` helper in `day15.go` generates all compositions of a sum into a fixed number of non-negative parts:

```go
func iteratePartitions(sum, k int, accumulated []int, callback func([]int)) {
	if k == 1 {
		callback(append(accumulated, sum))
		return
	}

	for i := range sum + 1 {
		iteratePartitions(sum-i, k-1, append(accumulated, i), callback)
	}
}
```

The callback receives a slice with the current composition. The comment at the top of the file warns that the slice should be treated as immutable: the same backing array may be reused across calls, and mutating it would corrupt subsequent iterations.

For each composition, we compute the four property sums and multiply them, treating any negative sum as zero:

```go
getMaxScore := func(parts []int) {
	caloriesPos := len(ingredients[0]) - 1
	sums := make([]int, caloriesPos)
	for i, coeff := range parts {
		for j, val := range ingredients[i][:caloriesPos] {
			sums[j] += coeff * val
		}
	}

	product := 1
	for _, sum := range sums {
		if sum < 0 {
			product = 0
			break
		}
		product *= sum
	}

	result = max(result, product)
}
```

The `break` on a negative sum skips the rest of the multiplication: if any property is negative, the whole score is zero by the puzzle's rule, and there is no point computing the rest.

See [part1.go](part1.go).

### Part Two

The only change is the calorie constraint. We compute the calorie total alongside the four property sums, and reject any composition whose calories are not exactly 500:

```go
if calories != nCalories {
	return
}
```

This is a filter, not a different algorithm. The enumeration is the same, and the only extra work per composition is one more multiply-add per ingredient. Since the enumeration itself dominates, the runtime is essentially unchanged.

See [part2.go](part2.go).

## Complexity

Let *k* be the number of ingredients and *S* the total number of spoons.

- **Time:** *O*(*C*(*S* + *k* − 1, *k* − 1) · *k* · *p*) ≈ *O*((*S* + 1)ᵏ⁻¹ · *k* · *p*), where *p* is the number of properties (four scoring properties, plus calories in Part Two). The recursion enumerates every composition of *S* into *k* parts, and for each leaf—each complete composition—the callback walks over all *k* ingredients and all *p* properties. The interior nodes of the recursion tree add a lower-order term, since their total count is bounded by the number of leaves times *k*. For *k* = 4 and *S* = 100 this is 176 851 × 4 × 4 ≈ 2.8 million operations in Part One, and 176 851 × 4 × 5 ≈ 3.5 million in Part Two.

- **Space:** *O*(*k*) for the recursion stack and the accumulated composition. The `sums` slice in the callback is *O*(*p*). The ingredient data is *O*(*k* · *p*), read once and kept for the whole run.

## Associations

- **The diet problem.** The puzzle is a constrained optimisation over a linear combination of ingredients, which is the same shape as the classical diet problem: choose quantities of foods to meet nutritional requirements at minimum cost. The diet problem was one of the first applications of linear programming, studied by George Stigler in 1945 and solved by George Dantzig in 1947 using the simplex method. Our puzzle differs in one crucial way: the objective is a product, not a sum, so the linear machinery does not apply. See [Diet problem](https://en.wikipedia.org/wiki/Stigler_diet).

- **Compositions and the stars and bars argument.** The number of ways to distribute *S* identical spoons among *k* distinct ingredients is the number of compositions of *S* into *k* non-negative parts, which by the stars and bars argument equals *C*(*S* + *k* − 1, *k* − 1). The same formula counts the number of multisets of size *S* drawn from *k* types, and it appears whenever identical objects are placed into distinct boxes. See [Stars and bars](https://en.wikipedia.org/wiki/Stars_and_bars_(combinatorics)).

## Notes

The measurements show Part Two running slightly faster than Part One, which is counter-intuitive: Part Two does the same enumeration plus the extra work of computing calories. The explanation is the early return in the callback: if the calorie total is not exactly 500, the function exits before computing the four property sums and their product. Most compositions are rejected this way, so the average work per iteration is lower in Part Two than in Part One, and the total time is slightly less despite the extra calorie arithmetic. The effect is small—both parts are dominated by process startup and input parsing—but the direction is consistent across runs.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~17 ms |
| 2 | ~14 ms |

Part Two is slightly faster, for the reason discussed in `Notes`: the early return on the calorie check skips most of the work per composition. The difference is small and close to run-to-run variation, but it is consistent.

---

[← Previous: Day 14](day14.md) · [Up: Contents](../README.md) · [Next: Day 16 →](day16.md)
