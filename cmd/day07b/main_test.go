package main

import (
	"adventofcode2024/internal/day"
	"testing"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay07b(day.WithInput("example.txt"))

	want := 3749
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay07b(day.WithInput("example.txt"))

	want := 11387
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
