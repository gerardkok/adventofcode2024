package main

import (
	"adventofcode2024/internal/day"
	"testing"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay03(day.WithInput("example1.txt"))

	want := 161
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay03(day.WithInput("example2.txt"))

	want := 48
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
