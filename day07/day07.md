# Day 7: Some Assembly Required

## The Puzzle

A circuit is described by a list of connections, one per line. Each line assigns a value to a wire, using a small set of bitwise operations: `AND`, `OR`, `NOT`, `LSHIFT`, `RSHIFT`. Operands are either other wires or 16-bit constants.

Part One asks for the signal on wire `a` after the whole circuit settles. Part Two overrides wire `b` with the answer from Part One and asks for the signal on wire `a` again.

## Input

Several hundred lines, each describing one connection. The forms are:

```
123 -> x
x AND y -> z
x OR y -> z
x LSHIFT 2 -> f
x RSHIFT 2 -> g
NOT x -> h
```

The left-hand side is an expression, the right-hand side is the target wire. Operands are either decimal constants (`0` to `65535`) or wire names, which are some lowercase letters—not necessarily a single character.

The circuit is acyclic. The puzzle says: "A gate provides no signal until all of its inputs have a signal." With a cycle, no gate in the cycle would ever get a signal, and the circuit would be undefined. So we can assume the circuit is a DAG, and we do not need to check for cycles. This is the one place where a structural property of the input matters for the algorithm.

## The Approach

The circuit is a DAG: each wire is a node, and each connection is an edge (or several) from its inputs to its output. Evaluating wire `a` means a depth-first traversal of that DAG, with memoisation so that each wire is computed once.

The structure I use is not an invention. It is the expression evaluator from *The Go Programming Language* by Donovan and Kernighan, adapted to this problem. The book presents a small calculator that parses arithmetic expressions into a tree and evaluates them; each node type carries its own `eval` method, and the evaluator does not know what operations exist. The same recipe fits here almost without change. I have not copied it exactly—the parsing is different, and the memoisation is added—but the shape is the book's, and the credit is theirs.

Each connection is an object that knows how to evaluate itself given the circuit. Wires and constants are also connections, but trivial ones. The evaluator does not know what operations exist; it just calls `eval` on whatever it has.

### Types

Two interfaces. A `connection` is anything that can produce a signal:

```go
type connection interface {
	eval(circ circuit) signal
}
```

A `source` is a connection that can appear on the left-hand side of an operation—a wire or a constant:

```go
type source interface {
	connection
	isSource()
}
```

The `isSource` method is a marker. It does not do anything; it exists so that the type system can tell a `wire` or a `signal` apart from a `unary` or `binary` at compile time. This idiom is not in the book—the calculator there does not need to distinguish operand types—so it is worth a sentence of explanation.

In Go, there is no inheritance, and there is no `sealed` keyword. If we want a set of types that belong together and that nobody outside the package can extend, the standard trick is to add an unexported method. Only types in this package can implement `isSource`, because only they can refer to the method. The compiler then rejects any attempt to pass an unrelated type where a `source` is expected. This gives us the closest thing Go has to a closed sum type, without a single line of runtime code. The method body is empty; its existence is the entire point.

A `signal` is a 16-bit value, a `wire` is a name, and `unary` and `binary` are the two compound operation types:

```go
type signal uint16
type wire string

type unary struct {
	op string
	x  source
}

type binary struct {
	op   string
	x, y source
}
```

Finally, an `element` holds the memoised value of a wire, plus the connection that defines it:

```go
type element struct {
	sig   signal
	isSet bool
	conn  connection
}

type circuit map[wire]*element
```

The circuit maps a wire name to its element. The elements are stored as pointers because we mutate `sig` and `isSet` during evaluation, and we want those mutations to be visible across the whole circuit.

### Evaluation

The evaluation is depth-first. A `signal` evaluates to itself. A `wire` looks itself up in the circuit and asks its element to evaluate. An `element` checks the memo flag first: if the value is already known, it returns it immediately. Only if the value is not known does the element ask its connection to evaluate, store the result, and return it.

The order matters. The memo check comes first, before any other logic. This is what makes the evaluation linear in the size of the circuit: a wire that has already been computed is never recomputed, no matter how many other wires depend on it.

```go
func (e *element) eval(circ circuit) signal {
	if e.isSet {
		return e.sig
	}
	e.sig = e.conn.eval(circ)
	e.isSet = true
	return e.sig
}
```

The `unary` and `binary` nodes do the actual arithmetic. They are only reached when the memo check has failed, which is exactly once per wire.

```go
func (u unary) eval(circ circuit) signal {
	x := u.x.eval(circ)
	switch u.op {
	case "NOT":
		return ^x
	default:
		return 0
	}
}
```

```go
func (b binary) eval(circ circuit) signal {
	x := b.x.eval(circ)
	y := b.y.eval(circ)
	switch b.op {
	case "AND":
		return x & y
	case "OR":
		return x | y
	case "LSHIFT":
		return x << y
	case "RSHIFT":
		return x >> y
	default:
		return 0
	}
}
```

