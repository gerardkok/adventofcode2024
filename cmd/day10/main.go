package main

import (
	"maps"
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

type grid [][]int

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

func parseGrid(input []string) grid {
	result := make(grid, len(input)+2)
	result[0] = slices.Repeat([]int{-1}, len(input[0])+2)
	result[len(input)+1] = slices.Repeat([]int{-1}, len(input[0])+2)
	for x, line := range input {
		result[x+1] = slices.Repeat([]int{-1}, len(input[0])+2)
		for y, c := range line {
			result[x+1][y+1] = int(c - '0')
		}
	}
	return result
}

func (g grid) height(p position) int {
	return g[p.x][p.y]
}

func (p position) to(d direction) position {
	return position{p.x + d.dx, p.y + d.dy}
}

func (g grid) ends(p position) map[position]struct{} {
	if g.height(p) == 9 {
		return map[position]struct{}{p: {}}
	}

	result := make(map[position]struct{})
	for _, dir := range directions {
		next := p.to(dir)
		if g.height(next) == g.height(p)+1 {
			n := g.ends(next)
			maps.Copy(result, n)
		}
	}

	return result
}

func (g grid) rating(p position) int {
	if g.height(p) == 9 {
		return 1
	}

	result := 0
	for _, dir := range directions {
		next := p.to(dir)
		if g.height(next) == g.height(p)+1 {
			result += g.rating(next)
		}
	}

	return result
}

func (d day10) Part1() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	result := 0

	for x := range len(grid) {
		for y := range len(grid[0]) {
			p := position{x, y}
			if grid.height(p) == 0 {
				n := grid.ends(p)
				result += len(n)
			}
		}
	}

	return result
}

func (d day10) Part2() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	result := 0

	for x := range len(grid) {
		for y := range len(grid[0]) {
			p := position{x, y}
			if grid.height(p) == 0 {
				result += grid.rating(p)
			}
		}
	}

	return result
}

func main() {
	d := NewDay10(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
