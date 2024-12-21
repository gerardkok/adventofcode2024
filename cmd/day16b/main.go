package main

import (
	"bytes"
	"fmt"
	"maps"
	"math"
	"os"
	"path/filepath"
	"runtime"

	pq "github.com/emirpasic/gods/queues/priorityqueue"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/grid"
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

func (d day16) paths(p map[state]map[state]struct{}, endStates []state) int {
	result := make(map[tile]struct{})
	for _, e := range endStates {
		maps.Copy(result, d.path(p, e))
	}
	return len(result)
}

func (d day16) path(p map[state]map[state]struct{}, e state) map[tile]struct{} {
	result := map[tile]struct{}{{e.x, e.y}: {}}

	if e == d.start {
		return result
	}

	prevs := p[e]
	fmt.Printf("l: %d\n", len(prevs))
	for k := range prevs {
		maps.Copy(result, d.path(p, k))
	}

	return result
}

// func (d day16) singlePathDijkstra2() (map[state]map[state]struct{}, []state) {
// 	q := newQueue()
// 	q.enqueue(d.start, 0)
// 	seen := make(map[state]struct{})
// 	dist := make(visited)
// 	dist[d.start] = 0
// 	prev := make(map[state]map[state]struct{})

// 	min := math.MaxInt
// 	var endStates []state

// 	for {
// 		s, cost := q.dequeue()
// 		if cost > min {
// 			return prev, endStates
// 		}
// 		if s.x == d.end.x && s.y == d.end.y {
// 			min = cost
// 			endStates = append(endStates, s)
// 		}

// 		seen[s] = struct{}{}

// 		for _, newNode := range d.moves(s) {
// 			if _, ok := seen[newNode.state]; ok {
// 				continue
// 			}

// 			c := cost + newNode.cost
// 			if c <= dist.get(newNode.state) {
// 				dist[newNode.state] = c
// 				q.enqueue(newNode.state, c)
// 				if _, ok := prev[newNode.state]; !ok {
// 					prev[newNode.state] = make(map[state]struct{})
// 				}
// 				prev[newNode.state][s] = struct{}{}
// 			}
// 		}
// 	}
// }

func (d day16) dijkstra() int {
	dist, _, _ := grid.ShortestPath(d.start, d.neighbours, d.isStop)
	return dist
}

func (d day16) isStop(s state) bool {
	return d.end.x == s.x && d.end.y == s.y
}

func (d day16) neighbours(s state) []grid.Edge[state] {
	var result []grid.Edge[state]

	if d.grid[s.x+s.direction.dx][s.y+s.direction.dy] != '#' {
		result = append(result, grid.Edge[state]{
			To:     state{s.x + s.direction.dx, s.y + s.direction.dy, s.direction},
			Weight: 1,
		})
	}

	for _, dir := range turnMap[s.direction] {
		if d.grid[s.x+dir.dx][s.y+dir.dy] != '#' {
			result = append(result, grid.Edge[state]{
				To:     state{s.x, s.y, dir},
				Weight: 1000,
			})
		}
	}

	return result
}

func (d day16) Part1() int {
	return d.dijkstra()
}

func (d day16) Part2() int {
	// p, e := d.singlePathDijkstra2()
	// fmt.Println(p)

	// path := d.paths(p, e)

	// fmt.Println(path)

	// return path
	return 0
}

func main() {
	d := NewDay16(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
