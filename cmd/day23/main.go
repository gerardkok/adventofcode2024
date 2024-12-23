package main

import (
	"fmt"
	"iter"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type connection struct {
	a, b string
}

type day23 struct {
	connections map[string]map[string]struct{}
}

func NewDay23(opts ...day.Option) day23 {
	input := day.NewDayInput(path, opts...)

	lines := input.ReadLines()

	connections := make(map[string]map[string]struct{})

	for _, line := range lines {
		a, b, _ := strings.Cut(line, "-")

		if _, ok := connections[a]; !ok {
			connections[a] = make(map[string]struct{})
		}
		if _, ok := connections[b]; !ok {
			connections[b] = make(map[string]struct{})
		}
		connections[a][b] = struct{}{}
		connections[b][a] = struct{}{}
	}

	return day23{connections}
}

func (d day23) networks() iter.Seq[[3]string] {
	return func(yield func([3]string) bool) {
		for a, connections := range d.connections {
			for b := range connections {
				if b > a {
					continue
				}
				for c := range d.connections[b] {
					if c > b {
						continue
					}
					if _, ok := d.connections[c][a]; ok {
						network := [3]string{a, b, c}
						if !yield(network) {
							return
						}
					}
				}
			}
		}
	}
}

func intersect(s, t map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})

	for k := range s {
		if _, ok := t[k]; ok {
			result[k] = struct{}{}
		}
	}

	return result
}

func include(s map[string]struct{}, v string) map[string]struct{} {
	result := map[string]struct{}{v: {}}

	maps.Copy(result, s)

	return result
}

func (d day23) bronKerbosch(R, P, X map[string]struct{}) iter.Seq[map[string]struct{}] {
	return func(yield func(map[string]struct{}) bool) {
		if len(P) == 0 && len(X) == 0 {
			if !yield(R) {
				return
			}
		}

		// stay on safe side, and construct new slice of keys of P, because P will be modified through the loop
		for _, v := range slices.Collect(maps.Keys(P)) {
			for c := range d.bronKerbosch(include(R, v), intersect(P, d.connections[v]), intersect(X, d.connections[v])) {
				if !yield(c) {
					return
				}
			}

			delete(P, v)
			X[v] = struct{}{}
		}
	}
}

func startsWithT(computer string) bool {
	return computer[0] == 't'
}

func hasStartsWithT(network [3]string) bool {
	return startsWithT(network[0]) || startsWithT(network[1]) || startsWithT(network[2])
}

func (d day23) Part1() int {
	result := 0

	for network := range d.networks() {
		if hasStartsWithT(network) {
			result++
		}
	}

	return result
}

func networkName(network map[string]struct{}) string {
	computers := slices.Collect(maps.Keys(network))
	sort.Strings(computers)

	return strings.Join(computers, ",")
}

func (d day23) Part2() string {
	R, P, X := make(map[string]struct{}), make(map[string]struct{}), make(map[string]struct{})
	for computer := range d.connections {
		P[computer] = struct{}{}
	}

	cliques := slices.Collect(d.bronKerbosch(R, P, X))

	max := slices.MaxFunc(cliques, func(a, b map[string]struct{}) int {
		return len(a) - len(b)
	})

	return networkName(max)
}

func main() {
	d := NewDay23(day.FromArgs(os.Args[1:]))

	fmt.Println(d.Part1())
	fmt.Println(d.Part2())
}
