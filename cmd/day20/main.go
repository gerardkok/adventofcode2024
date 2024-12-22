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

func (d day20) neighbours(p grid.Point) []grid.Edge[grid.Point] {
	var result []grid.Edge[grid.Point]

	for _, dir := range []grid.Direction{{Dx: 0, Dy: 1}, {Dx: 1, Dy: 0}, {Dx: 0, Dy: -1}, {Dx: -1, Dy: 0}} {
		x, y := p.X+dir.Dx, p.Y+dir.Dy
		if d.track[x][y] == '#' {
			continue
		}

		result = append(result, grid.Edge[grid.Point]{To: grid.Point{X: x, Y: y}, Weight: 1})
	}

	return result
}

func (d *day20) cheatablePaths(distStart, distEnd map[grid.Point]int, maxCheat int) map[cheat]int {
	result := make(map[cheat]int)

	for cheatStart, distCheatStart := range distStart {
		for cheatEnd, distCheatEnd := range distEnd {
			cheat := cheat{cheatStart, cheatEnd}
			cheatDist := conv.Abs(cheatEnd.X-cheatStart.X) + conv.Abs(cheatEnd.Y-cheatStart.Y)
			pathDist := distCheatStart + distCheatEnd + cheatDist
			if cheatDist <= maxCheat {
				result[cheat] = pathDist
			}
		}
	}

	return result
}

func (d day20) Part1() int {
	distStart, _ := grid.Dijkstra(d.start, d.neighbours)
	distEnd, _ := grid.Dijkstra(d.end, d.neighbours)

	shortest := distEnd[d.start]

	paths := d.cheatablePaths(distStart, distEnd, 2)

	sum := 0

	for _, length := range paths {
		if length <= shortest-d.minSaving {
			sum++
		}
	}

	return sum
}

func (d day20) Part2() int {
	distStart, _ := grid.Dijkstra(d.start, d.neighbours)
	distEnd, _ := grid.Dijkstra(d.end, d.neighbours)

	shortest := distEnd[d.start]

	paths := d.cheatablePaths(distStart, distEnd, 20)

	sum := 0

	for _, length := range paths {
		if length <= shortest-d.minSaving {
			sum++
		}
	}

	return sum
}

func main() {
	d := NewDay20(100, day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
