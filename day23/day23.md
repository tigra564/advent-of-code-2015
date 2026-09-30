# Day 23: Opening the Turing Lock

## The Puzzle

You are given a small program in a made-up assembly language, and a machine with two registers, `a` and `b`, that start at 0. The instruction set is tiny: `hlf` halves a register, `tpl` triples it, `inc` increments it, `jmp` jumps to a relative offset, and `jie` and `jio` jump conditionally if a register is even or equal to 1.

Part One asks for the value in register `b` when the program halts, starting with both registers at 0.

Part Two asks the same, but with register `a` starting at 1. The program is the same; only the initial state changes.

## Input

About fifty lines, one per instruction:

```
jio a, +18
inc a
tpl a
inc a
tpl a
...
```

Each line is one instruction, in one of six forms. Three of them take a register and nothing else (`hlf`, `tpl`, `inc`), one takes only an offset (`jmp`), and two take a register and an offset (`jie`, `jio`). The offsets are signed integers, always written with an explicit `+` or `-`.

The program is small—fifty instructions—and there is nothing in its structure that constrains the algorithm. There are no loops that cannot be unrolled, no self-modifying code, no input. The machine is deterministic: given the initial state, the entire execution is fixed, and the only question is how long it runs and what it leaves in `b`.

The same program is used for both parts. The only difference is the initial value of register `a`, which is 0 in Part One and 1 in Part Two.

## The Approach

This is an emulator: a small virtual machine with two registers, a program counter, and six instructions. There is no search, no optimisation, no state space—just a `for` loop that reads instructions and applies them until the program counter runs off the end.

The state is a struct with two parts: an array of general-purpose registers and the program counter.

```go
type registers struct {
	// General purpose registers: empty (for ops without registers), a, b
	gp [3]int
	// program counter
	pc int
}
```

The array has three slots, not two, and the zeroth slot is unused.

Each instruction is compiled into a small struct: a function pointer, a register number, and an offset.

```go
type microCode func(regs *registers, regNum int, offset int) int

type operation struct {
	microcode microCode
	regNum    int
	offset    int
}
```

The `microCode` function takes the whole machine state, the register number, and the offset, and returns the number of instructions to advance the program counter by. For most instructions this is 1; for jumps it is the offset. This design means the main loop does not need to know which instruction it is executing—it just calls the function and adds the result to `pc`.

The six instructions are six functions with the same signature:

```go
func opHlf(regs *registers, regNum int, offset int) int {
	regs.gp[regNum] /= 2
	return 1
}

func opTpl(regs *registers, regNum int, offset int) int {
	regs.gp[regNum] *= 3
	return 1
}

func opInc(regs *registers, regNum int, offset int) int {
	regs.gp[regNum]++
	return 1
}

func opJmp(regs *registers, regNum int, offset int) int {
	return offset
}

func opJie(regs *registers, regNum int, offset int) int {
	if 0 == regs.gp[regNum] % 2 {
		return offset
	}
	return 1
}

func opJio(regs *registers, regNum int, offset int) int {
	if 1 == regs.gp[regNum] {
		return offset
	}
	return 1
}
```

The conditional jumps are the only ones that look at a register value; `jie` checks for evenness, `jio` checks for equality with 1. The offset is returned only when the condition holds; otherwise the instruction falls through to the next one.

The parsing is a single regex, and the instruction name is mapped to the corresponding function by a `switch`:

```go
var reOp = regexp.MustCompile(`^([a-z]+) (?:([a-z]?)(?:, )?)([-+][0-9]+)?$`)
```

The regex captures three groups: the mnemonic, an optional register letter, and an optional signed offset. The register letter is mapped to 1 for `a` and 2 for `b`, with 0 for the absent case. The offset is parsed by `strconv.Atoi`, which handles the leading `+` or `-` without complaint.

The main loop is the whole execution model:

```go
func run(regs *registers, ops []operation, trace bool) {
	lastOp := len(ops)

	for regs.pc < lastOp {
		op := ops[regs.pc]
		regs.pc += op.microcode(regs, op.regNum, op.offset)
		if trace {
			fmt.Printf("%s %2d %3d => %v\n", funcName(op.microcode), op.regNum, op.offset, *regs)
		}
	}
}
```

The loop condition is `regs.pc < lastOp`. Since the program counter starts at 0 and all offsets are bounded, the counter cannot go negative on a valid input, and the program terminates when it runs past the last instruction. The `trace` flag is there for debugging: when set, the loop prints each instruction and the resulting state. It is off in both parts.

The `funcName` helper is worth a note, because it is not the kind of code that usually appears in a puzzle solution:

```go
func funcName(f any) string {
	if f == nil {
		return "<nil>"
	}
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
```

