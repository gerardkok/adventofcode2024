package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)

	boolValue = map[bool]int{false: 0, true: 1}
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

func (m *memo1) possible(design string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.HasPrefix(design, pattern) {
			tail := design[len(pattern):]

			if _, ok := (*m)[tail]; !ok {
				(*m)[tail] = m.possible(tail, patterns)
			}

			if (*m)[tail] {
				return true
			}
		}
	}

	return false
}

func (m *memo2) countWays(design string, patterns []string) int {
	result := 0

	for _, pattern := range patterns {
		if strings.HasPrefix(design, pattern) {
			tail := design[len(pattern):]

			if _, ok := (*m)[tail]; !ok {
				(*m)[tail] = m.countWays(tail, patterns)
			}

			result += (*m)[tail]
		}
	}

	return result
}

func (d day19) Part1() int {
	memo := memo1{"": true}

	return conv.SumFunc(d.designs, func(design string) int {
		return boolValue[memo.possible(design, d.patterns)]
	})
}

func (d day19) Part2() int {
	memo := memo2{"": 1}

	return conv.SumFunc(d.designs, func(design string) int {
		return memo.countWays(design, d.patterns)
	})
}

func main() {
	d := NewDay19(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
