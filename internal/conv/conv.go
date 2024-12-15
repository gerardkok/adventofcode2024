package conv

// convenience

import "strconv"

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
