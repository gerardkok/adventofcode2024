package main

import (
	"bytes"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
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

type edge struct {
	label    string
	from, to string
}

type graph struct {
	nodes map[string]string
	edges map[string]edge
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

func makeInputEdges(name string, inputMapping, outputMapping map[string]map[string]struct{}) []edge {
	var result []edge

	for i := range inputMapping[name] {
		for j := range outputMapping[name] {
			result = append(result, edge{name, j, i})
		}
	}

	return result
}

func makeOutputEdges(name string, inputMapping, outputMapping map[string]map[string]struct{}) []edge {
	var result []edge

	for i := range inputMapping[name] {
		for j := range outputMapping[name] {
			result = append(result, edge{name, j, i})
		}
	}

	return result
}

func (d day24) makeGraph() graph {
	inputMapping := make(map[string]map[string]struct{})
	outputMapping := make(map[string]map[string]struct{})
	nodes := make(map[string]string)

	for inputs, gates := range d.gates {
		i1 := inputs[0]
		i2 := inputs[1]
		if i2 < i1 {
			i1, i2 = i2, i1
		}
		if i1[0] == 'x' || i1[0] == 'y' {
			nodes[i1] = i1
			if _, ok := outputMapping[i1]; !ok {
				outputMapping[i1] = make(map[string]struct{})
			}
			outputMapping[i1][i1] = struct{}{}
		}
		if i2[0] == 'x' || i2[0] == 'y' {
			nodes[i2] = i2
			if _, ok := outputMapping[i2]; !ok {
				outputMapping[i2] = make(map[string]struct{})
			}
			outputMapping[i2][i2] = struct{}{}
		}
		for _, gate := range gates {
			name := fmt.Sprintf("%s %s %s", i1, gate.operator, i2)
			if _, ok := inputMapping[i1]; !ok {
				inputMapping[i1] = make(map[string]struct{})
			}
			inputMapping[i1][name] = struct{}{}
			if _, ok := inputMapping[i2]; !ok {
				inputMapping[i2] = make(map[string]struct{})
			}
			inputMapping[i2][name] = struct{}{}
			nodes[name] = gate.operator
			if gate.output[0] == 'z' {
				nodes[gate.output] = gate.output
				if _, ok := outputMapping[gate.output]; !ok {
					inputMapping[gate.output] = make(map[string]struct{})
				}
				inputMapping[gate.output][gate.output] = struct{}{}
			}
			if _, ok := outputMapping[gate.output]; !ok {
				outputMapping[gate.output] = make(map[string]struct{})
			}
			outputMapping[gate.output][name] = struct{}{}
		}
	}

	edges := make(map[string]edge)

	for inputs, gates := range d.gates {
		for i, e := range makeInputEdges(inputs[0], inputMapping, outputMapping) {
			n := fmt.Sprintf("%s%02d", inputs[0], i)
			edges[n] = e
		}
		for i, e := range makeInputEdges(inputs[1], inputMapping, outputMapping) {
			n := fmt.Sprintf("%s%02d", inputs[1], i)
			edges[n] = e
		}
		for _, gate := range gates {
			for i, e := range makeOutputEdges(gate.output, inputMapping, outputMapping) {
				n := fmt.Sprintf("%s%02d", gate.output, i)
				edges[n] = e
			}
		}
	}

	fmt.Println("nodes:")
	for node, label := range nodes {
		fmt.Printf("[%s] %s\n", node, label)
	}
	fmt.Println("edges:")
	for name, edge := range edges {
		fmt.Printf("(%s) [%s] -- %s --> [%s]\n", name, edge.from, edge.label, edge.to)
	}

	return graph{nodes, edges}
}

func (d day24) value(b byte) int {
	result := 0

	for i := range maxWire(b, slices.Collect(maps.Keys(d.wires))) + 1 {
		zWire := fmt.Sprintf("%c%02d", b, i)
		result += int(d.wires[zWire]) << i
	}

	return result
}

func (d day24) Part1() int {
	return d.simulate()
}

func (d day24) Part2() int {
	dot := d.makeGraph()

	file, err := os.Create("graph.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString("digraph {\n    rankdir=\"TB\"\n\n")
	ordered := slices.Collect(maps.Keys(dot.nodes))
	sort.Strings(ordered)
	for _, node := range ordered {
		file.WriteString(fmt.Sprintf("    \"%s\" [label=\"%s\"];\n", node, dot.nodes[node]))
	}
	for _, edge := range dot.edges {
		file.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [label=\"%s\"];\n", edge.from, edge.to, edge.label))
	}
	file.WriteString("}\n")

	file.Sync()

	x := d.value('x')
	y := d.value('y')
	z := d.value('z')
	fmt.Printf("%d + %d = %d, adder says: %d\n", x, y, x+y, z)

	return 0
}

func main() {
	d := NewDay24(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
