package main

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type connection interface {
	eval(circ circuit) signal
}

type source interface {
	connection
	isSource()
}

type signal uint16

func (signal) isSource() {}

type wire string

func (wire) isSource() {}

type unary struct {
	op string
	x  source
}

type binary struct {
	op   string
	x, y source
}

type element struct {
	sig   signal
	isSet bool
	conn  connection
}
type circuit map[wire]*element

func (s signal) eval(_ circuit) signal {
	return s
}

func (w wire) eval(circ circuit) signal {
	return circ[w].eval(circ)
}

func (u unary) eval(circ circuit) signal {
	x := u.x.eval(circ)
	switch u.op {
	case "NOT":
		return ^x
	default:
		return 0
	}
}

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

func (e *element) eval(circ circuit) signal {
	if e.isSet {
		return e.sig
	}
	e.sig = e.conn.eval(circ)
	e.isSet = true
	return e.sig
}

func parseSource(s string) source {
	if num, err := strconv.Atoi(s); err == nil {
		return signal(num)
	}
	return wire(s)
}

var (
	reId     = regexp.MustCompile(`^([0-9a-z]+) -> ([a-z]+)$`)
	reUnary  = regexp.MustCompile(`^([A-Z]+) ([0-9a-z]+) -> ([a-z]+)$`)
	reBinary = regexp.MustCompile(`^([0-9a-z]+) ([A-Z]+) ([0-9a-z]+) -> ([a-z]+)$`)
)

func parseConnection(s string) (connection, wire) {
	var m []string

	m = reId.FindStringSubmatch(s)
	if m != nil {
		return parseSource(m[1]), wire(m[2])
	}

	m = reUnary.FindStringSubmatch(s)
	if m != nil {
		return unary{m[1], parseSource(m[2])}, wire(m[3])
	}

	m = reBinary.FindStringSubmatch(s)
	if m != nil {
		return binary{m[2], parseSource(m[1]), parseSource(m[3])}, wire(m[4])
	}

	panic("invalid input: " + s)
}

func (c circuit) String() string {
	keys := slices.Sorted(maps.Keys(c))
	if len(keys) == 0 {
		return "{}"
	}

	var b strings.Builder
	b.WriteString("{\n")
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("   %s: ", k))
		if c[k] != nil {
			b.WriteString(fmt.Sprintf("%v", *c[k]))
		} else {
			b.WriteString("nil")
		}
		b.WriteString("\n")
	}
	b.WriteString("}")
	return b.String()
}

func parseCircuit(lines []string) circuit {
	circ := make(circuit)
	for _, line := range lines {
		conn, wire := parseConnection(line)
		circ[wire] = &element{conn: conn}
	}
	return circ
}
