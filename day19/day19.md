# Day 19: Medicine for Rudolph

## The Puzzle

Rudolph has a rare illness, and you need to synthesise the cure: a specific molecule described by the input. The rules for building it are a set of replacements—each one says that a single atom or a pair of atoms can be replaced by a longer sequence—and a target molecule.

Part One asks how many distinct molecules can be produced by applying exactly one replacement rule, anywhere in the target, exactly once.

Part Two asks the opposite: how many replacements are needed to build the target molecule starting from a single electron `e`, if each rule is applied in reverse—that is, replacing a substring of the target with the shorter left-hand side—until only `e` remains. Part Two asks for the minimum number of such reverse replacements.

## Input

The input has two parts. The first is a list of replacement rules, one per line:

```
H => HO
H => OH
O => HH
```

Each rule has a left-hand side (the substring to be replaced) and a right-hand side (what it becomes), separated by ` => `. The left-hand side is always exactly one atom; the right-hand side is a longer sequence of atoms.

The second part is a single line: the target molecule, separated from the rules by a blank line. It is a long string of atoms, where each atom is either a single uppercase letter or an uppercase letter followed by a lowercase one—`H`, `O`, `Ca`, `Mg`, and so on. The distinction matters: a naive character-by-character scan would split `Ca` into `C` and `a`, which are not atoms.

The number of rules is in the dozens, and the target molecule is a few hundred atoms long. The same left-hand side may appear in several rules with different right-hand sides, but the right-hand sides are distinct, which Part Two relies on. For Part One, the answer counts *distinct* molecules, not distinct rules applied, so two different rules that happen to produce the same string would count once; in this input, that does not happen.

For Part Two, the rules are applied in reverse, and the puzzle guarantees that a solution exists—the target can be reduced to `e`. The optimal solution may not be unique, but the question asks for the number of steps, and that number is well-defined.

## The Approach

This is the first day where the puzzles get genuinely hard. Part One is a straightforward search: apply every rule at every possible position, collect the results in a set, and count the set. Part Two is the reverse problem, and I did not find a rigorous algorithm that guarantees the minimum in the general case. What I have is a heuristic that works on the given input—and the reason a general solution is hard is a topic for `Notes`.

### Part One

The search is over positions. For each rule, we scan the target molecule for occurrences of the rule's left-hand side, and for each occurrence, build the new molecule by replacing that one occurrence with the right-hand side:

```go
newMolecules := make(map[string]common.Empty)
for _, prod := range prods {
	offset := 0
	for {
		idx := strings.Index(molecule[offset:], prod.from)
		if idx == -1 {
			break
		}
		offset += idx
		nm := molecule[:offset] + prod.to + molecule[offset + len(prod.from):]
		offset++
		newMolecules[nm] = common.Empty{}
	}
}
```

The result is a set, not a list, because two different rules can produce the same molecule, and the puzzle asks for distinct molecules. The set is a `map[string]common.Empty`—the value type carries no information, only the key matters. This is the same idiom used elsewhere in the event for sets.

The `offset++` after each match is what allows overlapping occurrences: if the left-hand side is `H` and the molecule contains `HHH`, we want to find all three positions, not just the first. Incrementing by 1 rather than by `len(prod.from)` is the simple way to do this.

See [part1.go](part1.go).

### Part Two

Part Two asks for the minimum number of reverse replacements from the target molecule down to `e`. The search space is enormous—every reverse replacement step produces a set of candidates, and a breadth-first search would explode immediately. What I use instead is a greedy heuristic: sort the rules so that the "most useful" ones come first, and apply them in that order until the molecule is reduced to `e`.

The sorting key is computed by `fillProductions`, which counts two things for each rule: how many atoms on the right-hand side are themselves left-hand sides of other rules (`numNonTerms`), and how many atoms the right-hand side contains in total (`numAtoms`):

```go
matches := reAtom.FindAllString(prod.to, -1)
prod.numAtoms = len(matches)

prod.numNonTerms = 0
for _, m := range matches {
	if _, ok := nonTerms[m]; ok {
		prod.numNonTerms++
	}
}
```

The rules are then sorted by `numNonTerms` descending, then by `numAtoms` descending, then by the right-hand side string descending. The intuition is that rules whose output contains many non-terminals are the ones that build up the complex structure of the target, and they should be applied first; rules that produce only atoms are the last step, reducing the molecule to `e`.

The main loop applies the first rule that matches, then restarts the scan:

```go
changed := true
result := 0
for changed {
	changed = false
	for _, prod := range prods {
		numSubst := strings.Count(molecule, prod.to)
		if 0 == numSubst {
			continue
		}
		molecule = strings.ReplaceAll(molecule, prod.to, prod.from)
		result += numSubst
		changed = true
		break
	}
}
```

