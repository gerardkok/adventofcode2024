package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day04 struct {
	day.DayInput
}

type wordSearch [][]byte

type direction struct {
	dRow, dColumn int
}

var directions = []direction{
	{0, 1},   // east
	{1, 1},   // southeast
	{1, 0},   // south
	{1, -1},  // southwest
	{0, -1},  // west
	{-1, -1}, // northwest
	{-1, 0},  // north
	{-1, 1},  // northeast
}

func NewDay04(inputFile string) Day04 {
	return Day04{day.DayInput(inputFile)}
}

func (w wordSearch) isMAS(r, c int, dir direction) bool {
	return w[r+dir.dRow][c+dir.dColumn] == 'M' &&
		w[r+2*dir.dRow][c+2*dir.dColumn] == 'A' &&
		w[r+3*dir.dRow][c+3*dir.dColumn] == 'S'
}

func (w wordSearch) isX(r, c int) bool {
	if w[r][c] != 'A' {
		return false
	}

	return (w[r-1][c-1] == 'M' && w[r+1][c+1] == 'S' ||
		w[r-1][c-1] == 'S' && w[r+1][c+1] == 'M') &&
		(w[r-1][c+1] == 'M' && w[r+1][c-1] == 'S' ||
			w[r-1][c+1] == 'S' && w[r+1][c-1] == 'M')
}

func (w wordSearch) countXMAS(r, c int) int {
	if w[r][c] != 'X' {
		return 0
	}

	result := 0

	for _, dir := range directions {
		if w.isMAS(r, c, dir) {
			result++
		}
	}

	return result
}

func (w wordSearch) print() {
	for _, r := range w {
		fmt.Println(string(r))
	}
}

func makeWordSearch(input []string) wordSearch {
	result := make(wordSearch, len(input)+2)
	result[0] = []byte(strings.Repeat(".", len(input[0])+2))
	result[len(input)+1] = result[0]
	for i, line := range input {
		result[i+1] = []byte("." + line + ".")
	}
	return result
}

func (d Day04) Part1() int {
	lines, _ := d.ReadLines()
	w := makeWordSearch(lines)

	xmas := 0

	for r := range len(w) {
		for c := range len(w[0]) {
			xmas += w.countXMAS(r, c)
		}
	}

	return xmas
}

func (d Day04) Part2() int {
	lines, _ := d.ReadLines()
	w := makeWordSearch(lines)

	xmas := 0

	for r := range len(w) {
		for c := range len(w[0]) {
			if w.isX(r, c) {
				xmas++
			}
		}
	}

	return xmas
}

func main() {
	d := NewDay04(filepath.Join(projectpath.Root, "cmd", "day04", "input.txt"))

	day.Solve(d)
}
