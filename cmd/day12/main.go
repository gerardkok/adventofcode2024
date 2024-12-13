package main

import (
	"bytes"
	"fmt"
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

type position struct {
	x, y int
}

type direction struct {
	dx, dy int
}

// all plots in a region with the number of neighbours
type region map[position]struct{}

type garden struct {
	regions          []region
	positionToRegion map[position]int
}

type square struct {
	nw, ne, se, sw position
}

var directions = []direction{{0, 1}, {1, 0}, {-1, 0}, {0, -1}}

func NewDay12(opts ...day.Option) day12 {
	return day12{day.NewDayInput(path, opts...)}
}

func (p position) to(d direction) position {
	return position{p.x + d.dx, p.y + d.dy}
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

func (g grid) plant(p position) byte {
	return g[p.x][p.y]
}

func (g grid) border(p position) bool {
	return g.plant(p) == '#'
}

func (g grid) neighbours(p position) iter.Seq[position] {
	return func(yield func(position) bool) {
		if g.border(p) {
			return
		}
		for _, dir := range directions {
			next := p.to(dir)
			if g.plant(next) == g.plant(p) && !yield(next) {
				return
			}
		}
	}
}

func (g grid) regionCost(start position) (map[position]struct{}, int) {
	// if g.border(start) {
	// 	return map[position]struct{}{}, 0
	// }

	seen := make(map[position]struct{})
	todo := []position{start}
	perimeter := 0

	for len(todo) > 0 {
		p := todo[0]
		todo = todo[1:]

		if _, ok := seen[p]; ok {
			continue
		}

		seen[p] = struct{}{}

		perimeter += 4
		for n := range g.neighbours(p) {
			todo = append(todo, n)
			perimeter--
		}
	}

	return seen, perimeter
}

func (g grid) borderRegion() region {
	result := make(region)

	for x := range len(g) {
		if x > 0 && x < len(g) {
			continue
		}
		for y := range len(g[0]) {
			if y > 0 && y < len(g[0]) {
				continue
			}

			p := position{x, y}
			result[p] = struct{}{}
		}
	}

	return result
}

func (g grid) regions() []region {
	result := []region{g.borderRegion()}

	seen := make(map[position]struct{})

	for x := range len(g) {
		for y := range len(g[0]) {
			p := position{x, y}
			if g.border(p) {
				continue
			}

			if _, ok := seen[p]; ok {
				continue
			}

			region, _ := g.regionCost(p)
			maps.Copy(seen, region)
			result = append(result, region)
		}
	}

	return result
}

func (g grid) cost() int {
	seen := make(map[position]struct{})
	result := 0

	for x := range len(g) {
		for y := range len(g[0]) {
			p := position{x, y}
			if g.border(p) {
				continue
			}
			if _, ok := seen[p]; ok {
				continue
			}

			region, perimeter := g.regionCost(p)
			maps.Copy(seen, region)
			result += len(region) * perimeter
		}
	}

	return result
}

func (g grid) makeGarden() garden {
	regions := g.regions()
	p2r := make(map[position]int)

	for i, region := range regions {
		for p := range region {
			p2r[p] = i
		}
	}

	return garden{regions, p2r}
}

func (s square) rotate() square {
	return square{s.sw, s.nw, s.ne, s.se}
}

func (g garden) convex(s square) bool {
	r := g.positionToRegion[s.nw] != g.positionToRegion[s.ne] &&
		g.positionToRegion[s.nw] != g.positionToRegion[s.sw] // &&
		// g.positionToRegion[s.nw] != g.positionToRegion[s.se]
	return r
}

func (g garden) concave(s square) bool {
	r := g.positionToRegion[s.nw] == g.positionToRegion[s.ne] &&
		g.positionToRegion[s.nw] == g.positionToRegion[s.sw] &&
		g.positionToRegion[s.nw] != g.positionToRegion[s.se]
	return r
}

func (g garden) corner(s square) bool {
	return g.convex(s) || g.concave(s)
}

func (p position) square() square {
	nw := p
	ne := position{p.x, p.y + 1}
	se := position{p.x + 1, p.y + 1}
	sw := position{p.x + 1, p.y}
	return square{nw, ne, se, sw}
}

func (s square) print() string {
	return fmt.Sprintf("nw: %v, ne: %v\nsw: %v, se: %v\n", s.nw, s.ne, s.sw, s.se)
}

func (g grid) corners() int {
	garden := g.makeGarden()
	corners := make(map[int]int)

	for x := range len(g) - 1 {
		for y := range len(g[0]) - 1 {
			p := position{x, y}

			sq := p.square()
			for range 4 {
				corner := garden.corner(sq)
				if corner {
					corners[garden.positionToRegion[sq.nw]]++
				}

				sq = sq.rotate()
			}
		}
	}

	cost := 0
	for i, region := range garden.regions {
		if i == 0 {
			continue
		}
		fmt.Printf("area: %d, corners: %d\n", len(region), corners[i])
		cost += len(region) * corners[i]
	}

	return cost
}

func (d day12) Part1() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	return grid.cost()
}

func (d day12) Part2() int {
	lines := d.ReadLines()

	grid := parseGrid(lines)

	fmt.Println(grid)

	return grid.corners()
}

func main() {
	d := NewDay12(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
