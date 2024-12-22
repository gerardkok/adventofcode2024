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

type secretNumber int

type monkey struct {
	secretNumber secretNumber
	price        int
	changes      [4]int
	prices       map[[4]int]int
}

type day22 struct {
	secretNumbers []secretNumber
}

func NewDay22(opts ...day.Option) day22 {
	input := day.NewDayInput(path, opts...)

	lines := input.ReadLines()

	result := make([]secretNumber, len(lines))

	for i, line := range lines {
		result[i] = secretNumber(conv.MustAtoi(line))
	}

	return day22{result}
}

func (s secretNumber) mix(value secretNumber) secretNumber {
	return s ^ value
}

func (s secretNumber) prune() secretNumber {
	return s % 16777216
}

func (s secretNumber) mul(value int) secretNumber {
	return s * secretNumber(value)
}

func (s secretNumber) div(value int) secretNumber {
	return s / secretNumber(value)
}

func (s secretNumber) step1() secretNumber {
	return s.mix(s.mul(64)).prune()
}

func (s secretNumber) step2() secretNumber {
	return s.mix(s.div(32)).prune()
}

func (s secretNumber) step3() secretNumber {
	return s.mix(s.mul(2048)).prune()
}

func (s secretNumber) next() secretNumber {
	return s.step1().step2().step3()
}

func (s secretNumber) loop(i int) secretNumber {
	result := s
	for range i {
		result = result.next()
	}
	return result
}

func (s secretNumber) price() int {
	return int(s) % 10
}

func valid(changes [4]int) bool {
	return changes[0] != -10
}

func (d day22) monkeys() []monkey {
	result := make([]monkey, len(d.secretNumbers))

	for i, secretNumber := range d.secretNumbers {
		result[i] = monkey{secretNumber, secretNumber.price(), [4]int{-10, -10, -10, -10}, make(map[[4]int]int)}
	}

	return result
}

func (m *monkey) next() {
	n := m.secretNumber.next()
	p := n.price()
	change := p - m.price
	m.secretNumber, m.price = n, p
	m.changes[0], m.changes[1], m.changes[2], m.changes[3] = m.changes[1], m.changes[2], m.changes[3], change
	if !valid(m.changes) {
		return
	}

	if _, ok := m.prices[m.changes]; !ok {
		m.prices[m.changes] = m.price
	}
}

func mergePrices(monkeys []monkey) map[[4]int]int {
	result := make(map[[4]int]int)

	for _, m := range monkeys {
		for change, price := range m.prices {
			result[change] += price
		}
	}

	return result
}

func (d day22) maxBananas() int {
	monkeys := d.monkeys()

	for range 2000 {
		for i := range monkeys {
			monkeys[i].next()
		}
	}

	prices := mergePrices(monkeys)

	return slices.Max(slices.Collect(maps.Values(prices)))
}

func (d day22) Part1() int {
	sum := secretNumber(0)

	for _, s := range d.secretNumbers {
		sum += s.loop(2000)
	}

	return int(sum)
}

func (d day22) Part2() int {
	return d.maxBananas()
}

func main() {
	d := NewDay22(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
