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

	mulRE     = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	enabledRE = regexp.MustCompile(`(?s)do\(\).*?don't\(\)`)
)

type day03 struct {
	day.DayInput
}

func NewDay03(opts ...day.Option) day03 {
	return day03{day.NewDayInput(path, opts...)}
}

func sumMuls(s string) int {
	result := 0

	muls := mulRE.FindAllStringSubmatch(s, -1)
	for _, mul := range muls {
		result += conv.MustAtoi(mul[1]) * conv.MustAtoi(mul[2])
	}

	return result
}

func (d day03) Part1() int {
	input := d.ReadInput()

	return sumMuls(string(input))
}

func (d day03) Part2() int {
	input := d.ReadInput()
	enabledMuls := enabledRE.FindAllString("do()"+string(input)+"don't()", -1)

	result := 0

	for _, muls := range enabledMuls {
		result += sumMuls(muls)
	}

	return result
}

func main() {
	d := NewDay03(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
