package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day07b struct {
	day.DayInput
}

type operator func(int, int) (int, bool)

func NewDay07b(opts ...day.Option) day07b {
	return day07b{day.NewDayInput(path, opts...)}
}

func parseLine(line string) (int, []int) {
	r, o, _ := strings.Cut(line, ": ")
	result := conv.MustAtoi(r)

	ops := strings.Split(o, " ")

	operands := make([]int, len(ops))

	for i, op := range ops {
		operands[i] = conv.MustAtoi(op)
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

func (d day07b) sumValid(operators []operator) int {
	lines := d.ReadLines()

	var sum atomic.Int64

	var wg sync.WaitGroup

	for _, line := range lines {
		wg.Add(1)

		go func() {
			target, operands := parseLine(line)
			if valid(target, operands, operators) {
				sum.Add(int64(target))
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return int(sum.Load())
}

func (d day07b) Part1() int {
	return d.sumValid([]operator{sub, div})
}

func (d day07b) Part2() int {
	return d.sumValid([]operator{sub, div, trimSuffix})
}

func main() {
	d := NewDay07b(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
