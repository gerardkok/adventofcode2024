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

type Grid[T comparable] struct {
	Points [][]T
}

func (g Grid[T]) At(p Point) T {
	return g.Points[p.X][p.Y]
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
			if next.X < 0 || next.X >= len(g.Points) || next.Y < 0 || next.Y >= len(g.Points[0]) {
				continue
			}

			if !yield(next) {
				return
			}
		}
	}
}

func (g Grid[T]) PointsIter() iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for x, row := range g.Points {
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
	p := make([][]T, len(points)+2)

	p[0] = slices.Repeat([]T{border}, len(points[0])+2)
	for i, l := range points {
		p[i+1] = slices.Concat([]T{border}, l, []T{border})
	}
	p[len(points)+1] = slices.Repeat([]T{border}, len(points[0])+2)

	return BorderedGrid[T]{Grid[T]{p}, border}
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

func (b BorderedGrid[T]) PointsIter() iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for x, row := range b.Points {
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