Two details matter here. First, `strings.ReplaceAll` replaces *all* occurrences of the right-hand side in one call, not just the first. This is a deliberate choice: if the rule matches several places, applying it everywhere at once is a valid sequence of reverse replacements, and it reduces the molecule faster. Second, after each successful application, the loop restarts from the first rule in the sorted list. This is because applying a rule changes the molecule, and a rule that did not match before may match now.

The loop terminates when no rule matches. At that point, the molecule should be exactly `"e"`. If it is not, the algorithm has failed, and the code panics:

```go
if "e" != molecule {
	panic(`Algorithm ended with "` + molecule + `", not with just "e"`)
}
```

The panic is there because the heuristic is not guaranteed to succeed. If the greedy order ever leads the molecule into a state where no rule matches but the molecule is not `e`, the algorithm has no way to backtrack, and the only honest thing to do is to report the failure rather than return a wrong answer.

See [part2.go](part2.go).

## Complexity

Let *R* be the number of rules, *L* the length of the target molecule in atoms, and *M* the length of the longest right-hand side.

- **Time:** *O*(*R* · *L* · *M*) for Part One. For each rule, we scan the molecule for occurrences of the left-hand side, and each occurrence produces a new string of length *O*(*L*). The number of occurrences per rule is at most *L*, so the worst case is *R* · *L* string constructions. In practice the left-hand sides are short, and `strings.Index` finds matches quickly.

  For Part Two, the complexity is harder to pin down because it depends on the heuristic. Each iteration of the outer loop applies one rule to all its occurrences, and the molecule shrinks. The number of iterations is bounded by the total number of replacements, which for the given input is in the low hundreds. Each iteration does *O*(*R* · *L*) work to scan for matches. The product is *O*(*R* · *L*²) in the worst case, but the constants are small and the actual runtime is milliseconds.

- **Space:** *O*(*L*) for the molecule strings, plus *O*(*R* · *M*) for the rule table. Part One additionally stores a set of distinct molecules, each of length *O*(*L*), so the space is *O*(*U* · *L*) where *U* is the number of distinct products—a few hundred for the given input.

## Associations

- **USPEX and evolutionary structure prediction.** The problem of finding the minimum number of reverse replacements is a search through a combinatorial space, and the heuristic I used is a form of greedy search. A much more sophisticated approach to a related problem is taken by USPEX (Universal Structure Predictor: Evolutionary Xtallography), an evolutionary algorithm for predicting crystal structures from chemical composition alone. USPEX maintains a population of candidate structures, evaluates each one by computing its free energy, and generates new candidates by applying heredity and mutation operators to the best parents. It has been applied to systems with up to 80 atoms in the unit cell and has a success rate near 100% for known structures. Like the heuristic in Part Two, it does not guarantee finding the global optimum, but it works well in practice. See [USPEX](https://uspex-team.org/en).

## Notes

The difficulty of Part Two comes from overlapping right-hand sides. The input does not contain a right-hand side that is a substring of another—if it did, the conflict could be resolved by preferring the longer match, the same way a shift-reduce parser prefers shift over reduce. But that preference gives no guarantee of optimality, and here the conflicts are of a different kind: two right-hand sides can overlap in the molecule without one containing the other, and the choice of which one to reduce first changes what is available next.

In my input, the fragment `CaSiTh` appears eight times, and it is a good illustration of the problem. It can be parsed in two ways. The rule `Si => CaSi` matches the first two atoms, giving `SiTh`, which in turn matches the rule `Ca => SiTh`. Alternatively, the rule `Ca => SiTh` matches the last two atoms, giving `CaSi`, which matches the rule `Si => CaSi`. Both paths are locally valid, both reduce the fragment by two atoms in two steps, and they lead to different intermediate molecules. Only the global structure of the target decides which path is correct. Your input will have its own overlapping fragments; the structure of the rule set makes them almost inevitable.

This is the same phenomenon that parser theory calls a shift-reduce conflict: the parser cannot decide, from local information alone, whether to reduce what it has already seen or to shift the next symbol and reduce later. Preferring shift resolves the local ambiguity but does not guarantee the minimum number of reductions—the choice that looks best now may lead to a dead end or a longer path later.

The heuristic sidesteps the problem by sorting the rules so that longer and more structured right-hand sides are tried first. When the molecule contains `CaSiTh`, the rule `Ca => SiTh` is tried before `Si => CaSi`, because its right-hand side contains more non-terminals. This resolves the conflict in the direction the input was constructed to favour, and on the given input, the greedy path leads to the minimum. There is no guarantee that it works in general, which is why the code panics if the molecule does not reduce to `e`—it is better to report a failure than to return a wrong answer.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~1 ms |
| 2 | ~2 ms |

Both parts are fast. Part One scans the molecule once per rule and builds a set of distinct products; Part Two runs the greedy reduction loop, which terminates after a few hundred iterations. Neither part does enough work to escape the noise of process startup, and the difference between them is within run-to-run variation.

---

[← Previous: Day 18](day18.md) · [Up: Contents](../README.md) · [Next: Day 20 →](day20.md)
