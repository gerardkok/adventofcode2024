package main

import (
	"adventofcode2024/internal/day"
	"testing"
)

func TestExample1Part1(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("example1.txt"))

	want := 140
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part1(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("example2.txt"))

	want := 772
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample3Part1(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("example3.txt"))

	want := 1930
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample1Part2(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("example1.txt"))

	want := 80
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part2(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("example2.txt"))

	want := 436
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample3Part2(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("example3.txt"))

	want := 1206
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample4Part2(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("example4.txt"))

	want := 236
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample5Part2(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("example5.txt"))

	want := 368
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestSmallPart2(t *testing.T) {
	t.Parallel()
	d := NewDay12(day.WithInput("small.txt"))

	// 6 * 3 + 4 * 1 = 22
	// ####
	// #AA#
	// #AB#
	// ####
	want := 22
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
