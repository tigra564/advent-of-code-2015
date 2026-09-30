# Day 22: Wizard Simulator 20XX

## The Puzzle

You are a wizard fighting a boss. You start with 50 hit points and 500 mana, and the boss has stats given in the input. On each turn you cast one spell, then the boss attacks you. The fight continues until one of you drops to 0 hit points or you run out of mana.

There are five spells: Magic Missile (instant damage), Drain (instant damage and healing), Shield (temporary armor boost), Poison (damage over time), and Recharge (mana regeneration). Each spell has a mana cost, and some have effects that last several turns. A spell with an ongoing effect cannot be cast again while its effect is still active.

Part One asks for the minimum amount of mana you must spend to win the fight.

Part Two adds a rule: at the start of every player turn, you lose 1 hit point. This is the "hard mode" variant, and it makes some otherwise winning strategies unviable. Part Two asks for the minimum mana under this rule.

## Input

Two lines, one per boss stat:

```
Hit Points: 55
Damage: 8
```

Each line is a label followed by a single integer, and the numbers are extracted with `common.ExtractNumbers`. The boss has no armor stat in this puzzle, which distinguishes it from Day 21—there, armor was part of the fight and the equipment could offset it; here, the only thing that matters is the boss's damage.

The player's starting stats are not part of the input. They are fixed by the puzzle statement at 50 hit points and 500 mana, and they are hard-coded in `solve`. The spell list is also fixed: five spells with their costs, durations, and effects, stored as a package-level array. Both are the same pattern as in previous days—data from the statement, kept in the code.

## The Approach

This is the first day where the fight is not a single calculation but a search over sequences of decisions. Day 21 resolved a fight arithmetically—the outcome depended only on the equipment, and there was exactly one fight per loadout. Here, the outcome depends on the order of spells, and the space of possible orders is huge. The fight is a game tree, and we need the cheapest path to a win.

The five spells are stored in a fixed array, and the spell indices are named constants:

```go
const numSpells = 5

var spells = [numSpells]spell{
	{53, 1, 4, 0, 0, 0},
	{73, 1, 2, 0, 2, 0},
	{113, 6, 0, 7, 0, 0},
	{173, 6, 3, 0, 0, 0},
	{229, 5, 0, 0, 0, 101},
}

const (
	spMagicMissile = iota
	spDrain
	spShield
	spPoison
	spRecharge
)
```

Instant spells—Magic Missile and Drain—are given a timer of 1. This is a small trick: instead of special-casing them in the round logic, the code treats them as spells with a one-turn duration, and the effect applies during the same round in which they are cast. The comment in the code says this changes nothing for Part One, and that is true—the timing is such that the instant effect lands on the boss before he attacks, which is where it would land if it were applied immediately.

The round is evaluated by `evalRound`, which takes a spell number and mutates the fighters in place:

```go
func evalRound(boss, player *fighter, spellNum int) roundResult {
	applyPreliminaryDamage(player)
	if player.hp <= 0 {
		return roundResult{ocLost, player.manaSpent}
	}

	applyActiveSpells(boss, player)

	if player.mana >= spells[spellNum].cost {
		player.spellTimers[spellNum] = spells[spellNum].timer
		player.mana -= spells[spellNum].cost
		player.manaSpent += spells[spellNum].cost
	} else {
		return roundResult{ocLost, player.manaSpent}
	}

	applyActiveSpells(boss, player)

	if boss.hp <= 0 {
		return roundResult{ocWon, player.manaSpent}
	}

	player.hp -= max(1, boss.damage - player.armor)

	if player.hp <= 0 {
		return roundResult{ocLost, player.manaSpent}
	}

	return roundResult{ocInProgress, player.manaSpent}
}
```

The three outcomes are encoded as integers: `ocLost`, `ocInProgress`, `ocWon`. This is simpler than an enum or a string, and it makes the recursion readable.

`applyActiveSpells` is called twice per round. The first call handles the effects that were active before the player's turn—shield, poison, recharge. The second call handles the effects that the player just activated, which for instant spells means applying the damage or healing immediately. For ongoing spells, the second call does nothing on the turn they are cast, because their timer was just set and the effect starts on the next round.

The recursion is where the search happens. `realFightImpl` tries every spell that can be cast, recurses into the resulting state, and keeps the cheapest winning path:

```go
for n := range spells {
	if player.spellTimers[n] > 1 {
		continue
	}
	newBoss := *boss
	newPlayer := *player
	rndRes := evalRound(&newBoss, &newPlayer, n)

	if rndRes.outcome == ocLost || rndRes.manaSpent > best.manaSpent {
		continue
	}
	if rndRes.outcome == ocInProgress {
		rndRes = realFightImpl(&newBoss, &newPlayer, memo)
	}

	if rndRes.outcome == ocWon && rndRes.manaSpent <= best.manaSpent {
		best = rndRes
	}
}
```

Two things are worth pointing out. First, the check `player.spellTimers[n] > 1` prevents casting a spell that is still active. The value 1 means the spell will wear off at the end of this round, so it can be recast next round; a value greater than 1 means the effect is still running. Second, the early exit `rndRes.manaSpent > best.manaSpent` prunes branches that are already more expensive than the best known solution. This is branch and bound, though without a lower-bound heuristic—the pruning is just against the current best, not against an optimistic estimate.

The memoisation is what makes the search tractable. Without it, the game tree explodes: five branches per turn, dozens of turns deep. With it, the number of distinct states is much smaller, because many different spell sequences lead to the same (boss HP, player HP, mana, timers) tuple. The cache key captures exactly that state:

```go
type cacheKey struct {
	playerMana, playerManaSpent int
	playerHP byte
	playerTimers [numSpells]byte
	bossHP byte
}
```

