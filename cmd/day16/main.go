package main

import (
	"bytes"
	"maps"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/grid"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)

	turns = map[direction][2]direction{
		{0, 1}:  {{-1, 0}, {1, 0}},
		{1, 0}:  {{0, -1}, {0, 1}},
		{0, -1}: {{-1, 0}, {1, 0}},
		{-1, 0}: {{0, -1}, {0, 1}},
	}
)

type tile struct {
	x, y int
}

type direction struct {
	dx, dy int
}

type edge struct {
	to     state
	weight int
}

type state struct {
	x, y      int
	direction direction
}

type day16 struct {
	grid       [][]byte
	start, end state
}

func (e edge) To() state {
	return e.to
}

func (e edge) Weight() int {
	return e.weight
}

func readInput(d day.DayInput) day16 {
	lines := d.ReadByteGrid()

	grid := make([][]byte, len(lines))
	var start, end state

	for i, line := range lines {
		grid[i] = line
		if j := bytes.IndexByte(line, 'S'); j != -1 {
			start = state{i, j, direction{0, 1}}
		}
		if j := bytes.IndexByte(line, 'E'); j != -1 {
			end = state{i, j, direction{}}
		}
	}

	return day16{grid, start, end}
}

func NewDay16(opts ...day.Option) day16 {
	input := day.NewDayInput(path, opts...)

	return readInput(input)
}

func (d day16) path(prev map[state]map[state]struct{}, e state) map[tile]struct{} {
	result := map[tile]struct{}{{e.x, e.y}: {}}

	if e == d.start {
		return result
	}

	for k := range prev[e] {
		maps.Copy(result, d.path(prev, k))
	}

	return result
}

func (d day16) allShortestPaths() int {
	dist, prev := grid.AllShortestPaths(d.start, d.neighbours)

	end := d.endState(dist)

	return len(d.path(prev, end))
}

func (d day16) endState(dist map[state]int) state {
	shortest := math.MaxInt
	var result state

	for dir := range maps.Keys(turns) {
		end := state{d.end.x, d.end.y, dir}
		if distEnd, ok := dist[end]; ok && distEnd < shortest {
			shortest = distEnd
			result = end
		}
	}

	return result
}

func (d day16) shortestPath() int {
	dist, _ := grid.ShortestPath(d.start, d.neighbours)

	return dist[d.endState(dist)]
}

func (s state) forward() state {
	return state{s.x + s.direction.dx, s.y + s.direction.dy, s.direction}
}

func (s state) turn(d direction) state {
	return state{s.x, s.y, d}
}

func (d day16) scan(s state, dir direction) byte {
	return d.grid[s.x+dir.dx][s.y+dir.dy]
}

func (d day16) neighbours(s state) []grid.Edge[state] {
	var result []grid.Edge[state]

	if d.scan(s, s.direction) != '#' {
		result = append(result, edge{s.forward(), 1})
	}

	for _, dir := range turns[s.direction] {
		if d.scan(s, dir) != '#' {
			result = append(result, edge{s.turn(dir), 1000})
		}
	}

	return result
}

func (d day16) Part1() int {
	return d.shortestPath()
}

func (d day16) Part2() int {
	return d.allShortestPaths()
}

func main() {
	d := NewDay16(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
