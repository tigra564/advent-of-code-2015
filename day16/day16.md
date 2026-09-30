# Day 16: Aunt Sue

## The Puzzle

You receive a gift from an unknown aunt, and the card is unsigned. The MFCSAM machine gives you a reading of the sender's "properties" from the card itself. Separately, you have hazy recollections of the same properties for hundreds of candidate aunts—each remembered only partially, a few values out of the full set.

Part One asks for the number of the aunt whose remembered properties all match the MFCSAM reading exactly.

Part Two changes the matching rules for a few properties: `cats` and `trees` are greater-than comparisons, `pomeranians` and `goldfish` are less-than, and the rest remain equality. Part Two asks for the aunt who matches under these rules.

## Input

The input has two parts. The first is a list of hundreds of lines, one per candidate aunt:

```
Sue 1: children: 1, cars: 8, vizslas: 7
Sue 2: akitas: 10, perfumes: 10, children: 5
```

Each line gives the aunt's number and a few of her properties, separated by commas. Not all properties are present in every line—only a subset, and which subset varies. The property names are lowercase words, and the values are non-negative integers.

The second part is the MFCSAM reading, which in the code is a small fixed map of ten properties and their values. It is not read from the input file: it comes from the puzzle statement, and we embed it directly in `day16.go` as a `map[string]int` with a comment pointing to its origin.

The asymmetry between the two parts matters for the algorithm. The candidate aunts are partial—each has only a few known properties—while the MFCSAM reading is complete. The matching rule is therefore "for every property the candidate remembers, its value must match the MFCSAM value", not the other way around.

## The Approach

Both parts scan the candidate list and return the first aunt whose remembered properties are consistent with the MFCSAM reading. The puzzle guarantees that exactly one candidate matches, so the first match is also the only one—the loop can return early without checking the rest. The only difference between the parts is the rule that decides consistency, so the scan is factored out into `findSue`, which takes a validation function as a parameter:

```go
func findSue(validate func(traitsMFCSAM, traitsSue map[string]int) bool) int {
	lines := common.ReadInputLines()
	for _, line := range lines {
		numSue, traitsSue := readSue(line)
		if validate(traitsMFCSAM, traitsSue) {
			return numSue
		}
	}

	return 0
}
```

The `validate` callback receives two maps: the full MFCSAM reading and the partial properties of the current candidate. It returns `true` if the candidate is consistent with the reading under its rule. The scan itself does not know what "consistent" means—it just iterates until a candidate passes, and returns the aunt's number.

The MFCSAM reading is a package-level variable, not a parameter threaded through every call. It is fixed data from the puzzle statement, and treating it as global keeps the call sites clean:

```go
var traitsMFCSAM = map[string]int{
	"children":    3,
	"cats":        7,
	"samoyeds":    2,
	"pomeranians": 3,
	"akitas":      0,
	"vizslas":     0,
	"goldfish":    5,
	"trees":       3,
	"cars":        2,
	"perfumes":    1,
}
```

Parsing a line is two regexes and a loop. The first regex extracts the aunt's number, the second pulls out the property–value pairs:

```go
var (
	reSue = regexp.MustCompile(`Sue ([0-9]+): `)
	reMap = regexp.MustCompile(`([a-z]+): ([0-9]+)`)
)
```

`readSue` splits the line at the number and hands the rest to `readMapStrInt`, which fills a `map[string]int`. The map is small—three entries per line—and the allocation is not worth avoiding.

### Part One

Part One requires exact equality on every remembered property. The validation function iterates over the candidate's map and compares each value to the MFCSAM value for the same key:

```go
result := findSue(func(traitsMFCSAM, traitsSue map[string]int) bool {
	for k, v := range traitsSue {
		if v != traitsMFCSAM[k] {
			return false
		}
	}
	return true
})
```

The iteration is over `traitsSue`, not over `traitsMFCSAM`, and this is the key detail. Since the candidate is partial, iterating over the full MFCSAM reading would check properties the candidate never mentioned and reject her for missing values. Iterating over the candidate's own map checks only what she actually remembers.

See [part1.go](part1.go).

### Part Two

