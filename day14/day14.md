# Day 14: Reindeer Olympics

## The Puzzle

A group of reindeer are competing in a race. Each reindeer has three numbers: a flying speed in km/s, how many seconds it can fly before tiring, and how many seconds it must rest before flying again. The pattern repeats for the whole race, which lasts 2503 seconds.

Part One asks for the distance travelled by the winning reindeer—the one that covers the most ground by the end of the race, regardless of what happens along the way.

Part Two changes the scoring. Instead of measuring distance at the finish line, the elves award one point at the end of every second to whichever reindeer is in the lead at that moment, or to several of them if their results are equal so far. Part Two asks for the highest total score after 2503 seconds.

## Input

A dozen lines, one per reindeer:

```
Vixen can fly 19 km/s for 7 seconds, but then must rest for 124 seconds.
```

Each line gives a name, a speed, a flight duration, and a rest duration. The numbers are extracted with `common.ExtractNumbers`, which pulls them out in order and ignores the surrounding words. The names are not used by either part—they are decoration.

The race duration is not part of the input; it is fixed at 2503 seconds by the puzzle statement. This is worth noting because both parts depend on it, and neither can be written without hard-coding the constant or passing it in.

## The Approach

Both parts simulate the race, but they simulate it differently. Part One needs only the final distance, and there is a closed-form expression for it—no simulation required. Part Two needs the leader at every second, and that is easier to get by ticking through the race one second at a time.

### Part One

A reindeer's motion is periodic: it flies for *f* seconds, rests for *r* seconds, and repeats. The period is *f* + *r*. In the full race of *T* seconds, the number of complete periods is `T / (f + r)`, and each complete period contributes `f * speed` kilometres. The remaining `T % (f + r)` seconds are either spent flying—if the remainder is less than *f*—or split between a full flight and the start of a rest:

```go
dist := (raceTime / (flightDuration + restDuration)) * flightDuration * speed
dist += min(raceTime % (flightDuration + restDuration), flightDuration) * speed
result = max(result, dist)
```

The `min` in the second line is what handles the partial period: the reindeer flies for at most `flightDuration` seconds of the remainder, and any leftover time is spent resting, which does not add distance.

See [part1.go](part1.go).

### Part Two

The scoring rule breaks the closed form. We need the leader at every second, and the leader depends on the state of every reindeer at that moment, not just on the final totals. There is no algebraic shortcut for "who was ahead at second 1374"; the simplest way is to advance the clock and watch.

Each reindeer is a small state machine. It has a `state` (flying or resting), a `timer` counting down the seconds remaining in the current phase, a running `distance`, and a `score`:

```go
type reindeer struct {
	speed int
	flightDuration int
	restDuration int

	state deerState
	timer int
	distance int
	score int
}
```

The simulation ticks once per second. On each tick, every reindeer advances by one second: if it is flying, it adds `speed` to its distance and decrements its timer; if the timer reaches zero, it switches to resting and resets the timer to `restDuration`. The resting case is symmetric—no distance is added, the timer runs down, and when it reaches zero the reindeer switches back to flying.

After all reindeer have advanced, we find the maximum distance among them, and award a point to each reindeer whose distance equals that maximum. The `maxDist` is computed in the same loop that advances the state, and the scoring is a second pass over the slice. This separation is not accidental: a reindeer can only be scored once the leader is known, and the leader is only known once everyone has moved.

```go
for range raceTime {
	maxDist := 0
	for i := range contestants {
		c := &contestants[i]
		switch c.state {
		case FLYING:
			c.distance += c.speed
			if c.timer--; c.timer == 0 {
				c.state = RESTING
				c.timer = c.restDuration
			}
		case RESTING:
			if c.timer--; c.timer == 0 {
				c.state = FLYING
				c.timer = c.flightDuration
			}
		}
		maxDist = max(maxDist, c.distance)
	}

	for i := range contestants {
		c := &contestants[i]
		if c.distance == maxDist {
			c.score++
		}
	}
}
```

The initial state of each reindeer is flying, with `timer` set to `flightDuration`. This is correct because every reindeer starts the race in the air—the puzzle does not mention a warm-up period.

The final answer is the maximum score among all reindeer. Note that the winner of Part Two is not necessarily the winner of Part One: a reindeer can lead for many seconds and still be overtaken at the finish, and vice versa.

See [part2.go](part2.go).

## Complexity

Let *n* be the number of reindeer and *T* the race duration.

- **Time:** *O*(*n*) for Part One—one pass over the input, constant work per reindeer. For Part Two it is *O*(*n* · *T*)—each of the *T* seconds iterates over all *n* reindeer twice, once to advance and once to score.

- **Space:** *O*(*n*) for the slice of reindeer in Part Two. Part One needs only a few scalars, so *O*(1) beyond the input.

## Associations

- **Closed-form versus simulation.** The two parts of this puzzle illustrate a general trade-off: when a process is periodic and we only need its state at a specific time, a closed-form expression is faster and simpler than stepping through it. When the question is about the history of the process—who was ahead at every intermediate moment—the simulation cannot be avoided. This is the same distinction as between evaluating a function at a point and plotting its entire graph. See [Closed-form expression](https://en.wikipedia.org/wiki/Closed-form_expression).

- **Scoring rules and the Condorcet paradox.** Part Two awards points second by second, which means the winner is the reindeer that led most often, not the one that finished first. This is a different aggregation of the same underlying data, and it can give a different winner—a small illustration of the general fact that the choice of scoring rule determines the outcome. The same phenomenon appears in voting theory, where different aggregation methods can rank the same set of candidates differently. See [Condorcet paradox](https://en.wikipedia.org/wiki/Condorcet_paradox).

## Notes

The runtime surprise from the measurements—Part One and Part Two taking roughly the same time—has a simple explanation: both are dominated by process startup, not by the algorithm. Part One does a dozen arithmetic operations; Part Two ticks a clock 2503 times over a dozen reindeer, which is about 60 000 iterations. On a modern CPU both are lost in the noise of loading a binary and reading the input file. The `Complexity` section says *O*(*n*) for Part One and *O*(*n* · *T*) for Part Two, and that is correct, but for this input both are effectively constant. The benchmarks confirm it: 25 ms and 23 ms of *user* time, with the difference well within run-to-run variation.

If the race were a hundred times longer, the second-by-second simulation would still be correct, but a different approach would pay off. Each reindeer changes state at predictable times, and between state changes the distance grows linearly—so the moment when one reindeer overtakes another can be computed directly, without stepping through every second. The events form a priority queue: at each step, extract the next overtake or state change, advance the clock, and award the accumulated points to whoever is leading during the interval. The complexity stays linear in the race length, but with a smaller constant: instead of *T* steps, we process only the moments when something actually changes. The saving depends on how often the leaders swap, which in a close race can be almost every second—so the improvement is real, but not dramatic.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~3 ms |
| 2 | ~2 ms |

Both parts are dominated by process startup and input parsing, not by the algorithm. Part One performs a dozen arithmetic operations per reindeer; Part Two ticks a clock 2503 times over a dozen reindeer, which is about 60 000 iterations. On a modern CPU both are lost in the noise of loading a binary and reading the input file. The difference between the two parts is within run-to-run variation.

---

[← Previous: Day 13](day13.md) · [Up: Contents](../README.md) · [Next: Day 15 →](day15.md)
