package main

import (
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day22b struct {
	secrets []int
}

func NewDay22b(opts ...day.Option) day22b {
	input := day.NewDayInput(path, opts...)

	lines := input.ReadLines()

	result := make([]int, len(lines))

	for i, line := range lines {
		result[i] = conv.MustAtoi(line)
	}

	return day22b{result}
}

func nextSecret(s int) int {
	s ^= (s << 6) & 16777215
	s ^= (s >> 5) & 16777215
	return (s ^ (s << 11)) & 16777215
}

func nextStep(secret, price, index int) (int, int, int) {
	s := nextSecret(secret)
	p := s % 10
	change := p - price
	return s, p, (index%(19*19*19))*19 + (change + 9)
}

func loop(s, i int) int {
	for range i {
		s = nextSecret(s)
	}
	return s
}

func (d day22b) maxPrice() int {
	seen := make(map[int]int)
	prices := make(map[int]int)

	for i, secret := range d.secrets {
		price := secret % 10
		index := 0

		for range 3 {
			secret, price, index = nextStep(secret, price, index)
		}

		for range 2000 - 3 {
			secret, price, index = nextStep(secret, price, index)

			if seen[index] != i+1 {
				seen[index] = i + 1
				prices[index] += price
			}
		}
	}

	return slices.Max(slices.Collect(maps.Values(prices)))
}

func (d day22b) Part1() int {
	sum := 0

	for _, s := range d.secrets {
		sum += loop(s, 2000)
	}

	return sum
}

func (d day22b) Part2() int {
	return d.maxPrice()
}

func main() {
	d := NewDay22b(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