This uses `reflect` and `runtime` to get the name of the function that implements the instruction, so that the trace output shows `main.opHlf` instead of a pointer. It is only used when `trace` is true, and both parts call `run` with `trace = false`. It is a debugging aid that survived into the final code because it costs nothing when disabled and is useful when something goes wrong.

### Part One

Part One starts with both registers at 0 and reads the result from register `b`:

```go
regs := registers{}
ops := compile(lines)
run(&regs, ops, false)

result := regs.gp[2]
```

Register `b` is at index 2, following the same convention as the instruction encoding: 0 for absent, 1 for `a`, 2 for `b`.

See [part1.go](part1.go).

### Part Two

Part Two sets register `a` to 1 before running the program:

```go
regs := registers{}
regs.gp[1] = 1
ops := compile(lines)
run(&regs, ops, false)
```

Everything else is the same. The program is identical, the result is read from the same register, and the only difference is the initial value.

See [part2.go](part2.go).

## Complexity

Let *I* be the number of instructions in the program, and *S* the number of steps the machine executes before halting.

- **Time:** *O*(*I*) for compilation and *O*(*S*) for execution. Each instruction is parsed once and mapped to a function pointer, which is a fixed amount of work. Execution then proceeds one instruction at a time, and the program counter advances by at least 1 on every step—jumps are always to a different instruction, never to the same one, and the only way to loop is to jump backwards to a point that has not yet been reached with the same register state. The program terminates, so *S* is finite, but it can be arbitrarily large: a loop with a counter can run for as many iterations as the counter allows, and the counter itself can grow without bound through `tpl` and `inc`. We will return to this in `Associations`.

- **Space:** *O*(*I*) for the compiled program—one `operation` struct per instruction—and *O*(1) for the machine state: two registers and a program counter. The `run` function's trace output is disabled, so no allocation happens during execution.

## Associations

There are no classical algorithms hiding in this puzzle—it is a small emulator, and the interesting parts are the design choices, not the mathematics. What follows are two notes on where the ideas come from and what the program actually computes.

- **Building a computer from scratch.** The design of this emulator—registers, a program counter, an instruction decoder, a dispatch loop—is a simplified version of how real CPUs work, and there is no better place to see the full picture than Ben Eater's video series on building an 8-bit breadboard computer. It starts with a clock module made of 555 timers and ends with a machine that runs programs from an EEPROM. Along the way it covers registers, the bus, the ALU, microcode, and the control logic that ties everything together. The series is long, but it is the closest thing to *Code* by Charles Petzold in video form, and it shows every wire. If you have ever wondered what happens between `pc++` and the next instruction, this is where to look. See [Ben Eater's 8-bit computer](https://www.youtube.com/watch?v=HyznrdDSSGM&list=PLowKtXNTBypGqImE405J2565dvjafglHU).

- **The Collatz conjecture.** The program in this puzzle is not arbitrary: it computes the Collatz sequence for the value in register `a` and leaves the length of the sequence in register `b`. The Collatz conjecture states that no matter what positive integer you start with, repeatedly applying "if even, halve; if odd, triple and add one" eventually reaches 1. It has been verified for numbers up to 2⁶⁸, but no proof exists, and Paul Erdős famously said that "mathematics is not yet ready for such problems". The program here is a direct encoding of that process, and it terminates because the starting value is small enough to have been checked. See [Collatz conjecture](https://en.wikipedia.org/wiki/Collatz_conjecture).

## Notes

The `trace` flag and the `funcName` helper are a small gift to myself. The flag turns the run loop into a disassembler-style trace: each step prints the name of the microcode function, its operands, and the full register state after the step. It is off in both parts, but it is the first thing I would turn on if something went wrong. Writing a debug mode that is never enabled in production is usually a waste—but for an emulator, where the state is small and the steps are discrete, the trace costs almost nothing and pays for itself the first time a program does not do what you expect.

The `funcName` function is the only place in the repository where `reflect` and `runtime` are used. It converts a function value into its fully qualified name, so that the trace shows `main.opHlf` instead of a hexadecimal pointer. The technique relies on `reflect.ValueOf(f).Pointer()`, which returns the code address of the function, and `runtime.FuncForPC`, which looks up the name for that address. It works for any function in the program, not just the microcode handlers, and it would work for methods as well if we ever needed to trace those. The cost is a small amount of reflection at startup, and the benefit is that the trace is readable without a symbol table.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | ~3 ms |
| 2 | ~2 ms |

Both parts compile the same fifty instructions and run them on the same emulator. The difference between them is one assignment—register `a` starts at 1 instead of 0—and the resulting program takes a slightly different path. The measured times are dominated by process startup and input parsing; the emulation itself is a few thousand steps, each one a function call and a handful of arithmetic operations. There is nothing to optimise here.

---

[← Previous: Day 22](day22.md) · [Up: Contents](../README.md) · [Next: Day 24 →](day24.md)
