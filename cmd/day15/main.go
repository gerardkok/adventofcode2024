package main

import (
	"bytes"
	"maps"
	"os"
	"path/filepath"
	"runtime"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day15 struct {
	gridInput [][]byte
	moves     []byte
}

type warehouse struct {
	grid  [][]byte
	moves []byte
	robot position
}

type position struct {
	x, y int
}

type direction struct {
	dx, dy int
}

type swap struct {
	p, q position
}

type swapSet map[swap]struct{}

var (
	moveMap = map[byte]direction{
		'^': {-1, 0},
		'>': {0, 1},
		'v': {1, 0},
		'<': {0, -1},
	}
	dyOtherHalf = map[byte]int{'[': 1, ']': -1}
)

func parseInput(lines [][]byte) ([][]byte, []byte) {
	var result [2][][]byte

	i := 0
	for _, line := range lines {
		if len(line) == 0 {
			i++
			continue
		}

		result[i] = append(result[i], line)
	}

	m := bytes.Join(result[1], []byte{})
	return result[0], m
}

func NewDay15(opts ...day.Option) day15 {
	input := day.NewDayInput(path, opts...)
	lines := input.ReadByteGrid()
	grid, moves := parseInput(lines)

	return day15{grid, moves}
}

func (p position) to(d direction) position {
	return position{p.x + d.dx, p.y + d.dy}
}

func (w *warehouse) swap(a, b position) {
	w.grid[a.x][a.y], w.grid[b.x][b.y] = w.grid[b.x][b.y], w.grid[a.x][a.y]
}

func merge(s, t []swapSet) []swapSet {
	if len(t) < len(s) {
		s, t = t, s
	}

	result := make([]swapSet, len(t))
	copy(result, t)

	for i := range s {
		maps.Copy(result[i], s[i])
	}

	return result
}

func (w warehouse) horizontalSwaps(p position, d direction) ([]swapSet, bool) {
	next := p.to(d)
	s := swap{p, next}
	c := w.grid[next.x][next.y]

	switch c {
	case '.':
		return []swapSet{{s: {}}}, true
	case 'O', '[', ']':
		swaps, ok := w.horizontalSwaps(next, d)
		return append(swaps, swapSet{s: {}}), ok
	default:
		return nil, false
	}
}

func (w warehouse) verticalSwaps(p position, d direction) ([]swapSet, bool) {
	next := p.to(d)
	s := swap{p, next}
	c := w.grid[next.x][next.y]

	switch c {
	case '.':
		return []swapSet{{s: {}}}, true
	case 'O':
		swaps, ok := w.verticalSwaps(next, d)
		return append(swaps, swapSet{s: {}}), ok
	case '[', ']':
		otherHalf := position{next.x, next.y + dyOtherHalf[c]}
		swaps, ok := w.wideBoxSwaps(next, otherHalf, d)
		return append(swaps, swapSet{s: {}}), ok
	default:
		return nil, false
	}
}

func (w warehouse) wideBoxSwaps(p, q position, d direction) ([]swapSet, bool) {
	swapsP, okP := w.verticalSwaps(p, d)
	swapsQ, okQ := w.verticalSwaps(q, d)
	return merge(swapsP, swapsQ), okP && okQ
}

func (w warehouse) move(p position, d direction) ([]swapSet, bool) {
	if d.dx == 0 {
		return w.horizontalSwaps(p, d)
	}

	c := w.grid[p.x][p.y]
	switch c {
	case '[', ']':
		otherHalf := position{p.x, p.y + dyOtherHalf[c]}
		return w.wideBoxSwaps(p, otherHalf, d)
	default:
		return w.verticalSwaps(p, d)
	}
}

func (w *warehouse) moveRobot(d direction) {
	if swaps, ok := w.move(w.robot, d); ok {
		for _, v := range swaps {
			for s := range v {
				w.swap(s.p, s.q)
			}
		}

		w.robot.x, w.robot.y = w.robot.x+d.dx, w.robot.y+d.dy
	}
}

func (w warehouse) sumBoxesGPS() int {
	result := 0

	for x, line := range w.grid {
		for y, c := range line {
			if c == 'O' || c == '[' {
				result += 100*x + y
			}
		}
	}

	return result
}

func (d day15) warehousePart1() warehouse {
	grid := make([][]byte, len(d.gridInput))
	var robot position

	for i, line := range d.gridInput {
		for j, c := range line {
			grid[i] = append(grid[i], c)
			if c == '@' {
				robot = position{i, j}
			}
		}
	}

	return warehouse{grid, d.moves, robot}
}

func (d day15) warehousePart2() warehouse {
	grid := make([][]byte, len(d.gridInput))
	var robot position

	for i, line := range d.gridInput {
		for j, c := range line {
			switch c {
			case '@':
				robot = position{i, 2 * j}
				grid[i] = append(grid[i], '@', '.')
			case 'O':
				grid[i] = append(grid[i], '[', ']')
			default:
				grid[i] = append(grid[i], c, c)
			}
		}
	}

	return warehouse{grid, d.moves, robot}
}

func (w warehouse) moveSequence() {
	for _, move := range w.moves {
		w.moveRobot(moveMap[move])
	}
}

func (d day15) Part1() int {
	w := d.warehousePart1()

	w.moveSequence()

	return w.sumBoxesGPS()
}

func (d day15) Part2() int {
	w := d.warehousePart2()

	w.moveSequence()

	return w.sumBoxesGPS()
}

func main() {
	d := NewDay15(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
