package main

import (
	"iter"
	"os"
	"path/filepath"
	"runtime"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/grid"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day10b struct {
	day.DayInput
}

func NewDay10b(opts ...day.Option) day10b {
	return day10b{day.NewDayInput(path, opts...)}
}

func up(g grid.BorderedGrid[byte], p grid.Point) iter.Seq[grid.Point] {
	return func(yield func(grid.Point) bool) {
		for n := range g.Neighbours4(p) {
			if g.At(p)+1 != g.At(n) {
				continue
			}

			if !yield(n) {
				return
			}
		}
	}
}

func rating(g grid.BorderedGrid[byte], p grid.Point) int {
	if g.At(p) == '9' {
		return 1
	}

	result := 0

	for n := range up(g, p) {
		result += rating(g, n)
	}

	return result
}

func peaks(g grid.BorderedGrid[byte], trailhead grid.Point) int {
	appender := func(p grid.Point) iter.Seq[grid.Point] {
		return up(g, p)
	}

	peaks := 0

	for p := range g.Bfs(trailhead, appender) {
		if g.At(p) == '9' {
			peaks++
		}
	}

	return peaks
}

func (d day10b) Part1() int {
	input := d.ReadByteGrid()
	g := grid.NewBorderedGrid[byte](input, byte(0))

	sum := 0

	for p := range g.PointsIter() {
		if g.At(p) == '0' {
			sum += peaks(g, p)
		}
	}

	return sum
}

func (d day10b) Part2() int {
	input := d.ReadByteGrid()
	g := grid.NewBorderedGrid[byte](input, byte(0))

	sum := 0

	for p := range g.PointsIter() {
		if g.At(p) == '0' {
			sum += rating(g, p)
		}
	}

	return sum
}

func main() {
	d := NewDay10b(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
