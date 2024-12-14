package main

import (
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type coord struct {
	w, h int
}

type robot struct {
	position, velocity coord
}

type day14 struct {
	width, height int
	robots        []robot
}

func parseCoord(s string) coord {
	_, c, _ := strings.Cut(s, "=")
	w, h, _ := strings.Cut(c, ",")
	return coord{conv.MustAtoi(w), conv.MustAtoi(h)}
}

func readRobots(d day.DayInput) []robot {
	lines := d.ReadLines()

	result := make([]robot, len(lines))

	for i, line := range lines {
		p, v, _ := strings.Cut(line, " ")
		position := parseCoord(p)
		velocity := parseCoord(v)
		result[i] = robot{position, velocity}
	}

	return result
}

func NewDay14(width, height int, opts ...day.Option) day14 {
	input := day.NewDayInput(path, opts...)

	robots := readRobots(input)

	return day14{width, height, robots}
}

func (d day14) robotPositions(seconds int) [][]int {
	result := make([][]int, d.height)
	for i := range d.height {
		result[i] = make([]int, d.width)
	}

	for _, robot := range d.robots {
		p, v := robot.position, robot.velocity
		w := ((p.w+seconds*v.w)%d.width + d.width) % d.width
		h := ((p.h+seconds*v.h)%d.height + d.height) % d.height
		result[h][w]++
	}

	return result
}

func (d day14) print(grid [][]int) string {
	result := ""
	for h := range d.height {
		for w := range d.width {
			if grid[h][w] == 0 {
				result += "."
			} else {
				result += fmt.Sprintf("%d", grid[h][w])
			}
		}
		result += "\n"
	}
	return result
}

func (d day14) quadrant(w, h int) [2]int {
	qw, qh := 0, 0

	switch {
	case h < d.height/2:
		qh = -1
	case h > d.height/2:
		qh = 1
	}

	switch {
	case w < d.width/2:
		qw = -1
	case w > d.width/2:
		qw = 1
	}

	return [2]int{qh, qw}
}

func (d day14) safetyFactor(grid [][]int) int {
	quadrants := map[[2]int]int{{-1, -1}: 0, {-1, 1}: 0, {1, -1}: 0, {1, 1}: 0}

	for h := range d.height {
		for w := range d.width {
			q := d.quadrant(w, h)
			if _, ok := quadrants[q]; !ok {
				continue
			}

			quadrants[q] += grid[h][w]
		}
	}

	return safetyFactor(quadrants)
}

func safetyFactor(quadrants map[[2]int]int) int {
	result := 1

	for _, nRobots := range quadrants {
		result *= nRobots
	}

	return result
}

func (d day14) nNeighbours(h, w int, grid [][]int) int {
	result := 0
	for i := h - 1; i <= h+1; i++ {
		if i < 0 || i >= d.height {
			continue
		}

		for j := w - 1; j <= w+1; j++ {
			if j < 0 || j >= d.width {
				continue
			}

			result += grid[i][j]
		}
	}

	return result
}

func (d day14) totalNeighbours(grid [][]int) int {
	result := 0

	for h := range d.height {
		for w := range d.width {
			if grid[h][w] > 0 {
				result += d.nNeighbours(h, w, grid)
			}
		}
	}

	return result
}

func (d day14) neighbours() iter.Seq[int] {
	return func(yield func(int) bool) {
		seconds := 0

		for {
			grid := d.robotPositions(seconds)
			n := d.totalNeighbours(grid)
			if !yield(n) {
				return
			}

			seconds++
		}
	}
}

func (d day14) Part1() int {
	grid := d.robotPositions(100)

	return d.safetyFactor(grid)
}

func (d day14) Part2() int {
	seconds := 0
	for a := range d.neighbours() {
		if a == 2346 {
			return seconds
		}

		seconds++
	}

	return 0
}

func main() {
	d := NewDay14(101, 103, day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
