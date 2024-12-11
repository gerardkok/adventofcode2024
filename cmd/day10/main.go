package main

import (
	"iter"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day10 struct {
	day.DayInput
}

type grid [][]byte

type position struct {
	x, y int
}

type direction struct {
	dx, dy int
}

var directions = []direction{{0, 1}, {1, 0}, {-1, 0}, {0, -1}}

func NewDay10(opts ...day.Option) day10 {
	return day10{day.NewDayInput(path, opts...)}
}

func (p position) to(d direction) position {
	return position{p.x + d.dx, p.y + d.dy}
}

func parseGrid(input []string) grid {
	result := make(grid, len(input)+2)
	result[0] = slices.Repeat([]byte{'#'}, len(input[0])+2)
	result[len(input)+1] = slices.Repeat([]byte{'#'}, len(input[0])+2)
	for x, line := range input {
		result[x+1] = []byte("#" + line + "#")
	}
	return result
}

func (g grid) height(p position) byte {
	return g[p.x][p.y]
}

func (g grid) trailhead(p position) bool {
	return g.height(p) == '0'
}

func (g grid) peak(p position) bool {
	return g.height(p) == '9'
}

func (g grid) neighbours(p position) iter.Seq[position] {
	return func(yield func(position) bool) {
		for _, dir := range directions {
			next := p.to(dir)
			if g.height(next) == g.height(p)+1 && !yield(next) {
				return
			}
		}
	}
}

func (g grid) trailheads() iter.Seq[position] {
	return func(yield func(position) bool) {
		for x := range len(g) {
			for y := range len(g[0]) {
				p := position{x, y}
				if g.trailhead(p) && !yield(p) {
					return
				}
			}
		}
	}
}

func (g grid) peaks(trailhead position) int {
	todo := []position{trailhead}
	seen := make(map[position]struct{})

	result := 0

	for len(todo) > 0 {
		p := todo[0]
		todo = todo[1:]

		if _, ok := seen[p]; ok {
			continue
		}

		seen[p] = struct{}{}

		if g.peak(p) {
			result++
		}

		todo = slices.AppendSeq(todo, g.neighbours(p))
	}

	return result
}

func (g grid) rating(p position) int {
	if g.peak(p) {
		return 1
	}

	result := 0

	for neighbour := range g.neighbours(p) {
		result += g.rating(neighbour)
	}

	return result
}

func (g grid) sum(fn func(position) int) int {
	result := 0

	for trailhead := range g.trailheads() {
		result += fn(trailhead)
	}

	return result
}

func (d day10) Part1() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	return grid.sum(grid.peaks)
}

func (d day10) Part2() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	return grid.sum(grid.rating)
}

func main() {
	d := NewDay10(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
