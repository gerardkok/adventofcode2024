package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day09 struct {
	day.DayInput
}

type block struct {
	id   int
	free bool
}

type disk []block

func NewDay09(opts ...day.Option) day09 {
	return day09{day.NewDayInput(path, opts...)}
}

func parseInput(input []byte) disk {
	var blocks []block

	for i, c := range bytes.TrimSpace(input) {
		id := 0
		if i%2 == 0 {
			id = i / 2
		}

		b := block{id, i%2 == 1}

		blocks = append(blocks, slices.Repeat([]block{b}, int(c-'0'))...)
	}

	return blocks
}

func (d disk) compact() {
	for l, r := 0, len(d)-1; ; l, r = l+1, r-1 {
		// find file
		for d[r].free {
			r--
		}

		// find free block
		for !d[l].free {
			l++
		}

		if l >= r {
			break
		}

		d[l], d[r] = d[r], d[l]
	}
}

func (d disk) checksum() int {
	result := 0

	for i, b := range d {
		result += i * b.id
	}

	return result
}

func (d disk) print() {
	fmt.Printf("free: ?\n")
	for i, b := range d {
		if b.id > 0 {
			fmt.Printf("[%4d] %d\n", i, b.id)
		}
	}
}

func (d day09) Part1() int {
	input := d.ReadInput()

	disk := parseInput(input)

	disk.compact()

	return disk.checksum()
}

func (d day09) Part2() int {
	return 0
}

func main() {
	d := NewDay09(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
