package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	pq "github.com/emirpasic/gods/queues/priorityqueue"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day18 struct {
	grid   [][]byte
	spots  []spot
	fallen int
}

type spot struct {
	r, c int
}

type queue struct {
	*pq.Queue
}

type node struct {
	spot spot
	cost int
}

type direction struct {
	dr, dc int
}

type visited map[spot]int

func (v visited) get(s spot) int {
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

func (q *queue) enqueue(s spot, cost int) {
	i := node{s, cost}
	q.Enqueue(i)
}

func (q *queue) dequeue() (spot, int) {
	i, _ := q.Dequeue()
	return i.(node).spot, i.(node).cost
}

func parseGrid(lines []string, size, fallenBytes int) ([][]byte, []spot) {
	result := make([][]byte, size)
	var spots []spot

	for i := range result {
		result[i] = bytes.Repeat([]byte{'.'}, size)
	}

	for i, line := range lines {
		c, r, _ := strings.Cut(line, ",")
		s := spot{conv.MustAtoi(r), conv.MustAtoi(c)}
		spots = append(spots, s)
		if i < fallenBytes {
			result[s.r][s.c] = '#'
		}
	}

	return result, spots
}

func NewDay18(size, fallenBytes int, opts ...day.Option) day18 {
	input := day.NewDayInput(path, opts...)

	lines := input.ReadLines()

	grid, spots := parseGrid(lines, size, fallenBytes)

	return day18{grid, spots, fallenBytes}
}

func (d day18) singlePathDijkstra() int {
	q := newQueue()
	q.enqueue(spot{0, 0}, 0)
	seen := make(map[spot]struct{})
	dist := make(visited)
	dist[spot{0, 0}] = 0

	for q.Size() > 0 {
		s, cost := q.dequeue()
		if s.r == len(d.grid)-1 && s.c == len(d.grid[0])-1 {
			return cost
		}

		seen[s] = struct{}{}

		for _, newNode := range d.moves(s) {
			if _, ok := seen[newNode.spot]; ok {
				continue
			}

			c := cost + newNode.cost
			if c < dist.get(newNode.spot) {
				dist[newNode.spot] = c
				q.enqueue(newNode.spot, c)
			}
		}
	}

	return -1
}

func (d day18) moves(s spot) []node {
	var result []node

	for _, dir := range []direction{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
		r := s.r + dir.dr
		c := s.c + dir.dc
		if r < 0 || r >= len(d.grid) || c < 0 || c >= len(d.grid[0]) {
			continue
		}

		if d.grid[r][c] == '#' {
			continue
		}

		result = append(result, node{spot{r, c}, 1})
	}

	return result
}

func (d day18) Part1() int {
	return d.singlePathDijkstra()
}

func (d day18) Part2() string {
	for _, s := range d.spots[d.fallen:] {
		d.grid[s.r][s.c] = '#'
		if d.singlePathDijkstra() == -1 {
			output := fmt.Sprintf("%d,%d", s.c, s.r)
			return output
		}
	}
	return "71,71"
}

func main() {
	d := NewDay18(71, 1024, day.FromArgs(os.Args[1:]))

	fmt.Println(d.Part1())
	fmt.Println(d.Part2())
}
