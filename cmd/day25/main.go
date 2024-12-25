package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day25 struct {
	locks, keys [][5]int
	height      int
}

func parseInput(input day.DayInput) ([][5]int, [][5]int, int) {
	parts := bytes.Split(input.ReadInput(), []byte{'\n', '\n'})

	var (
		locks, keys [][5]int
		height      int
	)

	for _, part := range parts {
		lines := bytes.Split(bytes.TrimSpace(part), []byte{'\n'})
		var (
			pins   [5]int
			isLock bool
		)
		for i, line := range lines {
			if i == 0 {
				isLock = bytes.Equal(line, []byte("#####"))
				continue
			}
			for j, c := range line {
				if isLock && c == '#' {
					pins[j] = i
				}
				if !isLock && c == '.' {
					pins[j] = i
				}
			}
		}
		height = len(lines) - 1
		if isLock {
			locks = append(locks, pins)
		} else {
			for i, pin := range pins {
				pins[i] = height - pin - 1
			}
			keys = append(keys, pins)
		}
	}

	return locks, keys, height
}

func NewDay25(opts ...day.Option) day25 {
	input := day.NewDayInput(path, opts...)

	locks, keys, height := parseInput(input)

	return day25{locks, keys, height}
}

func fits(lock, key [5]int, height int) bool {
	for i := range 5 {
		if lock[i]+key[i] >= height {
			return false
		}
	}
	return true
}

func (d day25) nFitting() int {
	result := 0
	for _, lock := range d.locks {
		for _, key := range d.keys {
			if fits(lock, key, d.height) {
				result++
			}
		}
	}
	return result
}

func (d day25) Part1() int {
	fmt.Println(d)

	return d.nFitting()
}

func (d day25) Part2() int {
	return 0
}

func main() {
	d := NewDay25(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
