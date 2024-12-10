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

type disk struct {
	blocks []int
	free   int
}

func NewDay09(opts ...day.Option) day09 {
	return day09{day.NewDayInput(path, opts...)}
}

func parseInput(input []byte) disk {
	var blocks []int
	free := 0

	for i, c := range bytes.TrimSpace(input) {
		var id int
		if i%2 == 0 {
			id = i / 2
		} else {
			id = -1
			free += int(c - '0')
		}
		blocks = append(blocks, slices.Repeat([]int{id}, int(c-'0'))...)
	}

	return disk{blocks, free}
}

func (d disk) compact() disk {
	for l, r := 0, len(d.blocks)-1; l < r; l, r = l+1, r-1 {
		// find file
		for d.blocks[r] == -1 {
			r--
		}

		// find free block
		for d.blocks[l] != -1 {
			l++
		}

		d.blocks[l] = d.blocks[r]
	}

	return disk{d.blocks[:len(d.blocks)-d.free], 0}
}

func (d disk) checksum() int {
	result := 0

	for i, b := range d.blocks {
		result += i * b
	}

	return result
}

func (d disk) print() {
	fmt.Printf("free: %d\n", d.free)
	for i, b := range d.blocks {
		fmt.Printf("[%5d] %d\n", i, b)
	}
}

func (d day09) Part1() int {
	input := d.ReadInput()

	disk := parseInput(input)

	compacted := disk.compact()

	return compacted.checksum()
}

func (d day09) Part2() int {
	return 0
}

func main() {
	d := NewDay09(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
