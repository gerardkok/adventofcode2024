package main

import (
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day02 struct {
	day.DayInput
}

type report []int

func NewDay02(inputFile string) Day02 {
	return Day02{day.DayInput(inputFile)}
}

func parseReport(line string) report {
	fields := strings.Fields(line)
	report := make(report, len(fields))

	for i, f := range fields {
		level, _ := strconv.Atoi(f)
		report[i] = level
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

func (r report) isSafe() bool {
	return r.isAscending() || r.reverse().isAscending()
}

func (r report) isAlmostSafe() bool {
	for i := range len(r) {
		d := slices.Concat(r[:i], r[i+1:])
		if d.isSafe() {
			return true
		}
	}
	return false
}

func (d Day02) Part1() int {
	lines, _ := d.ReadLines()

	safe := 0

	for _, line := range lines {
		report := parseReport(line)
		if report.isSafe() {
			safe++
		}
	}

	return safe
}

func (d Day02) Part2() int {
	lines, _ := d.ReadLines()

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
	d := NewDay02(filepath.Join(projectpath.Root, "cmd", "day02", "input.txt"))

	day.Solve(d)
}
