package main

import (
	"bytes"
	"fmt"
	"iter"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
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

	numericalGap   = key{3, 0}
	directionalGap = key{0, 0}

	moves = map[int]map[int][]string{
		-3: {-2: {"<<^^^", "^^^<<"}, -1: {"<^^^", "^^^<"}, 0: {"^^^"}, 1: {">^^^", "^^^>"}, 2: {">>^^^", "^^^>>"}},
		-2: {-2: {"<<^^", "^^<<"}, -1: {"<^^", "^^<"}, 0: {"^^"}, 1: {">^^", "^^>"}, 2: {">>^^", "^^>>"}},
		-1: {-2: {"<<^", "^<<"}, -1: {"<^", "^<"}, 0: {"^"}, 1: {">^", "^>"}, 2: {">>^", "^>>"}},
		0:  {-2: {"<<"}, -1: {"<"}, 0: {""}, 1: {">"}, 2: {">>"}},
		1:  {-2: {"<<v", "v<<"}, -1: {"<v", "v<"}, 0: {"v"}, 1: {">v", "v>"}, 2: {">>v", "v>>"}},
		2:  {-2: {"<<vv", "vv<<"}, -1: {"<vv", "vv<"}, 0: {"vv"}, 1: {">vv", "vv>"}, 2: {">>vv", "vv>>"}},
		3:  {-2: {"<<vvv", "vvv<<"}, -1: {"<vvv", "vvv<"}, 0: {"vvv"}, 1: {">vvv", "vvv>"}, 2: {">>vvv", "vvv>>"}},
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

func (k keypad) keys() iter.Seq[key] {
	return func(yield func(key) bool) {
		for r, row := range k {
			for c, ch := range row {
				if ch == ' ' {
					continue
				}

				if !yield(key{r, c}) {
					return
				}
			}
		}
	}
}

func (kpt keypadType) overGap(from, to key) bool {
	if kpt == numerical {
		return key{from.r, to.c} == numericalGap || key{to.r, from.c} == numericalGap
	}

	return key{from.r, to.c} == directionalGap || key{to.r, from.c} == directionalGap
}

// func (kpt keypadType) optionsMap() map[move][]string {
// 	k := keypads[kpt]

// 	result := make(map[move][]string)

// 	for f := range k.keys() {
// 		for t := range k.keys() {
// 			if f == t {
// 				continue
// 			}

// 			from := k[f.r][f.c]
// 			to := k[t.r][t.c]

// 			move := move{from, to}

// 			dr := t.r - f.r
// 			dc := t.c - f.c

// 			result[move] = moves[dr][dc]

// 			// for i, opt := range moves[dr][dc] {
// 			// 	if kpt.overGap(f, t) && i == 0 {
// 			// 		continue
// 			// 	}

// 			// 	result[move] = append(result[move], opt)
// 			// }
// 		}
// 	}

// 	return result
// }

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

	// fmt.Printf("from: %c, to: %c\n", from, to)
	// for i, opt := range result {
	// 	fmt.Printf("[%d] %s\n", i, string(opt))
	// }

	return result
}

func useKeypad(level int) keypadType {
	if level == 0 {
		return numerical
	}

	return directional
}

func (d day21) length(sequence []byte, level, nRobots int) string {
	//fmt.Printf("sequence: %s, level: %d\n", string(sequence), level)
	if level > nRobots {
		// human operator
		//fmt.Printf("human operator: %s\n", string(sequence))
		return string(sequence)
	}

	keypad := useKeypad(level)

	seq := ""
	prev := byte('A')
	for _, ch := range sequence {
		//fmt.Printf("considering char: %c\n", ch)
		shortest := math.MaxInt
		shortestOpt := ""
		for _, option := range keypad.options(prev, ch) {
			optA := append(option, 'A')
			//fmt.Printf("optA: %s\n", string(optA))
			optSeq := d.length(optA, level+1, nRobots)
			l := len(optSeq)
			if l < shortest {
				shortest = l
				shortestOpt = optSeq
			}
		}
		seq += shortestOpt
		prev = ch
	}

	//fmt.Printf("seq found: %s\n", string(seq))
	return seq
}

func codeToInt(code string) int {
	return conv.MustAtoi(code[:len(code)-1])
}

func (d day21) Part1() int {
	sum := 0
	for _, code := range d.codes {
		//		fmt.Printf("code: %s\n", code)
		l := d.length([]byte(code), 0, 2)
		fmt.Printf("code: %d, len: %d\n", codeToInt(code), len(l))
		sum += codeToInt(code) * len(l)
		//		fmt.Println(l)
		//		fmt.Printf("len: %d\n", len(l))
	}
	return sum
}

func (d day21) Part2() int {
	sum := 0
	for _, code := range d.codes {
		l := d.length([]byte(code), 0, 25)
		fmt.Printf("code: %d, len: %d\n", codeToInt(code), len(l))
		sum += codeToInt(code) * len(l)
	}
	return sum
}

func main() {
	d := NewDay21(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
