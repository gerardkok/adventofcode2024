package main

import (
	"path/filepath"
	"testing"

	"adventofcode2024/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay03b(filepath.Join(projectpath.Root, "cmd", "day03b", "example1.txt"))

	want := 161
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay03b(filepath.Join(projectpath.Root, "cmd", "day03b", "example2.txt"))

	want := 48
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
