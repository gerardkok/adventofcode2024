package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
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
	switch d.comboOperand(operand) % 8 {
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	case 3:
		return '3'
	case 4:
		return '4'
	case 5:
		return '5'
	case 6:
		return '6'
	case 7:
		return '7'
	default:
		return '8'
	}
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
		//fmt.Printf("A: %d, B: %d, C: %d, operator: %c, combo-operand: %d\n", d.register['A'], d.register['B'], d.register['C'], d.program[pointer], d.comboOperand(operand))
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

func print(output []int) string {
	r := make([]string, len(output))

	for i, o := range output {
		r[i] = strconv.Itoa(o)
	}
	return strings.Join(r, ",")

}

func octalToDecimal(n int) int {
	dec := 0

	// Initializing base value to 1, i.e 8^0
	base := 1

	for n > 0 {
		lastDigit := n % 10
		n /= 10
		dec += lastDigit * base
		base *= 8
	}

	return dec
}

func (d day17) equalTail(program []byte) bool {
	for i, j := len(program)-1, len(d.program)-1; i >= 0; i, j = i-1, j-1 {
		if program[i] != d.program[j] {
			return false
		}
	}

	return true
}

func (d day17) recurPart2(program []byte, a int) []byte {
	fmt.Printf("recur(%s, %d)\n", string(program), a)
	if slices.Equal(d.program, program) {
		return program
	}

	fmt.Printf("tail: %s, program: %s\n", string(d.program[len(d.program)-len(program):]), string(program))
	a *= 8
	l := len(program) + 1
	fmt.Printf("l: %d\n", l)
	for i := a; i < a+1024; i++ {
		d.register['A'] = i
		d.register['B'] = 0
		d.register['C'] = 0
		output := d.execute()
		fmt.Printf("output: %s\n", string(output))
		fmt.Printf("len output: %d\n", len(output))
		if len(output) < l {
			continue
		}
		p := output[len(output)-l:]
		fmt.Printf("tail: %v, p: %v\n", d.program[len(d.program)-l:], p)
		if !d.equalTail(p) {
			fmt.Println("not equal")
			continue
		}
		fmt.Println("going recur")
		if d.recurPart2(p, i) != nil {
			return p
		}
		fmt.Println("---")
	}

	return nil
}

func (d day17) Part1() string {
	output := d.execute()
	return string(output)
}

func (d day17) Part2() string {
	program := d.recurPart2([]byte{}, 0)

	println(string(program))

	return ""
}

func main() {
	d := NewDay17(day.FromArgs(os.Args[1:]))

	//fmt.Println(d.Part1())
	fmt.Println(d.Part2())
}
