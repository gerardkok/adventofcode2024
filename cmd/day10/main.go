package main

import (
	"bytes"
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
	grid       [][]byte
	trailheads []position
}

type position struct {
	x, y int
}

type direction struct {
	dx, dy int
}

func readGrid(d day.DayInput) [][]byte {
	input := d.ReadByteGrid()

	grid := make([][]byte, len(input)+2)

	grid[0] = bytes.Repeat([]byte{0}, len(input[0])+2)
	for i, line := range input {
		grid[i+1] = slices.Concat([]byte{0}, line, []byte{0})
	}
	grid[len(input)+1] = bytes.Repeat([]byte{0}, len(input[0])+2)

	return grid
}

func trailheads(grid [][]byte) []position {
	var result []position

	for x, line := range grid {
		for y, v := range line {
			if v == '0' {
				result = append(result, position{x, y})
			}
		}
	}

	return result
}

func NewDay10(opts ...day.Option) day10 {
	input := day.NewDayInput(path, opts...)

	grid := readGrid(input)
	trailheads := trailheads(grid)

	return day10{grid, trailheads}
}

func (p position) to(d direction) position {
	return position{p.x + d.dx, p.y + d.dy}
}

func (d day10) height(p position) byte {
	return d.grid[p.x][p.y]
}

func (d day10) peak(p position) bool {
	return d.height(p) == '9'
}

func (d day10) neighbours(p position) []position {
	var result []position

	for _, dir := range []direction{{0, 1}, {1, 0}, {-1, 0}, {0, -1}} {
		next := p.to(dir)
		if d.height(next) == d.height(p)+1 {
			result = append(result, next)
		}
	}

	return result
}

func (d day10) peaks(trailhead position) int {
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

		if d.peak(p) {
			result++
		}

		todo = append(todo, d.neighbours(p)...)
	}

	return result
}

func (d day10) rating(p position) int {
	if d.peak(p) {
		return 1
	}

	result := 0

	for _, neighbour := range d.neighbours(p) {
		result += d.rating(neighbour)
	}

	return result
}

func (d day10) Part1() int {
	return day.SumFunc(d.trailheads, d.peaks)
}

func (d day10) Part2() int {
	return day.SumFunc(d.trailheads, d.rating)
}

func main() {
	d := NewDay10(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
