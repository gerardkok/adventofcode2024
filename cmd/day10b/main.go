package main

import (
	"iter"
	"os"
	"path/filepath"
	"runtime"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
	"adventofcode2024/internal/grid"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day10b struct {
	grid       grid.Grid[byte]
	trailheads []grid.Point
}

func trailheads(g grid.Grid[byte]) []grid.Point {
	var result []grid.Point

	for x, line := range g {
		for y, v := range line {
			if v == '0' {
				result = append(result, grid.Point{X: x, Y: y})
			}
		}
	}

	return result
}

func NewDay10b(opts ...day.Option) day10b {
	input := day.NewDayInput(path, opts...)

	grid := input.ReadByteGrid()
	trailheads := trailheads(grid)

	return day10b{grid, trailheads}
}

func (d day10b) up(p grid.Point) iter.Seq[grid.Point] {
	return func(yield func(grid.Point) bool) {
		for n := range d.grid.Neighbours4(p) {
			if d.grid.At(p)+1 != d.grid.At(n) {
				continue
			}

			if !yield(n) {
				return
			}
		}
	}
}

func (d day10b) rating(p grid.Point) int {
	if d.grid.At(p) == '9' {
		return 1
	}

	result := 0

	for n := range d.up(p) {
		result += d.rating(n)
	}

	return result
}

func (d day10b) peaks(trailhead grid.Point) int {
	appender := func(p grid.Point) iter.Seq[grid.Point] {
		return d.up(p)
	}

	peaks := 0

	for p := range d.grid.Bfs(trailhead, appender) {
		if d.grid.At(p) == '9' {
			peaks++
		}
	}

	return peaks
}

func (d day10b) Part1() int {
	return conv.SumFunc(d.trailheads, d.peaks)
}

func (d day10b) Part2() int {
	return conv.SumFunc(d.trailheads, d.rating)
}

func main() {
	d := NewDay10b(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
