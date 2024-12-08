package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day07b struct {
	day.DayInput
}

type operator func(int, int) (int, bool)

func NewDay07b(inputFile string) Day07b {
	return Day07b{day.DayInput(inputFile)}
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

func divmod(numerator, denominator int) (int, int) {
	quotient := numerator / denominator
	remainder := numerator % denominator
	return quotient, remainder
}

func sub(a, b int) (int, bool) {
	return a - b, a >= b
}

func div(a, b int) (int, bool) {
	if b == 0 {
		return 0, a == 0
	}
	q, r := divmod(a, b)
	return q, r == 0
}

func trimSuffix(a, b int) (int, bool) {
	pow := 10
	for b >= pow {
		pow *= 10
	}
	q, r := divmod(a-b, pow)
	return q, r == 0
}

func valid(target int, operands []int, operators []operator) bool {
	if len(operands) == 1 {
		return target == operands[0]
	}

	if target == 0 && operands[len(operands)-1] == 0 {
		// if target and last operand are both 0, the equation can always
		// be made valid by putting a '*' before last operand
		return true
	}

	operand, remaining := operands[len(operands)-1], operands[:len(operands)-1]
	for _, operator := range operators {
		newTarget, possible := operator(target, operand)
		if possible && valid(newTarget, remaining, operators) {
			return true
		}
	}
	return false
}

func (d Day07b) Part1() int {
	lines, _ := d.ReadLines()

	sum := 0

	for _, line := range lines {
		target, operands := parseLine(line)
		if valid(target, operands, []operator{sub, div}) {
			sum += target
		}
	}

	return sum
}

func (d Day07b) Part2() int {
	lines, _ := d.ReadLines()

	sum := 0

	for _, line := range lines {
		target, operands := parseLine(line)
		if valid(target, operands, []operator{sub, div, trimSuffix}) {
			sum += target
		}
	}

	return sum
}

func main() {
	d := NewDay07b(filepath.Join(projectpath.Root, "cmd", "day07b", "input.txt"))

	day.Solve(d)
}