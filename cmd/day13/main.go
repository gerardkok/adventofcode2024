package main

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)

	machineRE = regexp.MustCompile(`(?s)Button A: X\+(\d+), Y\+(\d+).*?Button B: X\+(\d+), Y\+(\d+).*?Prize: X=(\d+), Y=(\d+)`)
)

type day13 struct {
	day.DayInput
}

func NewDay13(opts ...day.Option) day13 {
	return day13{day.NewDayInput(path, opts...)}
}

func divmod(numerator, denominator int) (int, int) {
	quotient := numerator / denominator
	remainder := numerator % denominator
	return quotient, remainder
}

func tokens(xa, ya, xb, yb, xp, yp int) (int, int) {
	b, r := divmod(xa*yp-xp*ya, xa*yb-xb*ya)
	if r != 0 {
		return 0, 0
	}

	a, r := divmod(yp-yb*b, ya)
	if r != 0 {
		return 0, 0
	}

	return a, b
}

func (d day13) Part1() int {
	input := d.ReadInput()

	result := 0

	machines := machineRE.FindAllStringSubmatch(string(input), -1)
	for _, machine := range machines {
		xa := conv.MustAtoi(machine[1])
		ya := conv.MustAtoi(machine[2])
		xb := conv.MustAtoi(machine[3])
		yb := conv.MustAtoi(machine[4])
		xp := conv.MustAtoi(machine[5])
		yp := conv.MustAtoi(machine[6])
		a, b := tokens(xa, ya, xb, yb, xp, yp)
		result += 3*a + b
	}

	return result
}

func (d day13) Part2() int {
	const prizeAddition = 10_000_000_000_000

	input := d.ReadInput()

	result := 0

	machines := machineRE.FindAllStringSubmatch(string(input), -1)
	for _, machine := range machines {
		xa := conv.MustAtoi(machine[1])
		ya := conv.MustAtoi(machine[2])
		xb := conv.MustAtoi(machine[3])
		yb := conv.MustAtoi(machine[4])
		xp := conv.MustAtoi(machine[5]) + prizeAddition
		yp := conv.MustAtoi(machine[6]) + prizeAddition
		a, b := tokens(xa, ya, xb, yb, xp, yp)
		result += 3*a + b
	}

	return result
}

func main() {
	d := NewDay13(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
