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

type day05 struct {
	day.DayInput
}

type rules map[int]map[int]struct{}

type page []int

func NewDay05(opts ...day.Option) day05 {
	return day05{day.NewDayInput(path, opts...)}
}

func parseRules(lines []string) rules {
	result := make(rules)

	for _, line := range lines {
		s, t, _ := strings.Cut(line, "|")
		x := conv.MustAtoi(s)
		y := conv.MustAtoi(t)
		if _, ok := result[x]; !ok {
			result[x] = make(map[int]struct{})
		}
		result[x][y] = struct{}{}
	}

	return result
}

func parsePages(lines []string) []page {
	result := make([]page, len(lines))

	for i, line := range lines {
		r := strings.Split(line, ",")
		for _, p := range r {
			result[i] = append(result[i], conv.MustAtoi(p))
		}
	}

	return result
}

func cmp(rules rules) func(a, b int) int {
	return func(a, b int) int {
		if _, ok := rules[a][b]; ok {
			return -1
		} else if _, ok := rules[b][a]; ok {
			return 1
		} else {
			return 0
		}
	}
}

func (p page) middle() int {
	return p[len(p)/2]
}

func (p page) isSorted(rules rules) bool {
	return slices.IsSortedFunc(p, cmp(rules))
}

func (p page) sort(rules rules) {
	slices.SortFunc(p, cmp(rules))
}

func parseInput(lines []string) (rules, []page) {
	result := [2][]string{make([]string, 0), make([]string, 0)}

	i := 0
	for _, line := range lines {
		if len(line) == 0 {
			i++
			continue
		}

		result[i] = append(result[i], line)
	}

	r := parseRules(result[0])
	p := parsePages(result[1])
	return r, p
}

func (d day05) Part1() int {
	lines := d.ReadLines()
	rules, pages := parseInput(lines)

	sum := 0

	for _, page := range pages {
		if page.isSorted(rules) {
			sum += page.middle()
		}
	}

	return sum
}

func (d day05) Part2() int {
	lines := d.ReadLines()
	rules, pages := parseInput(lines)

	sum := 0

	for _, page := range pages {
		if !page.isSorted(rules) {
			page.sort(rules)
			sum += page.middle()
		}
	}

	return sum
}

func main() {
	d := NewDay05(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
