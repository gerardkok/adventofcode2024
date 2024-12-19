package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day19 struct {
	patterns, designs []string
}

type memo1 map[string]bool

type memo2 map[string]int

func parseInput(input day.DayInput) ([]string, []string) {
	lines := input.ReadInput()

	parts := bytes.Split(lines, []byte{'\n', '\n'})

	patterns := strings.Split(string(parts[0]), ", ")
	designs := strings.Split(strings.TrimSpace(string(parts[1])), "\n")

	return patterns, designs
}

func NewDay19(opts ...day.Option) day19 {
	input := day.NewDayInput(path, opts...)

	patterns, designs := parseInput(input)

	return day19{patterns, designs}
}

func (m *memo1) memoPossible(design string, patterns []string) bool {
	if _, ok := (*m)[design]; !ok {
		(*m)[design] = m.possible(design, patterns)
	}

	return (*m)[design]
}

func (m *memo1) possible(design string, patterns []string) bool {
	if design == "" {
		return true
	}

	for _, pattern := range patterns {
		if strings.HasPrefix(design, pattern) && m.memoPossible(design[len(pattern):], patterns) {
			return true
		}
	}

	return false
}

func (m *memo2) memoCountWays(design string, patterns []string) int {
	if _, ok := (*m)[design]; !ok {
		(*m)[design] = m.countWays(design, patterns)
	}

	return (*m)[design]
}

func (m *memo2) countWays(design string, patterns []string) int {
	if design == "" {
		return 1
	}

	result := 0

	for _, pattern := range patterns {
		if strings.HasPrefix(design, pattern) {
			result += m.memoCountWays(design[len(pattern):], patterns)
		}
	}

	return result
}

func (d day19) Part1() int {
	sum := 0

	memo := make(memo1)

	for _, design := range d.designs {
		if memo.possible(design, d.patterns) {
			sum++
		}
	}

	return sum
}

func (d day19) Part2() int {
	sum := 0

	memo := make(memo2)

	for _, design := range d.designs {
		sum += memo.countWays(design, d.patterns)
	}

	return sum
}

func main() {
	d := NewDay19(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
