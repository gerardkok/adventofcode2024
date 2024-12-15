package main

import (
	"bytes"
	"fmt"
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
	robot     position
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

var moveMap = map[byte]direction{
	'^': {-1, 0},
	'>': {0, 1},
	'v': {1, 0},
	'<': {0, -1},
}

func parseInput(lines [][]byte) ([][]byte, []byte, position) {
	var (
		result [2][][]byte
		robot  position
	)

	i := 0
	for r, line := range lines {
		if len(line) == 0 {
			i++
			continue
		}

		result[i] = append(result[i], line)

		c := bytes.IndexByte(line, '@')
		if c != -1 {
			robot = position{r, c}
		}
	}

	m := bytes.Join(result[1], []byte{})
	return result[0], m, robot
}

func NewDay15(opts ...day.Option) day15 {
	input := day.NewDayInput(path, opts...)
	lines := input.ReadByteGrid()
	grid, moves, robot := parseInput(lines)

	return day15{grid, moves, robot}
}

func (p position) to(d direction) position {
	return position{p.x + d.dx, p.y + d.dy}
}

// func (d day15) offGrid(p position) bool {
// 	return p.x < 0 || p.x >= len(w.grid) || p.y < 0 || p.y >= len(w.grid[0])
// }

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

func (w warehouse) movePositions(dir direction) int {
	for i := 1; ; i++ {
		x, y := w.robot.x+i*dir.dx, w.robot.y+i*dir.dy
		//		fmt.Printf("checking (%d,%d): %c\n", x, y, w.grid[x][y])
		switch w.grid[x][y] {
		case '.':
			return i
		case '#':
			return 0
		}
	}
}

func (w warehouse) canMoveVertically(p position, d direction) bool {
	next := position{p.x + d.dx, p.y + d.dy}
	switch w.grid[next.x][next.y] {
	case '.':
		return true
	case 'O':
		return w.canMoveVertically(next, d)
	case '[':
		right := position{next.x, next.y + 1}
		return w.canMoveVertically(next, d) && w.canMoveVertically(right, d)
	case ']':
		left := position{next.x, next.y - 1}
		return w.canMoveVertically(next, d) && w.canMoveVertically(left, d)
	default: // including '#'
		return false
	}
}

func (w warehouse) canMoveHorizontally(p position, d direction) bool {
	next := position{p.x + d.dx, p.y + d.dy}
	switch w.grid[next.x][next.y] {
	case '.':
		return true
	case 'O', '[', ']':
		return w.canMoveHorizontally(next, d)
	default:
		return false
	}
}

func (w warehouse) canMove(p position, d direction) bool {
	if d.dx == 0 {
		return w.canMoveHorizontally(p, d)
	}

	switch w.grid[p.x][p.y] {
	case '[':
		right := position{p.x, p.y + 1}
		return w.canMoveVertically(p, d) && w.canMoveVertically(right, d)
	case ']':
		left := position{p.x, p.y - 1}
		return w.canMoveVertically(p, d) && w.canMoveVertically(left, d)
	default:
		return w.canMoveVertically(p, d)
	}
}

func (w *warehouse) swap(a, b position) {
	w.grid[a.x][a.y], w.grid[b.x][b.y] = w.grid[b.x][b.y], w.grid[a.x][a.y]
}

func (w *warehouse) moveHorizontally(p position, d direction) {
	next := position{p.x + d.dx, p.y + d.dy}
	switch w.grid[next.x][next.y] {
	case '.':
		w.swap(p, next)
	case 'O', '[', ']':
		w.moveHorizontally(next, d)
		w.swap(p, next)
	}
}

func (w *warehouse) moveVertically(p position, d direction) {
	next := position{p.x + d.dx, p.y + d.dy}
	// fmt.Printf("p: %v: %c, next: %v: %c\n", p, w.grid[p.x][p.y], next, w.grid[next.x][next.y])
	switch w.grid[next.x][next.y] {
	case '.':
		//	fmt.Println("swapping")
		w.swap(p, next)
	case 'O':
		w.moveVertically(next, d)
		w.swap(p, next)
	case '[':
		//fmt.Println("dealing with [")
		rightOfNext := position{next.x, next.y + 1}
		//	fmt.Printf("next: %v: %c, right of next: %v: %c\n", next, w.grid[next.x][next.y], rightOfNext, w.grid[rightOfNext.x][rightOfNext.y])
		w.moveVertically(next, d)
		w.moveVertically(rightOfNext, d)
		w.swap(next, p)
	case ']':
		//fmt.Println("dealing with ]")
		leftOfNext := position{next.x, next.y - 1}
		w.moveVertically(next, d)
		w.moveVertically(leftOfNext, d)
		w.swap(next, p)
	}
}

func (w *warehouse) move(p position, d direction) {
	if d.dx == 0 {
		w.moveHorizontally(p, d)
		return
	}

	switch w.grid[p.x][p.y] {
	case '[':
		right := position{p.x, p.y + 1}
		w.moveVertically(p, d)
		w.moveVertically(right, d)
	case ']':
		left := position{p.x, p.y - 1}
		w.moveVertically(p, d)
		w.moveVertically(left, d)
	default:
		w.moveVertically(p, d)
	}
}

func (w warehouse) canMoveRobot(d direction) bool {
	return w.canMove(w.robot, d)
}

func (w *warehouse) moveRobot(d direction) {
	w.move(w.robot, d)
	w.robot.x, w.robot.y = w.robot.x+d.dx, w.robot.y+d.dy
}

func (w warehouse) sumBoxesGPS(edge byte) int {
	result := 0

	for x, line := range w.grid {
		for y, c := range line {
			if c == edge {
				result += 100*x + y
			}
		}
	}

	return result
}

func (d day15) warehousePart1() warehouse {
	return warehouse{d.gridInput, d.moves, d.robot}
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
		fmt.Printf("moving %v %c\n", w.robot, move)
		if w.canMoveRobot(moveMap[move]) {
			// fmt.Println("can move")
			w.moveRobot(moveMap[move])
		}

		fmt.Println(w)
		fmt.Println("---")
	}
}

func (d day15) Part1() int {
	w := d.warehousePart1()

	//w.sequence()

	return w.sumBoxesGPS('O')
}

func (d day15) Part2() int {
	w := d.warehousePart2()

	// fmt.Println(w)
	w.sequence()

	return w.sumBoxesGPS('[')
}

func main() {
	d := NewDay15(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
