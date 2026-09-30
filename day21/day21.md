# Day 21: RPG Simulator 20XX

## The Puzzle

You are playing a simplified RPG. Your character has 100 hit points, no damage, and no armor; the boss has stats given in the input. Before the fight, you buy equipment from a shop: exactly one weapon, at most one piece of armor, and up to two rings. Each item has a cost, a damage bonus, and an armor bonus.

The fight is turn-based. On each turn, the attacker deals `max(1, attacker.damage - defender.armor)` damage to the defender. You attack first, then the boss, and so on until one of you drops to 0 hit points.

Part One asks for the minimum amount of gold you can spend and still win the fight.

Part Two asks for the maximum amount of gold you can spend and still *lose* the fight—the most expensive loadout that guarantees defeat.

## Input

Three lines, one per boss stat:

```
Hit Points: 103
Damage: 9
Armor: 2
```

Each line is a label followed by a single integer, and the numbers are extracted with `common.ExtractNumbers`. The boss stats are fixed for the whole puzzle; there is no randomness in the fight, and the outcome is determined entirely by the equipment you choose.

The shop contents are not part of the input. They are given in the puzzle statement as fixed lists of items, and they are embedded in `day21.go` as package-level slices: `weapons`, `armor`, and `rings`. This is the same pattern as the MFCSAM reading in Day 16—data from the statement that never changes, kept in the code rather than in a separate file.

The shop lists are small. There are five weapons, six armor options (including the "no armor" entry), and eight ring entries (including two "no ring" entries, so that one or two rings can be selected). The total number of loadouts is the product of these counts: 5 × 6 × C(8, 2) = 840. That is small enough to enumerate exhaustively, which is what both parts do.

## The Approach

Both parts enumerate the same set of loadouts and evaluate the fight for each. The only difference is which loadout wins: Part One wants the cheapest winning loadout, Part Two wants the most expensive losing one. The shared logic lives in `solve`, which takes a flag indicating whether we are looking for a win or a loss.

The shop is modelled with three slices of items, each item carrying cost, damage, and armor:

```go
type item struct {
	cost, damage, armor int
}
```

The "optional" slots—armor and rings—are handled by including placeholder items with all stats equal to zero. For armor, there is one such entry; for rings, there are two, so that the enumeration can pick zero, one, or two rings from the same list without special-casing the count.

The enumeration is four nested loops: one over weapons, one over armor, and two over rings. The ring loops use the `rings[n+1:]` idiom to avoid picking the same ring twice and to avoid counting the same pair in both orders:

```go
playerInv := make([]item, 4, 4)
for _, wp := range weapons {
	playerInv[0] = wp
	for _, ar := range armor {
		playerInv[1] = ar
		for n, rn1 := range rings {
			for _, rn2 := range rings[n+1:] {
				playerInv[2] = rn1
				playerInv[3] = rn2
				cost, win := evalFight(boss, player, playerInv)
				if mustWin == win {
					result = evalFn(result, cost)
				}
			}
		}
	}
}
```

The `playerInv` slice is reused across all iterations. This is safe because `evalFight` only reads from it—it does not store a reference or modify the items. Reusing the slice avoids 840 small allocations, which at this scale is not a performance concern but is a habit worth keeping.

`evalFight` sums the stats of the loadout and resolves the fight:

```go
gear := item{0, 0, 0}
for _, elem := range playerInv {
	gear.cost += elem.cost
	gear.damage += elem.damage
	gear.armor += elem.armor
}

player.damage += gear.damage
player.armor += gear.armor

damageToPlayer := max(1, boss.damage - player.armor)
damageToBoss := max(1, player.damage - boss.armor)
```

The `max(1, ...)` is the puzzle's rule that armor cannot reduce incoming damage below 1. Without it, a high-armor character could take zero damage and the fight would never end.

The fight is resolved arithmetically, without simulating turn by turn. If each fighter deals a constant amount of damage per turn, the number of turns each can survive is `(hp - 1) / damagePerTurn`, rounded down:

```go
playerRoundsToDeath :=  (player.hp - 1) / damageToPlayer
bossRoundsToDeath := (boss.hp - 1) / damageToBoss

return gear.cost, playerRoundsToDeath >= bossRoundsToDeath
```

