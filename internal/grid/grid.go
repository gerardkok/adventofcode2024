package grid

import (
	"container/heap"
	"iter"
	"math"
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

func Bfs[T comparable](start T, neighbours func(T) []T) iter.Seq[[2]T] {
	return func(yield func([2]T) bool) {
		todo := []T{start}
		parent := make(map[T]T)

		for len(todo) > 0 {
			p := todo[0]
			todo = todo[1:]

			if !yield([2]T{p, parent[p]}) {
				return
			}

			for _, n := range neighbours(p) {
				if _, ok := parent[n]; ok {
					continue
				}
				parent[n] = p
				todo = append(todo, n)
			}
		}
	}
}

func get[T comparable](d map[T]int, p T) int {
	if v, ok := d[p]; ok {
		return v
	}

	return math.MaxInt
}

func ShortestPath[T comparable](source T, neighbours func(T) []Edge[T]) (map[T]int, map[T]T) {
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

func AllShortestPaths[T comparable](source T, neighbours func(T) []Edge[T]) (map[T]int, map[T]map[T]struct{}) {
	dist := map[T]int{source: 0}
	prev := make(map[T]map[T]struct{})

	pq := make(PriorityQueue[T], 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item[T]{vertex: source, dist: 0})

	for pq.Len() > 0 {
		u := heap.Pop(&pq).(*Item[T]).vertex

		for _, edge := range neighbours(u) {
			v := edge.To
			weight := edge.Weight
			if dist[u]+weight <= get(dist, v) {
				dist[v] = dist[u] + weight
				heap.Push(&pq, &Item[T]{vertex: v, dist: dist[v]})
				if _, ok := prev[v]; !ok {
					prev[v] = make(map[T]struct{})
				}
				prev[v][u] = struct{}{}
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
