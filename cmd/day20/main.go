package main

import (
	"bytes"
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

type day20 struct {
	track               [][]byte
	start, end          grid.Point
	minSaving, maxCheat int
}

type cheat struct {
	start, end grid.Point
}

func readInput(d day.DayInput) ([][]byte, grid.Point, grid.Point) {
	lines := d.ReadByteGrid()

	track := make([][]byte, len(lines))
	var start, end grid.Point

	for i, line := range lines {
		track[i] = line
		if j := bytes.IndexByte(line, 'S'); j != -1 {
			start = grid.Point{X: i, Y: j}
		}
		if j := bytes.IndexByte(line, 'E'); j != -1 {
			end = grid.Point{X: i, Y: j}
		}
	}

	return track, start, end
}

func NewDay20(minSaving int, opts ...day.Option) day20 {
	input := day.NewDayInput(path, opts...)

	track, start, end := readInput(input)

	return day20{track, start, end, minSaving, 0}
}

func (d day20) neighbours(from grid.Point) []grid.Point {
	var result []grid.Point

	for _, dir := range []grid.Direction{{Dx: 0, Dy: 1}, {Dx: 1, Dy: 0}, {Dx: 0, Dy: -1}, {Dx: -1, Dy: 0}} {
		to := from.To(dir)
		if d.track[to.X][to.Y] == '#' {
			continue
		}

		result = append(result, to)
	}

	return result
}

func (d day20) bfs(start grid.Point) map[grid.Point]int {
	dist := map[grid.Point]int{start: 0}

	for p := range grid.Bfs(start, d.neighbours) {
		if p[0] == start {
			continue
		}

		dist[p[0]] = dist[p[1]] + 1
	}

	return dist
}

func (d day20) cheatablePaths(distStart, distEnd map[grid.Point]int, maxCheat, maxLength int) int {
	result := 0

	for start, sDist := range distStart {
		for dx := -maxCheat; dx <= maxCheat; dx++ {
			for dy := -maxCheat + conv.Abs(dx); dy <= maxCheat-conv.Abs(dx); dy++ {
				end := start.To(grid.Direction{Dx: dx, Dy: dy})
				if eDist, ok := distEnd[end]; ok {
					dist := sDist + eDist + conv.Abs(dx) + conv.Abs(dy)
					if dist <= maxLength {
						result++
					}
				}
			}
		}
	}

	return result
}

func (d day20) Part1() int {
	distStart := d.bfs(d.start)
	distEnd := d.bfs(d.end)

	shortest := distEnd[d.start]

	return d.cheatablePaths(distStart, distEnd, 2, shortest-d.minSaving)
}

func (d day20) Part2() int {
	distStart := d.bfs(d.start)
	distEnd := d.bfs(d.end)

	shortest := distEnd[d.start]

	return d.cheatablePaths(distStart, distEnd, 20, shortest-d.minSaving)
}

func main() {
	d := NewDay20(100, day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