The `-1` in the numerator is the standard trick for "how many hits can I take before reaching 0 HP": if I have 100 HP and take 9 damage per turn, I die on the 12th hit (100 → 91 → 82 → ... → 1 → dead), and `(100 − 1) / 9 = 11`, which means 11 hits survived, so the 12th is fatal. The player wins if the boss dies first or on the same turn, which is why the comparison is `>=`, not `>`.

### Part One

Part One looks for the cheapest winning loadout. The initial value is `math.MaxInt`, and the comparison keeps the minimum:

```go
result := math.MaxInt
evalFn := func (a, b int) int { return min(a, b) }
```

The call from `part1.go` is `solve(true)`.

See [part1.go](part1.go).

### Part Two

Part Two looks for the most expensive losing loadout. The initial value is `math.MinInt`, and the comparison keeps the maximum:

```go
result = math.MinInt
evalFn = func (a, b int) int { return max(a, b) }
```

The call from `part2.go` is `solve(false)`.

See [part2.go](part2.go).

## Complexity

Let *W* be the number of weapons, *A* the number of armor options, and *R* the number of ring entries (including the two "no ring" placeholders). Let *k* be the number of slots in the loadout—four in this puzzle: one weapon, one armor, two rings.

- **Time:** *O*(*W* · *A* · *R*²) for both parts. The enumeration is four nested loops: *W* choices of weapon, *A* choices of armor, and *R*(*R*−1)/2 choices of ring pair. The factor of 1/2 comes from the `rings[n+1:]` idiom, which avoids counting each pair twice. For *W* = 5, *A* = 6, *R* = 8, the count is 5 × 6 × 28 = 840 loadouts, and each loadout costs *O*(*k*) to evaluate. The total is *O*(*W* · *A* · *R*² · *k*), which for these numbers is a few thousand operations.

- **Space:** *O*(*k*) for the `playerInv` slice and the `gear` accumulator. The shop lists are *O*(*W* + *A* + *R*) and live in package-level variables, allocated once at program start. No part of the algorithm allocates per loadout—the slice is reused—so the total working memory beyond the shop lists is constant.

## Associations

- **The George Business.** Roger Zelazny's 1980 short story tells the origin of St. George—not as a dragon-slayer, but as a businessman. A knight and a dragon named Dart, instead of fighting to the death, strike a mutually beneficial deal: the knight wins public glory, the dragon keeps a share of the treasure and the reputation. A good game often allows this kind of expansion—cooperation instead of combat, negotiation instead of slaughter, as in the *Civilization* series, where you can win by diplomacy or culture rather than by arms. Simulating that in Advent of Code would be a different and much harder problem; here we are limited to buying gear and rolling dice. Still, the story is a useful reminder that "win" does not always mean "defeat the other side".

## Notes

The four nested loops could be replaced with a single loop over a bitmask if the shop rules were messier. Each item would get one bit, and the mask would encode the whole loadout; the item lists would be indexed by bit position, and validity would be checked by counting the bits in each category. This is the same technique used in Day 17, and it is occasionally the right one—for example, when the rules forbid certain combinations or allow arbitrary numbers of items. Here the rules are simple enough that the loops are shorter and clearer than the mask arithmetic would be.

The bitmask approach has a long history that predates the puzzle's era. When CPUs had no instruction for counting set bits—no `POPCNT`, no `bits.OnesCount`—programmers built lookup tables by hand. A table of 16 entries would give the popcount of every possible nibble, and a larger count would be assembled from four nibble lookups. The technique is described in detail in *Hacker's Delight* by Henry S. Warren, and it remains relevant on architectures where popcount is still slow or absent. The moral is the same as in the loops-versus-mask question: the right representation depends on what operations the hardware and the language make cheap.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~3 ms |
| 2 | ~2 ms |

Both parts enumerate the same 840 loadouts and evaluate a fight for each. The arithmetic per fight is a handful of additions, one division, and a comparison—microseconds in total. The measured times are dominated by process startup and input parsing, and the difference between the two parts is within run-to-run variation. There is nothing to optimise here: the work is too small to measure meaningfully.

---

[← Previous: Day 20](day20.md) · [Up: Contents](../README.md) · [Next: Day 22 →](day22.md)
