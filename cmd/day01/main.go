package main

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day01 struct {
	day.DayInput
}

func NewDay01(inputFile string) Day01 {
	return Day01{day.DayInput(inputFile)}
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func readLine(line string) (int, int) {
	fields := strings.Fields(line)
	left, _ := strconv.Atoi(fields[0])
	right, _ := strconv.Atoi(fields[1])
	return left, right
}

func parseInput(lines []string) ([]int, []int) {
	var left, right []int

	for _, line := range lines {
		l, r := readLine(line)
		left = append(left, l)
		right = append(right, r)
	}

	return left, right
}

func histogram(list []int) map[int]int {
	result := make(map[int]int)

	for _, l := range list {
		result[l]++
	}

	return result
}

func (d Day01) Part1() int {
	lines, _ := d.ReadLines()

	left, right := parseInput(lines)

	sort.Ints(left)
	sort.Ints(right)

	sum := 0

	for i, l := range left {
		sum += abs(right[i] - l)
	}

	return sum
}

func (d Day01) Part2() int {
	lines, _ := d.ReadLines()

	l, r := parseInput(lines)

	left := histogram(l)
	right := histogram(r)

	sum := 0

	for k, v := range left {
		sum += k * v * right[k]
	}

	return sum
}

func main() {
	d := NewDay01(filepath.Join(projectpath.Root, "cmd", "day01", "input.txt"))

	day.Solve(d)
}
