package main

import (
	"bytes"
	"fmt"
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

func (w warehouse) String() string {
	result := ""

	for _, line := range w.grid {
		result += fmt.Sprintf("%s\n", string(line))
	}
	result += fmt.Sprintln()
	// result += fmt.Sprintf("%s\n", string(w.moves))
	// result += fmt.Sprintln()
	result += fmt.Sprintf("robot: (%d,%d)\n", w.robot.x, w.robot.y)

	return result
}

func (w *warehouse) swap(a, b position) {
	w.grid[a.x][a.y], w.grid[b.x][b.y] = w.grid[b.x][b.y], w.grid[a.x][a.y]
}

func mergeS(s, t map[swap]struct{}) map[swap]struct{} {
	// result := make(map[swap]struct{})
	// maps.Copy(result, s)
	// maps.Copy(result, t)
	// return result
	maps.Copy(s, t)
	return s
}

func merge(s, t []map[swap]struct{}) []map[swap]struct{} {
	if len(s) == 0 || len(t) == 0 {
		return nil
	}

	var result []map[swap]struct{}

	for i := range min(len(s), len(t)) {
		maps.Copy(s[i], t[i])
		result = append(result, s[i])
	}
	if len(s) < len(t) {
		result = append(result, t[len(s):]...)
	}
	if len(s) > len(t) {
		result = append(result, s[len(t):]...)
	}

	return result
}

func swapAppend(swaps []map[swap]struct{}, e map[swap]struct{}) []map[swap]struct{} {
	if len(swaps) == 0 {
		return nil
	}

	return append(swaps, e)
}

func (w warehouse) moveHorizontally2(p position, d direction) []map[swap]struct{} {
	next := position{p.x + d.dx, p.y + d.dy}
	s := swap{p, next}
	switch w.grid[next.x][next.y] {
	case '.':
		return []map[swap]struct{}{map[swap]struct{}{s: struct{}{}}}
	case 'O', '[', ']':
		return swapAppend(w.moveHorizontally2(next, d), map[swap]struct{}{s: struct{}{}})
	default:
		return nil
	}
}

func (w warehouse) moveVertically2(p position, d direction) []map[swap]struct{} {
	next := position{p.x + d.dx, p.y + d.dy}
	s := swap{p, next}
	switch w.grid[next.x][next.y] {
	case '.':
		return []map[swap]struct{}{map[swap]struct{}{s: struct{}{}}}
	case 'O':
		return swapAppend(w.moveVertically2(next, d), map[swap]struct{}{s: struct{}{}})
	case '[':
		rightOfNext := position{next.x, next.y + 1}
		r := swapAppend(merge(w.moveVertically2(next, d), w.moveVertically2(rightOfNext, d)), map[swap]struct{}{s: struct{}{}})
		// fmt.Printf("[ swaps: %v\n", r)
		return r
	case ']':
		leftOfNext := position{next.x, next.y - 1}
		r := swapAppend(merge(w.moveVertically2(next, d), w.moveVertically2(leftOfNext, d)), map[swap]struct{}{s: struct{}{}})
		// fmt.Printf("] swaps: %v\n", r)
		return r
	default:
		return nil
	}
}

func (w warehouse) move2(p position, d direction) []map[swap]struct{} {
	if d.dx == 0 {
		return w.moveHorizontally2(p, d)
	}

	switch w.grid[p.x][p.y] {
	case '[':
		right := position{p.x, p.y + 1}
		return merge(w.moveVertically2(p, d), w.moveVertically2(right, d))
	case ']':
		left := position{p.x, p.y - 1}
		return merge(w.moveVertically2(p, d), w.moveVertically2(left, d))
	default:
		return w.moveVertically2(p, d)
	}
}

func (w *warehouse) moveRobot2(d direction) {
	swaps := w.move2(w.robot, d)
	// fmt.Println(swaps)

	if len(swaps) == 0 {
		return
	}

	for _, v := range swaps {
		for s := range v {
			w.swap(s.p, s.q)
		}
	}

	w.robot.x, w.robot.y = w.robot.x+d.dx, w.robot.y+d.dy
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

func (w warehouse) sequence() {
	for _, move := range w.moves {
		// if w.canMoveRobot(moveMap[move]) {
		// 	w.moveRobot(moveMap[move])
		// }
		// fmt.Printf("%c\n", move)
		// fmt.Println(w)
		w.moveRobot2(moveMap[move])
		// fmt.Println(w)
		// fmt.Println("---")
	}
}

func (d day15) Part1() int {
	w := d.warehousePart1()

	w.sequence()

	return w.sumBoxesGPS()
}

func (d day15) Part2() int {
	w := d.warehousePart2()

	w.sequence()

	return w.sumBoxesGPS()
}

func main() {
	d := NewDay15(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
