# Day 12: JSAbacusFramework.io

## The Puzzle

Santa's accounting software stores its data as JSON, and the elves need a total. Part One asks for the sum of all numbers in the document—integers, positive and negative, wherever they appear.

Part Two adds a rule: any object (a JSON map) that contains the value `"red"` is removed entirely, along with everything inside it. Arrays are not affected by the rule—only objects. Then the sum is recomputed.

## Input

A single line of JSON, about 40 KB long. It contains nested objects and arrays, integers, and strings. The structure is arbitrary: there is no fixed schema, and the depth of nesting is not specified.

Every number in the document is a JSON number—an integer, positive or negative. No digits appear inside string literals, and no number is used as an object key. This matters for Part One: it means we can pull integers out of the raw text with a simple pattern, without parsing the JSON at all. In fact, we already have a tool for that: `common.ExtractNumbers` does exactly what is needed!

## The Approach

### Part One

The simplest approach is to pull all integers out of the line with a regex and sum them. `common.ExtractNumbers` does exactly this—it finds every sequence of digits, with an optional leading minus, and returns them as a slice:

```go
func countNumbers(str string) int64 {
	result := int64(0)
	for _, n := range common.ExtractNumbers(str) {
		result += int64(n)
	}

	return result
}
```

This works because of the property of the input we noted above: every number in the document is a JSON number, and no digits appear inside string literals or object keys. If any of those conditions failed, the regex would pick up false positives or miss negatives, and the answer would be wrong.

All that remains is to call the function:

```go
lines := common.ReadInputLines()

result := countNumbers(lines[0])
```

See [part1.go](part1.go).

### Part Two

A regex would not be the natural tool here. With backreferences it could in principle match a balanced object and check its contents, but Go's RE2 engine does not support them. Another option is to emulate a stack by hand: scan for `{`, walk forward to the matching `}`, and if a `"red"` value appears along the way, drop the whole span. This is doable, but the bookkeeping—tracking nesting depth, handling strings that contain braces, and remembering where each object started—is no simpler than reaching for the standard library. So I parsed the JSON instead.

Go's `encoding/json` package does the parsing. The result is a generic tree of `interface{}` values:

```go
var data []interface{}
json.Unmarshal([]byte(s), &data)
```

The `processValue` function walks this tree recursively:

```go
func processValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case []interface{}:
		var newSlice []interface{}
		for _, elem := range val {
			if newElem := processValue(elem); newElem != nil {
				newSlice = append(newSlice, newElem)
			}
		}
		if len(newSlice) > 0 {
			return newSlice
		} else {
			return nil
		}

	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range val {
			if value == "red" {
				return nil
			}
			if newValue := processValue(value); newValue != nil {
				newMap[key] = newValue
			}
		}
		if len(newMap) > 0 {
			return newMap
		} else {
			return nil
		}
	}

	return v
}
```

The rule about `"red"` is checked only in the `map[string]interface{}` branch: if any value in the object equals the string `"red"`, the whole object is dropped. Arrays are not checked, because the puzzle says only objects are affected.

The `return nil` when a container becomes empty is a small design decision. A JSON object with all of its values removed is semantically the same as an empty object, and an empty object contributes nothing to the sum. Returning `nil` for it simply prunes the tree a little, and saves the marshalling step from producing `{}` and `[]` in the output. The behavior would be identical if we returned the empty containers instead—the sum would not change—but the code would be marginally noisier.

After the tree is rebuilt, it is marshalled back to a string, and the same `countNumbers` from Part One is applied:

```go
str, _ := json.Marshal(newdata)
result := countNumbers(str)
```

See [part2.go](part2.go).

## Complexity

Let *n* be the size of the input JSON.

- **Time:** *O*(*n*) for both parts. Part One scans the string once with a regex. Part Two parses the JSON (*O*(*n*)), walks the resulting tree (*O*(*n*)), marshals it back (*O*(*n*)), and scans the result with the same regex as Part One (*O*(*n*)). All steps are linear, though the constants differ.

- **Space:** *O*(*n*) for Part Two, where the parsed tree and the rebuilt tree both live in memory simultaneously. Part One uses only the slice of numbers returned by `ExtractNumbers`, which is *O*(*n*) in the worst case as well.

## Associations

- **Subgroup analysis and multiple comparisons.** The two parts of the puzzle compute a statistic—a function of a sample—first over the entire document, then over the subset of objects that survive the `"red"` rule. The statistic itself is artificial, but the structure is the same as in clinical trials, where an effect is first estimated on the whole population and then recomputed on a subgroup. The danger in real applications is that testing many subgroups inflates the chance of finding a spurious effect—a problem known as multiple comparisons. Regulators like NICE require subgroups to be specified before data collection for exactly this reason. Bayesian methods address it by shrinking subgroup estimates toward the overall mean when data are sparse. See [Subgroup analysis](https://www.ncbi.nlm.nih.gov/books/NBK602660/).

## Notes

The runtime surprise is real and worth a closer look. Part One does more work than Part Two by the naive count: it sums 2533 numbers, while Part Two sums only 1200. Yet Part One is slower—by roughly a factor of two.

The reason is that `ExtractNumbers` uses `regexp.FindAllString`, which walks the entire string and collects all matches into a slice before returning. For a 40 KB input, the regex engine must examine every byte, and the allocation of a 2533-element slice adds up. `ExtractNumbers` grows its result slice with `append`, which means a handful of reallocations and copies before the final size is reached. In Part One the slice holds 2533 numbers, in Part Two only 1200, so Part One pays a bit more for this. The difference is small compared to the cost of the regex itself, but it is there. Part Two, by contrast, allocates a whole tree of `interface{}` values during `json.Unmarshal`—far more allocations than a single integer slice—and still runs faster. This reinforces the main point: the regex engine, not the allocator, dominates Part One.

The lesson is not that regexes are slow in general, but that Go's `regexp` package is a linear-time engine with a fair amount of per-byte overhead. When the input is larger than a few kilobytes and the pattern is simple, a hand-written parser or a standard library decoder is often faster. The puzzle's input is small enough that the difference is milliseconds, but the direction of the effect is clear.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~10 ms |
| 2 | ~6 ms |

Part Two is faster, despite doing more work: parsing the JSON, walking the tree, marshalling it back, and then running the same `ExtractNumbers` over the result. The reason is discussed in `Notes`: the regex engine in Part One costs more than the whole JSON round-trip in Part Two.

---

[← Previous: Day 11](day11.md) · [Up: Contents](../README.md) · [Next: Day 13 →](day13.md)
