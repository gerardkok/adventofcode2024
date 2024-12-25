package main

import (
	"adventofcode2024/internal/day"
	"testing"
)

func TestExample1Part1(t *testing.T) {
	t.Parallel()
	d := NewDay24(day.WithInput("example1.txt"))

	want := 4
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part1(t *testing.T) {
	t.Parallel()
	d := NewDay24(day.WithInput("example2.txt"))

	want := 2024
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay24(day.WithInput("example1.txt"))

	want := 0
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