The `default` branches are unreachable for valid input—the parser only produces these operation names—but Go requires the switch to return a value, and I prefer an explicit `return 0` over a panic. If the parser ever produces an unknown operation, the circuit will silently produce zero for that wire, and the answer will be wrong. This is the one place where I traded safety for brevity.

## Complexity

Let *n* be the number of connections in the circuit.

- **Time:** *O*(*n*) for the first evaluation of a wire. Each wire is computed exactly once, thanks to the memoisation in `element.eval`. The recursive descent visits each node at most once, and each edge (the inputs of a node) at most once. The depth of recursion is bounded by the longest path in the DAG, which is *O*(*n*).

  Part Two evaluates wire `a` twice: once in the original circuit, then again after overriding `b`. The second evaluation is also *O*(*n*), because we reset the memo flags and recompute everything. So the total is still *O*(*n*).

- **Space:** *O*(*n*) for the circuit map and the recursion stack. The map holds one element per wire, and each element holds one connection.

## Associations

- **Lambda calculus and the Y combinator.** A circuit with a cycle would be a self-referential definition: wire `a` depends on wire `b`, which depends on `a`. In the lambda calculus, such a definition is not inherently meaningless—it can be expressed, and the Y combinator exists precisely to allow it. A cyclic circuit would be an attempt at a recursive definition of a signal, not a syntax error. The puzzle rules this out by fiat: "A gate provides no signal until all of its inputs have a signal." What looks like a technical restriction is really a statement about the model: we are not allowed to express recursion, only finite acyclic composition.

- **Lem's *The Star Diaries* and the sepulki.** The puzzle's insistence on acyclicity is reminiscent of Stanisław Lem's story about Ijon Tichy and the sepulki—objects that cannot be defined without reference to themselves. Tichy looks them up in the *Cosmic Encyclopedia* and finds a perfectly circular definition:

  > SEPULKI – an important element of the civilization of the Ardrites (see) from the planet Enteropia (see). See SEPULKARIA.
  >
  > I followed this advice and read:
  >
  > SEPULKARIA – devices for sepuling (see).
  >
  > I looked under "Sepuling"; there it said:
  >
  > SEPULING – an activity of the Ardrites (see) from the planet Enteropia (see). See SEPULKI.

  For readers of the Russian and Polish originals, the terms are more familiar as *sepulki*, *sepulkarii*, and *sepulenie*. The English translation by Michael Kandel renders them as *scrupts*, *scrupture*, and *scruption*—a clever choice, because the English words carry the same sense of a closed, self-referential loop. But the Polish *sepulki* is the one Lem actually wrote, and it is the one that stuck as a metaphor for circular definitions.

  In a cyclic circuit, every gate would be waiting for every other gate, and the whole construction would be a sepulka: defined, described, and completely inert. The puzzle avoids the problem by requiring a DAG, which is a polite way of saying "no sepulki allowed".

## Notes

The instruction set is small and well-behaved. `AND`, `OR`, `NOT`, `LSHIFT`, and `RSHIFT` are all defined on 16-bit unsigned values, and the puzzle guarantees that every wire is defined exactly once and that the circuit is acyclic. There is no need for bounds checks, no need to guard against overflows, no need to handle a wire that is read before it is written. The evaluator can be simple because the problem is simple.

This is not true of every puzzle. In Day 22, the state space is large enough that a naive search must be pruned and memoised, and the code carries several guards that exist only to keep the search tractable. The trade-off there is different: more code, more checks, more places for bugs.

I will use the same simple structure in Day 23, where the instruction set is even smaller and the machine model is even more constrained. The pattern is worth recognising: when the problem guarantees a well-behaved input, the code can be straightforward, and adding defensive checks only obscures the algorithm. When the guarantees are weak, the checks are the algorithm.

## Benchmarks

Measured on a compiled binary (`go build`), ten runs, *user* time.

| Part | Time (user, per run) |
|------|----------------------|
| 1 | 2.7 ms |
| 2 | 1.7 ms |

Part One parses the circuit and evaluates wire `a` once. Part Two parses the same circuit, evaluates `a`, resets the memo flags, overrides `b`, and evaluates `a` again. So Part Two does the same work plus one extra evaluation—but the measurement shows Part Two as *faster*, which is backwards.

The explanation is that the two parts are separate binaries. The measurement is dominated by process startup and parsing, not by the evaluation itself. The evaluation is a few hundred nodes deep and takes microseconds; the parsing and setup take milliseconds. The 1 ms difference between the two parts is noise, not signal.

---

[← Previous: Day 6](day06.md) · [Up: Contents](../README.md) · [Next: Day 8 →](day08.md)
