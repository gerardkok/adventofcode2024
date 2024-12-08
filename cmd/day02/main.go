package main

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day02 struct {
	day.DayInput
}

type report []int

func NewDay02(opts ...day.Option) day02 {
	return day02{day.NewDayInput(path, opts...)}
}

func parseReport(line string) report {
	fields := strings.Fields(line)
	report := make(report, len(fields))

	for i, f := range fields {
		report[i] = conv.MustAtoi(f)
	}

	return report
}

func (r report) reverse() report {
	slices.Reverse(r)
	return r
}

func (r report) isAscending() bool {
	for i := range len(r) - 1 {
		diff := r[i+1] - r[i]
		if diff < 1 || diff > 3 {
			return false
		}
	}
	return true
}

func (r report) isDescending() bool {
	return r.reverse().isAscending()
}

func (r report) isSafe() bool {
	return r.isAscending() || r.isDescending()
}

func (r report) isAlmostSafe() bool {
	for i := range r {
		d := slices.Concat(r[:i], r[i+1:])
		if d.isSafe() {
			return true
		}
	}
	return false
}

func (d day02) Part1() int {
	lines := d.ReadLines()

	safe := 0

	for _, line := range lines {
		report := parseReport(line)
		if report.isSafe() {
			safe++
		}
	}

	return safe
}

func (d day02) Part2() int {
	lines := d.ReadLines()

	safe := 0

	for _, line := range lines {
		report := parseReport(line)
		if report.isAlmostSafe() {
			safe++
		}
	}

	return safe
}

func main() {
	d := NewDay02(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
