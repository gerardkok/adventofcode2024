package grid

import (
	"iter"
	"slices"
)

type Direction struct {
	Dx, Dy int
}

type Point struct {
	X, Y int
}

type Grid interface {
	Neighbours(Point) iter.Seq[Point]
	Points() iter.Seq[Point]
}

func (p Point) To(d Direction) Point {
	return Point{p.X + d.Dx, p.Y + d.Dy}
}

func Bfs(g Grid, start Point) iter.Seq[Point] {
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

			todo = slices.AppendSeq(todo, g.Neighbours(p))
		}
	}
}

type BorderedGrid[T comparable] struct {
	points [][]T
	border T
}

func MakeBorderedGrid[T comparable](points [][]T, border T) BorderedGrid[T] {
	p := make([][]T, len(points)+2)

	p[0] = slices.Repeat([]T{border}, len(points[0])+2)
	for i, l := range points {
		p[i] = slices.Concat([]T{border}, l, []T{border})
	}
	p[len(points)+1] = slices.Repeat([]T{border}, len(points[0])+2)

	return BorderedGrid[T]{p, border}
}

func (b BorderedGrid[T]) Neighbours(p Point) iter.Seq[Point] {
	return func(yield func(Point) bool) {
		if b.points[p.X][p.Y] == b.border {
			return
		}

		for _, dir := range []Direction{{0, 1}, {1, 0}, {-1, 0}, {0, -1}} {
			next := p.To(dir)
			if b.points[next.X][next.Y] != b.border && !yield(next) {
				return
			}
		}
	}
}

func (b BorderedGrid[T]) Points() iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for x := range len(b.points) {
			for y := range len(b.points[0]) {
				p := Point{x, y}
				if b.points[x][y] == b.border {
					continue
				}

				if !yield(p) {
					return
				}
			}
		}
	}
}
