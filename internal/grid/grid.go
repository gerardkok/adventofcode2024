package grid

import (
	"iter"
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

type Grid[T comparable] [][]T

func (g Grid[T]) At(p Point) T {
	return g[p.X][p.Y]
}

func (g Grid[T]) Bfs(start Point, appendSeq func(p Point) iter.Seq[Point]) iter.Seq[Point] {
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
