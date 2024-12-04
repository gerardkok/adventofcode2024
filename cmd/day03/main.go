package main

import (
	"path/filepath"
	"regexp"
	"strconv"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day03 struct {
	day.DayInput
}

var (
	mulRE     = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	enabledRE = regexp.MustCompile(`(?s)do\(\).*?don't\(\)`)
)

func NewDay03(inputFile string) Day03 {
	return Day03{day.DayInput(inputFile)}
}

func sumMuls(s string) int {
	result := 0

	muls := mulRE.FindAllStringSubmatch(s, -1)
	for _, mul := range muls {
		left, _ := strconv.Atoi(mul[1])
		right, _ := strconv.Atoi(mul[2])
		result += left * right
	}

	return result
}

func (d Day03) Part1() int {
	input, _ := d.ReadFile()

	return sumMuls(string(input))
}

func (d Day03) Part2() int {
	input, _ := d.ReadFile()
	enabledMuls := enabledRE.FindAllString("do()"+string(input)+"don't()", -1)

	result := 0

	for _, muls := range enabledMuls {
		result += sumMuls(muls)
	}

	return result
}

func main() {
	d := NewDay03(filepath.Join(projectpath.Root, "cmd", "day03", "input.txt"))

	day.Solve(d)
}
