# Advent of Code 2015

## The Goal

This is a continuation of [my Advent of Code 2025 repository](https://github.com/tigra564/advent-of-code-2025/), where I solved the puzzles in C++. That project was, among other things, a way to relearn a language I had not touched seriously in fifteen years. This one has a different premise: I had never written Go before I started it. A month ago I had not touched the language at all. The choice was not deep—Go is popular, relatively simple, and its standard library is a good fit for the kind of small programs Advent of Code asks for. But it was a real choice, and it made the first few days slower than they would have been in a language I knew.

What I did expect was that the write-ups would take time. What I did not expect was that they would take more time than the solutions themselves—even with AI generating the drafts under my direction and control. The code is here, but the explanations are the point. If you want to see the solutions, they are in the day folders. If you want to see the thinking, they are in the `dayNN.md` files.

If you want the long version of why I started doing Advent of Code at all, the 2025 README has it. The short version is that I spent twenty years writing code and then ten years managing people who write code, and I wanted to make sure I could still do the former. The 2025 project proved that I could, in C++. This one proves that I can learn a new language from scratch, and finish all twenty-five days while doing it.

## The Repository

Each day lives in its own directory: `day01/`, `day02/`, and so on, up to `day25/`. Inside each, you will find:

- `part1.go` and `part2.go` — the two solutions. Some days have only one, either because the puzzle has no second part or because both parts share a single binary.
- `dayNN.go` — shared code used by both parts, when there is any.
- `dayNN.md` — the write-up. This is the part I actually care about.
- `input.txt` — not present. Advent of Code asks that you do not publish your puzzle input, and I follow that rule. The write-ups describe the inputs I received, and the code works on any valid input.

The shared utilities live in `common/common.go`. It contains a handful of functions used across the days: reading input from a file or standard input, extracting integers from a line, iterating over permutations, and a small structure for undirected weighted graphs. I do not describe `common` as a package anywhere; each function is introduced in the day where it first appears.

The write-ups follow a consistent structure. `The Puzzle` states the problem in a few sentences. `Input` describes the shape of the input and, more importantly, the properties that affect the algorithm. `The Approach` explains the solution, part by part, with code excerpts. `Complexity` gives the time and space bounds. `Associations` links the puzzle to something outside it—a classical algorithm, a real-world application, a story. `Notes` collects the dead ends, the counterexamples, and the things I would have done differently. `Benchmarks` gives measured runtimes, when measuring makes sense.

The code is not the point. It is correct, and it is clean enough to read, but the algorithms are the same ones you would find in any textbook. What I tried to do is explain *why* each choice was made—why brute force here, why dynamic programming there, why a particular data structure rather than another. Some days are tuned further than others; Day 20, for instance, went through several rounds of optimisation before I was satisfied. If you are looking for a reference solution to copy, this is not it. If you are looking for a worked example of how someone thinks through a problem, it might be.

## Daily Write-ups

[Day 1](day01/day01.md) · [Day 2](day02/day02.md) · [Day 3](day03/day03.md) · [Day 4](day04/day04.md) · [Day 5](day05/day05.md) · [Day 6](day06/day06.md) · [Day 7](day07/day07.md) · [Day 8](day08/day08.md) · [Day 9](day09/day09.md) · [Day 10](day10/day10.md) · [Day 11](day11/day11.md) · [Day 12](day12/day12.md) · [Day 13](day13/day13.md) · [Day 14](day14/day14.md) · [Day 15](day15/day15.md) · [Day 16](day16/day16.md) · [Day 17](day17/day17.md) · [Day 18](day18/day18.md) · [Day 19](day19/day19.md) · [Day 20](day20/day20.md) · [Day 21](day21/day21.md) · [Day 22](day22/day22.md) · [Day 23](day23/day23.md) · [Day 24](day24/day24.md) · [Day 25](day25/day25.md)

## Building and Running

The code is written in Go and requires Go 1.26 or newer. There are no external dependencies: everything is in the standard library and in `common/common.go`.

You do not have to build anything to try a solution. `go run` compiles and executes in one step:

```bash
cd day25
go run part1.go input.txt
```

For a single run this is the easiest way, and the difference in speed is negligible—a few milliseconds of compilation, once. If you want to run the same binary many times, or if you want to measure its performance, build it first:

```bash
cd day25
go build part1.go
./part1 input.txt
```

This produces a binary named `part1` in the same directory. If no argument is given, the program reads from standard input instead. This is convenient for testing on the sample inputs from the puzzle statements, which are small enough to paste into a terminal:

```bash
./part1 < sample.txt
```

Some days have a shared file (`dayNN.go`) that both parts depend on. In those cases, pass it to `go run` or `go build` together with the part:

```bash
cd day07
go run part1.go day07.go input.txt
```

The write-ups mention when a day has a shared file and what it contains.

There is no top-level `Makefile` or build script. Twenty-five days, each with one or two small binaries, do not need one. If you want to build everything at once, `go build ./...` from the repository root will do it, but it will not produce binaries in the day directories unless you run it from there.

## A note on AI

I used AI in this project, and I am not going to hide it. There are three places where it helped.

The first is the write-ups themselves. The text you are reading was drafted by an AI under my direction. I decided what to say, which examples to use, which associations to make, and where to be honest about mistakes. The AI turned those decisions into prose. This is different from the 2025 repository, where I wrote everything myself and only ran it through an AI at the end to catch inconsistencies. Here I delegated the drafting from the start, which saved time and made the descriptions more consistent.

The second is code review. After solving a day and getting the right answer, I always asked an AI to look at the code and point out anything that looked off. This was part of the process, not an occasional extra: I had no colleague to review this project, so the AI was the reviewer. It usually flagged naming, leftover debug output, and unnecessary copies. Nothing substantial—no algorithmic errors, no missed edge cases—but the code became a little cleaner each time. In Day 24 the discussion went further: the AI pointed out that my solution ignored the balancing of the remaining compartments, and that disagreement is why the `Notes` section for that day is the longest in the repository.

The third is `count-dissections.go` in Day 24. That program exists because I solved the puzzle incorrectly, got the right answer anyway, and wanted to understand why. Rather than write the exhaustive checker myself, I described what it should do and let an AI generate it. It worked on the first prompt. This is the only piece of code in the repository that I did not write.

What I did not use AI for is the solutions themselves. Every `part1.go`, `part2.go`, and `dayNN.go` was written by hand, and every algorithm choice was mine. The AI can generate code that passes the tests—if the goal had been to solve everything as fast as possible, it could have done most of it, with corrections here and there. But that was not the goal. I wanted to do it myself, and the write-ups are the record of that work. Since the explanations are the point of this repository, the code had to be mine.
