package main

import (
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

type day08 struct {
	day.DayInput
}

func NewDay08(opts ...day.Option) day08 {
	return day08{day.NewDayInput(path, opts...)}
}

type location [2]int

type city struct {
	cityMap  [][]byte
	antennae map[byte][]location
}

type antinodesFn func(location, location) []location

func parseCity(input []string) city {
	cityMap := make([][]byte, len(input))
	antennae := make(map[byte][]location)

	for i, line := range input {
		cityMap[i] = []byte(line)
		for j, c := range cityMap[i] {
			if c != '.' {
				antennae[c] = append(antennae[c], location{i, j})
			}
		}
	}

	return city{cityMap, antennae}
}

func (c city) covers(l location) bool {
	return l[0] >= 0 && l[0] < len(c.cityMap) &&
		l[1] >= 0 && l[1] < len(c.cityMap[0])
}

func (c city) antinodesPart1(a, b location) []location {
	var result []location

	dr := b[0] - a[0]
	dc := b[1] - a[1]
	n := location{a[0] - dr, a[1] - dc}
	if c.covers(n) {
		result = append(result, n)
	}
	m := location{b[0] + dr, b[1] + dc}
	if c.covers(m) {
		result = append(result, m)
	}

	return result
}

func (ct city) antinodesPart2(a, b location) []location {
	var result []location

	dr := b[0] - a[0]
	dc := b[1] - a[1]
	for r, c := a[0], a[1]; ct.covers(location{r, c}); r, c = r-dr, c-dc {
		result = append(result, location{r, c})
	}
	for r, c := b[0], b[1]; ct.covers(location{r, c}); r, c = r+dr, c+dc {
		result = append(result, location{r, c})
	}

	return result
}

func (c city) antennaAntinodes(antennae []location, antinodes antinodesFn) map[location]struct{} {
	result := make(map[location]struct{})

	for i, a := range antennae[:len(antennae)-1] {
		for _, b := range antennae[i+1:] {
			for _, antinode := range antinodes(a, b) {
				result[antinode] = struct{}{}
			}
		}
	}

	return result
}

func (c city) antinodes(antinodes antinodesFn) map[location]struct{} {
	result := make(map[location]struct{})

	for _, antennae := range c.antennae {
		antinodes := c.antennaAntinodes(antennae, antinodes)
		maps.Copy(result, antinodes)
	}

	return result
}

func (d day08) Part1() int {
	lines := d.ReadLines()

	city := parseCity(lines)

	antinodes := city.antinodes(city.antinodesPart1)

	return len(antinodes)
}

func (d day08) Part2() int {
	lines := d.ReadLines()

	city := parseCity(lines)

	antinodes := city.antinodes(city.antinodesPart2)

	return len(antinodes)
}

func main() {
	d := NewDay08(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
