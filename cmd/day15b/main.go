package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day15b struct {
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

var (
	moves = map[byte]direction{
		'^': {-1, 0},
		'>': {0, 1},
		'v': {1, 0},
		'<': {0, -1},
	}
	otherHalf = map[byte]direction{'[': {0, 1}, ']': {0, -1}}
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

func NewDay15b(opts ...day.Option) day15b {
	input := day.NewDayInput(path, opts...)
	lines := input.ReadByteGrid()
	grid, moves := parseInput(lines)

	return day15b{grid, moves}
}

func (p position) to(d direction) position {
	return position{p.x + d.dx, p.y + d.dy}
}

func (w warehouse) at(p position) byte {
	return w.grid[p.x][p.y]
}

func (w *warehouse) swap(p position, d direction) {
	// q is always '.'
	q := p.to(d)
	w.grid[q.x][q.y] = w.grid[p.x][p.y]
	w.grid[p.x][p.y] = '.'
}

func (w warehouse) horizontalSwaps(from position, d direction) []position {
	var result []position

	for {
		next := from.to(d)

		result = append(result, from)

		c := w.at(next)

		switch c {
		case '.':
			return result
		case 'O', '[', ']':
			from = next
		default:
			return nil
		}
	}
}

func (w warehouse) verticalSwaps(from position, d direction) []position {
	seen := make(map[int]struct{})
	todo := []position{from}
	c := w.at(from)
	if c == '[' || c == ']' {
		other := from.to(otherHalf[c])
		todo = append(todo, other)
	}

	var result []position

	for len(todo) > 0 {
		p := todo[0]
		next := p.to(d)
		todo = todo[1:]

		key := p.x*len(w.grid[0]) + p.y

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}

		result = append(result, p)

		c := w.at(next)

		switch c {
		case '.':
		case 'O':
			todo = append(todo, next)
		case '[', ']':
			otherHalf := next.to(otherHalf[c])
			todo = append(todo, next, otherHalf)
		default:
			return nil
		}
	}

	return result
}

func (w warehouse) move(from position, d direction) []position {
	if d.dx == 0 {
		return w.horizontalSwaps(from, d)
	}

	return w.verticalSwaps(from, d)
}

func (w *warehouse) moveRobot(d direction) {
	swaps := w.move(w.robot, d)
	for i := len(swaps) - 1; i >= 0; i-- {
		w.swap(swaps[i], d)
	}

	if len(swaps) > 0 {
		w.robot = w.robot.to(d)
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

func (d day15b) warehousePart1() warehouse {
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

func (d day15b) warehousePart2() warehouse {
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
		w.moveRobot(moves[move])
	}
}

func (d day15b) Part1() int {
	w := d.warehousePart1()

	w.moveSequence()

	return w.sumBoxesGPS()
}

func (d day15b) Part2() int {
	w := d.warehousePart2()

	w.moveSequence()

	return w.sumBoxesGPS()
}

func main() {
	d := NewDay15b(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
