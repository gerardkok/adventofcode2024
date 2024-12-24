package main

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type wire int

type operator int

type gate struct {
	operator string
	output   string
}

var (
	operatorMap = map[string]map[[2]wire]wire{
		"AND": {
			{wire(0), wire(0)}: wire(0),
			{wire(0), wire(1)}: wire(0),
			{wire(1), wire(0)}: wire(0),
			{wire(1), wire(1)}: wire(1),
		},
		"OR": {
			{wire(0), wire(0)}: wire(0),
			{wire(0), wire(1)}: wire(1),
			{wire(1), wire(0)}: wire(1),
			{wire(1), wire(1)}: wire(1),
		},
		"XOR": {
			{wire(0), wire(0)}: wire(0),
			{wire(0), wire(1)}: wire(1),
			{wire(1), wire(0)}: wire(1),
			{wire(1), wire(1)}: wire(0),
		},
	}
)

type day24 struct {
	wires map[string]wire
	gates map[[2]string][]gate
}

func parseWire(s string) wire {
	return wire(conv.MustAtoi(s))
}

func parseWires(lines string) map[string]wire {
	result := make(map[string]wire)

	for _, w := range strings.Split(strings.TrimSpace(lines), "\n") {
		name, value, _ := strings.Cut(w, ": ")
		result[name] = parseWire(value)
	}

	return result
}

func parseGates(lines string) map[[2]string][]gate {
	result := make(map[[2]string][]gate)

	for _, g := range strings.Split(strings.TrimSpace(lines), "\n") {
		input, output, _ := strings.Cut(g, " -> ")
		i := strings.Split(input, " ")
		result[[2]string{i[0], i[2]}] = append(result[[2]string{i[0], i[2]}], gate{i[1], output})
	}

	return result
}

func parseInput(input day.DayInput) (map[string]wire, map[[2]string][]gate) {
	lines := input.ReadInput()

	parts := bytes.Split(lines, []byte{'\n', '\n'})

	wires := parseWires(string(parts[0]))
	gates := parseGates(string(parts[1]))

	return wires, gates
}

func NewDay24(opts ...day.Option) day24 {
	input := day.NewDayInput(path, opts...)

	wires, gates := parseInput(input)

	return day24{wires, gates}
}

func maxWire(suffix byte, wires []string) int {
	max := 0

	for _, wire := range wires {
		if wire[0] == suffix {
			z := conv.MustAtoi(wire[1:])
			if z > max {
				max = z
			}
		}
	}

	return max
}

func (d day24) simulate() int {
	done := false

	for !done {
		done = true
		for input, gates := range d.gates {
			for _, gate := range gates {
				i1, ok1 := d.wires[input[0]]
				i2, ok2 := d.wires[input[1]]
				_, ok3 := d.wires[gate.output]
				if ok1 && ok2 && !ok3 {
					d.wires[gate.output] = operatorMap[gate.operator][[2]wire{i1, i2}]
					done = false
				}
			}
		}
	}

	result := 0

	for i := range maxWire('z', slices.Collect(maps.Keys(d.wires))) + 1 {
		zWire := fmt.Sprintf("z%02d", i)
		result += int(d.wires[zWire]) << i
	}

	return result
}

func (d day24) Part1() int {
	return d.simulate()
}

func (d day24) Part2() int {
	return 0
}

func main() {
	d := NewDay24(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
