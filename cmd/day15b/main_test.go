package main

import (
	"adventofcode2024/internal/day"
	"testing"
)

func TestLargeExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay15b(day.WithInput("large.txt"))

	want := 10092
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestSmallExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay15b(day.WithInput("small.txt"))

	want := 2028
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestLargeExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay15b(day.WithInput("large.txt"))

	want := 9021
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