The `manaSpent` field is part of the key as a safety measure: two states that look identical in HP, mana, and timers may have arrived there at different costs, and their optimal continuations may differ. Whether such states actually occur is a question I did not try to settle—including the field costs nothing, and proving it unnecessary would take more effort than it saves.

The types in the key are chosen for compactness: `byte` for HP and timers, because they never exceed 255, and `int` for mana, because 500 and the cost of spells fit comfortably. The key is compared by value, so no hashing is needed beyond what Go's map does.

### Part One

Part One is the fight as described, with no additional rule. The function `applyPreliminaryDamage` does nothing:

```go
func applyPreliminaryDamage(player *fighter) {
	// do nothing
}
```

See [part1.go](part1.go).

### Part Two

Part Two adds the "hard mode" rule: at the start of every player turn, the player loses 1 hit point. This is implemented by a single function:

```go
func applyPreliminaryDamage(player *fighter) {
	player.hp--
}
```

The call is the first thing `evalRound` does, before checking whether the player is still alive. If the player's HP reaches 0 from this damage alone, the round is lost immediately. The rest of the algorithm is unchanged—the same recursion, the same memoisation, the same pruning.

See [part2.go](part2.go).

## Complexity

Let *S* be the number of distinct states, and *K* the number of spells (five in this puzzle).

- **Time:** *O*(*S* · *K*) for both parts. The recursion visits each distinct state at most once, thanks to memoisation, and from each state it tries every spell. The value of *S* is bounded by the product of the state dimensions: boss HP (up to a few hundred), player HP (at most 50), player mana (at most 500), mana spent (up to the mana budget), and the timer array (six values per spell, five spells). The naive product is enormous, but most combinations are unreachable—timers only take values that correspond to spells actually being active, and mana is consumed discretely, so not every intermediate value is attainable. In my input, the search visited 3701 states in Part One and 3836 in Part Two—a few thousand, which is why the whole thing finishes in milliseconds.

- **Space:** *O*(*S*) for the cache, plus *O*(*K*) for the recursion stack. The cache holds one entry per distinct state, and each entry is a small struct. The recursion depth is bounded by the maximum fight length, which is on the order of a few hundred rounds: the boss takes at least 2 damage from Drain every few turns, and the player's HP puts a hard ceiling on how long the fight can last before someone dies.

## Associations

- **The Accountant.** Robert Sheckley's short story "The Accountant" (1954) tells of Morton Dee, the son of a family of powerful wizards, who chooses to become an accountant. His parents are horrified—accounting is dull, mechanical, and utterly beneath a sorcerer. But the family's demon tutor sees it differently: accounting is a form of magic, and a more powerful one than the old kind. "Evasion of income tax as a road to damnation," he says, "that's real magic." The story is a good companion to this puzzle, where the winning strategy is not the flashiest spell sequence but the cheapest one. Brute-force search over every possible fight is the "old magic" of Day 21; the pruned, memoised search here is the accounting version—less dramatic, but it wins.

## Notes

- **The mistakes I made.** This was the hardest puzzle in the event, and I got Part One right twice for the wrong reasons. The first bug was in the check that prevents recasting an active spell. I wrote `if player.spellTimers[n] > 0 { continue }`, which skips a spell as long as its timer is positive. But a timer of 1 means the spell wears off at the end of this round, so it can be recast next round—the condition should have been `> 1`. The fix is explained in `The Approach`, but the point is that the wrong condition did not affect the answer to Part One: the optimal sequence happened to avoid recasting spells early, so the extra restriction never bound. The second bug came from an old habit. In many games, a damage-over-time effect is a reason to wait: let the poison do its work, save your mana, and cast nothing until the enemy is dead or you need to renew the effect. I read the puzzle through that lens and assumed the player could skip a turn. The statement says the opposite: on each turn you *must* cast a spell, and if you cannot afford any, you lose. The habit did not affect Part One either, because the optimal sequence never skips. But in Part Two, where the hard-mode damage makes the fight tighter, the ability to skip turns changes the answer, and my algorithm returned a value that was too low. It took several hours and a careful rereading to find both bugs. The lesson is that a correct answer is not the same as a correct algorithm: a bug can be invisible if the optimal path happens to avoid it, and then surface on a different input where the path changes.

- **The upper bound on rounds.** A fight cannot go on forever, and it is worth working out why. The player starts with 50 HP, and in Part One the boss deals `max(1, 8 − armor)` damage per round. With Shield active, armor is 7 and the damage is 1; without it, the damage is 8. So without healing, the player survives either 50 rounds (with Shield up) or 7 rounds (without). Drain heals 2 HP per cast, and Poison deals 3 damage per round for 6 rounds. My input has 55 HP for the boss. Killing it with Poison alone takes 4 cycles of 6 rounds, or 24 rounds; killing it with Drain alone takes 30 casts, which is also at least 30 rounds. Maintaining Shield and Recharge costs about 342 mana per 6 rounds, and Recharge generates 505, leaving roughly 163 mana per 6 rounds for offense. That is enough for one Poison every 6 rounds, or two Drain casts. In either case, the boss is dead within a few dozen rounds, and the player is kept alive by Shield and Drain. The bound is not tight, but it is finite: a few dozen rounds at most, and in practice the optimal fight in my input lasted well under that.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~7 ms |
| 2 | ~4 ms |

Both parts run the same memoised search; the only difference is the hard-mode damage in Part Two, which prunes some branches earlier and leaves fewer states to explore. The measured times are a few milliseconds each, dominated by process startup. The search itself visits a few thousand states, and each state costs a handful of operations—well within the noise of loading the binary.

---

[← Previous: Day 21](day21.md) · [Up: Contents](../README.md) · [Next: Day 23 →](day23.md)
