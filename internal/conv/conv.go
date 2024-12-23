package conv

// convenience

import (
	"iter"
	"strconv"
)

func MustAtoi(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return result
}

func SumFunc[T any](s []T, fn func(T) int) int {
	result := 0

	for _, e := range s {
		result += fn(e)
	}

	return result
}

func Sum(seq iter.Seq[int]) int {
	result := 0

	for i := range seq {
		result += i
	}

	return result
}

func Abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
