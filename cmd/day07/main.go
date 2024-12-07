package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day07 struct {
	day.DayInput
}

func NewDay07(inputFile string) Day07 {
	return Day07{day.DayInput(inputFile)}
}

func parseLine(line string) (int, []int) {
	r, o, _ := strings.Cut(line, ":")
	result, _ := strconv.Atoi(r)

	ops := strings.Split(o, " ")

	operands := make([]int, len(ops))

	for i, op := range ops {
		p, _ := strconv.Atoi(op)
		operands[i] = p
	}

	return result, operands
}

func doublePipe(a, b int) int {
	r, _ := strconv.Atoi(fmt.Sprintf("%d%d", a, b))
	return r
}

func isValid(result, intermediate int, operands []int) bool {
	if len(operands) == 0 {
		return result == intermediate
	}

	operand := operands[0]
	remaining := operands[1:]
	return isValid(result, intermediate+operand, remaining) ||
		isValid(result, intermediate*operand, remaining)
}

func isValidPart2(result, intermediate int, operands []int) bool {
	if len(operands) == 0 {
		return result == intermediate
	}

	operand := operands[0]
	remaining := operands[1:]
	return isValidPart2(result, intermediate+operand, remaining) ||
		isValidPart2(result, intermediate*operand, remaining) ||
		isValidPart2(result, doublePipe(intermediate, operand), remaining)
}

func (d Day07) Part1() int {
	lines, _ := d.ReadLines()

	sum := 0

	for _, line := range lines {
		result, operands := parseLine(line)
		if isValid(result, 0, operands) {
			sum += result
		}
	}

	return sum
}

func (d Day07) Part2() int {
	lines, _ := d.ReadLines()

	sum := 0

	for _, line := range lines {
		result, operands := parseLine(line)
		if isValidPart2(result, 0, operands) {
			sum += result
		}
	}

	return sum
}

func main() {
	d := NewDay07(filepath.Join(projectpath.Root, "cmd", "day07", "input.txt"))

	day.Solve(d)
}
