package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day11 struct {
	day.DayInput
}

type stone int

type stones struct {
	count map[stone]int
	memo  map[stone][]stone
}

func NewDay11(opts ...day.Option) day11 {
	return day11{day.NewDayInput(path, opts...)}
}

func parseInput(lines []string) stones {
	// only read first line
	result := make(map[stone]int)

	for _, f := range strings.Fields(lines[0]) {
		result[stone(conv.MustAtoi(f))] += 1
	}

	return stones{result, make(map[stone][]stone)}
}

func (s stone) nDigits() int {
	count := 0
	for s != 0 {
		s /= 10
		count++
	}
	return count
}

func divmod(numerator, denominator int) (int, int) {
	quotient := numerator / denominator
	remainder := numerator % denominator
	return quotient, remainder
}

func (s stone) split() []stone {
	n := s.nDigits() / 2
	pow := 1
	for range n {
		pow *= 10
	}
	q, r := divmod(int(s), pow)
	return []stone{stone(q), stone(r)}
}

func (s stone) transform() []stone {
	switch {
	case s == stone(0):
		return []stone{1}
	case s.nDigits()%2 == 0:
		return s.split()
	default:
		return []stone{s * 2024}
	}
}

func (s stones) blink() stones {
	result := make(map[stone]int)

	for st, c := range s.count {
		if _, ok := s.memo[st]; !ok {
			t := st.transform()
			s.memo[st] = t
		}

		for _, u := range s.memo[st] {
			result[u] += c
		}
	}

	return stones{result, s.memo}
}

func (s stones) length() int {
	l := 0
	for _, n := range s.count {
		l += n
	}
	return l
}

func (d day11) Part1() int {
	input := d.ReadLines()
	stones := parseInput(input)

	for range 25 {
		stones = stones.blink()
	}

	return stones.length()
}

func (d day11) Part2() int {
	input := d.ReadLines()
	stones := parseInput(input)

	for range 75 {
		stones = stones.blink()
	}

	return stones.length()
}

func main() {
	d := NewDay11(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
