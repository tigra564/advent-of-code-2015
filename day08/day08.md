# Day 8: Matchsticks

## The Puzzle

The input is a list of string literals, with the usual C-like escape sequences: `\"`, `\\`, `\x` followed by two hex digits, and so on. Part One asks how many characters of code are saved by removing the quotes and interpreting the escapes. Part Two asks the opposite: how many characters are added by re-encoding each literal.

## Input

A list of string literals, one per line:

```
""
"abc"
"aaa\"aaa"
"\x27"
```

Each line is a quoted literal with escape sequences. The escapes are the familiar ones from C and its descendants—`\"`, `\\`, `\x` followed by two hex digits—nothing exotic.

It happens that Go's standard library already knows how to parse and produce these literals: `strconv.Unquote` and `strconv.Quote`. The Go designers did not invent a new escape syntax; they reused the one that was already everywhere. That turns each part of this puzzle into a two-liner, which is a happy coincidence rather than a property of the puzzle.

## The Approach

Both parts are two-liners: one call into the standard library, one arithmetic operation. Go's `strconv` package already knows how to parse and produce these literals, and the puzzle is mostly about noticing that.

### Part One

We need the difference between the length of the literal and the length of the string it represents. `strconv.Unquote` parses a quoted literal with escapes and returns the actual string:

```go
unquoted, _ := strconv.Unquote(line)
result += len(line) - len(unquoted)
```

`len` here is the byte length. The escapes produce ASCII characters, so byte length and character count coincide for this input.

See [part1.go](part1.go).

### Part Two

We need the difference between the length of the encoded literal and the length of the original string. `strconv.Quote` does the encoding:

```go
quoted := strconv.Quote(line)
result += len(quoted) - len(line)
```

See [part2.go](part2.go).

## Complexity

Let *n* be the number of lines and *m* the length of a line.

- **Time:** *O*(*n* · *m*). Each line is unquoted or quoted once, and both operations scan the string from start to finish.

- **Space:** *O*(*m*) per line for the intermediate result. The `common.ReadInputLines` slice holds the whole input, but that is a property of the input reader, not of this algorithm.

## Associations

- **Escape sequences as a language.** The backslash is the oldest and most persistent piece of syntax in computing. It predates C: it appears in ALGOL 60, in BCPL, and in the teleprinter conventions of the 1950s, where `\n` was a way to tell a mechanical device to move the carriage. Every language since has inherited some version of it, and every language has had to make the same choice—which characters to escape, and how. Go chose to reuse the C set rather than invent a new one, and that choice is what makes this puzzle a two-liner. See [Escape sequence](https://en.wikipedia.org/wiki/Escape_sequence).

- **Titlo: marking a letter as a number.** The idea that a symbol needs a special mark to signal "this is not a letter, this is a number" is old. In medieval Greek, Latin, and Slavic manuscripts, a titlo—a wavy or zigzag line placed above a letter—indicated that the letter stood for a numeral rather than a sound. The Greeks placed a horizontal stroke over letters used as numbers: alpha with a bar meant 1, beta with a bar meant 2, and so on through the alphabet. The convention was inherited by the Slavs, and the same technique appears in the Cyrillic tradition as the titlo. It is the same problem the escape sequence solves: a symbol's meaning depends on context, and a marker disambiguates. See [Titlo](https://en.wikipedia.org/wiki/Titlo).

## Notes

Most standard libraries in mainstream languages have an equivalent of `Quote` and `Unquote`. Go's are in `strconv`. Python has `codecs.escape_decode` and `unicode_escape`. C has neither as a function, so a C version would need some code of its own—but not much.

The puzzle can also be solved in bash, and the two parts together fit on two lines:

Part One:

```bash
a=$(sed -E -e 's,\\(x..|\\|"),_,g' -e 's,",,g' < input.txt| wc -c); b=$(wc -c < input.txt); echo $((b-a))
```

Part Two:

```bash
a=$(sed -e 's,[\\"],__,g' -e 's,^,"",g' < input.txt | wc -c); b=$(wc -c < input.txt); echo $((a-b))
```

In Part One, the first `-e` expression in `sed` replaces every escape sequence with a single underscore, which has the same length as the two-character escape it replaces but represents one character in the output. The second `-e` expression in `sed` strips the surrounding quotes. In Part Two, the first `-e` in `sed` replaces each backslash and each quote with two underscores, and the second `-e` in `sed` adds two quotea at the beginning of each line: the place does not matter, only character count. The arithmetic is the same in both cases—the difference between two byte counts—and the direction of the subtraction is the only thing that changes.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | 1.2 ms |
| 2 | 1.6 ms |

Both parts are dominated by process startup. The algorithm itself is a single pass over each line, and the input is small. The difference between the two parts is within noise.

---

[← Previous: Day 7](day07.md) · [Up: Contents](../README.md) · [Next: Day 9 →](day09.md)
