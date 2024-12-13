package main

import (
	"bytes"
	"iter"
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

type day12 struct {
	day.DayInput
}

type grid [][]byte

type plot struct {
	x, y int
}

type direction struct {
	dx, dy int
}

type region struct {
	plots           map[plot]struct{}
	area, perimeter int
}

type regionIDs [][]int

type square struct {
	nw, ne, se, sw int
}

func NewDay12(opts ...day.Option) day12 {
	return day12{day.NewDayInput(path, opts...)}
}

func (p plot) to(d direction) plot {
	return plot{p.x + d.dx, p.y + d.dy}
}

func parseGrid(input []string) grid {
	result := make(grid, len(input)+2)

	result[0] = bytes.Repeat([]byte{'#'}, len(input[0])+2)
	for x, line := range input {
		result[x+1] = []byte("#" + line + "#")
	}
	result[len(input)+1] = result[0]

	return result
}

func (g grid) plant(p plot) byte {
	return g[p.x][p.y]
}

func (g grid) border(p plot) bool {
	return g.plant(p) == '#'
}

func (g grid) neighbours(p plot) []plot {
	var result []plot

	if g.border(p) {
		return result
	}

	for _, dir := range []direction{{0, 1}, {1, 0}, {-1, 0}, {0, -1}} {
		next := p.to(dir)
		if g.plant(next) == g.plant(p) {
			result = append(result, next)
		}
	}

	return result
}

func (g grid) region(start plot) region {
	seen := make(map[plot]struct{})
	todo := []plot{start}
	nNeighbours := 0

	for len(todo) > 0 {
		p := todo[0]
		todo = todo[1:]

		if _, ok := seen[p]; ok {
			continue
		}

		seen[p] = struct{}{}

		neighbours := g.neighbours(p)
		todo = append(todo, neighbours...)
		nNeighbours += len(neighbours)
	}

	area := len(seen)
	perimeter := area*4 - nNeighbours

	return region{seen, area, perimeter}
}

func (g grid) plots() iter.Seq[plot] {
	return func(yield func(plot) bool) {
		for x := range len(g) {
			for y := range len(g[0]) {
				p := plot{x, y}
				if g.border(p) {
					continue
				}

				if !yield(p) {
					return
				}
			}
		}
	}
}

func (g grid) regions() []region {
	result := []region{{}} // fake region for border
	seen := make(map[plot]struct{})

	for p := range g.plots() {
		if _, ok := seen[p]; ok {
			continue
		}

		region := g.region(p)
		maps.Copy(seen, region.plots)
		result = append(result, region)
	}

	return result
}

func (g grid) costPart1() int {
	result := 0

	for _, r := range g.regions() {
		result += r.area * r.perimeter
	}

	return result
}

func (g grid) costPart2() int {
	regions := g.regions()

	regionIDs := g.regionIDs(regions)

	corners := regionIDs.corners()

	result := 0

	for i, region := range regions {
		result += region.area * corners[i]
	}

	return result
}

func (g grid) regionIDs(regions []region) regionIDs {
	result := make(regionIDs, len(g))

	for i := range len(g) {
		result[i] = make([]int, len(g[0]))
	}

	for j, r := range regions {
		for p := range r.plots {
			result[p.x][p.y] = j
		}
	}

	return result
}

func (r regionIDs) squares() iter.Seq[square] {
	return func(yield func(square) bool) {
		for x := range len(r) - 1 {
			for y := range len(r[0]) - 1 {
				sq := r.square(x, y)
				if sq.inside() {
					continue
				}

				for _, s := range sq.rotations() {
					if !yield(s) {
						return
					}
				}
			}
		}
	}
}

func (r regionIDs) square(x, y int) square {
	nw, ne, se, sw := r[x][y], r[x][y+1], r[x+1][y+1], r[x+1][y]
	return square{nw, ne, se, sw}
}

func (r regionIDs) corners() map[int]int {
	result := make(map[int]int)

	for sq := range r.squares() {
		if sq.corner() {
			result[sq.nw]++
		}
	}

	return result
}

func (s square) rotations() []square {
	result := []square{s}

	for range 3 {
		s = s.rotate()

		result = append(result, s)
	}

	return result
}

func (s square) inside() bool {
	return s.nw == s.ne && s.ne == s.se && s.se == s.sw
}

func (s square) rotate() square {
	return square{s.sw, s.nw, s.ne, s.se}
}

func (s square) convex() bool {
	return s.nw != s.ne && s.nw != s.sw
}

func (s square) concave() bool {
	return s.nw == s.ne && s.nw == s.sw && s.nw != s.se
}

func (s square) corner() bool {
	return s.convex() || s.concave()
}

func (d day12) Part1() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	return grid.costPart1()
}

func (d day12) Part2() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	return grid.costPart2()
}

func main() {
	d := NewDay12(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
