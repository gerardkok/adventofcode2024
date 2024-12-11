package main

import (
	"bytes"
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

type file struct {
	start, length int
	id            int
}

type disk struct {
	blocks    []block
	files     []file
	freeSpace []file
}

func NewDay09(opts ...day.Option) day09 {
	return day09{day.NewDayInput(path, opts...)}
}

func isFree(i int) bool {
	return i%2 == 1
}

func id(i int) int {
	if isFree(i) {
		return 0
	}

	return i / 2
}

func parseInput(input []byte) disk {
	var (
		blocks    []block
		files     []file
		freeSpace []file
	)

	for i, c := range bytes.TrimSpace(input) {
		id := id(i)
		free := isFree(i)
		length := int(c - '0')
		f := file{len(blocks), length, id}
		if free {
			freeSpace = append(freeSpace, f)
		} else {
			files = append([]file{f}, files...) // reversed
		}

		b := block{id, free}
		blocks = append(blocks, slices.Repeat([]block{b}, length)...)
	}

	return disk{blocks, files, freeSpace}
}

func (d disk) blockCompact() {
	for l, r := 0, len(d.blocks)-1; ; l, r = l+1, r-1 {
		// find file
		for d.blocks[r].free {
			r--
		}

		// find free block
		for !d.blocks[l].free {
			l++
		}

		if l >= r {
			break
		}

		d.blocks[l], d.blocks[r] = d.blocks[r], d.blocks[l]
	}
}

func (d disk) findFreeSpace(f file) int {
	for i, s := range d.freeSpace {
		if s.start >= f.start {
			return -1
		}

		if s.length >= f.length {
			return i
		}
	}

	return -1
}

func (d disk) fileCompact() {
	for _, f := range d.files {
		s := d.findFreeSpace(f)

		if s == -1 {
			continue
		}

		for i := 0; i < f.length; i++ {
			d.blocks[d.freeSpace[s].start+i], d.blocks[f.start+i] = d.blocks[f.start+i], d.blocks[d.freeSpace[s].start+i]
		}
		d.freeSpace[s].start += f.length
		d.freeSpace[s].length -= f.length
	}
}

func (d disk) checksum() int {
	result := 0

	for i, b := range d.blocks {
		result += i * b.id
	}

	return result
}

func (d day09) Part1() int {
	input := d.ReadInput()

	disk := parseInput(input)

	disk.blockCompact()

	return disk.checksum()
}

func (d day09) Part2() int {
	input := d.ReadInput()

	disk := parseInput(input)

	disk.fileCompact()

	return disk.checksum()
}

func main() {
	d := NewDay09(day.FromArgs(os.Args[1:]))

	day.Solve(d)
}
