package main

import (
	"adventofcode2024/internal/day"
	"testing"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay20(6, day.WithInput("example.txt"))

	want := 16
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
