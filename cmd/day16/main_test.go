package main

import (
	"adventofcode2024/internal/day"
	"testing"
)

func TestExample1Part1(t *testing.T) {
	t.Parallel()
	d := NewDay16(day.WithInput("example1.txt"))

	want := 7036
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part1(t *testing.T) {
	t.Parallel()
	d := NewDay16(day.WithInput("example2.txt"))

	want := 11048
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample1Part2(t *testing.T) {
	t.Parallel()
	d := NewDay16(day.WithInput("example1.txt"))

	want := 45
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part2(t *testing.T) {
	t.Parallel()
	d := NewDay16(day.WithInput("example2.txt"))

	want := 64
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
