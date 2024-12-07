package main

import (
	"bytes"
	"path/filepath"
	"strings"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day06 struct {
	day.DayInput
}

type direction [2]int

type position struct {
	row, column int
	face        direction
}

type patrolMap [][]byte

var (
	rotateMap = map[direction]direction{
		{0, 1}:  {1, 0},
		{1, 0}:  {0, -1},
		{0, -1}: {-1, 0},
		{-1, 0}: {0, 1},
	}
)

func NewDay06(inputFile string) Day06 {
	return Day06{day.DayInput(inputFile)}
}

func (p *position) rotate() {
	p.face = rotateMap[p.face]
}

func (p *position) move() {
	p.row += p.face[0]
	p.column += p.face[1]
}

func (p position) lookRight() position {
	return position{p.row, p.column, rotateMap[p.face]}
}

func (p position) ahead() position {
	return position{p.row + p.face[0], p.column + p.face[1], p.face}
}

func (p position) on(m patrolMap) bool {
	return m[p.row][p.column] != 'O'
}

func (p position) blocked(m patrolMap) bool {
	return m[p.row+p.face[0]][p.column+p.face[1]] == '#'
}

func parsePatrolMap(input []string) (patrolMap, position) {
	result := make(patrolMap, len(input)+2)
	var pos position

	result[0] = []byte(strings.Repeat("O", len(input[0])+2))
	result[len(input)+1] = result[0]
	for i, line := range input {
		result[i+1] = []byte("O" + line + "O")
		c := bytes.IndexByte(result[i+1], '^')
		if c != -1 {
			pos = position{i + 1, c, [2]int{-1, 0}}
		}
	}
	return result, pos
}

func (p position) visits(m patrolMap) map[direction]struct{} {
	visited := make(map[direction]struct{})

	for p.on(m) {
		visited[[2]int{p.row, p.column}] = struct{}{}

		for p.blocked(m) {
			p.rotate()
		}

		p.move()
	}

	return visited
}

func (p position) loop(m patrolMap) bool {
	turns := make(map[position]struct{})

	for p.on(m) {
		if _, ok := turns[p]; ok {
			return true
		}

		if p.blocked(m) {
			turns[p] = struct{}{}
			for p.blocked(m) {
				p.rotate()
			}
		}

		p.move()
	}

	return false
}

func (d Day06) Part1() int {
	lines, _ := d.ReadLines()

	patrolMap, guard := parsePatrolMap(lines)

	visited := guard.visits(patrolMap)

	return len(visited)
}

func (d Day06) Part2() int {
	lines, _ := d.ReadLines()

	patrolMap, guard := parsePatrolMap(lines)

	visited := guard.visits(patrolMap)
	delete(visited, direction{guard.row, guard.column})

	obstruct := make(map[direction]struct{})

	for v := range visited {
		patrolMap[v[0]][v[1]] = '#'
		if guard.loop(patrolMap) {
			obstruct[v] = struct{}{}
		}
		patrolMap[v[0]][v[1]] = '.'
	}

	return len(obstruct)
}

func main() {
	d := NewDay06(filepath.Join(projectpath.Root, "cmd", "day06", "input.txt"))

	day.Solve(d)
}
