package main

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"runtime"

	pq "github.com/emirpasic/gods/queues/priorityqueue"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)

	turnMap = map[direction][2]direction{
		{0, 1}:  [2]direction{{-1, 0}, {1, 0}},
		{1, 0}:  [2]direction{{0, -1}, {0, 1}},
		{0, -1}: [2]direction{{-1, 0}, {1, 0}},
		{-1, 0}: [2]direction{{0, -1}, {0, 1}},
	}
)

type tile struct {
	x, y int
}

type direction struct {
	dx, dy int
}

type state struct {
	x, y      int
	direction direction
}

type queue struct {
	*pq.Queue
}

type node struct {
	state state
	cost  int
}

type day16 struct {
	grid       [][]byte
	start, end state
}

type visited map[state]int

func (v visited) get(s state) int {
	if v, ok := v[s]; ok {
		return v
	}
	return math.MaxInt
}

func newQueue() queue {
	q := pq.NewWith(func(a, b any) int {
		return a.(node).cost - b.(node).cost
	})
	return queue{q}
}

func (q *queue) enqueue(s state, cost int) {
	i := node{s, cost}
	q.Enqueue(i)
}

func (q *queue) dequeue() (state, int) {
	i, _ := q.Dequeue()
	return i.(node).state, i.(node).cost
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

func (d day16) singlePathDijkstra() int {
	q := newQueue()
	q.enqueue(d.start, 0)
	dist := make(visited)
	dist[d.start] = 0

	for {
		s, cost := q.dequeue()
		if s.x == d.end.x && s.y == d.end.y {
			return cost
		}

		for _, newNode := range d.moves(s) {
			c := cost + newNode.cost
			if dist.get(newNode.state) <= c {
				continue
			}
			dist[newNode.state] = c
			q.enqueue(newNode.state, c)
		}
	}
}

func (d day16) moves(s state) []node {
	var result []node

	if d.grid[s.x+s.direction.dx][s.y+s.direction.dy] != '#' {
		result = append(result, node{state{s.x + s.direction.dx, s.y + s.direction.dy, s.direction}, 1})
	}

	for _, dir := range turnMap[s.direction] {
		if d.grid[s.x+dir.dx][s.y+dir.dy] != '#' {
			result = append(result, node{state{s.x, s.y, dir}, 1000})
		}
	}

	return result
}

func (d day16) Part1() int {
	cost := d.singlePathDijkstra()

	return cost
}

func (d day16) Part2() int {
	return 0
}

func main() {
	d := NewDay16(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
