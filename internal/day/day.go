package day

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Day interface {
	Part1() int
	Part2() int
}

type DayInput struct {
	Input string
}

type Option func(*DayInput)

func (d DayInput) ReadLines() []string {
	file, err := os.Open(d.Input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	result := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func (d DayInput) ReadInput() []byte {
	input, err := os.ReadFile(d.Input)
	if err != nil {
		log.Fatal(err)
	}

	return input
}

func NewDayInput(path string, opts ...Option) DayInput {
	d := DayInput{
		Input: filepath.Join(path, "input.txt"),
	}
	for _, opt := range opts {
		opt(&d)
	}
	return d
}

func FromArgs(args []string) Option {
	return func(d *DayInput) {
		fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		fset.StringVar(&d.Input, "i", d.Input, "input file")
		if err := fset.Parse(args); err != nil {
			if err == flag.ErrHelp {
				os.Exit(0)
			}
			log.Fatal(err)
		}
	}
}

func WithInput(input string) Option {
	return func(d *DayInput) {
		d.Input = input
	}
}

func Solve(p Day) {
	fmt.Println(p.Part1())
	fmt.Println(p.Part2())
}
