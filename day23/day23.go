package main

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
)

type registers struct {
	// General purpose registers: empty (for ops without registers), a, b
	gp [3]int
	// program counter
	pc int
}

type microCode func(regs *registers, regNum int, offset int) int

type operation struct {
	microcode microCode
	regNum    int
	offset    int
}

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

func funcName(f any) string {
	if f == nil {
		return "<nil>"
	}
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

var reOp = regexp.MustCompile(`^([a-z]+) (?:([a-z]?)(?:, )?)([-+][0-9]+)?$`)

func compile(lines []string) []operation {
	ops := make([]operation, 0, len(lines))

	for _, line := range lines {
		m := reOp.FindStringSubmatch(line)

		var reg int
		switch m[2] {
		case "a": reg = 1
		case "b": reg = 2
		default: reg = 0
		}

		offset, _ := strconv.Atoi(m[3])

		var mcode microCode
		switch m[1] {
		case "hlf": mcode = opHlf
		case "tpl": mcode = opTpl
		case "inc": mcode = opInc
		case "jmp": mcode = opJmp
		case "jie": mcode = opJie
		case "jio": mcode = opJio
		}

		op := operation{mcode, reg, offset}
		ops = append(ops, op)
	}

	return ops
}

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
