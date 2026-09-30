# Day 4: The Ideal Stocking Stuffer

## The Puzzle

Santa needs help mining AdventCoins. The mining algorithm is simple: take a secret key (the input line), append a decimal number to it without leading zeroes, hash the result with MD5, and check the hexadecimal representation of the hash. A hash is considered valid if it starts with a certain number of zeroes.

Part One asks for the smallest number that produces a hash starting with at least five zeroes. Part Two asks for the smallest number that produces a hash starting with exactly six zeroes.

## Input

A single line—the secret key. For the real input it is a short lowercase string like `abcdef`, but the puzzle guarantees nothing about its length or content beyond it being a non-empty string.

Unlike previous days, the input is not consumed character by character. We take the whole line as a prefix and append numbers to it.

## The Approach

Both parts are a brute-force search: try 1, 2, 3, ... until the hash of prefix + number meets the condition. There is no clever way to invert MD5, and the puzzle is designed around this fact.

The search itself is a single loop. The difference between the two parts is only the condition on the hash. I keep that condition inline rather than extracting it into a function, because the two parts are separate programs and the condition is the only thing that changes.

### Part One

We need at least five hex zeroes, which is twenty bits. The MD5 digest is a `[16]byte` array. The first two bytes must be zero, and the high nibble of the third byte must also be zero. The low nibble of the third byte is unconstrained: five zeroes in hex correspond to twenty bits, and twenty bits is exactly two and a half bytes.

Checking the bits directly is faster than formatting the hash as a hex string and comparing strings. The condition uses a bitmask to isolate the high nibble:

```go
sample := []byte(prefix + strconv.Itoa(result))
digest := md5.Sum(sample)
if 0 == digest[0] && 0 == digest[1] && 0 == digest[2]&0b11110000 {
	break
}
result++
```

See [part1.go](part1.go).

### Part Two

Now we need exactly six hex zeroes, which is twenty-four bits—exactly three bytes. The condition changes: the first three bytes must be zero, and the high nibble of the fourth byte must not be zero. The second part matters: without it, a hash with seven zeroes would also be accepted, and the number we are looking for would not be the smallest one with exactly six zeroes.

```go
sample := []byte(prefix + strconv.Itoa(result))
digest := md5.Sum(sample)
if 0 == digest[0] && 0 == digest[1] && 0 == digest[2] && 0 != digest[3]&0b11110000 {
	break
}
result++
```

See [part2.go](part2.go).

## Complexity

Let *k* be the number of required hex zeroes (*k* = 5 for Part One, *k* = 6 for Part Two). Let *N* be the answer—the smallest number that produces a valid hash.

- **Time:** *O*(*N*) MD5 computations. Each computation is *O*(1) in the input size: the prefix is short, and the number we append has *O*(log *N*) digits.

  In expectation, *N* is about 16*ᵏ*: a random hash has each hex digit uniformly distributed, so the probability that the first *k* digits are zero is 16⁻*ᵏ*. For *k* = 5 that is one in 2²⁰ ≈ 10⁶; for *k* = 6, one in 2²⁴ ≈ 1.7×10⁷. The measured values (117 946 and 3 938 038) are smaller than these expectations—both parts got lucky.

- **Space:** *O*(1). A fixed-size digest buffer and a few scalars.

## Associations

- **MD5 and cryptographic hash functions.** MD5 is a cryptographic hash function: it maps an arbitrary input to a fixed 128-bit output, and it is designed so that the output looks uniformly random. That uniformity is exactly what makes the brute-force search in this puzzle work: each hex digit of the output behaves like an independent coin flip with sixteen outcomes. See [MD5](https://en.wikipedia.org/wiki/MD5).

- **The birthday attack.** MD5 is considered cryptographically broken. The best known collision attack finds two different inputs with the same hash in about 2⁶⁴ operations—far fewer than the 2¹²⁸ that a perfect hash function would require. This is a consequence of the birthday paradox: in a set of *n* random hashes, a collision appears after roughly √*n* samples, not *n*. The puzzle does not require collisions, but the same probabilistic reasoning explains why the search for a prefix terminates so quickly. See [Birthday attack](https://en.wikipedia.org/wiki/Birthday_attack).

## Notes

The requirement of "exactly six zeroes" in Part Two, rather than "at least six", does not change the complexity estimate: the search still terminates after roughly 16⁶ hashes. It can, however, change the result. If the first hash with at least six zeroes also has seven, then it is not valid for Part Two, and we must keep searching. In my input this did not happen—the first candidate with six zeroes had exactly six—but your mileage may vary.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~23 ms |
| 2 | ~643 ms |

The ratio is about 28×, larger than the theoretical 16×, because both parts deviated from the average—Part One found its answer earlier, Part Two later.

---

[← Previous: Day 3](day03.md) · [Up: Contents](../README.md) · [Next: Day 5 →](day05.md)
