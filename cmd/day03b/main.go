package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"

	"adventofcode2024/internal/day"
	"adventofcode2024/internal/projectpath"
)

type Day03b struct {
	day.DayInput
}

var (
	mulRE     = regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)$`)
	tokenRE   = regexp.MustCompile(`do\(\)|don't\(\)|mul\(\d{1,3},\d{1,3}\)`)
	enabledRE = regexp.MustCompile(`(?s)do\(\).*?don't\(\)`)
)

func NewDay03b(inputFile string) Day03b {
	return Day03b{day.DayInput(inputFile)}
}

func parseNumberBackwards(data []byte, startAt int) (int, int) {
	result := 0
	for i, j := startAt, 1; i >= 0 && j <= 1000; i, j = i-1, j*10 {
		c := data[i]
		if c < '0' || c > '9' {
			//fmt.Printf("returning number: %d\n", result)
			return result, i
		}
		result += int(c-'0') * j
	}
	return 0, 0
}

func mul(data []byte) int {
	right, i := parseNumberBackwards(data, len(data)-2)

	if data[i] != ',' {
		return 0
	}

	left, i := parseNumberBackwards(data, i-1)

	if i-3 >= 0 && bytes.Equal(data[i-3:i+1], []byte("mul(")) {
		// fmt.Printf("bytes: %s\n", string(data[i-3:i+1]))
		//fmt.Printf("returning %d * %d\n", left, right)
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
	if i := bytes.IndexByte(data, byte(')')); i >= 0 {
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
		// line2 := make([]byte, len(line))
		// copy(line2, line)

		// m1 := 0

		// if match := mulRE.Find(line); match != nil {
		// 	m := string(match[4 : len(match)-1])

		// 	l, r, _ := strings.Cut(m, ",")
		// 	left, _ := strconv.Atoi(l)
		// 	right, _ := strconv.Atoi(r)
		// 	fmt.Printf("m1 left: %d, m1 right: %d\n", left, right)
		// 	m1 = left * right
		// }
		sum += mul(line)
	}

	return sum
}

func (d Day03b) Part2() int {
	file, _ := os.Open(string(d.DayInput))
	defer file.Close()

	// f, _ := os.Create("myprogram.prof")
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

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
			// if match := mulRE.Find(line); match != nil {
			// 	m := string(match[4 : len(match)-1])

			// 	l, r, _ := strings.Cut(m, ",")
			// 	left, _ := strconv.Atoi(l)
			// 	right, _ := strconv.Atoi(r)
			// 	sum += left * right
			// }
			sum += mul(line)
		}
	}

	// 	switch {
	// 	case match == "do()":
	// 		enabled = true
	// 	case match == "don't()":
	// 		enabled = false
	// 	case enabled:
	// 		m := match[4 : len(match)-1]

	// 		l, r, _ := strings.Cut(m, ",")
	// 		left, _ := strconv.Atoi(l)
	// 		right, _ := strconv.Atoi(r)
	// 		sum += left * right
	// 	}
	// }

	return sum
}

func main() {
	d := NewDay03b(filepath.Join(projectpath.Root, "cmd", "day03b", "aoc-2024-day-03-challenge-2.txt"))

	day.Solve(d)
}