Part Two changes the comparison for four properties: `cats` and `trees` must be strictly greater than the MFCSAM value, `pomeranians` and `goldfish` must be strictly less. The remaining properties keep the equality rule:

```go
result := findSue(func(traitsMFCSAM, traitsSue map[string]int) bool {
	for k, v := range traitsSue {
		switch k {
		case "cats", "trees":
			if v <= traitsMFCSAM[k] {
				return false
			}
		case "pomeranians", "goldfish":
			if v >= traitsMFCSAM[k] {
				return false
			}
		default:
			if v != traitsMFCSAM[k] {
				return false
			}
		}
	}
	return true
})
```

The switch is the whole difference between the two parts. The scan, the parsing, and the MFCSAM data are identical.

See [part2.go](part2.go).

## Complexity

Let *n* be the number of candidate aunts and *p* the number of properties per candidate.

- **Time:** *O*(*n* · *p*) for both parts. Each line is parsed with two regexes, which take *O*(*p*) time, and then the validation callback iterates over the candidate's properties, which is also *O*(*p*). The MFCSAM map is looked up once per property, and each lookup is *O*(1) on average. The number of properties a candidate can have is bounded by the size of the MFCSAM reading, which is a constant, so *p* is effectively constant as well—but the algorithm does not depend on that, and would scale linearly if the reading ever grew.

- **Space:** *O*(*p*) per line for the parsed map, plus *O*(*P*) for the global MFCSAM map. The input is read line by line by `common.ReadInputLines`, which materialises the whole file as a slice of strings—that is *O*(*n* · *L*) in the input, but it is a property of the reader, not of this algorithm.

## Associations

- **Holmes and the method of elimination.** The puzzle follows a pattern Sherlock Holmes described as his method: "when you have eliminated the impossible, whatever remains, however improbable, must be the truth" (Chapter VI, *The Sign of the Four*). The MFCSAM reading is the evidence; each candidate aunt is a hypothesis; and matching every remembered property against the reading eliminates the hypotheses that contradict it. The puzzle's input is built so that exactly one aunt survives all the filters. The same principle appears in a Russian joke about a young lecturer from the Department of General Surgery who reconstructs an entire evening from a student's damaged gradebook, eliminating water, wine, and vodka as causes and concluding that it must have been pure alcohol—a substance he happens to know something about. The joke is too long to reproduce here, and it is in Russian, but machine translation handles it well enough: [История №40566](https://www.anekdot.ru/id/40566/).

- **Partial records and indexed lookup.** Each candidate aunt is a sparse record: three properties out of ten, with different properties remembered for different aunts. In a relational database, this would be a table with ten nullable columns, one per property, and the search would be a `SELECT` with clauses like `(children = 3 OR children IS NULL)` for every property in the reading. The `OR ... IS NULL` is what makes the query work on partial records: a candidate who did not mention a property is not required to match it. Indexes on the property columns would let the query planner prune candidates without scanning the whole table—the same filtering we do by hand, but expressed declaratively. For 500 rows, the database would almost certainly be slower: connection setup, query parsing, and planning would dominate the actual lookup. But the structure is identical, and the association is worth noting because it shows the same problem from a different angle.

## Notes

The second part's comparison rules are not symmetric: `cats` and `trees` are "more than the reading", `pomeranians` and `goldfish` are "fewer than", and everything else is "equal to". If there were more rules, or if the rules were more complex but still one per property, the natural refactoring would be to store comparators in the map instead of plain integers—each property pointing to a function that takes the remembered value and the MFCSAM value and returns a boolean. This is the same shape as Day 5 and Day 11, where each rule became its own small predicate and the puzzle reduced to combining them. For ten properties and two puzzles, the `switch` is shorter and easier to read, and a table of functions would be a solution in search of a problem.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~3 ms |
| 2 | ~2 ms |

Both parts are dominated by process startup and input parsing. The scan itself is a few hundred lines, each parsed with two regexes and checked with a handful of map lookups—microseconds of work that disappear next to the cost of loading the binary and reading the file. The difference between the two parts is within run-to-run variation.

---

[← Previous: Day 15](day15.md) · [Up: Contents](../README.md) · [Next: Day 17 →](day17.md)
