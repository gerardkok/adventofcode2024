package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day03b struct {
	day.DayInput
}

func NewDay03b(inputFile string) Day03b {
	return Day03b{day.DayInput(inputFile)}
}

func parseNumberBackwards(data []byte, startAt int) (int, int) {
	result := 0

	pow := 1
	for i := startAt; i >= max(0, startAt-3); i-- {
		c := data[i]
		if c < '0' || c > '9' {
			return result, i
		}
		result += int(c-'0') * pow
		pow *= 10
	}

	return 0, 0
}

func mul(data []byte) int {
	right, i := parseNumberBackwards(data, len(data)-2)

	if data[i] != ',' {
		return 0
	}

	left, i := parseNumberBackwards(data, i-1)

	if i >= 3 && bytes.Equal(data[i-3:i+1], []byte("mul(")) {
		return left * right
	}

	return 0
}

// from https://gist.github.com/guleriagishere/8185da56df6d64c2ab652a59808c1011
func splitAtClosingBracket(data []byte, atEOF bool) (advance int, token []byte, err error) {
	dataLen := len(data)

	// Return Nothing if at the end of file or no data passed.
	if atEOF && dataLen == 0 {
		return 0, nil, nil
	}

	// Find next separator and return token.
	if i := bytes.IndexByte(data, ')'); i >= 0 {
		return i + 1, data[:i+1], nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return dataLen, data, nil
	}

	// Request more data.
	return 0, nil, nil
}

func (d Day03b) Part1() int {
	file, _ := os.Open(string(d.DayInput))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(splitAtClosingBracket)

	sum := 0

	for scanner.Scan() {
		line := scanner.Bytes()

		sum += mul(line)
	}

	return sum
}

func (d Day03b) Part2() int {
	file, _ := os.Open(string(d.DayInput))
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(splitAtClosingBracket)

	sum := 0

	enabled := true

	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.HasSuffix(line, []byte("do()")) {
			enabled = true
		} else if bytes.HasSuffix(line, []byte("don't()")) {
			enabled = false
		} else if enabled {
			sum += mul(line)
		}
	}

	return sum
}

func main() {
	d := NewDay03b(filepath.Join(projectpath.Root, "cmd", "day03b", "input.txt"))

	day.Solve(d)
}
