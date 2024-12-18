package main

import (
	"adventofcode2024/internal/day"
	"testing"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay17(day.WithInput("example.txt"))

	want := "4,6,3,5,6,3,5,2,1,0"
	got := d.Part1()
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func Test1Part1(t *testing.T) {
	t.Parallel()
	d := day17{[]byte{'2', '6'}, map[byte]int{'A': 0, 'B': 0, 'C': 9}}

	d.execute()
	want := 1
	if d.register['B'] != want {
		t.Errorf("want %d, got %d", want, d.register['B'])
	}
}

func Test2Part1(t *testing.T) {
	t.Parallel()
	d := day17{[]byte{'5', '0', '5', '1', '5', '4'}, map[byte]int{'A': 10, 'B': 0, 'C': 0}}

	want := "0,1,2"
	got := d.Part1()
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func Test3Part1(t *testing.T) {
	t.Parallel()
	d := day17{[]byte{'0', '1', '5', '4', '3', '0'}, map[byte]int{'A': 2024, 'B': 0, 'C': 0}}

	want := "4,2,5,6,7,7,7,7,3,1,0"
	got := d.Part1()
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}

	wantA := 0
	if d.register['A'] != wantA {
		t.Errorf("want %d, got %d", wantA, d.register['A'])
	}
}

func Test4Part1(t *testing.T) {
	t.Parallel()
	d := day17{[]byte{'1', '7'}, map[byte]int{'A': 0, 'B': 29, 'C': 0}}

	d.execute()
	want := 26
	if d.register['B'] != want {
		t.Errorf("want %d, got %d", want, d.register['B'])
	}
}

func Test5Part1(t *testing.T) {
	t.Parallel()
	d := day17{[]byte{'4', '0'}, map[byte]int{'A': 0, 'B': 2024, 'C': 43690}}

	d.execute()
	want := 44354
	if d.register['B'] != want {
		t.Errorf("want %d, got %d", want, d.register['B'])
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay17(day.WithInput("example.txt"))

	want := ""
	got := d.Part2()
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}
