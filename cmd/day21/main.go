package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"

	"github.com/cespare/xxhash/v2"
)

type keypad [][]byte

type key struct {
	r, c int
}

type move struct {
	from, to byte
}

type day21 struct {
	codes []string
}

type seq struct {
	xxh   uint64
	level int
}

type memo map[seq]int

type keypadType int

const (
	numerical keypadType = iota
	directional
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)

	keypads = map[keypadType]keypad{
		numerical: {
			{'7', '8', '9'},
			{'4', '5', '6'},
			{'1', '2', '3'},
			{' ', '0', 'A'},
		},
		directional: {
			{' ', '^', 'A'},
			{'<', 'v', '>'},
		},
	}
	keyMap = map[keypadType]map[byte]key{
		numerical: {
			'A': {3, 2},
			'0': {3, 1},
			'1': {2, 0},
			'2': {2, 1},
			'3': {2, 2},
			'4': {1, 0},
			'5': {1, 1},
			'6': {1, 2},
			'7': {0, 0},
			'8': {0, 1},
			'9': {0, 2},
		},
		directional: {
			'>': {1, 2},
			'v': {1, 1},
			'<': {1, 0},
			'A': {0, 2},
			'^': {0, 1},
		},
	}
)

func NewDay21(opts ...day.Option) day21 {
	input := day.NewDayInput(path, opts...)

	codes := input.ReadLines()

	return day21{codes}
}

func printKeypadOpts(keypadOpts map[move][]string) {
	for k, v := range keypadOpts {
		opts := strings.Join(v, ", ")
		fmt.Printf("%c -> %c: %s\n", k.from, k.to, opts)
	}
}

func horizontal(dc int) []byte {
	c := byte('<')
	if dc > 0 {
		c = '>'
	}

	return bytes.Repeat([]byte{c}, conv.Abs(dc))
}

func vertical(dr int) []byte {
	c := byte('^')
	if dr > 0 {
		c = 'v'
	}

	return bytes.Repeat([]byte{c}, conv.Abs(dr))
}

func (kpt keypadType) options(from, to byte) [][]byte {
	f := keyMap[kpt][from]
	t := keyMap[kpt][to]

	dr := t.r - f.r
	dc := t.c - f.c

	if dr == 0 { // horizontal
		return [][]byte{horizontal(dc)}
	}
	if dc == 0 { //vertical
		return [][]byte{vertical(dr)}
	}

	var result [][]byte

	if keypads[kpt][t.r][f.c] != ' ' { // don't go horizontal first
		result = append(result, slices.Concat(vertical(dr), horizontal(dc)))
	}
	if keypads[kpt][f.r][t.c] != ' ' { // don't go vertical first}
		result = append(result, slices.Concat(horizontal(dc), vertical(dr)))
	}

	return result
}

func useKeypad(level int) keypadType {
	if level == 0 {
		return numerical
	}

	return directional
}

func (m *memo) length(sequence []byte, level, nRobots int) int {
	if level > nRobots {
		return len(sequence)
	}

	keypad := useKeypad(level)

	sum := 0
	prev := byte('A')
	for _, ch := range sequence {
		shortest := math.MaxInt
		for _, option := range keypad.options(prev, ch) {
			newSeq := append(option, 'A')
			xxh := xxhash.Sum64(newSeq)
			s := seq{xxh, level + 1}
			if _, ok := (*m)[s]; !ok {
				(*m)[s] = m.length(append(option, 'A'), level+1, nRobots)
			}

			if (*m)[s] < shortest {
				shortest = (*m)[s]
			}
		}
		sum += shortest
		prev = ch
	}

	return sum
}

func codeToInt(code string) int {
	return conv.MustAtoi(code[:len(code)-1])
}

func (d day21) Part1() int {
	memo := memo{}
	sum := 0
	for _, code := range d.codes {
		length := memo.length([]byte(code), 0, 2)
		sum += codeToInt(code) * length
	}
	return sum
}

func (d day21) Part2() int {
	memo := memo{}
	sum := 0
	for _, code := range d.codes {
		length := memo.length([]byte(code), 0, 25)
		sum += codeToInt(code) * length
	}
	return sum
}

func main() {
	d := NewDay21(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
