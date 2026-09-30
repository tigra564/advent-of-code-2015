# Day 13: Knights of the Dinner Table

## The Puzzle

You are hosting a family dinner and need to seat the guests around a table. The input gives, for every ordered pair of people, how many happiness units the first person gains or loses by sitting next to the second. Happiness is not symmetric: Alice may gain 54 by sitting next to Bob, while Bob loses 79 by sitting next to Alice.

Part One asks for the seating arrangement that maximizes the total happiness. Every person sits at a round table, so each person has exactly two neighbours, and the total is the sum of all pairwise contributions along the circle.

Part Two adds you to the table. You have no opinion about anyone and nobody has an opinion about you, so all edges involving you contribute zero. But your presence changes the seating: the table now has one more chair, and the circle has one more person on it.

## Input

About a dozen lines, one per ordered pair of people:

```
Carol would lose 62 happiness units by sitting next to Alice.
Carol would gain 60 happiness units by sitting next to Bob.
```

Each line gives a directed edge: the person who gains or loses, the sign, the number of units, and the neighbour. The happiness is not symmetric—the graph is directed, and the edge from A to B is separate from the edge from B to A.

The input describes a relation for every ordered pair: each pair has both directions specified, one for each of the two people. This means any seating arrangement is valid and can be scored, and there is no need to check whether a particular pair of neighbours is allowed.

## The Approach

This puzzle is a disguise of Day 9. There we had a set of locations and distances, and we looked for the shortest and longest routes through all of them. Here we have a set of people and happiness values, and we look for the seating that maximizes the total. The two are the same problem with different names.

The reduction is straightforward. In Day 9, an edge was undirected: the distance from A to B was the same as from B to A, and `MakeEdge` normalised the pair. Here, the graph is directed, but every pair has two contributions—one in each direction—so we can collapse them into a single undirected edge whose weight is the sum of the two:

```go
hapUnits[common.MakeEdge(from, to)] += weight
```

This `+=` is the whole trick. When the input gives us `Alice would gain 54 by sitting next to Bob`, we add 54 to the edge `{Alice, Bob}`. When it gives us `Bob would lose 79 by sitting next to Alice`, we add −79 to the same edge. The final weight of the edge is the total happiness that Alice and Bob contribute when they sit next to each other—regardless of which side each of them is on, since a seat at a round table is symmetric.

The symmetric seating means the directed graph collapses into an undirected one, and the puzzle becomes exactly Day 9—except that we want the maximum, not the minimum, and the path is closed into a cycle. The `FindBestHamiltonian` function from `common.go` already handles both: the `cycle` flag closes the path, and the `cmp` callback chooses between minimisation and maximisation. The only change from Day 9 is passing `true` for `cycle` and reversing the comparison.

### Part One

We build the graph, call `FindBestHamiltonian` with `cycle` set to `true` and the comparison reversed:

```go
result, _ := common.FindBestHamiltonian(people, hapUnits, true, math.MinInt, func(a, b int) bool {
	return a > b
})
```

The initial value is `math.MinInt`, so any real seating beats it, and the comparison `a > b` keeps the maximum.

See [part1.go](part1.go).

### Part Two

The change is one line: we add ourselves to the list of people, and rely on the map returning zero for edges that do not exist. No explicit initialisation is needed—the missing edges simply contribute nothing.

```go
people = append(people, "I")

result, _ := common.FindBestHamiltonian(people, hapUnits, true, math.MinInt, func(a, b int) bool {
	return a > b
})
```

Since the number of people increases by one, the number of circular arrangements grows accordingly, and the search takes longer. But the algorithm itself is unchanged.

See [part2.go](part2.go).

## Complexity

Let *n* be the number of people.

- **Time:** *O*(*n*! · *n*). We enumerate all circular arrangements by permuting the list and summing the weights along the cycle. The factorial dominates; the extra *n* factor is the sum over the edges. For *n* = 8 this is small enough to run in milliseconds; for Part Two the number of people grows by one, and the work grows by roughly a factor of *n* + 1.

- **Space:** *O*(*n*) for the current permutation and the best arrangement. The recursion in `Permute16` uses *O*(*n*) stack frames, and the accumulator is a fixed-size array. As usual, we do not count the input itself.

## Associations

- **Speaking names.** The names in the puzzle are not arbitrary. Alice and Bob are the canonical placeholder characters in cryptography, and Mallory is the active attacker—the one who does not just eavesdrop but interferes. In my input, Mallory gives away less happiness than anyone else except Carol and Eric, which is a fitting temperament for someone who spends her time tampering with other people's messages. Whether the puzzle author chose the names on purpose or by habit, the resonance is there. See [Alice and Bob](https://en.wikipedia.org/wiki/Alice_and_Bob).

- **The marriage theorem and seating.** A related classical problem asks whether a set of people can be seated so that no one sits next to an enemy. This is a matching problem, solvable in polynomial time by Hall's marriage theorem. Our puzzle is different: every seating is allowed, and the question is which one is best, not whether one exists. The contrast is worth noting—the feasibility version is easy, the optimisation version is hard. See [Hall's marriage theorem](https://en.wikipedia.org/wiki/Hall%27s_marriage_theorem).

## Notes

The reduction to Day 9 is the part worth remembering. A seating problem with directed happiness values sounds different from a routing problem with undirected distances, but the `+=` on a symmetric edge key collapses the two directions into one weight, and the round table turns the path into a cycle. Once that is done, the same `FindBestHamiltonian` handles both puzzles, with only the `cycle` flag and the comparison callback changing.

The same function handles incomplete input as well, with a small trick. If the data for some pair is missing, we can read that as "these two must never sit together" and fill in the missing edges with a large negative value. If a valid seating still exists—some arrangement that avoids all the forbidden adjacencies—then the maximisation will find it, because any seating that uses a forbidden pair will be pushed below every valid one. The value should be large enough to dominate any real sum but small enough to avoid overflow when several such edges appear in the same seating: `math.MinInt / 2` is a reasonable choice. The result is the best allowed arrangement, not just the best unrestricted one.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~30 ms |
| 2 | ~254 ms |

Part Two is about eight times slower, and the reason is purely factorial. Eight people give 8! = 40 320 permutations; adding one more gives 9! = 362 880, roughly nine times as many. The measured ratio is close to that, which confirms what the `Complexity` section already says: the enumeration is the whole cost, and everything else—parsing, map lookups, summing—is lost in the noise.

---

[← Previous: Day 12](day12.md) · [Up: Contents](../README.md) · [Next: Day 14 →](day14.md)
