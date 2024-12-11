package main

import (
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

func (g grid) walkNeighbours(p position, fn func(position)) {
	for _, dir := range directions {
		next := p.to(dir)
		if g.height(next) == g.height(p)+1 {
			fn(next)
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

		g.walkNeighbours(p, func(n position) {
			todo = append(todo, n)
		})
	}

	return result
}

func (g grid) rating(p position) int {
	if g.peak(p) {
		return 1
	}

	result := 0

	g.walkNeighbours(p, func(n position) {
		result += g.rating(n)
	})

	return result
}

func (g grid) walk(fn func(position)) {
	for x := range len(g) {
		for y := range len(g[0]) {
			p := position{x, y}
			fn(p)
		}
	}
}

func (d day10) Part1() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	result := 0

	grid.walk(func(p position) {
		if grid.trailhead(p) {
			result += grid.peaks(p)
		}
	})

	return result
}

func (d day10) Part2() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	result := 0

	grid.walk(func(p position) {
		if grid.trailhead(p) {
			result += grid.rating(p)
		}
	})

	return result
}

func main() {
	d := NewDay10(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
