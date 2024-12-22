package main

import (
	"log"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
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
	s ^= (s << 6) % 16777216
	s ^= (s >> 5) % 16777216
	return (s ^ (s << 11)) % 16777216
}

func nextStep(secret, price int) (int, int, int) {
	s := nextSecret(secret)
	p := s % 10
	return s, p, p - price
}

func loop(s, i int) int {
	for range i {
		s = nextSecret(s)
	}
	return s
}

func (d day22b) maxPrice() int {
	seen := make(map[[4]int]int)
	prices := make(map[[4]int]int)

	for i, secret := range d.secrets {
		var changes [4]int
		price := secret % 10

		for j := range 3 {
			secret, price, changes[j] = nextStep(secret, price)
		}

		for range 2000 - 3 {
			changes[0], changes[1], changes[2] = changes[1], changes[2], changes[3]
			secret, price, changes[3] = nextStep(secret, price)

			if seen[changes] != i+1 {
				seen[changes] = i + 1
				prices[changes] += price
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

	f, err := os.Create("myprogram.prof")
	if err != nil {
		log.Fatal(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	return d.maxPrice()
}

func main() {
	d := NewDay22b(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
