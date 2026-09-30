# Day 11: Corporate Policy

## The Puzzle

Santa's password has expired, and the elves have a set of rules for the new one. It must be exactly eight lowercase letters, contain at least one increasing straight of three letters (`abc`, `bcd`, ..., `xyz`), contain at least two different non-overlapping pairs of the same letter (`aa`, `bb`, ..., `zz`), and contain none of the letters `i`, `o`, `l`.

Part One asks for the next valid password after the current one, in lexicographic order. Part Two asks for the next valid password after the answer to Part One—that is, the one after that.

## Input

A single line of eight lowercase letters, like `abcdefgh`. The password is short, and the rules are restrictive enough that the next valid password is usually close: for the given inputs, the answers are found within a few thousand increments.

## The Approach

The brute-force approach is the obvious one: increment the password like an odometer, check the rules, and stop at the first valid string. The rules are cheap to check, and the answers lie close enough to the starting points that the loop does not run for long.

The `next` function increments the password in place, with `'z'` wrapping to `'a'`:

```go
func next(passwd []byte) {
	// We won't get more than 8 bytes anyway, so no overflow at the first byte
	for i := len(passwd) - 1; i >= 0; i-- {
		c := &passwd[i]
		if *c++; *c <= 'z' {
			break
		}
		*c = 'a'
	}
}
```

The loop goes from right to left, incrementing each character and stopping as soon as the new value fits within `'a'`–`'z'`. If it does not, the character wraps back to `'a'` and the carry propagates to the next position. The comment about overflow is a note to myself: the passwords are always eight characters long, and the answers never come close to `zzzzzzzz`, so the loop never runs out of positions to carry into. If it did, the carry would simply fall off the end and the leftmost character would stay at `'a'`—a silent wrap that the puzzle input never triggers.

Three predicates encode the rules. Each is a single pass over the eight bytes:

```go
func has3Straight(passwd []byte) bool {
	for i := 0; i < len(passwd)-2; i++ {
		if passwd[i]+1 == passwd[i+1] && passwd[i+1]+1 == passwd[i+2] {
			return true
		}
	}
	return false
}
```

```go
func hasNoIllegals(passwd []byte) bool {
	for i := 0; i < len(passwd); i++ {
		if c := passwd[i]; c == 'i' || c == 'o' || c == 'l' {
			return false
		}
	}
	return true
}
```

```go
func hasDifferentPairs(passwd []byte) bool {
	i := 0
	for ; i < len(passwd)-1; i++ {
		if passwd[i+1] == passwd[i] {
			break
		}
	}
	for j := i + 1; j < len(passwd)-1; j++ {
		if passwd[j+1] == passwd[j] && passwd[j] != passwd[i] {
			return true
		}
	}
	return false
}
```

The first two are straightforward. The third finds the first pair of consecutive identical letters, then starts from the next position and looks for a pair of a different letter. If the second scan hits the same letter as the first pair, it silently skips it and keeps going until the end of the string.

The main loop pulls it together:

```go
passwd := []byte(lines[0])
next(passwd)
for !isValid(passwd) {
	next(passwd)
}
result := string(passwd)
```

The first `next` is called before the loop, because we want the *next* valid password, not the current one. Then the loop keeps incrementing until `isValid` returns true.

See [part1.go](part1.go).

### Part One

The answer is the first valid password after the input.

### Part Two

Part Two asks for the next valid password after the answer to Part One. The cleanest way to express this is to run the same loop twice—once for each part—reusing the function. The second call takes the result of the first as its input.

In the repository, Part Two is not a separate file: the same binary is invoked with two different inputs. The first input holds the original password; the second holds the answer to Part One. The elf sitting at the keyboard is the one connecting the two parts.

## Complexity

Let *n* be the number of candidates tested before the first valid password is found. Since the password is always eight characters, each rule check is *O*(1), and so is each increment. The total time is *O*(*n*), with a small constant. For the given inputs, *n* is in the low thousands.

- **Time:** *O*(*n*), where *n* is the number of candidates skipped.
- **Space:** *O*(1). The password is a fixed eight-byte buffer, mutated in place.

## Associations

- **Password strength policies.** The rules in this puzzle are a parody of password policies: minimum length, required character classes, forbidden sequences. Real policies that look like this are widely considered counterproductive. The classic reference is NIST Special Publication 800-63B, which since 2017 recommends dropping forced composition rules (mandatory digits, symbols, mixed case) in favour of length and a blocklist of compromised passwords. The puzzle's rules are even worse than the ones NIST warns against: "at least one straight of three letters" and "no `i`, `o`, `l`" have no security value whatsoever. See [NIST SP 800-63B](https://pages.nist.gov/800-63-3/sp800-63b.html).

- **The elves' minimal effort.** We have learned by now that the elves are cunning and a little lazy—they prefer the smallest amount of work that satisfies the rules. There is a well-documented tendency in humans to do the same: given a composition policy, users pick the shortest string that passes it, not the strongest one. The elves are no different. For these rules, the minimal password is `aaaaabcc`: it has the straight `abc`, two different pairs (`aa` and `cc`), and no forbidden letters. A real attacker would try it first. See [Password strength](https://en.wikipedia.org/wiki/Password_strength).

## Notes

The password length is fixed at eight, and the answers are close to the inputs—a few thousand increments at most. This makes brute force the right choice: there is no need for an analytic solution or a clever search order, and the simple odometer is easy to reason about and hard to get wrong. The mechanism is the same as in a mechanical odometer, where rotating drums advance one position at a time and carry over to the next drum when they wrap around. See [Odometer](https://en.wikipedia.org/wiki/Odometer).

The puzzle is very similar to Day 5: the same idea of checking strings against a handful of rules, and the same style of predicates. The difference is that here we generate the input ourselves—one increment at a time—while in Day 5 the candidates were handed to us in a list. The rules are also a bit different, but the shape of the solution is the same.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~33 ms |
| 2 | ~61 ms |

Both parts run the same loop; the only difference is the starting point. Part Two begins from the answer to Part One, which is already deeper into the search space, so it skips more candidates before reaching a valid password. The measured ratio (61 / 33 ≈ 1.8) reflects that difference, not any change in the algorithm.

---

[← Previous: Day 10](day10.md) · [Up: Contents](../README.md) · [Next: Day 12 →](day12.md)
