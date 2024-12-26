package main

import (
	"bytes"
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
	grid   [][]byte
	spots  []spot
	fallen int
}

type spot struct {
	r, c int
}

type direction struct {
	dr, dc int
}

func parseGrid(lines []string, size, fallenBytes int) ([][]byte, []spot) {
	result := make([][]byte, size)
	var spots []spot

	for i := range result {
		result[i] = bytes.Repeat([]byte{'.'}, size)
	}

	for i, line := range lines {
		c, r, _ := strings.Cut(line, ",")
		s := spot{conv.MustAtoi(r), conv.MustAtoi(c)}
		spots = append(spots, s)
		if i < fallenBytes {
			result[s.r][s.c] = '#'
		}
	}

	return result, spots
}

func NewDay18(size, fallenBytes int, opts ...day.Option) day18 {
	input := day.NewDayInput(path, opts...)

	lines := input.ReadLines()

	grid, spots := parseGrid(lines, size, fallenBytes)

	return day18{grid, spots, fallenBytes}
}

func (s spot) to(d direction) spot {
	return spot{s.r + d.dr, s.c + d.dc}
}

func (d day18) offGrid(s spot) bool {
	return s.r < 0 || s.r >= len(d.grid) || s.c < 0 || s.c >= len(d.grid[0])
}

func (d day18) shortestPath() int {
	start := spot{0, 0}
	end := spot{len(d.grid) - 1, len(d.grid[0]) - 1}

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

		if d.grid[to.r][to.c] == '#' {
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
	for _, s := range d.spots[d.fallen:] {
		d.grid[s.r][s.c] = '#'
		if d.shortestPath() == -1 {
			output := fmt.Sprintf("%d,%d", s.c, s.r)
			return output
		}
	}
	return "71,71"
}

func main() {
	d := NewDay18(71, 1024, day.FromArgs(os.Args[1:]))

	fmt.Println(d.Part1())
	fmt.Println(d.Part2())
}
