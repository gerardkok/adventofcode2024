package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"adventofcode2024/internal/conv"
	"adventofcode2024/internal/day"
)

var (
	_, caller, _, _ = runtime.Caller(0)
	path            = filepath.Dir(caller)
)

type day17 struct {
	program  []byte
	register map[byte]int
}

func parseProgram(line string) []byte {
	_, p, _ := strings.Cut(line, ": ")

	var program []byte

	for _, c := range []byte(p) {
		if c == ',' {
			continue
		}

		program = append(program, c)
	}

	return program
}

func parseRegisters(lines []string) map[byte]int {
	result := make(map[byte]int)

	for _, line := range lines {
		decl, value, _ := strings.Cut(line, ": ")
		register := decl[len(decl)-1]
		result[register] = conv.MustAtoi(value)
	}

	return result
}

func parseInput(lines []string) ([]byte, map[byte]int) {
	var result [2][]string

	i := 0
	for _, line := range lines {
		if len(line) == 0 {
			i++
			continue
		}

		result[i] = append(result[i], line)
	}

	return parseProgram(result[1][0]), parseRegisters(result[0])
}

func NewDay17(opts ...day.Option) day17 {
	input := day.NewDayInput(path, opts...)

	lines := input.ReadLines()

	program, register := parseInput(lines)

	return day17{program, register}
}

func (d day17) literalOperand(b byte) int {
	return int(b - '0')
}

func (d day17) comboOperand(c byte) int {
	switch c {
	case '0', '1', '2', '3':
		return d.literalOperand(c)
	case '4', '5', '6':
		return d.register['A'-'4'+c]
	default:
		return -1
	}
}

func (d day17) div(operand, register byte) {
	numerator := d.register['A']
	denominator := 1 << d.comboOperand(operand)
	d.register[register] = numerator / denominator
}

func (d day17) adv(operand byte) {
	d.div(operand, 'A')
}

func (d day17) bxl(operand byte) {
	d.register['B'] ^= d.literalOperand(operand)
}

func (d day17) bst(operand byte) {
	d.register['B'] = d.comboOperand(operand) % 8
}

func (d day17) jnz(operand byte) int {
	if d.register['A'] == 0 {
		return -1
	}

	return d.literalOperand(operand)
}

func (d day17) bxc(_ byte) {
	d.register['B'] ^= d.register['C']
}

func (d day17) out(operand byte) byte {
	result := d.comboOperand(operand) % 8
	return byte(result + '0')
}

func (d day17) bdv(operand byte) {
	d.div(operand, 'B')
}

func (d day17) cdv(operand byte) {
	d.div(operand, 'C')
}

func (d day17) execute() []byte {
	pointer := 0
	var result []byte

	for pointer < len(d.program)-1 {
		operand := d.program[pointer+1]
		switch d.program[pointer] {
		case '0':
			d.adv(operand)
			pointer += 2
		case '1':
			d.bxl(operand)
			pointer += 2
		case '2':
			d.bst(operand)
			pointer += 2
		case '3':
			jump := d.jnz(operand)
			if jump == -1 {
				pointer += 2
			} else {
				pointer = d.jnz(operand)
			}
		case '4':
			d.bxc(operand)
			pointer += 2
		case '5':
			result = append(result, d.out(operand))
			pointer += 2
		case '6':
			d.bdv(operand)
			pointer += 2
		case '7':
			d.cdv(operand)
			pointer += 2
		default:
			pointer += 2
		}
	}

	return result
}

func join(b []byte, c byte) []byte {
	return bytes.Join(bytes.Split(b, []byte{}), []byte{c})

}

func tail(b []byte, n int) []byte {
	return b[len(b)-n:]
}

func (d day17) traceBack(pl, a int) int {
	if pl == len(d.program) {
		return a
	}

	a *= 8
	pl++
	for i := a; i < a+1024; i++ {
		d.register['A'] = i
		output := d.execute()

		if len(output) < pl || !bytes.Equal(tail(d.program, pl), tail(output, pl)) {
			continue
		}

		return d.traceBack(pl, i)
	}

	return -1
}

func (d day17) Part1() string {
	output := d.execute()
	return string(join(output, ','))
}

func (d day17) Part2() int {
	a := d.traceBack(0, 0)

	return a
}

func main() {
	d := NewDay17(day.FromArgs(os.Args[1:]))

	fmt.Println(d.Part1())
	fmt.Println(d.Part2())
}
