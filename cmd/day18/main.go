package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
	"adventofcode2024/internal/grid"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day18 struct {
	spots        []spot
	corrupted    map[spot]bool
	size, fallen int
}

type spot struct {
	r, c int
}

type direction struct {
	dr, dc int
}

func parseGrid(lines []string, fallen int) ([]spot, map[spot]bool) {
	spots := make([]spot, len(lines))
	corrupted := make(map[spot]bool, len(lines))

	for i, line := range lines {
		c, r, _ := strings.Cut(line, ",")
		s := spot{conv.MustAtoi(r), conv.MustAtoi(c)}
		spots[i] = s
		corrupted[s] = i < fallen
	}

	return spots, corrupted
}

func NewDay18(size, fallen int, opts ...day.Option) day18 {
	input := day.NewDayInput(path, opts...)

	lines := input.ReadLines()

	spots, corrupted := parseGrid(lines, fallen)

	return day18{spots, corrupted, size, fallen}
}

func (s spot) to(d direction) spot {
	return spot{s.r + d.dr, s.c + d.dc}
}

func (d day18) offGrid(s spot) bool {
	return s.r < 0 || s.r >= d.size || s.c < 0 || s.c >= d.size
}

func (d day18) shortestPath() int {
	start := spot{0, 0}
	end := spot{d.size - 1, d.size - 1}

	dist := map[spot]int{start: 0}

	for p := range grid.Bfs(start, d.neighbours) {
		if p[0] == start {
			continue
		}
		dist[p[0]] = dist[p[1]] + 1
		if p[0] == end {
			return dist[p[0]]
		}
	}

	return -1
}

func (d day18) neighbours(s spot) []spot {
	var result []spot

	for _, dir := range []direction{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
		to := s.to(dir)
		if d.offGrid(to) {
			continue
		}

		if d.corrupted[to] {
			continue
		}

		result = append(result, to)
	}

	return result
}

func (d day18) Part1() int {
	return d.shortestPath()
}

func (d day18) Part2() string {
	lo, hi := d.fallen, len(d.spots)
	for lo < hi {
		t := (lo+hi)/2 + 1
		for _, s := range d.spots[lo:t] {
			d.corrupted[s] = true
		}
		if d.shortestPath() == -1 {
			hi = t - 1
		} else {
			lo = t
		}
		for _, s := range d.spots[lo:t] {
			d.corrupted[s] = false
		}
	}
	return fmt.Sprintf("%d,%d", d.spots[lo].c, d.spots[lo].r)
}

func main() {
	d := NewDay18(71, 1024, day.FromArgs(os.Args[1:]))

	fmt.Println(d.Part1())
	fmt.Println(d.Part2())
}
