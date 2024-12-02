package main

import (
	"path/filepath"
	"testing"

	"adventofcode2024/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay02(filepath.Join(projectpath.Root, "cmd", "day02", "example.txt"))

	want := 2
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay02(filepath.Join(projectpath.Root, "cmd", "day02", "example.txt"))

	want := 4
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
