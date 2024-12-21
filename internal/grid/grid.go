package grid

import (
	"container/heap"
	"iter"
	"math"
	"slices"
)

type Point struct {
	X, Y int
}

type Direction struct {
	Dx, Dy int
}

func (p Point) To(d Direction) Point {
	return Point{p.X + d.Dx, p.Y + d.Dy}
}

type Edge[T any] struct {
	To     T
	Weight int
}

type Item[T any] struct {
	vertex T
	dist   int
	index  int
}

type PriorityQueue[T any] []*Item[T]

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*Item[T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

type Grid[T comparable] [][]T

func (g Grid[T]) At(p Point) T {
	return g[p.X][p.Y]
}

func Bfs(start Point, appendSeq func(p Point) iter.Seq[Point]) iter.Seq[Point] {
	return func(yield func(Point) bool) {
		seen := make(map[Point]struct{})
		todo := []Point{start}

		for len(todo) > 0 {
			p := todo[0]
			todo = todo[1:]

			if _, ok := seen[p]; ok {
				continue
			}

			if !yield(p) {
				return
			}

			seen[p] = struct{}{}

			todo = slices.AppendSeq(todo, appendSeq(p))
		}
	}
}

func Bfs2[T comparable](start T, neighbours func(T) []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		todo := []T{start}
		seen := make(map[T]struct{})

		for len(todo) > 0 {
			p := todo[0]
			todo = todo[1:]

			if _, ok := seen[p]; ok {
				continue
			}

			if !yield(p) {
				return
			}

			seen[p] = struct{}{}

			todo = append(todo, neighbours(p)...)
		}
	}
}

func get[T comparable](d map[T]int, p T) int {
	if v, ok := d[p]; ok {
		return v
	}

	return math.MaxInt
}

func ShortestPath[T comparable](source T, neighbours func(T) []Edge[T], isEnd func(T) bool) (int, map[T]T, T) {
	dist := make(map[T]int)
	dist[source] = 0

	pq := make(PriorityQueue[T], 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item[T]{vertex: source, dist: 0})

	prev := make(map[T]T)

	for {
		u := heap.Pop(&pq).(*Item[T]).vertex
		if isEnd(u) {
			return dist[u], prev, u
		}

		for _, edge := range neighbours(u) {
			v := edge.To
			weight := edge.Weight
			if dist[u]+weight < get(dist, v) {
				dist[v] = dist[u] + weight
				heap.Push(&pq, &Item[T]{vertex: v, dist: dist[v]})
				prev[v] = u
			}
		}
	}
}

func Dijkstra[T comparable](source T, neighbours func(T) []Edge[T]) (map[T]int, map[T]T) {
	dist := map[T]int{source: 0}
	prev := make(map[T]T)

	pq := make(PriorityQueue[T], 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item[T]{vertex: source, dist: 0})

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Item[T]).vertex

		for _, edge := range neighbours(u) {
			v := edge.To
			weight := edge.Weight
			if dist[u]+weight < get(dist, v) {
				dist[v] = dist[u] + weight
				heap.Push(&pq, &Item[T]{vertex: v, dist: dist[v]})
				prev[v] = u
			}
		}
	}

	return dist, prev
}

func (g Grid[T]) Neighbours4(p Point) iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for _, dir := range []Direction{{0, 1}, {1, 0}, {-1, 0}, {0, -1}} {
			next := p.To(dir)
			if next.X < 0 || next.X >= len(g) || next.Y < 0 || next.Y >= len(g[0]) {
				continue
			}

			if !yield(next) {
				return
			}
		}
	}
}

func (g Grid[T]) Points() iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for x, row := range g {
			for y := range row {
				if !yield(Point{x, y}) {
					return
				}
			}
		}
	}
}

type BorderedGrid[T comparable] struct {
	Grid[T]
	border T
}

func NewBorderedGrid[T comparable](points [][]T, border T) BorderedGrid[T] {
	g := make(Grid[T], len(points)+2)

	g[0] = slices.Repeat([]T{border}, len(points[0])+2)
	for i, l := range points {
		g[i+1] = slices.Concat([]T{border}, l, []T{border})
	}
	g[len(points)+1] = slices.Repeat([]T{border}, len(points[0])+2)

	return BorderedGrid[T]{g, border}
}

func (b BorderedGrid[T]) Neighbours4(p Point) iter.Seq[Point] {
	return func(yield func(Point) bool) {
		if b.At(p) == b.border {
			return
		}

		for _, dir := range []Direction{{0, 1}, {1, 0}, {-1, 0}, {0, -1}} {
			next := p.To(dir)
			if b.At(next) == b.border {
				continue
			}

			if !yield(next) {
				return
			}
		}
	}
}

func (b BorderedGrid[T]) Points() iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for x, row := range b.Grid {
			for y, p := range row {
				if p == b.border {
					continue
				}

				if !yield(Point{x, y}) {
					return
				}
			}
		}
	}
}
