# Day 2: I Was Told There Would Be No Math

## The Puzzle

Elves need wrapping paper and ribbon for a list of presents. Each present is a rectangular box with dimensions *l* × *w* × *h*. Part One asks for the total area of wrapping paper required, including a little extra for each present. Part Two asks for the total length of ribbon required, including a bow.

## Input

Each line describes one present as three positive integers separated by `x`:

```
2x3x4
1x1x10
```

The dimensions are not sorted: the smallest side is not necessarily the first one. This matters because both parts need the smallest pair of dimensions, and we cannot assume it is in any particular position.

The numbers are extracted with `common.ExtractNumbers`. It pulls out all integers from a string and ignores the separators between them, which is exactly what we need here: the `x` characters carry no information, only the numbers do. We will keep using this helper throughout the event, whenever the structure of the line matters less than the numbers on it.

## The Approach

Both parts reduce to arithmetic on the three dimensions of a present, with one shared observation: we need the two smallest dimensions. The third one is used only for the volume in Part Two.

We sort the three dimensions once and reuse the sorted slice. This makes the "smallest pair" unambiguous: after sorting, it is always the first two elements.

### Part One

We need the surface area of the box, plus the area of the smallest side as an extra.

The surface area is the sum of the areas of all six faces. Since opposite faces are equal, it is *2·(l·w + w·h + h·l)*. After sorting, *l* and *w* are the two smallest dimensions, so the extra is *l·w*.

```go
dims := common.ExtractNumbers(l)
sort.Ints(dims)
result += 2*(dims[0]*dims[1]+dims[1]*dims[2]+dims[2]*dims[0]) +
	dims[0]*dims[1]
```

See [part1.go](part1.go).

### Part Two

We need the length of the ribbon: the perimeter of the smallest face, plus the volume of the box for the bow.

The smallest perimeter is *2·(l + w)*, where *l* and *w* are the two smallest dimensions. The volume is *l·w·h*. Note that the ribbon wraps around the smallest face—not the two smallest sides individually.

```go
dims := common.ExtractNumbers(l)
sort.Ints(dims)
result += 2*(dims[0]+dims[1]) +
	dims[0]*dims[1]*dims[2]
```

See [part2.go](part2.go).

## Complexity

Let *n* be the number of presents.

- **Time:** *O*(*n*). Each line takes constant time: we extract three integers, sort them, and do a fixed amount of arithmetic. Sorting three elements is a constant cost.
- **Space:** *O*(1).

## Associations

- **Odysseus feigning madness.** The elves, it turns out, are cunning folk—much like Odysseus, who pretended to be insane to avoid sailing to Troy. He yoked an ox and a horse to a plough and sowed his fields with salt. Palamedes saw through the trick by placing baby Telemachus in front of the plough; Odysseus swerved, and his sanity was proven. The elves show a similar kind of practical cunning: they avoid complicated gift shapes by insisting on boxes.

- **The box as a convex hull.** A gift could be a bottle, a Klein bottle, or something worse. But the elves, being lazy, accept only convex shapes. In mathematical terms, the wrapping paper covers the convex hull of the object—or, more precisely, its bounding box, the smallest axis-aligned cuboid that contains it. For a box-shaped gift, the hull and the bounding box coincide. See [Convex hull](https://en.wikipedia.org/wiki/Convex_hull).

- **Klein bottle.** A famous example of a non-orientable surface with no inside and no outside. Wrapping one in paper would be a nightmare: the paper would have to pass through itself. It is fortunate the elves did not ask for that. See [Klein bottle](https://en.wikipedia.org/wiki/Klein_bottle).

## Notes

There is no benchmark for this day. The algorithm is *O*(*n*) with a tiny constant, and the input is short. A benchmark would mostly measure process startup, not the computation. Benchmarks become meaningful from Day 4 onward.

---

[← Previous: Day 1](day01.md) · [Up: Contents](../README.md) · [Next: Day 3 →](day03.md)
