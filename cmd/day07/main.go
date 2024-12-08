package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day07 struct {
	day.DayInput
}

type operator func(int, int) int

func NewDay07(opts ...day.Option) day07 {
	return day07{day.NewDayInput(path, opts...)}
}

func parseLine(line string) (int, []int) {
	r, o, _ := strings.Cut(line, ": ")
	result, _ := strconv.Atoi(r)

	ops := strings.Split(o, " ")

	operands := make([]int, len(ops))

	for i, op := range ops {
		p, _ := strconv.Atoi(op)
		operands[i] = p
	}

	return result, operands
}

func add(a, b int) int {
	return a + b
}

func mul(a, b int) int {
	return a * b
}

func concat(a, b int) int {
	pow := 10
	for b >= pow {
		pow *= 10
	}
	return a*pow + b
}

func valid(target, intermediate int, operands []int, operators []operator) bool {
	if len(operands) == 0 {
		return target == intermediate
	}

	if intermediate > target {
		return false
	}

	op, remaining := operands[0], operands[1:]
	for _, operator := range operators {
		if valid(target, operator(intermediate, op), remaining, operators) {
			return true
		}
	}
	return false
}

func (d day07) Part1() int {
	lines := d.ReadLines()

	sum := 0

	for _, line := range lines {
		target, operands := parseLine(line)
		if valid(target, operands[0], operands[1:], []operator{add, mul}) {
			sum += target
		}
	}

	return sum
}

func (d day07) Part2() int {
	lines := d.ReadLines()

	sum := 0

	for _, line := range lines {
		target, operands := parseLine(line)
		if valid(target, operands[0], operands[1:], []operator{add, mul, concat}) {
			sum += target
		}
	}

	return sum
}

func main() {
	d := NewDay07(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
